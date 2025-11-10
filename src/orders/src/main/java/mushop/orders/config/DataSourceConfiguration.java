
/**
 ** Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 ** Licensed under the Universal Permissive License v 1.0 as shown at http://oss.oracle.com/licenses/upl.
 **/
package  mushop.orders.config;

import com.zaxxer.hikari.HikariDataSource;
import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.actuate.metrics.jdbc.DataSourcePoolMetrics;
import org.springframework.boot.jdbc.DataSourceBuilder;
import org.springframework.boot.jdbc.metadata.DataSourcePoolMetadataProvider;
import org.springframework.boot.jdbc.metadata.HikariDataSourcePoolMetadata;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;

import javax.sql.DataSource;

@Configuration
public class DataSourceConfiguration {


    @Autowired
    private Environment env;

    @Value("${mushop.orders.postgres.host}")
    private String db_host;

    @Value("${mushop.orders.postgres.port}")
    private String db_port;

    @Value("${mushop.orders.postgres.database}")
    private String db_database;

    @Value("${mushop.orders.postgres.user}")
    private String db_user;

    @Value("${mushop.orders.postgres.password}")
    private String db_pass;

    @Bean
    public DataSource getDataSource(MeterRegistry registry) {
        DataSourceBuilder dataSourceBuilder = DataSourceBuilder.create();
        //
        if("mock".equalsIgnoreCase(db_host)) {
            dataSourceBuilder.driverClassName("org.h2.Driver");
            dataSourceBuilder.url("jdbc:h2:mem:test");
            dataSourceBuilder.username("SA");
            dataSourceBuilder.password("");
        }else{
            dataSourceBuilder.driverClassName("org.postgresql.Driver");
            dataSourceBuilder.url(String.format("jdbc:postgresql://%s:%s/%s", db_host, db_port, db_database));
            dataSourceBuilder.username(db_user);
            dataSourceBuilder.password(db_pass);
        }
        DataSource ordersDS = dataSourceBuilder.build();
        DataSourcePoolMetadataProvider provider = dataSource -> new HikariDataSourcePoolMetadata((HikariDataSource) dataSource);
        new DataSourcePoolMetrics(ordersDS,provider,"orders_data_source",null).bindTo(registry);
        return ordersDS;
    }
}