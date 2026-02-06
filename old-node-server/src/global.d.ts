/* eslint-disable no-var */
import { Logger } from "./lib/logger";
import { AdminService } from "./service/admin-service";
import { BucketService } from "./service/bucket-service";
import { DirectoryService } from "./service/directory-service";
import { FileService } from "./service/file-service";
import { SessionService } from "./service/session-service";
import { UserService } from "./service/user-service";
import { AuthService } from "./service/auth-service";
import { BlobService } from "./service/blob-service";
import { DatabaseEngine } from "./lib/database-engine.js";
import { MetricsService } from "./service/metrics-service.js";
import { Config } from "./lib/config.js";

declare global {
  var logger: Logger;
  var dispatch: {
    db: DatabaseEngine;
    config: Config;
    userService: UserService;
    sessionService: SessionService;
    adminService: AdminService;
    bucketService: BucketService;
    directoryService: DirectoryService;
    fileService: FileService;
    authService: AuthService;
    blobService: BlobService;
    metricsService: MetricsService;
  };
}

// We type-alias any as Generic to easily mark improvement scopes without adding comments
type Generic = any;

type SerializedError = {
  code: string;
  message: string;
  details: any;
};
