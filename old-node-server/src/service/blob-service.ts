import Nedb from "@seald-io/nedb";
import { ReadStream, WriteStream } from "fs";
import collections from "../constant/collections.js";
import { Generic } from "../global.js";
import { BlobStorage } from "../lib/blob-storage.js";
import { DatabaseEngine } from "../lib/database-engine.js";
import { Blob } from "../model/core-db-entities.js";
import { UserError } from "../utility/coded-error.js";
import { createSizeLimiterPassthroughStream } from "../utility/stream-utils.js";

export class BlobService {
  db: Nedb;
  blobStorage: BlobStorage;

  constructor(dbEngine: DatabaseEngine, blobStorage: BlobStorage) {
    this.db = dbEngine.connection;
    this.blobStorage = blobStorage;
  }

  createStandardSizeLimiter(startingOffset: number = 0, byteCountWrapper: Generic = { byteCount: 0 }) {
    return createSizeLimiterPassthroughStream(
      dispatch.config.blobStorage.maxFileSizeBytes - startingOffset, byteCountWrapper
    );
  }

  async createInProgressBlob(bucketId: string, fileId: string, cryptoMetaHeaderContent: string, createdByUserId: string):
    Promise<{ blob: Blob, stream: WriteStream }> {
    let data: Blob = {
      _id: undefined,
      bucketId,
      fileId,
      cryptoMetaHeaderContent,
      startedAt: Date.now(),
      finishedAt: null,
      status: "started",
      createdByUserIdentifier: `${createdByUserId}@.`,
      createdAt: Date.now(),
      updatedAt: Date.now(),
    };

    let blob: Generic = await this.db.insertAsync({
      collection: collections.BLOB,
      ...data
    });

    let stream = this.blobStorage.createWritableStream(blob._id);

    return { blob, stream };
  }

  async getInProgressBlob(bucketId: string, fileId: string, blobId: string, startOffset: number):
    Promise<{ blob: Blob, stream: WriteStream }> {
    let blob = await this.db
      .findOneAsync({
        collection: collections.BLOB,
        bucketId,
        fileId,
        _id: blobId,
        status: "started",
      });
    if (!blob) {
      throw new UserError("BLOB_INVALID", "No in-progress blob found with the given ID");
    }
    let stream = this.blobStorage.createWritableStream(blobId, startOffset);
    return { blob, stream };
  }

  async markBlobAsErroneous(bucketId: string, fileId: string, blobId: string) {
    return await this.db.updateAsync(
      {
        collection: collections.BLOB,
        _id: blobId,
        bucketId,
        fileId,
      },
      {
        $set: {
          status: "error",
          updatedAt: Date.now()
        },
      }
    );
  }

  async markBlobAsFinished(bucketId: string, fileId: string, blobId: string) {
    return await this.db.updateAsync(
      {
        collection: collections.BLOB,
        _id: blobId,
        bucketId,
        fileId,
      },
      {
        $set: {
          status: "finished",
          finishedAt: Date.now(),
          updatedAt: Date.now()
        },
      }
    );
  }

  async findBlobByBucketIdAndFileId(bucketId: string, fileId: string): Promise<Blob> {
    let list = await this.db
      .findAsync({
        collection: collections.BLOB,
        bucketId,
        fileId,
      })
      .sort({ finishedAt: -1 });

    let doc = list.length ? list[0] : null;

    return doc;
  }

  async createReadableStreamFromBlobId(blobId: string): Promise<{ readStream: ReadStream, sizeOfStream: number }> {
    let sizeOfStream = await this.blobStorage.queryBlobSize(blobId);
    let readStream = this.blobStorage.createReadableStream(blobId);
    return { sizeOfStream, readStream }
  }

  async removeAllOtherBlobs(bucketId: string, fileId: string, blobId: string) {
    let list = await this.db.findAsync({
      collection: collections.BLOB,
      bucketId,
      fileId,
      _id: { $ne: blobId },
    });

    for (let blob of list) {
      try {
        await this.blobStorage.removeByBlobId(blob._id);
      } catch (ex) {
        if (
          ex &&
          typeof ex === "object" &&
          "code" in ex &&
          (<Generic>ex).code === "ENOENT"
        ) {
          // the file not being there is not a catastrophe in this case
          ("pass");
        } else {
          throw ex;
        }
      }
    }

    return await this.db.removeAsync(
      {
        collection: collections.BLOB,
        bucketId,
        fileId,
        _id: { $ne: blobId },
      },
      { multi: false } // FIXME why multi: false?
    );
  }

  async removeAllBlobsOfFile(bucketId: string, fileId: string) {
    let list = await this.db.findAsync({
      collection: collections.BLOB,
      bucketId,
      fileId,
    });

    for (let blob of list) {
      try {
        await this.blobStorage.removeByBlobId(blob._id);
      } catch (ex) {
        if (
          ex &&
          typeof ex === "object" &&
          "code" in ex &&
          (<Generic>ex).code === "ENOENT"
        ) {
          // the file not being there is not a catastrophe in this case
          ("pass");
        } else {
          throw ex;
        }
      }
    }

    return await this.db.removeAsync(
      {
        collection: collections.BLOB,
        bucketId,
        fileId,
      },
      { multi: true }
    );
  }
}
