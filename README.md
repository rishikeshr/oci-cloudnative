# ![MuShop Logo](./images/logo.png#gh-light-mode-only)![MuShop Logo - Dark Mode](./images/logo-inverse.png#gh-dark-mode-only)

MuShop is a complete e-commerce platform built as a set of microservices, demonstrating modern cloud-native application development practices. The application has been refactored to be **cloud-agnostic**, using open-source technologies instead of proprietary cloud services.

| ![home](./images/screenshot/mushop.home.png) | ![browse](./images/screenshot/mushop.browse.png) | ![cart](./images/screenshot/mushop.cart.png) | ![about](./images/screenshot/mushop.about.png) |
|---|---|---|---|

## 🚀 Quick Start - Cloud-Agnostic Deployment

**New!** MuShop now runs on any cloud platform or locally using standard open-source technologies:

- **PostgreSQL** instead of Oracle Autonomous Database
- **MinIO** instead of OCI Object Storage
- **Apache Kafka** instead of OCI Streaming
- **NATS** for messaging
- **Mailhog** for SMTP (development)

### Deploy Locally with Docker Compose

The easiest way to run MuShop locally is using Docker Compose, which sets up the complete application stack including all microservices and infrastructure components.

#### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) 20.10 or later
- [Docker Compose](https://docs.docker.com/compose/install/) 1.29 or later
- At least 8GB of RAM allocated to Docker
- 10GB of free disk space

#### Quick Start

```bash
# Clone the repository
git clone https://github.com/rishikeshr/oci-cloudnative.git
cd oci-cloudnative

# Start all services (builds images and starts containers)
docker-compose up -d

# Wait for services to be ready (may take 2-3 minutes)
docker-compose ps

# Access the application
open http://localhost:8086
```

#### What Gets Deployed

The Docker Compose setup includes:

**Application Microservices:**
- **Storefront** (http://localhost:8086) - React-based web UI
- **API Gateway** (http://localhost:8080) - NGINX reverse proxy
- **User Service** (http://localhost:3000) - User authentication & management
- **Catalogue Service** (http://localhost:8081) - Product catalog
- **Carts Service** (http://localhost:8082) - Shopping cart management
- **Orders Service** (http://localhost:8083) - Order processing
- **Payment Service** (http://localhost:8084) - Payment processing
- **Fulfillment Service** (http://localhost:8085) - Order fulfillment
- **Assets Service** (http://localhost:3001) - Product images & static assets
- **Events Service** (http://localhost:8087) - Event streaming integration
- **Newsletter Service** (http://localhost:3002) - Newsletter subscriptions

**Infrastructure Services:**
- **PostgreSQL** (ports 5432-5435) - 4 database instances for different services
- **MinIO** (http://localhost:9000, Console: http://localhost:9001) - S3-compatible object storage
- **Apache Kafka** (localhost:9092) - Event streaming
- **Zookeeper** (localhost:2181) - Kafka coordination
- **NATS** (localhost:4222) - Messaging for orders/fulfillment
- **Mailhog** (UI: http://localhost:8025, SMTP: 1025) - Email testing

#### Verify Deployment

```bash
# Check all services are running
docker-compose ps

# View logs for all services
docker-compose logs -f

# View logs for a specific service
docker-compose logs -f storefront

# Check service health
curl http://localhost:8086  # Storefront
curl http://localhost:3000/health  # User service
curl http://localhost:3002/health  # Newsletter service
```

#### Using the Application

1. **Browse Products**: Navigate to http://localhost:8086
2. **Create Account**: Click "Sign In" → "Register" to create a user
3. **Add to Cart**: Browse products and add items to your cart
4. **Checkout**: Complete the purchase flow
5. **Subscribe to Newsletter**: Enter email in the newsletter form
6. **View Emails**: Check http://localhost:8025 (Mailhog) to see subscription emails

#### Accessing Infrastructure Components

**PostgreSQL Databases:**
```bash
# User database
psql -h localhost -p 5432 -U mushop -d mushop_user

# Catalogue database
psql -h localhost -p 5433 -U mushop -d mushop_catalogue

# Carts database
psql -h localhost -p 5434 -U mushop -d mushop_carts

# Orders database
psql -h localhost -p 5435 -U mushop -d mushop_orders

# Password for all: mushop
```

**MinIO Object Storage:**
```bash
# Access MinIO Console
open http://localhost:9001

# Credentials:
# Username: minioadmin
# Password: minioadmin
```

**Mailhog (Email Testing):**
```bash
# View sent emails
open http://localhost:8025
```

**Kafka Topics:**
```bash
# List topics
docker-compose exec kafka kafka-topics \
  --bootstrap-server localhost:9092 --list

# Consume events
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic mushop-events --from-beginning
```

#### Managing the Deployment

```bash
# Stop all services (keeps data)
docker-compose stop

# Start services again
docker-compose start

# Restart a specific service
docker-compose restart storefront

# View resource usage
docker stats

# Rebuild a service after code changes
docker-compose build user
docker-compose up -d user

# Scale a service (e.g., run 3 catalogue instances)
docker-compose up -d --scale catalogue=3
```

#### Cleanup

```bash
# Stop and remove all containers
docker-compose down

# Remove containers, networks, and volumes (complete cleanup)
docker-compose down -v

# Remove all images as well
docker-compose down -v --rmi all
```

#### Troubleshooting

**Services not starting:**
```bash
# Check logs for errors
docker-compose logs

# Restart all services
docker-compose restart

# Rebuild and restart
docker-compose up -d --build
```

**Database connection errors:**
```bash
# Check PostgreSQL is ready
docker-compose exec postgres-user pg_isready -U mushop

# Restart database and dependent service
docker-compose restart postgres-user user
```

**Port conflicts:**
```bash
# If ports are already in use, modify docker-compose.yml
# Change the first port number (host port) to something else
# Example: "8086:80" → "8090:80"
```

**Out of memory:**
```bash
# Increase Docker memory allocation in Docker Desktop preferences
# Recommended: At least 8GB RAM
```

**📖 See [DOCKER-COMPOSE-README.md](./DOCKER-COMPOSE-README.md) for complete local development guide and advanced topics**

### Deploy to Kubernetes

```bash
# Using Helm (recommended)
helm install mushop ./deploy/complete/helm-chart/mushop \
  --set global.postgres.host=your-postgres-host \
  --set global.postgres.password=your-password

# Or using kubectl
kubectl apply -f deploy/complete/kubernetes/
```

**📖 See [MIGRATION-SUMMARY.md](./MIGRATION-SUMMARY.md) for complete migration details**

---

## Legacy OCI Deployment Options

MuShop can also be deployed on [Oracle Cloud Infrastructure][oci] using OCI-managed services. Both deployment models can be used with trial subscriptions. However, [Oracle Cloud Infrastructure][oci] offers an *Always Free* tier with resources that can be used indefinitely.

| [Basic: `deploy/basic`](#Getting-Started-with-MuShop-Basic) | [Complete: `deploy/complete`](#Getting-Started-with-MuShop-Complete) |
|---|---|
| Simplified runtime utilizing **only** [Always Free](https://www.oracle.com/cloud/free/) eligible resources. <br/><br/> Deploy using: <br/>  1. [Terraform][tf] <br/>  2. [Resource Manager][orm_landing] following the steps below <br/>  3. (Recommended) Button below - launches in Resource Manager directly | Polyglot set of micro-services deployed on [Kubernetes](https://kubernetes.io/), showcasing [Oracle Cloud Native](https://www.oracle.com/cloud/cloud-native/) technologies and backing services. <br/><br/> Deploy using: <br/>  1. [Helm](https://helm.sh) <br/> 2. [Terraform][tf] <br/>  3. [Resource Manager][orm_landing] <br/>  4. (Recommended) Button below - launches in Resource Manager directly |
| [![Deploy to Oracle Cloud][magic_button]][magic_mushop_basic_stack]| [![Deploy to Oracle Cloud][magic_button]][magic_mushop_stack]|

```text
mushop
└── deploy
    ├── basic
    └── complete
```

## Getting Started with MuShop Basic

This is a Terraform configuration that deploys the MuShop basic sample application on [Oracle Cloud Infrastructure][oci] and is designed to run using only the Always Free tier resources.

The repository contains the application code as well as the [Terraform][tf] code to create a [Resource Manager][orm] stack, that creates all the required resources and configures the application on the created resources. To simplify getting started, the Resource Manager Stack is created as part of each [release](https://github.com/oracle-quickstart/oci-cloudnative/releases)

The steps below guide you through deploying the application on your tenancy using the OCI Resource Manager.

1. Download the latest [`mushop-basic-stack-latest.zip`](../../releases/latest/download/mushop-basic-stack-latest.zip) file.
2. [Login](https://cloud.oracle.com/resourcemanager/stacks/create) to Oracle Cloud Infrastructure to import the stack
    > `Home > Developer Services > Resource Manager > Stacks > Create Stack`
3. Upload the `mushop-basic-stack-latest.zip` file that was downloaded earlier, and provide a name and description for the stack
4. Configure the stack
   1. **Database Name** - You can choose to provide a database name (optional)
   2. **Node Count** - Select if you want to deploy one or two application instances.
   3. **SSH Public Key** - (Optional) Provide a public SSH key if you wish to establish SSH access to the compute node(s).
5. Review the information and click Create button.
   > The upload can take a few seconds, after which you will be taken to the newly created stack
6. On Stack details page, click on `Terraform Actions > Apply`

All the resources will be created, and the URL to the load balancer will be displayed as `lb_public_url` as in the example below.
> The same information is displayed on the Application Information tab

```text
Outputs:

autonomous_database_password = <generated>

comments = The application URL will be unavailable for a few minutes after provisioning, while the application is configured

dev = Made with ❤ by Oracle Developers

lb_public_url = http://xxx.xxx.xxx.xxx
```

> The application is being deployed to the compute instances asynchronously, and it may take a couple of minutes for the URL to serve the application.

### Cleanup

Even though it is Always Free, you will likely want to terminate the demo application
in your Oracle Cloud Infrastructure tenancy. With the use of Terraform, the [Resource Manager][orm]
stack is also responsible for terminating the application.

Follow these steps to completely remove all provisioned resources:

1. Return to the Oracle Cloud Infrastructure [Console](https://cloud.oracle.com/resourcemanager/stacks)

  > `Home > Developer Services > Resource Manager > Stacks`

1. Select the stack created previously to open the Stack Details view
1. From the Stack Details, select `Terraform Actions > Destroy`
1. Confirm the **Destroy** job when prompted

  > The job status will be **In Progress** while resources are terminated

1. Once the destroy job has succeeded, return to the Stack Details page
1. Click `Delete Stack` and confirm when prompted

#### Basic Topology

The following diagram shows the topology created by this stack.

![MuShop Basic Infra](./images/basic/00-Free-Tier.png)

---

## Getting Started with MuShop Complete

MuShop Complete is a polyglot micro-services application built to showcase a cloud native approach to application development on [Oracle Cloud Infrastructure][oci] using Oracle's [cloud native](https://www.oracle.com/cloud/cloud-native/) services. MuShop Complete uses a Kubernetes cluster, and can be deployed using the provided `helm` charts (preferred), or Kubernetes manifests. It is recommended to use an Oracle Container Engine for Kubernetes cluster, however other Kubernetes distributions will also work.

The [helm chart documentation][chartdocs] walks through the deployment process and various options for customizing the deployment.

### Complete Topology

The following diagram shows the topology created by this stack.

![MuShop Complete Infra](./images/complete/00-Topology.png)

### [![Deploy to Oracle Cloud][magic_button]][magic_mushop_stack]

## Contributing

This project welcomes contributions from the community. Before submitting a pull request, please [review our contribution guide](./CONTRIBUTING.md)

## Security

Please consult the [security guide](./SECURITY.md) for our responsible security vulnerability disclosure process

## Questions

If you have an issue or a question, please take a look at our [FAQs](./deploy/basic/FAQs.md) or [open an issue](https://github.com/oracle-quickstart/oci-cloudnative/issues/new).

[oci]: https://cloud.oracle.com/en_US/cloud-infrastructure
[orm]: https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm
[tf]: https://www.terraform.io
[orm_landing]:https://www.oracle.com/cloud/systems-management/resource-manager/
[chartdocs]: https://github.com/oracle-quickstart/oci-cloudnative/tree/master/deploy/complete/helm-chart#setup
[magic_button]: https://oci-resourcemanager-plugin.plugins.oci.oraclecloud.com/latest/deploy-to-oracle-cloud.svg
[magic_mushop_basic_stack]: https://cloud.oracle.com/resourcemanager/stacks/create?zipUrl=https://github.com/oracle-quickstart/oci-cloudnative/releases/latest/download/mushop-basic-stack-latest.zip
[magic_mushop_stack]: https://cloud.oracle.com/resourcemanager/stacks/create?zipUrl=https://github.com/oracle-quickstart/oci-cloudnative/releases/latest/download/mushop-stack-latest.zip

## License

Copyright (c) 2019 Oracle and/or its affiliates.

Released under the Universal Permissive License v1.0 as shown at
<https://oss.oracle.com/licenses/upl/>.
