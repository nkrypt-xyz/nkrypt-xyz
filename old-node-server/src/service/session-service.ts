import Nedb from "@seald-io/nedb";
import collections from "../constant/collections.js";
import constants from "../constant/common-constants.js";
import { miscConstants } from "../constant/misc-constants.js";
import { Generic } from "../global";
import { DatabaseEngine } from "../lib/database-engine.js";
import { Session } from "../model/core-db-entities.js";
import {
  DeveloperError,
  throwOnFalsy,
  UserError,
} from "../utility/coded-error.js";
import { generateRandomString } from "../utility/string-utils.js";

const LOGOUT_MESSAGE_PREFIX = "Logout: ";
const FORCE_LOGOUT_MESSAGE_PREFIX = "ForceLogout: ";

export class SessionService {
  db: Nedb;

  constructor(dbEngine: DatabaseEngine) {
    this.db = dbEngine.connection;
  }

  async createNewUniqueSession(user: Generic): Promise<{ session: Session, apiKey: string }> {
    let apiKey;
    let session: Session, exists;
    let safetyCap = constants.std.SAFETY_CAP;
    do {
      apiKey = generateRandomString(constants.iam.API_KEY_LENGTH);

      let exists = await this.db.findOneAsync({
        collection: collections.SESSION,
        apiKey,
      });
      if (!exists) {
        let data: Session = {
          _id: undefined,
          userId: user._id,
          apiKey,
          hasExpired: false,
          expiredAt: null,
          expireReason: null,
          createdAt: Date.now(),
          updatedAt: Date.now()
        }
        session = await this.db.insertAsync({
          collection: collections.SESSION,
          ...data
        });
      }
      throwOnFalsy(
        DeveloperError,
        safetyCap--,
        "API_KEY_CREATION_FAILED",
        "Timed out"
      );
    } while (exists);

    return { session: session!, apiKey };
  }


  async listSessionsByUserIdOrFail(userId: string): Promise<Session[]> {
    let sessionList = await this.db.findAsync({
      collection: collections.SESSION,
      userId: userId,
    }).sort({ createdAt: -1 }).limit(miscConstants.SESSION_LIST_ALL_COUNT_LIMIT);
    return sessionList;
  }

  async getSessionByIdOrFail(_id: string): Promise<Session> {
    let session = await this.db.findOneAsync({
      collection: collections.SESSION,
      _id,
    });

    throwOnFalsy(
      UserError,
      session,
      "SESSION_NOT_FOUND",
      "The requested session could not be found."
    );

    return session;
  }

  async getSessionByApiKey(apiKey: string): Promise<Session> {
    return await this.db.findOneAsync({
      collection: collections.SESSION,
      apiKey,
    });
  }

  async expireSessionById(_id: string, message: string) {
    return await this.db.updateAsync(
      {
        collection: collections.SESSION,
        _id,
      },
      {
        $set: {
          hasExpired: true,
          expireReason: `${LOGOUT_MESSAGE_PREFIX}${message}`,
          expiredAt: Date.now(),
          updatedAt: Date.now(),
        },
      }
    );
  }

  async expireAllSessionByUserId(userId: string, message: string) {
    return await this.db.updateAsync(
      {
        collection: collections.SESSION,
        hasExpired: false,
        userId,
      },
      {
        $set: {
          hasExpired: true,
          expireReason: `${FORCE_LOGOUT_MESSAGE_PREFIX}${message}`,
          expiredAt: Date.now(),
          updatedAt: Date.now(),
        },
      },
      {
        multi: true,
      }
    );
  }
}
