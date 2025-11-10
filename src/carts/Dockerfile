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
RUN mvn package

# ------------
# Stage 2 : Application container
#
FROM eclipse-temurin:11-jre-jammy
ARG APPLICATION_NAME
ARG VERSION

COPY --from=appbuild /usr/src/app/target/${APPLICATION_NAME}-${VERSION}.jar /app/${APPLICATION_NAME}-${VERSION}.jar
COPY --from=appbuild /usr/src/app/target/libs /app/libs
EXPOSE 80
WORKDIR /app

ENV APPLICATION_NAME=${APPLICATION_NAME}
ENV VERSION=${VERSION}

ENTRYPOINT java $JAVA_OPTS  -jar /app/${APPLICATION_NAME}-${VERSION}.jar
