import Nedb from "@seald-io/nedb";
import collections from "../constant/collections.js";
import constants from "../constant/common-constants.js";
import { ensureDir, resolvePath } from "../utility/file-utils.js";
import { Config } from "./config.js";

class DatabaseEngine {
  public config: Config;
  private _dir: string;
  private _dbFilePath: string;
  connection!: Nedb;

  constructor(config: Config) {
    this.config = config;
    this._dir = resolvePath(config.database.dir);
    this._dbFilePath = resolvePath(
      config.database.dir,
      constants.database.CORE_FILE_NAME
    );
  }

  async init() {
    ensureDir(this._dir);
    await this.backup();

    this.connection = new Nedb({ filename: this._dbFilePath });

    await this.connection.loadDatabaseAsync();

    await this.setupInternalData();
  }

  async setupInternalData() {
    await this.connection.compactDatafileAsync();
    await this.connection.updateAsync(
      { collection: collections.SYSTEM },
      { $inc: { timesApplicationRan: 1 } },
      { upsert: true }
    );
  }

  async backup() {
    // TODO backup
  }
}

export { DatabaseEngine };
