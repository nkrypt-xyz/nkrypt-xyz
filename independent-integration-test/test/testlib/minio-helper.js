import * as Minio from 'minio';
import { Readable } from 'stream';

/**
 * MinIO helper for test verification
 * This allows tests to verify that blobs are correctly stored in MinIO
 */
export class MinIOHelper {
  constructor(config) {
    this.client = new Minio.Client({
      endPoint: config.endpoint || 'localhost',
      port: config.port || 9000,
      useSSL: config.useSSL || false,
      accessKey: config.accessKey || 'minioadmin',
      secretKey: config.secretKey || 'minioadmin',
    });
    this.bucketName = config.bucketName || 'nkrypt-blobs';
  }

  /**
   * Ensure the bucket exists
   */
  async ensureBucket() {
    const exists = await this.client.bucketExists(this.bucketName);
    if (!exists) {
      await this.client.makeBucket(this.bucketName);
    }
  }

  /**
   * Get a blob from MinIO as a Buffer
   * @param {string} blobId - The blob ID
   * @returns {Promise<Buffer>} The blob data
   */
  async getBlob(blobId) {
    const objectKey = `blobs/${blobId}`;
    
    return new Promise((resolve, reject) => {
      const chunks = [];
      
      this.client.getObject(this.bucketName, objectKey, (err, stream) => {
        if (err) {
          return reject(err);
        }

        stream.on('data', (chunk) => chunks.push(chunk));
        stream.on('end', () => resolve(Buffer.concat(chunks)));
        stream.on('error', reject);
      });
    });
  }

  /**
   * Get a blob from MinIO as a readable stream
   * @param {string} blobId - The blob ID
   * @returns {Promise<Readable>} The blob stream
   */
  async getBlobStream(blobId) {
    const objectKey = `blobs/${blobId}`;
    
    return new Promise((resolve, reject) => {
      this.client.getObject(this.bucketName, objectKey, (err, stream) => {
        if (err) {
          return reject(err);
        }
        resolve(stream);
      });
    });
  }

  /**
   * Get blob metadata (size, etc.)
   * @param {string} blobId - The blob ID
   * @returns {Promise<Object>} The blob stat info
   */
  async getBlobStat(blobId) {
    const objectKey = `blobs/${blobId}`;
    return await this.client.statObject(this.bucketName, objectKey);
  }

  /**
   * Check if a blob exists
   * @param {string} blobId - The blob ID
   * @returns {Promise<boolean>} True if blob exists
   */
  async blobExists(blobId) {
    try {
      await this.getBlobStat(blobId);
      return true;
    } catch (err) {
      if (err.code === 'NotFound') {
        return false;
      }
      throw err;
    }
  }

  /**
   * Delete a blob from MinIO
   * @param {string} blobId - The blob ID
   */
  async deleteBlob(blobId) {
    const objectKey = `blobs/${blobId}`;
    await this.client.removeObject(this.bucketName, objectKey);
  }

  /**
   * List all blobs with a given prefix
   * @param {string} prefix - The prefix to filter by (default: 'blobs/')
   * @returns {Promise<Array>} Array of object names
   */
  async listBlobs(prefix = 'blobs/') {
    return new Promise((resolve, reject) => {
      const objects = [];
      const stream = this.client.listObjects(this.bucketName, prefix, true);
      
      stream.on('data', (obj) => objects.push(obj));
      stream.on('end', () => resolve(objects));
      stream.on('error', reject);
    });
  }

  /**
   * Clean up all test blobs
   */
  async cleanupAllBlobs() {
    const objects = await this.listBlobs('blobs/');
    const objectNames = objects.map(obj => obj.name);
    
    if (objectNames.length > 0) {
      await this.client.removeObjects(this.bucketName, objectNames);
    }
  }
}

/**
 * Create a MinIO helper with default test configuration
 * Can be overridden with environment variables:
 * - MINIO_ENDPOINT
 * - MINIO_PORT
 * - MINIO_ACCESS_KEY
 * - MINIO_SECRET_KEY
 * - MINIO_BUCKET_NAME
 */
export function createTestMinIOHelper() {
  return new MinIOHelper({
    endpoint: process.env.MINIO_ENDPOINT || 'localhost',
    port: parseInt(process.env.MINIO_PORT || '9000'),
    useSSL: process.env.MINIO_USE_SSL === 'true',
    accessKey: process.env.MINIO_ACCESS_KEY || 'minioadmin',
    secretKey: process.env.MINIO_SECRET_KEY || 'minioadmin',
    bucketName: process.env.MINIO_BUCKET_NAME || 'nkrypt-blobs',
  });
}
