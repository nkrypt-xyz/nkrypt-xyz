import Nedb from "@seald-io/nedb";
import * as ExpressCore from "express-serve-static-core";
import constants from "../constant/common-constants.js";
import { DatabaseEngine } from "../lib/database-engine.js";
import {
  DeveloperError,
  throwOnFalsy,
  throwOnTruthy,
  UserError
} from "../utility/coded-error.js";

export class AuthService {
  db: Nedb;

  constructor(dbEngine: DatabaseEngine) {
    this.db = dbEngine.connection;
  }

  async authenticate(_expressRequest: ExpressCore.Request) {
    const authorizationHeader = String(
      _expressRequest.headers["authorization"] || ""
    );

    throwOnFalsy(
      UserError,
      authorizationHeader.length > 0,
      "AUTHORIZATION_HEADER_MISSING",
      "Authorization header is missing"
    );

    let parts = authorizationHeader.split(" ");
    throwOnFalsy(
      UserError,
      parts.length === 2 &&
      parts[0].toLowerCase().indexOf("bearer") === 0 &&
      parts[1].length === constants.iam.API_KEY_LENGTH,
      "AUTHORIZATION_HEADER_MALFORMATTED",
      "Authorization header is malformatted"
    );
    let apiKey = parts.pop() as string;

    let session = await dispatch.sessionService.getSessionByApiKey(
      apiKey
    );

    throwOnFalsy(
      UserError,
      session,
      "API_KEY_NOT_FOUND",
      "Your session could not be found. Login again."
    );

    throwOnTruthy(
      UserError,
      session.hasExpired,
      "API_KEY_EXPIRED",
      "Your session has expired. Login again."
    );

    let hasExpired =
      Date.now() - session.createdAt >
      constants.iam.SESSION_VALIDITY_DURATION_MS;
    throwOnTruthy(
      UserError,
      hasExpired,
      "API_KEY_EXPIRED",
      "Your session has expired. Login again."
    );

    let { userId, _id: sessionId } = session;

    let user = await dispatch.userService.findUserByIdOrFail(userId);

    return { apiKey, userId, sessionId, user };
  }
}
