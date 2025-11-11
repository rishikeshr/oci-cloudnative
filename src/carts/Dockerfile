##### Application Metadata
ARG APPLICATION_NAME="carts"
ARG VERSION="1.1.0-SNAPSHOT"

# ------------
# Stage 1 : Build the application using Maven
# 
FROM maven:3.6-jdk-11 as appbuild

# Create source folder
RUN mkdir -p /usr/src/app
COPY pom.xml /usr/src/app/
WORKDIR      /usr/src/app/
RUN mvn dependency:go-offline -B

COPY src     /usr/src/app/src
RUN mvn package -DskipTests

# ------------
# Stage 2 : Application container
#
FROM eclipse-temurin:11-jre-jammy
ARG APPLICATION_NAME
ARG VERSION
ARG ELASTIC_OTEL_AGENT_VERSION=1.7.0

RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Download Elastic OTel Java Agent (EDOT)
RUN mkdir -p /app && \
    curl -L -o /app/elastic-otel-javaagent.jar \
    https://repo1.maven.org/maven2/co/elastic/otel/elastic-otel-javaagent/${ELASTIC_OTEL_AGENT_VERSION}/elastic-otel-javaagent-${ELASTIC_OTEL_AGENT_VERSION}.jar

COPY --from=appbuild /usr/src/app/target/${APPLICATION_NAME}-${VERSION}.jar /app/${APPLICATION_NAME}-${VERSION}.jar
COPY --from=appbuild /usr/src/app/target/libs /app/libs
EXPOSE 80
WORKDIR /app

ENV APPLICATION_NAME=${APPLICATION_NAME}
ENV VERSION=${VERSION}

# EDOT (Elastic Distribution of OpenTelemetry) configuration
ENV OTEL_SERVICE_NAME=mushop-carts
ENV OTEL_RESOURCE_ATTRIBUTES=deployment.environment=production,service.version=${VERSION}
ENV OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT:-http://otel-collector:4318}
ENV OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
ENV OTEL_LOGS_EXPORTER=otlp
ENV OTEL_METRICS_EXPORTER=otlp
ENV OTEL_TRACES_EXPORTER=otlp

ENTRYPOINT java $JAVA_OPTS -javaagent:/app/elastic-otel-javaagent.jar -jar /app/${APPLICATION_NAME}-${VERSION}.jar
