# Docker Hub Setup Instructions

This document explains how to configure your GitHub repository to automatically push Docker images to Docker Hub.

## Overview

The GitHub Actions workflow in `.github/workflows/build-and-push.yml` builds all MuShop microservice images and pushes them to Docker Hub at:

**https://hub.docker.com/repositories/rishikeshr**

## Prerequisites

1. A Docker Hub account
2. Admin access to the GitHub repository
3. Docker Hub repository namespace: `rishikeshr`

## Step 1: Create Docker Hub Access Token

1. Log in to Docker Hub at https://hub.docker.com
2. Click on your profile icon → **Account Settings**
3. Navigate to **Security** tab
4. Click **New Access Token**
5. Configure the token:
   - **Access Token Description**: `GitHub Actions - MuShop`
   - **Access permissions**: `Read, Write, Delete`
6. Click **Generate**
7. **IMPORTANT**: Copy the token immediately - it won't be shown again!

## Step 2: Configure GitHub Secrets

1. Go to your GitHub repository: https://github.com/rishikeshr/oci-cloudnative
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**

### Add DOCKER_USERNAME Secret:
- **Name**: `DOCKER_USERNAME`
- **Value**: `rishikeshr`
- Click **Add secret**

### Add DOCKER_TOKEN Secret:
- Click **New repository secret** again
- **Name**: `DOCKER_TOKEN`
- **Value**: Paste the Docker Hub access token you copied earlier
- Click **Add secret**

## Step 3: Verify Configuration

1. Go to **Actions** tab in your GitHub repository
2. Select the **Build and Push Docker Images** workflow
3. Click **Run workflow** → **Run workflow**
4. Wait for the workflow to complete
5. Check Docker Hub: https://hub.docker.com/repositories/rishikeshr
6. Verify that images are appearing with tags

## Images Published

The workflow will publish 11 images:

| Service | Docker Hub Image |
|---------|------------------|
| User | `rishikeshr/mushop-user` |
| Catalogue | `rishikeshr/mushop-catalogue` |
| Carts | `rishikeshr/mushop-carts` |
| Orders | `rishikeshr/mushop-orders` |
| Payment | `rishikeshr/mushop-payment` |
| Fulfillment | `rishikeshr/mushop-fulfillment` |
| Events | `rishikeshr/mushop-events` |
| Assets | `rishikeshr/mushop-assets` |
| Newsletter | `rishikeshr/mushop-newsletter` |
| API | `rishikeshr/mushop-api` |
| Storefront | `rishikeshr/mushop-storefront` |

## Image Tags

Each image will have multiple tags:
- `latest` - Latest build from main/master branch
- `<branch-name>` - Branch-specific builds
- `<branch>-<git-sha>` - Specific commit builds
- `v1.2.3` - Semantic version tags (if you create version tags)

## Using the Images

### Pull an Image:
```bash
docker pull rishikeshr/mushop-user:latest
docker pull rishikeshr/mushop-catalogue:latest
docker pull rishikeshr/mushop-storefront:latest
```

### Use in Docker Compose:
```yaml
services:
  user:
    image: rishikeshr/mushop-user:latest
    environment:
      - POSTGRES_HOST=postgres-user
      # ... other env vars

  catalogue:
    image: rishikeshr/mushop-catalogue:latest
    environment:
      - POSTGRES_HOST=postgres-catalogue
      # ... other env vars
```

### Use in Kubernetes/Helm:
```yaml
user:
  image:
    repository: rishikeshr/mushop-user
    tag: latest

catalogue:
  image:
    repository: rishikeshr/mushop-catalogue
    tag: latest
```

## Troubleshooting

### Authentication Failed
- Verify secrets are correctly configured in GitHub
- Check that the Docker Hub token hasn't expired
- Ensure token has write permissions

### Build Failed
- Check the GitHub Actions logs for specific errors
- Verify Dockerfiles are valid
- Ensure all dependencies are available

### Images Not Appearing
- Check if the workflow completed successfully
- Verify you're logged into the correct Docker Hub account
- Check repository visibility (public vs private)

### Rate Limiting
- Docker Hub free tier: 100 pulls/6 hours (anonymous), 200 pulls/6 hours (authenticated)
- If hitting limits, consider Docker Hub Pro plan
- Use Docker Hub authentication in your deployments

## Security Best Practices

1. **Never commit tokens**: Always use GitHub Secrets
2. **Rotate tokens regularly**: Change your Docker Hub access token periodically
3. **Use least privilege**: Only grant necessary permissions to tokens
4. **Enable 2FA**: Enable two-factor authentication on Docker Hub
5. **Monitor usage**: Check Docker Hub analytics for unusual activity

## Additional Resources

- [Docker Hub Documentation](https://docs.docker.com/docker-hub/)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Workflow README](.github/workflows/README.md)
