import { DatabaseEngine } from "./database-engine.js";
import { AdminService } from "../service/admin-service.js";
import { BucketService } from "../service/bucket-service.js";
import { DirectoryService } from "../service/directory-service.js";
import { FileService } from "../service/file-service.js";
import { SessionService } from "../service/session-service.js";
import { UserService } from "../service/user-service.js";
import { AuthService } from "../service/auth-service.js";
import { BlobService } from "../service/blob-service.js";
import { BlobStorage } from "./blob-storage.js";
import { MetricsService } from "../service/metrics-service.js";

export const prepareServiceDispatch = async (
  db: DatabaseEngine,
  blobStorage: BlobStorage
) => {
  global.dispatch = {
    db,
    config: db.config,
    bucketService: new BucketService(db),
    directoryService: new DirectoryService(db),
    fileService: new FileService(db),
    userService: new UserService(db),
    sessionService: new SessionService(db),
    adminService: new AdminService(db),
    authService: new AuthService(db),
    blobService: new BlobService(db, blobStorage),
    metricsService: new MetricsService(db, blobStorage),
  };
};
