import Nedb from "@seald-io/nedb";
import { BucketPermission } from "../constant/bucket-permission.js";
import collections from "../constant/collections.js";
import { miscConstants } from "../constant/misc-constants.js";
import { Generic } from "../global.js";
import { DatabaseEngine } from "../lib/database-engine.js";
import { Bucket } from "../model/core-db-entities.js";

export class BucketService {
  db: Nedb;

  constructor(dbEngine: DatabaseEngine) {
    this.db = dbEngine.connection;
  }

  async findBucketById(id: string): Promise<Bucket> {
    let doc = await this.db.findOneAsync({
      collection: collections.BUCKET,
      _id: id,
    });
    return doc;
  }

  async findBucketByName(name: string): Promise<Bucket> {
    let doc = await this.db.findOneAsync({
      collection: collections.BUCKET,
      name,
    });
    return doc;
  }

  async createBucket(
    name: string,
    cryptSpec: string,
    cryptData: string,
    metaData: Record<string, never>,
    createdByUserId: string
  ): Promise<Bucket> {
    let data: Bucket = {
      _id: undefined,
      name,
      cryptSpec,
      cryptData,
      metaData,
      bucketAuthorizations: [{
        userId: createdByUserId,
        notes: miscConstants.BUCKET_CREATOR_AUTHORIZATION_MESSAGE,
        permissions: this.createNewBucketPermissionAllAllowed(),
      },],
      createdByUserIdentifier: `${createdByUserId}@.`,
      createdAt: Date.now(),
      updatedAt: Date.now(),
    };

    return await this.db.insertAsync({
      collection: collections.BUCKET,
      ...data
    });
  }

  async listAllBuckets(): Promise<Bucket[]> {
    let list = await this.db.findAsync({
      collection: collections.BUCKET,
    });
    return list;
  }

  async setBucketName(bucketId: string, name: string) {
    return await this.db.updateAsync(
      {
        collection: collections.BUCKET,
        _id: bucketId,
      },
      {
        $set: {
          name,
          updatedAt: Date.now()
        },
      }
    );
  }

  async setBucketMetaData(bucketId: string, metaData: Generic) {
    return await this.db.updateAsync(
      {
        collection: collections.BUCKET,
        _id: bucketId,
      },
      {
        $set: {
          metaData,
          updatedAt: Date.now()
        },
      }
    );
  }

  async removeBucket(bucketId: string) {
    return await this.db.removeAsync(
      {
        collection: collections.BUCKET,
        _id: bucketId,
      },
      { multi: false }
    );
  }

  createNewBucketPermissionAllAllowed() {
    return Object.keys(BucketPermission).reduce((map: Generic, key) => {
      map[key] = true;
      return map;
    }, {});
  }

  createNewBucketPermissionAllForbidden() {
    return Object.keys(BucketPermission).reduce((map: Generic, key) => {
      map[key] = false;
      return map;
    }, {});
  }

  async authorizeUserWithAllPermissionsForbidden(
    bucketId: string,
    userId: string,
    notes: string
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.BUCKET,
        _id: bucketId,
      },
      {
        $push: {
          bucketAuthorizations: {
            userId,
            notes,
            permissions: this.createNewBucketPermissionAllForbidden(),
          },
        },
      }
    );
  }

  async setAuthorizationList(bucketId: string, bucketAuthorizations: Generic) {
    return await this.db.updateAsync(
      {
        collection: collections.BUCKET,
        _id: bucketId,
      },
      {
        $set: {
          bucketAuthorizations,
        },
      }
    );
  }
}
