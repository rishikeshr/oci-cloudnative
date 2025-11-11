# Application Metadata
ARG APPLICATION_NAME="orders"
ARG VERSION="0.0.1-SNAPSHOT"

# ------------
# Stage 1 : Setup the build environment
FROM gradle:7.6-jdk17 as buildenv

# create source folder
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Refresh Deps.
COPY settings.gradle /usr/src/app
# copy buildscript and cache all dependencies
COPY build.gradle /usr/src/app
#
# ------------

# ------------
# Stage 2 : Build the application
#
FROM buildenv as appbuild
ARG APPLICATION_NAME
ARG VERSION
# Copy the source code.
# This layer is recreated only when there are actual source chnages 
COPY src /usr/src/app/src

# Install the application
RUN gradle clean bootJar -x test
RUN ls -ltr /usr/src/app/build/libs
# ------------

# ------------
# Stage 3 : Application container
#
FROM eclipse-temurin:17-jre-jammy
ARG APPLICATION_NAME
ARG VERSION
ARG ELASTIC_OTEL_AGENT_VERSION=1.7.0

RUN apt-get update && apt-get install -y procps curl && rm -rf /var/lib/apt/lists/*

# Download Elastic OTel Java Agent (EDOT)
RUN curl -L -o /app/elastic-otel-javaagent.jar \
    https://repo1.maven.org/maven2/co/elastic/otel/elastic-otel-javaagent/${ELASTIC_OTEL_AGENT_VERSION}/elastic-otel-javaagent-${ELASTIC_OTEL_AGENT_VERSION}.jar

# copy the generated application distribution
COPY --from=appbuild /usr/src/app/build/libs/${APPLICATION_NAME}-${VERSION}.jar /app/${APPLICATION_NAME}-${VERSION}.jar

EXPOSE 80
WORKDIR /app
ENV APPLICATION_NAME=${APPLICATION_NAME}
ENV VERSION=${VERSION}
ENV TNS_ADMIN=/wallet/

# EDOT (Elastic Distribution of OpenTelemetry) configuration
ENV OTEL_SERVICE_NAME=mushop-orders
ENV OTEL_RESOURCE_ATTRIBUTES=deployment.environment=production,service.version=${VERSION}
ENV OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT:-http://otel-collector:4318}
ENV OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
ENV OTEL_LOGS_EXPORTER=otlp
ENV OTEL_METRICS_EXPORTER=otlp
ENV OTEL_TRACES_EXPORTER=otlp

ENTRYPOINT java $JAVA_OPTS -javaagent:/app/elastic-otel-javaagent.jar -jar /app/${APPLICATION_NAME}-${VERSION}.jar --port=80
#
# ------------

