import * as ExpressCore from "express-serve-static-core";
import Joi, { func } from "joi";
import stream from "stream";
import { promisify } from "util";
import { BucketPermission } from "../../constant/bucket-permission.js";
import constants from "../../constant/common-constants.js";
import {
  ensureFileBelongsToBucket,
  requireBucketAuthorizationByBucketId,
} from "../../utility/access-control-utils.js";
import { CodedError, UserError } from "../../utility/coded-error.js";
import { prepareAndSendCustomApiErrorResponse } from "../../utility/error-response-utils.js";
import {
  detectHttpStatusCode,
  stringifyErrorObject,
} from "../../utility/error-utils.js";
import { createSizeLimiterPassthroughStream } from "../../utility/stream-utils.js";
import { validators } from "../../validators.js";

const pipeline = promisify(stream.pipeline);

export const blobReadApiPath = "/api/blob/read/:bucketId/:fileId";

let schema = Joi.object().required().keys({
  bucketId: validators.id,
  fileId: validators.id,
});

export const blobReadApiHandler = async (
  req: ExpressCore.Request,
  res: ExpressCore.Response
) => {
  try {
    let { apiKey, userId, sessionId, user } =
      await dispatch.authService.authenticate(req);

    let { bucketId, fileId } = await schema.validateAsync(req.params);

    await ensureFileBelongsToBucket(bucketId, fileId);

    await requireBucketAuthorizationByBucketId(
      userId,
      bucketId,
      BucketPermission.VIEW_CONTENT
    );

    let blob = await dispatch.blobService.findBlobByBucketIdAndFileId(
      bucketId,
      fileId
    );

    if (!blob) {
      throw new UserError("BLOB_NOT_FOUND", "Desired blob could not be found");
    }

    res.setHeader("Access-Control-Expose-Headers", constants.webServer.BLOB_API_CRYPTO_META_HEADER_NAME);
    res.setHeader("Content-Type", "application/octet-stream");
    res.setHeader(constants.webServer.BLOB_API_CRYPTO_META_HEADER_NAME, blob.cryptoMetaHeaderContent);

    let { readStream: stream, sizeOfStream } = await dispatch.blobService.createReadableStreamFromBlobId(blob._id!);
    res.setHeader("Content-Length", sizeOfStream);

    await pipeline(stream, res);
  } catch (ex) {
    prepareAndSendCustomApiErrorResponse(ex, req, res);
  }
};
