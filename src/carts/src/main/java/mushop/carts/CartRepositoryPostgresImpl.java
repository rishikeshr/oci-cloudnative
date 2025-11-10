package mushop.carts;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

import javax.json.bind.Jsonb;
import javax.json.bind.JsonbBuilder;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;

import io.helidon.config.Config;

/**
 * Implements CartRepository using PostgreSQL with JSONB support.
 * This replaces the Oracle SODA implementation with PostgreSQL native JSONB.
 */
public class CartRepositoryPostgresImpl implements CartRepository {

    /** The name of the backing table */
    private final String tableName;

    /** Pool of reusable database connections */
    protected HikariDataSource dataSource;

    /** Used to automatically convert a Cart object to and from JSON */
    private Jsonb jsonb;

    private final static Logger log = Logger.getLogger(CartService.class.getName());

    public CartRepositoryPostgresImpl(Config config) {
        try {
            String host = config.get("POSTGRES_HOST").asString().orElse("localhost");
            String port = config.get("POSTGRES_PORT").asString().orElse("5432");
            String database = config.get("POSTGRES_DB").asString().orElse("mushop_carts");
            String user = config.get("POSTGRES_USER").asString().orElse("mushop");
            String password = config.get("POSTGRES_PASSWORD").asString().orElse("mushop");

            tableName = config.get("POSTGRES_CARTS_TABLE").asString().orElse("carts");

            String jdbcUrl = String.format("jdbc:postgresql://%s:%s/%s", host, port, database);

            HikariConfig hikariConfig = new HikariConfig();
            hikariConfig.setJdbcUrl(jdbcUrl);
            hikariConfig.setUsername(user);
            hikariConfig.setPassword(password);
            hikariConfig.setMaximumPoolSize(10);
            hikariConfig.setConnectionTimeout(30000);

            dataSource = new HikariDataSource(hikariConfig);

            // Create the carts table if it does not exist
            createTableIfNotExists();

            log.info("Connected to PostgreSQL database: " + database);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        jsonb = JsonbBuilder.create();
    }

    /**
     * Create the carts table if it doesn't exist
     */
    private void createTableIfNotExists() throws SQLException {
        String createTableSQL = String.format(
            "CREATE TABLE IF NOT EXISTS %s (" +
            "  id VARCHAR(255) PRIMARY KEY," +
            "  json_document JSONB NOT NULL," +
            "  version UUID DEFAULT gen_random_uuid()," +
            "  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
            "  created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP" +
            ")", tableName
        );

        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement(createTableSQL)) {
            stmt.execute();
        }

        // Create index on customerId within JSONB for efficient querying
        String createIndexSQL = String.format(
            "CREATE INDEX IF NOT EXISTS idx_%s_customer_id ON %s ((json_document->>'customerId'))",
            tableName, tableName
        );

        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement(createIndexSQL)) {
            stmt.execute();
        }
    }

    @Override
    public boolean deleteCart(String id) {
        String sql = String.format("DELETE FROM %s WHERE id = ?", tableName);
        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement(sql)) {
            stmt.setString(1, id);
            int rowsAffected = stmt.executeUpdate();
            return rowsAffected > 0;
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public Cart getById(String id) {
        String sql = String.format("SELECT json_document FROM %s WHERE id = ?", tableName);
        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement(sql)) {
            stmt.setString(1, id);
            try (ResultSet rs = stmt.executeQuery()) {
                if (rs.next()) {
                    String jsonString = rs.getString("json_document");
                    return jsonb.fromJson(jsonString, Cart.class);
                }
                return null;
            }
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public List<Cart> getByCustomerId(String custId) {
        if (custId == null) {
            throw new IllegalArgumentException("The customer id must be specified");
        }

        String sql = String.format(
            "SELECT json_document FROM %s WHERE json_document->>'customerId' = ?",
            tableName
        );

        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement(sql)) {
            stmt.setString(1, custId);

            List<Cart> result = new ArrayList<>();
            try (ResultSet rs = stmt.executeQuery()) {
                while (rs.next()) {
                    String jsonString = rs.getString("json_document");
                    Cart cart = jsonb.fromJson(jsonString, Cart.class);
                    result.add(cart);
                }
            }
            return result;
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void save(Cart cart) {
        String sql = String.format(
            "INSERT INTO %s (id, json_document, version, last_modified, created_on) " +
            "VALUES (?, ?::jsonb, gen_random_uuid(), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) " +
            "ON CONFLICT (id) DO UPDATE SET " +
            "  json_document = EXCLUDED.json_document, " +
            "  version = gen_random_uuid(), " +
            "  last_modified = CURRENT_TIMESTAMP",
            tableName
        );

        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement(sql)) {
            stmt.setString(1, cart.getId());
            stmt.setString(2, jsonb.toJson(cart));
            stmt.executeUpdate();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public boolean healthCheck() {
        try (Connection con = dataSource.getConnection();
             PreparedStatement stmt = con.prepareStatement("SELECT 1")) {
            try (ResultSet rs = stmt.executeQuery()) {
                return rs.next();
            }
        } catch (SQLException e) {
            log.log(Level.SEVERE, "DB health-check failed.", e);
            return false;
        }
    }
}
