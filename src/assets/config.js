/**
 * Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
require('dotenv').config();

const {
  MINIO_ENDPOINT,
  MINIO_PORT,
  MINIO_ACCESS_KEY,
  MINIO_SECRET_KEY,
  MINIO_BUCKET,
  MINIO_USE_SSL,
} = process.env;

function getMinIOUrl() {
  if (MINIO_ENDPOINT && MINIO_BUCKET) {
    const protocol = MINIO_USE_SSL === 'true' ? 'https' : 'http';
    const port = MINIO_PORT || '9000';
    return `${protocol}://${MINIO_ENDPOINT}:${port}/${MINIO_BUCKET}`;
  }
  return null;
}

module.exports = {
  env: process.env,
  dist: 'dist',
  minioEndpoint: MINIO_ENDPOINT,
  minioPort: MINIO_PORT || '9000',
  minioAccessKey: MINIO_ACCESS_KEY,
  minioSecretKey: MINIO_SECRET_KEY,
  minioBucket: MINIO_BUCKET || 'mushop-assets',
  minioUseSSL: MINIO_USE_SSL === 'true',
  bucketUrl: getMinIOUrl(),
  cache: {
    maxAge: 604800
  },
};
