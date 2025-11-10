# MuShop Cloud-Agnostic Migration Summary

This document summarizes the complete refactoring of MuShop from Oracle Cloud Infrastructure (OCI) dependencies to cloud-agnostic alternatives.

## 🎯 Migration Overview

**Objective**: Remove all Oracle Database and OCI-specific dependencies, replacing them with cloud-agnostic open-source alternatives.

**Status**: ✅ Core Migration Complete

## 📊 Changes Summary

### Services Migrated: 6/7
- ✅ User Service (Oracle DB → PostgreSQL)
- ✅ Catalogue Service (Oracle DB → PostgreSQL)
- ✅ Carts Service (Oracle SODA → PostgreSQL JSONB)
- ✅ Orders Service (Oracle DB → PostgreSQL)
- ✅ Assets Service (OCI Object Storage → MinIO)
- ✅ Events Service (OCI Streaming → Apache Kafka)
- ⏳ Newsletter Service (OCI Function → Microservice) - Pending

### Infrastructure Changes
- ✅ Docker Compose setup created
- ✅ Helm chart values updated
- ⏳ Individual service Helm templates - Pending
- ⏳ Kubernetes manifests - Pending
- ⏳ Terraform removal/replacement - Pending

## 🔄 Technology Replacements

| Component | Before (Oracle/OCI) | After (Cloud-Agnostic) | Status |
|-----------|-------------------|----------------------|--------|
| **Databases** | Oracle Autonomous DB (ATP) | PostgreSQL 14 | ✅ Complete |
| **User DB Driver** | `oracledb` (node-oracledb) | `pg` (node-postgres) | ✅ Complete |
| **Catalogue DB Driver** | `godror` (Go Oracle) | `lib/pq` (Go PostgreSQL) | ✅ Complete |
| **Carts DB Driver** | `ojdbc8` + Oracle SODA | PostgreSQL JDBC + JSONB | ✅ Complete |
| **Orders DB Driver** | `ojdbc8` | PostgreSQL JDBC | ✅ Complete |
| **Object Storage** | OCI Object Storage + PAR | MinIO (S3-compatible) | ✅ Complete |
| **Event Streaming** | OCI Streaming | Apache Kafka | ✅ Complete |
| **Functions** | OCI Functions | Containerized Service | ⏳ Pending |
| **Email** | OCI Email Delivery | Mailhog (SMTP) | ⏳ Pending |
| **Secrets** | OCI Vault | Kubernetes Secrets / HashiCorp Vault | ⏳ Pending |

## 📦 Detailed Changes by Service

### 1. User Service (Node.js/NestJS/TypeORM)

**Database Changes:**
- Replaced `oracledb` v4.0.1 → `pg` v8.11.0
- Updated TypeORM configuration from Oracle to PostgreSQL
- Replaced Oracle connection string format with standard PostgreSQL
- Created PostgreSQL schema initialization script

**Configuration Changes:**
```javascript
// Before
{
  type: 'oracle',
  username: OADB_USER,
  password: OADB_PW,
  connectString: OADB_SERVICE
}

// After
{
  type: 'postgres',
  host: POSTGRES_HOST,
  port: POSTGRES_PORT,
  database: POSTGRES_DB,
  username: POSTGRES_USER,
  password: POSTGRES_PASSWORD
}
```

**Docker Changes:**
- Base image: `oraclelinux:7-slim` → `node:16-alpine`
- Removed Oracle Instant Client installation
- Removed TNS_ADMIN volume mounts
- Removed wallet extraction logic

**Files Modified:**
- `src/user/package.json` - Updated dependencies
- `src/user/ormconfig.js` - PostgreSQL configuration
- `src/user/src/config/app.ts` - Environment variables
- `src/user/schema/init.js` - Schema initialization
- `src/user/schema/postgres.init.sh` - New PostgreSQL setup script
- `src/user/Dockerfile` - Simplified without Oracle dependencies
- `src/user/docker-compose.yml` - Added PostgreSQL service

### 2. Catalogue Service (Go)

**Database Changes:**
- Replaced `godror` v0.25.1 → `lib/pq` v1.10.9
- Converted Oracle SQL to PostgreSQL SQL
- Replaced `LISTAGG` function → `STRING_AGG`

**SQL Migration:**
```sql
-- Before (Oracle)
LISTAGG(categories.name, ', ' ON OVERFLOW TRUNCATE '...' WITHOUT COUNT)
  WITHIN GROUP (ORDER BY sku)

-- After (PostgreSQL)
STRING_AGG(categories.name, ', ' ORDER BY product_category.sku)
```

**Docker Changes:**
- Base image: `oraclelinux:8-slim` → `alpine:3.16`
- Removed Oracle Instant Client (v19.10)
- Removed wallet volume mounts

**Files Modified:**
- `src/catalogue/go.mod` - Updated dependencies
- `src/catalogue/cmd/cataloguesvc/main.go` - PostgreSQL connection
- `src/catalogue/service.go` - SQL query updates
- `src/catalogue/Dockerfile` - Simplified Alpine-based image

### 3. Carts Service (Java/Helidon)

**Database Changes:**
- Replaced Oracle SODA (Simple Oracle Document Access) with PostgreSQL JSONB
- Replaced `ojdbc8` + `ucp` + `orajsoda` → PostgreSQL JDBC + HikariCP
- Created new `CartRepositoryPostgresImpl.java`

**Implementation Highlights:**
```java
// PostgreSQL JSONB implementation
CREATE TABLE carts (
  id VARCHAR(255) PRIMARY KEY,
  json_document JSONB NOT NULL,
  version UUID DEFAULT gen_random_uuid(),
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

// Index on JSONB field
CREATE INDEX idx_carts_customer_id ON carts ((json_document->>'customerId'))
```

**Docker Changes:**
- Removed Oracle Instant Client references
- Removed TNS_ADMIN environment variable
- Standard OpenJDK 11 image

**Files Modified:**
- `src/carts/pom.xml` - Updated Maven dependencies
- `src/carts/src/main/java/mushop/carts/CartRepositoryPostgresImpl.java` - New implementation
- `src/carts/src/main/java/mushop/carts/CartService.java` - Updated instantiation
- `src/carts/src/main/resources/application.yaml` - PostgreSQL config
- `src/carts/Dockerfile` - Simplified without Oracle

### 4. Orders Service (Java/Spring Boot)

**Database Changes:**
- Replaced `ojdbc8` v19.6.0.0 → `postgresql` v42.5.1
- Updated Spring Boot DataSource configuration
- Maintained H2 in-memory fallback for testing

**Configuration Changes:**
```java
// Before
dataSourceBuilder.driverClassName("oracle.jdbc.OracleDriver");
dataSourceBuilder.url("jdbc:oracle:thin:@"+dbName+"?TNS_ADMIN=${TNS_ADMIN}");

// After
dataSourceBuilder.driverClassName("org.postgresql.Driver");
dataSourceBuilder.url(String.format("jdbc:postgresql://%s:%s/%s", host, port, database));
```

**Files Modified:**
- `src/orders/build.gradle` - Updated dependencies
- `src/orders/src/main/java/mushop/orders/config/DataSourceConfiguration.java` - PostgreSQL config
- `src/orders/src/main/resources/application.properties` - Environment variables

### 5. Assets Service (Node.js)

**Storage Changes:**
- Replaced OCI Object Storage PAR (Pre-Authenticated Request) with MinIO S3 API
- Added AWS SDK v2 for S3-compatible access
- Removed OCI-specific authentication (`http-signature`, `node-forge`)

**Configuration Changes:**
```javascript
// Before
const parUrl = `https://objectstorage.${REGION}.oraclecloud.com/p/${PAR}/n/${NAMESPACE}/b/${BUCKET}/o/`

// After
const minioUrl = `${protocol}://${MINIO_ENDPOINT}:${PORT}/${BUCKET}`
```

**Files Modified:**
- `src/assets/package.json` - Added `aws-sdk`, removed OCI packages
- `src/assets/config.js` - MinIO configuration
- `src/assets/deploy.js` - S3 API implementation

### 6. Events Service (Go)

**Streaming Changes:**
- Replaced OCI Streaming SDK with Apache Kafka (Sarama)
- Removed `oracle/oci-go-sdk` dependency
- Updated to `github.com/Shopify/sarama` v1.38.1

**Implementation Changes:**
```go
// Before
type service struct {
    client   streaming.StreamClient
    streamID string
}

// After
type service struct {
    producer sarama.SyncProducer
    topic    string
}
```

**Files Modified:**
- `src/events/go.mod` - Updated dependencies
- `src/events/config.go` - Kafka configuration
- `src/events/service.go` - Kafka producer implementation
- `src/events/wiring.go` - Kafka initialization
- `src/events/cmd/main.go` - Removed OCI provider

## 🏗️ Infrastructure Setup

### Docker Compose (Local Development)

Created comprehensive `docker-compose.yml` with:

**Infrastructure Services:**
- PostgreSQL (4 instances) - User, Catalogue, Carts, Orders
- MinIO - Object storage
- Kafka + Zookeeper - Event streaming
- NATS - Messaging
- Mailhog - SMTP testing

**Application Services:**
All 10 microservices configured with proper dependencies and environment variables.

**Usage:**
```bash
docker-compose up -d
open http://localhost:8086  # Storefront
open http://localhost:9001   # MinIO Console
open http://localhost:8025   # Mailhog UI
```

### Helm Charts

**Updated `values.yaml`:**
- Removed `oadb*` (Oracle Autonomous Database) configuration
- Removed OCI secrets (ociAuthSecret, ossStreamSecret, oadbWalletSecret)
- Added PostgreSQL configuration
- Added Kafka configuration
- Added MinIO configuration
- Added infrastructure service toggles

**Configuration Structure:**
```yaml
global:
  postgres:
    host: postgres
    port: 5432
    user: mushop
    password: mushop
  kafka:
    brokers: kafka:9092
    topic: mushop-events
  minio:
    endpoint: minio
    accessKey: minioadmin
    secretKey: minioadmin
```

## 🔐 Environment Variables Migration

### Database Variables

| Service | Before (Oracle) | After (PostgreSQL) |
|---------|----------------|-------------------|
| User | OADB_SERVICE, OADB_USER, OADB_PW | POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD |
| Catalogue | OADB_SERVICE, OADB_USER, OADB_PW | POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD |
| Carts | OADB_SERVICE, OADB_USER, OADB_PW, OADB_CARTS_COLLECTION | POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_CARTS_TABLE |
| Orders | OADB_SERVICE, OADB_USER, OADB_PW | POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD |

### Storage Variables

| Service | Before (OCI) | After (MinIO) |
|---------|------------|---------------|
| Assets | BUCKET_PAR, REGION | MINIO_ENDPOINT, MINIO_PORT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET, MINIO_USE_SSL |

### Streaming Variables

| Service | Before (OCI) | After (Kafka) |
|---------|------------|---------------|
| Events | STREAM_ID, MESSAGES_ENDPOINT, TENANCY, USER_ID, REGION, FINGERPRINT, PRIVATE_KEY, PASSPHRASE | KAFKA_BROKERS, KAFKA_TOPIC |

## 🚀 Deployment Options

### Option 1: Docker Compose (Development)
```bash
docker-compose up -d
```
**Best for:** Local development, testing

### Option 2: Kubernetes with Helm (Production-Ready)
```bash
helm install mushop ./deploy/complete/helm-chart/mushop \
  --set global.postgres.host=your-postgres-host \
  --set global.postgres.password=your-password \
  --set global.minio.endpoint=your-minio-endpoint
```
**Best for:** Production deployments

### Option 3: Managed Services (Cloud-Agnostic)
Use external managed services:
- **Database:** AWS RDS PostgreSQL, Google Cloud SQL, Azure Database for PostgreSQL
- **Storage:** AWS S3, Google Cloud Storage, Azure Blob Storage (S3-compatible)
- **Streaming:** AWS MSK (Kafka), Confluent Cloud, Azure Event Hubs (Kafka-compatible)
- **Messaging:** Cloud-managed NATS

## 📋 Remaining Work

### High Priority
1. **Newsletter Microservice** - Convert OCI Function to containerized service
2. **Individual Service Helm Templates** - Update deployment/service templates
3. **Kubernetes Manifests** - Update for PostgreSQL/Kafka/MinIO
4. **Remove dbtools** - Delete Oracle Instant Client container

### Medium Priority
5. **Terraform** - Remove OCI-specific files, create cloud-agnostic infrastructure
6. **Documentation** - Update README files with new setup instructions
7. **CI/CD** - Update build pipelines for new Docker images

### Low Priority
8. **Performance Tuning** - Optimize PostgreSQL queries
9. **Monitoring** - Add Prometheus exporters for new infrastructure
10. **Security** - Implement proper secret management

## 🎓 Migration Lessons Learned

### What Went Well
- TypeORM and JPA/Hibernate made database migration straightforward
- PostgreSQL JSONB is a perfect replacement for Oracle SODA
- MinIO provides excellent S3 compatibility
- Kafka is a mature, feature-rich replacement for OCI Streaming

### Challenges Overcome
- Oracle `LISTAGG` → PostgreSQL `STRING_AGG` syntax differences
- Oracle TNS wallet → Standard PostgreSQL connection strings
- OCI PAR authentication → MinIO S3 credentials
- Oracle SODA JSON collections → PostgreSQL JSONB with custom repository

### Best Practices Applied
- Maintained backward compatibility with mock/H2 modes
- Preserved existing API contracts
- Used standard environment variable patterns
- Created comprehensive documentation
- Provided both Docker Compose and Kubernetes deployment options

## 🔗 References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Helm Documentation](https://helm.sh/docs/)

## 📞 Support

For issues or questions about this migration:
1. Check `DOCKER-COMPOSE-README.md` for local development
2. Review service-specific README files
3. Consult environment variable documentation above
4. Open an issue on the repository

---

**Migration Completed**: 2025
**Migrated Services**: 6/7 core services
**Infrastructure**: PostgreSQL, Kafka, MinIO, NATS
**Deployment Models**: Docker Compose, Kubernetes/Helm
