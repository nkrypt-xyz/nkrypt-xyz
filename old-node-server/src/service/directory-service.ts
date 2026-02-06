import Nedb from "@seald-io/nedb";
import collections from "../constant/collections.js";
import { Generic } from "../global.js";
import { DatabaseEngine } from "../lib/database-engine.js";
import { Directory } from "../model/core-db-entities.js";

export class DirectoryService {
  db: Nedb;

  constructor(dbEngine: DatabaseEngine) {
    this.db = dbEngine.connection;
  }

  async listChildrenOfDirectory(bucketId: string, parentDirectoryId: string): Promise<Directory[]> {
    let list = await this.db.findAsync({
      collection: collections.DIRECTORY,
      bucketId,
      parentDirectoryId,
    });
    return list;
  }

  async findDirectoryById(bucketId: string, directoryId: string): Promise<Directory> {
    let doc = await this.db.findOneAsync({
      collection: collections.DIRECTORY,
      bucketId,
      _id: directoryId,
    });
    return doc;
  }

  async findDirectoryByNameAndParent(
    name: string,
    bucketId: string,
    parentDirectoryId: string
  ): Promise<Directory> {
    let doc = await this.db.findOneAsync({
      collection: collections.DIRECTORY,
      name,
      bucketId,
      parentDirectoryId,
    });
    return doc;
  }

  async createDirectory(
    name: string,
    bucketId: string,
    metaData: Generic,
    encryptedMetaData: string,
    createdByUserId: string,
    parentDirectoryId: string | null
  ): Promise<Directory> {
    let data: Directory = {
      _id: undefined,
      name,
      metaData,
      encryptedMetaData,
      bucketId,
      parentDirectoryId,
      createdByUserIdentifier: `${createdByUserId}@.`,
      createdAt: Date.now(),
      updatedAt: Date.now(),
    };

    return await this.db.insertAsync({
      collection: collections.DIRECTORY,
      ...data
    });
  }

  async listRootDirectoriesByBucketIdList(idList: string[]) {
    let list = await this.db.findAsync({
      collection: collections.DIRECTORY,
      bucketId: { $in: idList },
      parentDirectoryId: null,
    });
    return list;
  }

  async findRootDirectoryByBucketId(bucketId: string): Promise<Directory> {
    let doc = await this.db.findOneAsync({
      collection: collections.DIRECTORY,
      bucketId,
      parentDirectoryId: null,
    });
    return doc;
  }

  async setDirectoryName(bucketId: string, directoryId: string, name: string) {
    return await this.db.updateAsync(
      {
        collection: collections.DIRECTORY,
        _id: directoryId,
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

  async setDirectoryEncryptedMetaData(
    bucketId: string,
    directoryId: string,
    encryptedMetaData: string
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.DIRECTORY,
        _id: directoryId,
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

  async setDirectoryMetaData(
    bucketId: string,
    directoryId: string,
    metaData: Generic
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.DIRECTORY,
        _id: directoryId,
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

  async deleteDirectory(bucketId: string, directoryId: string) {
    return await this.db.removeAsync(
      {
        collection: collections.DIRECTORY,
        _id: directoryId,
        bucketId,
      },
      { multi: false }
    );
  }

  async moveDirectory(
    bucketId: string,
    directoryId: string,
    newParentDirectoryId: string,
    newName: string
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.DIRECTORY,
        _id: directoryId,
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

  async deleteDirectoryAndChildrenInTheBackground(bucketId: string, parentDirectory: Directory) {
    let fileList = await dispatch.fileService.listFilesUnderDirectory(bucketId, parentDirectory._id!);
    for (let file of fileList) {
      await dispatch.fileService.deleteFile(bucketId, file._id);
      await dispatch.blobService.removeAllBlobsOfFile(bucketId, file._id);
    }

    let directoryList = await dispatch.directoryService.listChildrenOfDirectory(bucketId, parentDirectory._id!);
    for (let directory of directoryList) {
      await dispatch.directoryService.deleteDirectory(bucketId, directory._id!);
      await this.deleteDirectoryAndChildrenInTheBackground(bucketId, directory);
    }
  }
}
