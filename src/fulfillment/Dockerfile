ARG APPLICATION_NAME="fulfillment"
ARG VERSION="0.0.1-SNAPSHOT"

# Stage 1 : Setup the build environment
FROM gradle:7.6-jdk11 as buildenv

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY settings.gradle /usr/src/app
# copy buildscript and cache all dependencies
COPY build.gradle /usr/src/app
COPY gradle.properties /usr/src/app
#RUN gradle --refresh-dependencies

# Stage 2 : Build the application
FROM buildenv as appbuild
# Copy the source code.
# This layer is recreated only when there are actual source chnages
COPY src /usr/src/app/src
# build an executable fat jar using shadowJar
RUN gradle clean shadowJar -x test

# Stage 3 : Application container using OpenJDK
FROM eclipse-temurin:11-jre-jammy
ARG APPLICATION_NAME
ARG VERSION

# Copy the generated shadow JAR (fat jar)
COPY --from=appbuild /usr/src/app/build/libs/${APPLICATION_NAME}-*-all.jar /app/${APPLICATION_NAME}.jar

EXPOSE 80
WORKDIR /app

ENV APPLICATION_NAME=${APPLICATION_NAME}
ENV VERSION=${VERSION}

ENTRYPOINT ["java", "-jar", "/app/fulfillment.jar"]

