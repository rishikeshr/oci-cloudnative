# GitHub Actions Workflows

## Build and Push Docker Images

This workflow automatically builds and pushes all MuShop microservice Docker images to Docker Hub (hub.docker.com).

### Triggers

The workflow runs on:

1. **Push to master/main branch** - Builds and pushes images with `latest` tag
2. **Tag creation** (v*.*.* pattern) - Builds and pushes images with semantic version tags
3. **Pull requests** - Builds images only (no push) to verify they compile
4. **Manual dispatch** - Can be triggered manually from the Actions tab

### Images Built

The workflow builds and pushes the following images:

| Service | Image Name | Description |
|---------|------------|-------------|
| User | `rishikeshr/mushop-user` | User authentication & management |
| Catalogue | `rishikeshr/mushop-catalogue` | Product catalog |
| Carts | `rishikeshr/mushop-carts` | Shopping cart management |
| Orders | `rishikeshr/mushop-orders` | Order processing |
| Payment | `rishikeshr/mushop-payment` | Payment processing |
| Fulfillment | `rishikeshr/mushop-fulfillment` | Order fulfillment |
| Events | `rishikeshr/mushop-events` | Event streaming |
| Assets | `rishikeshr/mushop-assets` | Static assets & images |
| Newsletter | `rishikeshr/mushop-newsletter` | Newsletter subscriptions |
| API | `rishikeshr/mushop-api` | Backend for Frontend (BFF) |
| Storefront | `rishikeshr/mushop-storefront` | React frontend |

### Image Tags

Images are automatically tagged based on the trigger:

- **Branch pushes**: `<branch-name>` and `<branch-name>-<git-sha>`
- **Latest tag**: `latest` (only on default branch)
- **Semantic versions**: `v1.2.3`, `v1.2`, `v1` (on version tags)
- **Pull requests**: `pr-<number>`

### Multi-Architecture Support

Most images are built for:
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit, including Apple Silicon)

**Note:** The fulfillment service is only built for `linux/amd64` due to GraalVM native image compilation limitations with QEMU emulation during cross-platform builds.

### Build Cache

The workflow uses GitHub Actions cache to speed up builds:
- Docker layer cache is automatically managed
- Subsequent builds are significantly faster

### Secrets Required

The workflow requires the following secrets to be configured in your GitHub repository:
- `DOCKER_USERNAME` - Your Docker Hub username (rishikeshr)
- `DOCKER_TOKEN` - Your Docker Hub access token or password

**To configure secrets:**
1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add `DOCKER_USERNAME` with value: `rishikeshr`
5. Add `DOCKER_TOKEN` with your Docker Hub access token

**To create a Docker Hub access token:**
1. Log in to https://hub.docker.com
2. Go to Account Settings → Security
3. Click "New Access Token"
4. Give it a description (e.g., "GitHub Actions")
5. Copy the token and add it to GitHub secrets

### Permissions Required

The workflow requires:
- `contents: read` - To checkout the repository
- `packages: write` - For workflow attestation (optional)

### Using the Images

After the workflow runs, images are available at:

```bash
docker pull rishikeshr/mushop-<service>:latest
```

For example:
```bash
docker pull rishikeshr/mushop-user:latest
docker pull rishikeshr/mushop-catalogue:latest
docker pull rishikeshr/mushop-storefront:latest
```

### Local Testing

To test the workflow locally before pushing:

```bash
# Install act (GitHub Actions local runner)
# https://github.com/nektos/act

# Run the workflow
act push -s GITHUB_TOKEN=<your-token>
```

### Deployment with Docker Compose

To use the pre-built images from Docker Hub, update your `docker-compose.yml`:

```yaml
services:
  user:
    image: rishikeshr/mushop-user:latest
    # Remove the 'build' section
    environment:
      - POSTGRES_HOST=postgres-user
      # ... other env vars

  catalogue:
    image: rishikeshr/mushop-catalogue:latest
    environment:
      - POSTGRES_HOST=postgres-catalogue
      # ... other env vars

  storefront:
    image: rishikeshr/mushop-storefront:latest
    # ... other env vars
```

### Deployment with Helm

Update Helm values to use the images:

```yaml
user:
  image:
    repository: rishikeshr/mushop-user
    tag: latest

catalogue:
  image:
    repository: rishikeshr/mushop-catalogue
    tag: latest

storefront:
  image:
    repository: rishikeshr/mushop-storefront
    tag: latest
```

### Manual Workflow Dispatch

You can manually trigger the workflow:

1. Go to the **Actions** tab in your GitHub repository
2. Select **Build and Push Docker Images** workflow
3. Click **Run workflow**
4. Choose the branch
5. Click **Run workflow** button

### Troubleshooting

**Authentication Issues:**
- Verify `DOCKER_USERNAME` and `DOCKER_TOKEN` secrets are configured
- Check that the Docker Hub access token has write permissions
- Ensure the token hasn't expired

**Build Failures:**
- Check the specific service logs in the workflow run
- Verify Dockerfiles are valid and dependencies are available
- Ensure base images are accessible

**Cache Issues:**
- GitHub Actions cache is limited to 10GB
- Old cache entries are automatically evicted
- You can clear cache from Settings > Actions > Caches

### Security

**Image Scanning:**
To add vulnerability scanning, add this step before pushing:

```yaml
- name: Scan image for vulnerabilities
  uses: aquasecurity/trivy-action@master
  with:
    image-ref: rishikeshr/mushop-${{ matrix.service.name }}:latest
    format: 'sarif'
    output: 'trivy-results.sarif'
```

**Image Signing:**
To sign images with Cosign:

```yaml
- name: Install Cosign
  uses: sigstore/cosign-installer@v3

- name: Sign image
  env:
    COSIGN_EXPERIMENTAL: "true"
  run: |
    cosign sign rishikeshr/mushop-${{ matrix.service.name }}@${{ steps.build.outputs.digest }}
```

### Cost Optimization

**Build Time Optimization:**
- Matrix strategy builds services in parallel
- Docker layer caching reduces build time by ~70%
- Multi-stage builds minimize image size

**Storage Optimization:**
- Docker Hub free tier includes unlimited public repositories
- Free tier has pull rate limits (100 pulls per 6 hours for anonymous users)
- Authenticated users get 200 pulls per 6 hours
- Old images can be cleaned up manually from Docker Hub UI
- Consider Docker Hub Pro ($5/month) for unlimited pulls and private repositories

### Related Documentation

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Hub](https://hub.docker.com)
- [Docker Hub Repositories](https://hub.docker.com/repositories/rishikeshr)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Docker Login Action](https://github.com/docker/login-action)
