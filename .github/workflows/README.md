# GitHub Actions Workflows

## Build and Push Docker Images

This workflow automatically builds and pushes all MuShop microservice Docker images to GitHub Container Registry (ghcr.io).

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
| User | `ghcr.io/<owner>/mushop-user` | User authentication & management |
| Catalogue | `ghcr.io/<owner>/mushop-catalogue` | Product catalog |
| Carts | `ghcr.io/<owner>/mushop-carts` | Shopping cart management |
| Orders | `ghcr.io/<owner>/mushop-orders` | Order processing |
| Payment | `ghcr.io/<owner>/mushop-payment` | Payment processing |
| Fulfillment | `ghcr.io/<owner>/mushop-fulfillment` | Order fulfillment |
| Events | `ghcr.io/<owner>/mushop-events` | Event streaming |
| Assets | `ghcr.io/<owner>/mushop-assets` | Static assets & images |
| Newsletter | `ghcr.io/<owner>/mushop-newsletter` | Newsletter subscriptions |
| API | `ghcr.io/<owner>/mushop-api` | Backend for Frontend (BFF) |
| Storefront | `ghcr.io/<owner>/mushop-storefront` | React frontend |

### Image Tags

Images are automatically tagged based on the trigger:

- **Branch pushes**: `<branch-name>` and `<branch-name>-<git-sha>`
- **Latest tag**: `latest` (only on default branch)
- **Semantic versions**: `v1.2.3`, `v1.2`, `v1` (on version tags)
- **Pull requests**: `pr-<number>`

### Multi-Architecture Support

Images are built for:
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit, including Apple Silicon)

### Build Cache

The workflow uses GitHub Actions cache to speed up builds:
- Docker layer cache is automatically managed
- Subsequent builds are significantly faster

### Permissions Required

The workflow requires:
- `contents: read` - To checkout the repository
- `packages: write` - To push images to GitHub Container Registry

### Using the Images

After the workflow runs, images are available at:

```bash
docker pull ghcr.io/<owner>/mushop-<service>:latest
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

To use the pre-built images from GitHub Container Registry, update your `docker-compose.yml`:

```yaml
services:
  user:
    image: ghcr.io/<owner>/mushop-user:latest
    # Remove the 'build' section
    environment:
      - POSTGRES_HOST=postgres-user
      # ... other env vars
```

### Deployment with Helm

Update Helm values to use the images:

```yaml
user:
  image:
    repository: ghcr.io/<owner>/mushop-user
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
- Ensure `packages: write` permission is granted
- Check that GITHUB_TOKEN is available (automatic in GitHub Actions)

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
    image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}-${{ matrix.service.name }}:latest
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
    cosign sign ${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}-${{ matrix.service.name }}@${{ steps.build.outputs.digest }}
```

### Cost Optimization

**Build Time Optimization:**
- Matrix strategy builds services in parallel
- Docker layer caching reduces build time by ~70%
- Multi-stage builds minimize image size

**Storage Optimization:**
- GitHub Container Registry is free for public repositories
- Private repositories have storage limits (see GitHub pricing)
- Old images can be cleaned up manually or with retention policies

### Related Documentation

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
