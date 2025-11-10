# MuShop Newsletter Service

A standalone microservice for handling newsletter subscriptions, replacing the OCI Functions implementation.

## Features

- RESTful API for newsletter subscriptions
- SMTP email delivery
- Health check endpoint
- Docker containerized
- Kubernetes ready

## Endpoints

### POST /subscribe
Subscribe to the newsletter

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Subscription successful! Check your email.",
  "messageId": "abc123"
}
```

### GET /health
Health check endpoint

**Response:**
```json
{
  "status": "healthy",
  "service": "newsletter",
  "timestamp": "2025-01-10T12:00:00.000Z"
}
```

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP port | `3000` |
| `SMTP_HOST` | SMTP server hostname | `localhost` |
| `SMTP_PORT` | SMTP server port | `1025` |
| `SMTP_USER` | SMTP username | Required |
| `SMTP_PASSWORD` | SMTP password | Required |
| `APPROVED_SENDER_EMAIL` | From email address | `noreply@mushop.com` |

## Development

### Local Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Start production server
npm start
```

### With Mailhog (SMTP Testing)

```bash
# Start Mailhog
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog

# Configure environment
export SMTP_HOST=localhost
export SMTP_PORT=1025
export SMTP_USER=test
export SMTP_PASSWORD=test

# Start service
npm start

# Subscribe
curl -X POST http://localhost:3000/subscribe \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'

# Check emails at http://localhost:8025
```

## Docker

### Build

```bash
docker build -t mushop/newsletter:latest .
```

### Run

```bash
docker run -d \
  -p 3000:3000 \
  -e SMTP_HOST=mailhog \
  -e SMTP_PORT=1025 \
  -e SMTP_USER=test \
  -e SMTP_PASSWORD=test \
  mushop/newsletter:latest
```

## Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: newsletter
spec:
  replicas: 2
  selector:
    matchLabels:
      app: newsletter
  template:
    metadata:
      labels:
        app: newsletter
    spec:
      containers:
      - name: newsletter
        image: mushop/newsletter:latest
        ports:
        - containerPort: 3000
        env:
        - name: SMTP_HOST
          value: "smtp.example.com"
        - name: SMTP_PORT
          value: "587"
        - name: SMTP_USER
          valueFrom:
            secretKeyRef:
              name: smtp-credentials
              key: username
        - name: SMTP_PASSWORD
          valueFrom:
            secretKeyRef:
              name: smtp-credentials
              key: password
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: newsletter
spec:
  selector:
    app: newsletter
  ports:
  - port: 80
    targetPort: 3000
```

## Migration from OCI Functions

This service replaces the OCI Functions implementation with a standard containerized microservice.

**Key differences:**
- Runs as HTTP service instead of serverless function
- No dependency on OCI Functions runtime
- Compatible with any SMTP server (not just OCI Email Delivery)
- Can be deployed to any Kubernetes cluster
- Supports Mailhog for local testing

**Backward compatibility:**
- Maintains same request/response format
- Supports both `/subscribe` and `/newsletter/subscribe` endpoints

## Production Considerations

### SMTP Providers

Use a production SMTP service:
- **SendGrid** - Reliable, good free tier
- **Mailgun** - Good API, email validation
- **Amazon SES** - Cost-effective, AWS integration
- **Postmark** - Transactional email focus

### Security

- Use TLS/SSL for SMTP connections in production
- Store SMTP credentials in Kubernetes Secrets
- Enable rate limiting to prevent abuse
- Add email validation and sanitization
- Implement unsubscribe functionality

### Monitoring

- Monitor email delivery success rates
- Track bounce and complaint rates
- Set up alerts for delivery failures
- Log all subscription attempts

## Testing

```bash
# Test subscription
curl -X POST http://localhost:3000/subscribe \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'

# Test health check
curl http://localhost:3000/health

# Invalid email
curl -X POST http://localhost:3000/subscribe \
  -H "Content-Type: application/json" \
  -d '{"email": "invalid-email"}'

# Missing email
curl -X POST http://localhost:3000/subscribe \
  -H "Content-Type: application/json" \
  -d '{}'
```

## License

Copyright © 2019, Oracle and/or its affiliates. All rights reserved.
Licensed under the Universal Permissive License v 1.0
