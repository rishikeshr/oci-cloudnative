/**
 * Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
const fs = require('fs');
const path = require('path');
const mime = require('mime');
const AWS = require('aws-sdk');
const config = require('./config');

// Configure AWS SDK for MinIO
const s3 = new AWS.S3({
  endpoint: config.minioEndpoint ? `${config.minioUseSSL ? 'https' : 'http'}://${config.minioEndpoint}:${config.minioPort}` : null,
  accessKeyId: config.minioAccessKey,
  secretAccessKey: config.minioSecretKey,
  s3ForcePathStyle: true, // Needed for MinIO
  signatureVersion: 'v4',
  sslEnabled: config.minioUseSSL,
});

const putImage = async (dir, img) => {
  const mType = mime.getType(img);
  const data = fs.readFileSync(path.join(dir, img));
  const fname = img;

  const params = {
    Bucket: config.minioBucket,
    Key: fname,
    Body: data,
    ContentType: mType,
    CacheControl: `max-age=${config.cache.maxAge}, public, no-transform`,
  };

  try {
    await s3.putObject(params).promise();
    console.log(`PUT Success: ${img}`);
  } catch (e) {
    console.error(`PUT Failed: ${img}`, e.toString());
  }
};

// Create bucket if it doesn't exist
async function ensureBucket() {
  try {
    await s3.headBucket({ Bucket: config.minioBucket }).promise();
    console.log(`Bucket ${config.minioBucket} exists`);
  } catch (err) {
    if (err.statusCode === 404) {
      console.log(`Creating bucket ${config.minioBucket}`);
      await s3.createBucket({ Bucket: config.minioBucket }).promise();
    } else {
      throw err;
    }
  }
}

// Main deployment logic
(async () => {
  if (!config.minioEndpoint || !config.minioAccessKey || !config.minioSecretKey) {
    console.log('MinIO configuration not provided, exiting with nothing to do');
    process.exit(0);
  }

  try {
    await ensureBucket();

    const dist = path.join(__dirname, config.dist);
    if (fs.existsSync(dist)) {
      const files = fs.readdirSync(dist);
      await Promise.all(files.map(f => putImage(dist, f)));
      console.log('Deployment complete!');
    } else {
      console.error(`Optimized image directory does not exist: ${config.dist}`);
      process.exit(1);
    }
  } catch (err) {
    console.error('Deployment failed:', err);
    process.exit(1);
  }
})();
