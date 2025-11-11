# Node version
ARG nodeVersion=16

# Node ENV
ARG nodeEnv

# PROXY
ARG httpProxy

################################
#    ----- Base Image -----    #
################################
FROM node:${nodeVersion}-alpine AS base

#################################
#    ----- Build Image -----    #
#################################
FROM base AS node-build

# Configure npm
ARG httpProxy
ENV HTTP_PROXY ${httpProxy:-""}
ENV HTTPS_PROXY ${httpProxy:-""}
RUN npm config set loglevel error \
  && npm set progress=false

# Install application dependencies
COPY package*.json /tmp/
RUN cd /tmp && npm ci --no-optional
RUN mkdir -p /usr/src/app && cp -a /tmp/node_modules /usr/src/app/
RUN rm -rf /tmp/node_modules

# Copy files
WORKDIR /usr/src/app
COPY *.js* ./
COPY src src

# Build
ARG nodeEnv
ENV NODE_ENV ${nodeEnv:-"production"}
RUN npm run build:docker

# Add schema for init operation
COPY schema schema

###################################
#    ----- Image Runtime -----    #
###################################
FROM base
ARG nodeEnv
ENV NODE_ENV ${nodeEnv:-"production"}
WORKDIR /usr/src/app
COPY --from=node-build /usr/src/app .

ENV PORT 3000
EXPOSE 3000

CMD [ "node", "dist/main.js" ]
