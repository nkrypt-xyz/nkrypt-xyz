import Nedb from "@seald-io/nedb";
import * as ExpressCore from "express-serve-static-core";
import Joi from "joi";
import { Generic, SerializedError } from "../global.js";
import {
  CodedError,
  DeveloperError,
  UserError,
} from "../utility/coded-error.js";
import {
  detectHttpStatusCode,
  stringifyErrorObject,
} from "../utility/error-utils.js";
import { Config } from "./config.js";
import { Server } from "./server.js";

const joiValidationOptions = {
  abortEarly: true,
  convert: true,
  allowUnknown: false,
};

abstract class AbstractApi {
  config: Config;
  db: Nedb<any>;
  interimData: {
    userId: string | null;
    user: any;
    apiKey: string | null;
    sessionId: string | null;
  };

  constructor(
    private apiPath: string,
    private server: Server,
    private networkDetails: { ip: string },
    private _expressRequest: ExpressCore.Request,
    private _expressResponse: ExpressCore.Response
  ) {
    this.apiPath = apiPath;
    this.server = server;

    this.db = this.server.db.connection;

    this.config = this.server.config;
    this.networkDetails = networkDetails;

    this.interimData = {
      userId: null,
      user: null,
      apiKey: null,
      sessionId: null,
    };
  }

  // ============================== region: properties - start ==============================

  abstract get isEnabled(): boolean;

  // This is a Joi Schema. If not null, request will be parsed and validated.
  abstract get requestSchema(): Joi.Schema;

  abstract get requiresAuthentication(): boolean;

  abstract handle(body: Generic): Promise<Generic>;

  // ============================== region: properties - end ==============================
  // ============================== region: request processing - start ==============================

  async _composeAndValidateSchema(body: Generic) {
    let schema = this.requestSchema;

    // Note: Throws ValidationError
    let validatedBody = await schema.validateAsync(body, joiValidationOptions);
    return validatedBody;
  }

  async _authenticate() {
    return await dispatch.authService.authenticate(this._expressRequest);
  }

  async _preHandleJsonPostApi(parsedJsonBody: Generic) {
    try {
      if (!this.isEnabled) {
        throw new DeveloperError(
          "API_DISABLED",
          "This action has been disabled by the developers."
        );
      }

      let body: Generic = {};
      if (this.requestSchema !== null) {
        body = await this._composeAndValidateSchema(parsedJsonBody);

        if (this.requiresAuthentication) {
          let authData = await this._authenticate();
          Object.assign(this.interimData, authData);
        }
      }

      let response = await this.handle(body);

      if (typeof response !== "object" || response === null) {
        throw new DeveloperError(
          "DEVELOPER_ERROR",
          "Expected response to be an object."
        );
      }

      // FIXME Better solution
      // eslint-disable-next-line
      // @ts-ignore
      response.hasError = false;

      this._sendResponse(200, response);
    } catch (ex: unknown) {
      // There is no need to log UserErrors since they are always logged as response.
      if (!(ex instanceof UserError) && !(ex instanceof Joi.ValidationError)) {
        logger.error(<Error>ex);
      }

      let [serializedError, errorName] = stringifyErrorObject(<Error>ex);
      let statusCode = detectHttpStatusCode(serializedError, errorName);

      this._sendResponse(statusCode, {
        hasError: true,
        error: serializedError,
      });
    }
  }

  _sendResponse(statusCode: number, data: Generic) {
    logger.log(
      statusCode,
      this.apiPath,
      (this._expressRequest as Generic).uuid,
      "\n" + JSON.stringify(data, null, 2)
    );
    this._expressResponse.status(statusCode).send(data);
  }

  // ============================== region: request processing - end ==============================
}

interface IAbstractApi {
  new(
    apiPath: string,
    server: Server,
    networkDetails: { ip: string },
    _expressRequest: ExpressCore.Request,
    _expressResponse: ExpressCore.Response
  ): AbstractApi;
}

export { AbstractApi, IAbstractApi };
