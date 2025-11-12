#
# Copyright (c) 2020-2021 Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at http://oss.oracle.com/licenses/upl.
#

##### Go Builder image
FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:1.20 AS go-builder

ARG TARGETARCH
ARG TARGETOS

WORKDIR /go/src/mushop/catalogue

# Support for Offline local image. Online image on OCI Object Storage
COPY images/ images/

# Catalogue Go Source
COPY cmd/cataloguesvc/*.go cmd/cataloguesvc/
COPY *.go .
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod tidy
RUN go mod download

# Build Catalogue service
RUN GOARCH=${TARGETARCH:-amd64} GOOS=${TARGETOS:-linux} \
  go build -a \
  -o /catalogue mushop/catalogue/cmd/cataloguesvc

##### Catalogue Service Image
FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.16

ARG TARGETPLATFORM

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

WORKDIR /app
COPY --from=go-builder --chown=app:app /catalogue /app/
COPY --chown=app:app images/ /app/images/

USER app

CMD ["/app/catalogue", "-port=8080"]
EXPOSE 8080

LABEL org.opencontainers.image.title="catalogue" \
    org.opencontainers.image.architecture="${TARGETPLATFORM}"
