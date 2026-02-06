import Nedb from "@seald-io/nedb";
import collections from "../constant/collections.js";
import { Generic } from "../global.js";
import { DatabaseEngine } from "../lib/database-engine.js";
import { File } from "../model/core-db-entities.js";

export class FileService {
  db: Nedb;

  constructor(dbEngine: DatabaseEngine) {
    this.db = dbEngine.connection;
  }

  async findFileById(bucketId: string, fileId: string): Promise<File> {
    let doc = await this.db.findOneAsync({
      collection: collections.FILE,
      bucketId,
      _id: fileId,
    });
    return doc;
  }

  async findFileByNameAndParent(
    name: string,
    bucketId: string,
    parentDirectoryId: string
  ): Promise<File> {
    let doc = await this.db.findOneAsync({
      collection: collections.FILE,
      name,
      bucketId,
      parentDirectoryId,
    });
    return doc;
  }

  async createFile(
    name: string,
    bucketId: string,
    metaData: Generic,
    encryptedMetaData: string,
    createdByUserId: string,
    parentDirectoryId: string
  ): Promise<File> {
    let data: File = {
      _id: undefined,
      name,
      metaData,
      encryptedMetaData,
      bucketId,
      parentDirectoryId,
      sizeAfterEncryptionBytes: 0,
      createdByUserIdentifier: `${createdByUserId}@.`,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      contentUpdatedAt: Date.now(),
    };

    return await this.db.insertAsync({
      collection: collections.FILE,
      ...data
    });
  }

  async setFileName(bucketId: string, fileId: string, name: string) {
    return await this.db.updateAsync(
      {
        collection: collections.FILE,
        _id: fileId,
        bucketId,
      },
      {
        $set: {
          name,
          updatedAt: Date.now()
        },
      }
    );
  }

  async setFileContentUpdateAt(bucketId: string, fileId: string, epoc: number) {
    return await this.db.updateAsync(
      {
        collection: collections.FILE,
        _id: fileId,
        bucketId,
      },
      {
        $set: {
          updatedAt: epoc,
        },
      }
    );
  }

  async setFileEncryptedMetaData(
    bucketId: string,
    fileId: string,
    encryptedMetaData: string
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.FILE,
        _id: fileId,
        bucketId,
      },
      {
        $set: {
          encryptedMetaData,
          updatedAt: Date.now()
        },
      }
    );
  }

  async setFileMetaData(bucketId: string, fileId: string, metaData: Generic) {
    return await this.db.updateAsync(
      {
        collection: collections.FILE,
        _id: fileId,
        bucketId,
      },
      {
        $set: {
          metaData,
          updatedAt: Date.now()
        },
      }
    );
  }

  async deleteFile(bucketId: string, fileId: string) {
    return await this.db.removeAsync(
      {
        collection: collections.FILE,
        _id: fileId,
        bucketId,
      },
      { multi: false }
    );
  }

  async moveFile(
    bucketId: string,
    fileId: string,
    newParentDirectoryId: string,
    newName: string
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.FILE,
        _id: fileId,
        bucketId,
      },
      {
        $set: {
          parentDirectoryId: newParentDirectoryId,
          name: newName,
          updatedAt: Date.now()
        },
      }
    );
  }

  async listFilesUnderDirectory(bucketId: string, parentDirectoryId: string) {
    let list = await this.db.findAsync({
      collection: collections.FILE,
      bucketId,
      parentDirectoryId,
    });
    return list;
  }
}
