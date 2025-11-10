# MuShop Local Development with Docker Compose

This Docker Compose configuration runs the complete MuShop application locally without any Oracle or OCI dependencies.

## Architecture

### Infrastructure Services
- **PostgreSQL** (4 instances) - Separate databases for User, Catalogue, Carts, and Orders services
- **MinIO** - S3-compatible object storage for assets
- **Kafka + Zookeeper** - Event streaming platform
- **NATS** - Messaging for Orders/Fulfillment communication
- **Mailhog** - SMTP server for email testing

### Application Services
- **User** - User authentication and profiles (Port 3001)
- **Catalogue** - Product catalog (Port 8080)
- **Carts** - Shopping cart with JSONB storage (Port 8081)
- **Orders** - Order processing (Port 8082)
- **Payment** - Payment processing (Port 8083)
- **Fulfillment** - Order fulfillment
- **Events** - Event streaming integration (Port 8084)
- **Assets** - Static asset serving (Port 8085)
- **API** - Backend for Frontend / API Gateway (Port 3000)
- **Storefront** - Frontend web application (Port 8086)

## Prerequisites

- Docker 20.10 or higher
- Docker Compose 2.0 or higher
- At least 8GB of available RAM
- At least 20GB of available disk space

## Quick Start

### 1. Start All Services

```bash
docker-compose up -d
```

This will start all infrastructure and application services in the background.

### 2. Verify Services Are Running

```bash
docker-compose ps
```

All services should show as "Up" or "healthy".

### 3. Initialize Databases (First Time Only)

The PostgreSQL databases will be created automatically. To initialize schemas:

```bash
# User service schema
docker-compose exec user npm run schema:sync

# Other services use JPA/Hibernate auto-DDL
```

### 4. Access the Application

- **Storefront**: http://localhost:8086
- **API Gateway**: http://localhost:3000
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)
- **Mailhog UI**: http://localhost:8025
- **NATS Monitoring**: http://localhost:8222

## Service-Specific Access

| Service | Port | URL |
|---------|------|-----|
| Storefront | 8086 | http://localhost:8086 |
| API (BFF) | 3000 | http://localhost:3000 |
| User | 3001 | http://localhost:3001 |
| Catalogue | 8080 | http://localhost:8080 |
| Carts | 8081 | http://localhost:8081 |
| Orders | 8082 | http://localhost:8082 |
| Payment | 8083 | http://localhost:8083 |
| Events | 8084 | http://localhost:8084 |
| Assets | 8085 | http://localhost:8085 |

## Database Access

### PostgreSQL Connections

```bash
# User DB
psql -h localhost -p 5432 -U mushop -d mushop_user

# Catalogue DB
psql -h localhost -p 5433 -U mushop -d mushop_catalogue

# Carts DB
psql -h localhost -p 5434 -U mushop -d mushop_carts

# Orders DB
psql -h localhost -p 5435 -U mushop -d mushop_orders
```

Password for all: `mushop`

## MinIO Object Storage

### Access Credentials
- **Endpoint**: http://localhost:9000
- **Access Key**: minioadmin
- **Secret Key**: minioadmin
- **Bucket**: mushop-assets

### Upload Assets

```bash
# From the assets service directory
cd src/assets
npm run build  # Optimize images
npm run deploy # Upload to MinIO
```

## Kafka Topics

The Events service uses Kafka for event streaming:

- **Topic**: mushop-events
- **Brokers**: localhost:29092 (external), kafka:9092 (internal)

### View Kafka Topics

```bash
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Consume Events

```bash
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic mushop-events \
  --from-beginning
```

## Common Operations

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f user
docker-compose logs -f catalogue
```

### Restart a Service

```bash
docker-compose restart user
```

### Rebuild and Restart

```bash
# Rebuild a specific service
docker-compose up -d --build user

# Rebuild all services
docker-compose up -d --build
```

### Stop All Services

```bash
docker-compose down
```

### Stop and Remove Volumes (Clean Start)

```bash
docker-compose down -v
```

## Development Workflow

### 1. Code Changes

After making code changes to a service:

```bash
# Rebuild and restart the service
docker-compose up -d --build <service-name>
```

### 2. Database Migrations

```bash
# User service (TypeORM)
docker-compose exec user npm run schema:sync

# Java services use auto-DDL (no manual migration needed)
```

### 3. View Application Logs

```bash
# Follow logs in real-time
docker-compose logs -f api user catalogue carts orders
```

## Troubleshooting

### Service Won't Start

```bash
# Check service logs
docker-compose logs <service-name>

# Check if ports are already in use
netstat -an | grep LISTEN | grep <port>
```

### Database Connection Issues

```bash
# Verify PostgreSQL is running
docker-compose ps postgres-user

# Check PostgreSQL logs
docker-compose logs postgres-user

# Test connection
docker-compose exec postgres-user psql -U mushop -d mushop_user -c "SELECT 1;"
```

### Kafka Connection Issues

```bash
# Verify Kafka and Zookeeper are running
docker-compose ps kafka zookeeper

# Check Kafka logs
docker-compose logs kafka

# List topics
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092
```

### MinIO Connection Issues

```bash
# Verify MinIO is running
docker-compose ps minio

# Check MinIO logs
docker-compose logs minio

# Access MinIO console
open http://localhost:9001
```

### Reset Everything

```bash
# Stop all services and remove volumes
docker-compose down -v

# Remove all images
docker-compose down --rmi all

# Start fresh
docker-compose up -d
```

## Environment Variables

Override defaults by creating a `.env` file:

```env
# PostgreSQL
POSTGRES_USER=mushop
POSTGRES_PASSWORD=your-secure-password

# MinIO
MINIO_ROOT_USER=your-access-key
MINIO_ROOT_PASSWORD=your-secret-key

# Kafka
KAFKA_BROKER_ID=1
```

## Performance Tuning

### Increase Resource Limits

Edit `docker-compose.yml` to add resource limits:

```yaml
services:
  postgres-user:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

### Use BuildKit for Faster Builds

```bash
DOCKER_BUILDKIT=1 docker-compose build
```

## Production Notes

⚠️ **This configuration is for local development only!**

For production:
- Use managed PostgreSQL (RDS, Cloud SQL, Azure Database)
- Use managed Kafka (Confluent Cloud, Amazon MSK)
- Use S3/GCS/Azure Blob instead of MinIO
- Configure proper security (TLS, authentication)
- Use Kubernetes for orchestration
- Implement monitoring and observability

## Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
