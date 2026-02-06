import pathlib from "path";
import { blobReadApiHandler, blobReadApiPath } from "./api/blob/read.js";
import { blobWriteApiHandler, blobWriteApiPath } from "./api/blob/write.js";
import { blobWriteQuantizedApiHandler, blobWriteQuantizedApiPath } from "./api/blob/write-quantized.js";
import constants from "./constant/common-constants.js";
import { BlobStorage } from "./lib/blob-storage.js";
import { Config } from "./lib/config.js";
import { DatabaseEngine } from "./lib/database-engine.js";
import { Logger } from "./lib/logger.js";
import { Server } from "./lib/server.js";
import { prepareServiceDispatch } from "./lib/service-dispatch.js";
import { apiNameList } from "./routes.js";
import { appRootDirPath, toFileUrl } from "./utility/file-utils.js";

// We initiate logger and inject it into global so that it is usable everywhere.
global.logger = new Logger({
  switches: {
    debug: true,
    log: true,
    important: true,
    warning: true,
    error: true,
    urgent: true,
  },
});
await logger.init();

export class NkWebServerProgram {
  db!: DatabaseEngine;
  server!: Server;
  blobStorage!: BlobStorage;
  config!: Config;

  async start(config: Config) {
    try {
      this.config = config;
      await this._initialize();
    } catch (ex) {
      logger.log("Error was propagated to root level. Throwing again.");
      throw ex;
    }
  }

  async _initialize() {
    this.db = new DatabaseEngine(this.config);
    await this.db.init();

    this.blobStorage = new BlobStorage(this.config);
    await this.blobStorage.init();

    await prepareServiceDispatch(this.db, this.blobStorage);

    await dispatch.adminService.createDefaultAdminAccountIfNotPresent();

    this.server = new Server(this.config, this.db);
    await this.server.prepare();
    await this._registerEndpoints();
    await this.server.start();
  }

  async _registerEndpoints() {
    logger.log("(server)> Dynamically registering APIs");

    await Promise.all(
      apiNameList.map(async (name) => {
        let path = toFileUrl(
          pathlib.join(appRootDirPath, constants.api.CORE_API_DIR, `${name}.js`)
        );

        let apiModule = await import(path);
        await this.server.registerJsonPostApi(name, apiModule.Api);
      })
    );

    await this.server.registerCustomHandler(
      blobWriteApiPath,
      blobWriteApiHandler
    );

    await this.server.registerCustomHandler(
      blobReadApiPath,
      blobReadApiHandler
    );

    await this.server.registerCustomHandler(
      blobWriteQuantizedApiPath,
      blobWriteQuantizedApiHandler
    );
  }
}

process.on("uncaughtException", function (err) {
  console.log("Suppressing uncaughtException");
  console.log("uncaughtException message:", JSON.stringify(err));
  console.log("uncaughtException stack:", err.stack);
  console.error(err);
});

