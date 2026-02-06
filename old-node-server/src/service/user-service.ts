import Nedb from "@seald-io/nedb";
import collections from "../constant/collections.js";
import { DatabaseEngine } from "../lib/database-engine.js";
import { throwOnFalsy, UserError } from "../utility/coded-error.js";

export class UserService {
  db: Nedb;

  constructor(dbEngine: DatabaseEngine) {
    this.db = dbEngine.connection;
  }

  async findUserByIdOrFail(_id: string) {
    let user = await this.db.findOneAsync({
      collection: collections.USER,
      _id,
    });
    throwOnFalsy(
      UserError,
      user,
      "USER_NOT_FOUND",
      "The requested user could not be found."
    );
    return user;
  }

  async findUserByUserName(userName: string) {
    let user = await this.db.findOneAsync({
      collection: collections.USER,
      userName,
    });
    return user;
  }

  async findUserOrFail(userName: string) {
    let user = await this.db.findOneAsync({
      collection: collections.USER,
      userName,
    });
    throwOnFalsy(
      UserError,
      user,
      "USER_NOT_FOUND",
      "The requested user could not be found."
    );
    return user;
  }

  async updateOwnCommonProperties(_id: string, displayName: string) {
    return await this.db.updateAsync(
      {
        collection: collections.USER,
        _id,
      },
      {
        $set: {
          displayName,
          updatedAt: Date.now(),
        },
      }
    );
  }

  async listAllUsers() {
    let userList = await this.db.findAsync({
      collection: collections.USER,
    });
    return userList;
  }

  async queryUsers(userIdList: string[], userNameList: string[]) {
    let query = {
      collection: collections.USER,
      $or: [
        {
          _id: { $in: userIdList },
        },
        {
          userName: { $in: userNameList },
        },
      ]
    };

    let userList = await this.db.findAsync(query);
    return userList;
  }

  async updateUserPassword(
    _id: string,
    newPassword: { hash: string; salt: string }
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.USER,
        _id,
      },
      {
        $set: {
          password: newPassword,
          updatedAt: Date.now(),
        },
      }
    );
  }

  async setBanningStatus(
    _id: string,
    isBanned: boolean
  ) {
    return await this.db.updateAsync(
      {
        collection: collections.USER,
        _id,
      },
      {
        $set: {
          isBanned,
          updatedAt: Date.now(),
        },
      }
    );
  }
}
