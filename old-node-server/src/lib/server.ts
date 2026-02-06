import bodyParser from "body-parser";
import express from "express";
import * as ExpressCore from "express-serve-static-core";
import http from "http";
import https from "https";
import pathlib from "path";
import constants from "../constant/common-constants.js";
import { ErrorCode } from "../constant/error-codes.js";
import { Generic } from "../global.js";
import { DeveloperError } from "../utility/coded-error.js";
import { prepareSslDetails } from "../utility/ssl-utils.js";
import { generateUuid } from "../utility/string-utils.js";
import { joinUrlParts } from "../utility/url-utils.js";
import { IAbstractApi } from "./abstract-api.js";
import { Config } from "./config.js";
import { DatabaseEngine } from "./database-engine.js";
import { registerStaticRequestHandlers } from "./static-content.js";

const jsonParser = bodyParser.json({
  limit: "100kb",
});

class Server {
  config: Config;
  db: DatabaseEngine;

  private _expressApp: ExpressCore.Express;
  private _nodeHttpWebServer!: http.Server;
  private _nodeHttpsWebServer!: http.Server;

  private _subContextPath: string;

  constructor(config: Config, db: DatabaseEngine) {
    this.config = config;
    this.db = db;

    this._expressApp = express();

    this._subContextPath = constants.api.CORE_API_SUBCONTEXT_PATH;
  }

  async prepare() {
    this._expressApp.settings["x-powered-by"] = false;
    this._expressApp.set("etag", false);
    this._expressApp.set("trust proxy", true);

    this._expressApp.use((req, res, next) => {
      // Assign UUID for effective tracking
      let uuid = generateUuid();
      (req as Generic).uuid = uuid;

      // Log including UUID
      let description = `HTTP ${req.method} ${req.url} ${
        (req as Generic).uuid
      }`;
      logger.log(description);
      return next();
    });

    // Enable CORS
    this._expressApp.use((req, res, next) => {
      res.header("Access-Control-Allow-Origin", "*");
      res.header(
        "Access-Control-Allow-Headers",
        `Origin, X-Requested-With, Content-Type, Accept, Authorization, ${constants.webServer.BLOB_API_CRYPTO_META_HEADER_NAME}`
      );
      return next();
    });

    // Accept OPTIONS to enable CORS
    this._expressApp.options("*", (req, res) => {
      return res.status(200).send("Ok");
    });
  }

  async start() {
    await this._initializeWebServer();
  }

  async _initializeWebServer() {
    // static content
    registerStaticRequestHandlers(this._expressApp);

    // Finally reject anything not supported
    this._expressApp.all("*", (req, res) => {
      let description = `REJECT ${req.method} ${req.url}`;
      logger.log(description);
      return res.status(400).send("Not supported");
    });

    if (this.config.webServer.http.enabled) {
      await new Promise<void>((accept, reject) => {
        let port = this.config.webServer.http.port;
        
        this._nodeHttpWebServer = http.createServer(this._expressApp);
        this._nodeHttpWebServer.listen(port, () => {
          logger.log("(server)> Http server is listening on port", port);
          return accept();
        });
      });
    }

    if (this.config.webServer.https.enabled) {
      await new Promise<void>((accept, reject) => {
        let port = this.config.webServer.https.port;
        let sslDetails = prepareSslDetails(this.config);

        this._nodeHttpsWebServer = https.createServer(
          sslDetails,
          this._expressApp
        );
        this._nodeHttpsWebServer.listen(port, () => {
          logger.log("(server)> Https server is listening on port", port);
          return accept();
        });
      });
    }
  }

  async registerCustomHandler(
    path: string,
    fn: (req: ExpressCore.Request, res: ExpressCore.Response) => void
  ) {
    this._expressApp.post(path, fn);
  }

  async registerJsonPostApi(path: string, ApiClass: IAbstractApi) {
    if (!ApiClass) {
      throw new DeveloperError(
        ErrorCode.DEVELOPER_ERROR,
        "Expected ApiClass to not be null/undefined."
      );
    }

    let apiPath = joinUrlParts(
      this.config.webServer.contextPath,
      this._subContextPath,
      path
    );

    this._expressApp.post(apiPath, jsonParser, (req, res) => {
      setTimeout(() => {
        logger.log(
          "POST",
          req.url,
          (req as Generic).uuid,
          "\n" + JSON.stringify(req.body, null, 2)
        );

        let api = new ApiClass(apiPath, this, { ip: req.ip }, req, res);

        api._preHandleJsonPostApi(req.body);
      }, constants.webServer.INTENTIONAL_REQUEST_DELAY_MS);
    });

    return { apiPath };
  }
}

export { Server };
