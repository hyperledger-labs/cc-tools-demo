import AWS = require('aws-sdk');
import env from './s3.env';

const s3Client = new AWS.S3({
  accessKeyId: env.AWS_ACCESS_KEY,
  secretAccessKey: env.AWS_SECRET_ACCESS_KEY
});

export default s3Client;
