import Nedb from "@seald-io/nedb";
import { BlobStorage } from "../lib/blob-storage.js";
import { DatabaseEngine } from "../lib/database-engine.js";

export class MetricsService {
  db: Nedb;
  blobStorage: BlobStorage;

  constructor(dbEngine: DatabaseEngine, blobStorage: BlobStorage) {
    this.db = dbEngine.connection;
    this.blobStorage = blobStorage;
  }

  async getDiskUsage(): Promise<{ usedBytes: number, totalBytes: number }> {
    let { free, size } = await this.blobStorage.queryUsage();
    return {
      totalBytes: size,
      usedBytes: size - free
    };
  }
}
