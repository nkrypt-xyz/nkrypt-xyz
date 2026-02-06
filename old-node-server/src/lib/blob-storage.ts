import checkDiskSpace from 'check-disk-space';
import { createReadStream, createWriteStream, promises } from "fs";
import fsPromises from "fs/promises";
import { ensureDir, getAbsolutePath, resolvePath } from "../utility/file-utils.js";
import { Config } from "./config.js";

class BlobStorage {
  config: Config;
  private _dir: string;

  constructor(config: Config) {
    this.config = config;
    this._dir = resolvePath(config.blobStorage.dir);
  }

  async init() {
    ensureDir(this._dir);
  }

  createWritableStream(blobId: string, start: number = 0) {
    let path = resolvePath(this._dir, blobId);
    let stream = createWriteStream(path, { start, flags: 'a+' });
    return stream;
  }

  async queryBlobSize(blobId: string) {
    let path = resolvePath(this._dir, blobId);
    let stats = (await fsPromises.stat(path))
    return stats.size;
  }

  createReadableStream(blobId: string) {
    let path = resolvePath(this._dir, blobId);
    let stream = createReadStream(path);
    return stream;
  }

  async removeByBlobId(blobId: string) {
    let path = resolvePath(this._dir, blobId);
    await promises.unlink(path);
  }

  async queryUsage(): Promise<{ free: number, size: number }> {
    return checkDiskSpace(getAbsolutePath(this._dir));
  }
}

export { BlobStorage };
