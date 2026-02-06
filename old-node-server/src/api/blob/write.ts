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
import { CodedError, DeveloperError, UserError } from "../../utility/coded-error.js";
import { prepareAndSendCustomApiErrorResponse } from "../../utility/error-response-utils.js";
import {
  detectHttpStatusCode,
  stringifyErrorObject,
} from "../../utility/error-utils.js";
import {
  createDelayerTransformStream,
} from "../../utility/stream-utils.js";
import { validators } from "../../validators.js";

const pipeline = promisify(stream.pipeline);

export const blobWriteApiPath = "/api/blob/write/:bucketId/:fileId";

let schema = Joi.object().required().keys({
  bucketId: validators.id,
  fileId: validators.id,
});

export const blobWriteApiHandler = async (
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
      BucketPermission.MANAGE_CONTENT
    );

    let cryptoMetaHeaderContent = req.headers[constants.webServer.BLOB_API_CRYPTO_META_HEADER_NAME];

    if (!cryptoMetaHeaderContent || (typeof cryptoMetaHeaderContent === "object" && Array.isArray(cryptoMetaHeaderContent))) {
      throw new DeveloperError("CRYPTO_META_HEADER_INVALID", `Provided ${constants.webServer.BLOB_API_CRYPTO_META_HEADER_NAME} header is invalid`);
    }

    let { blob, stream: fileStream } =
      await dispatch.blobService.createInProgressBlob(bucketId, fileId, cryptoMetaHeaderContent, userId);

    try {
      await pipeline(
        req,
        dispatch.blobService.createStandardSizeLimiter(),
        fileStream
      );

      await dispatch.blobService.markBlobAsFinished(bucketId, fileId, blob._id!);

      await dispatch.fileService.setFileContentUpdateAt(bucketId, fileId, Date.now());

      await dispatch.blobService.removeAllOtherBlobs(
        bucketId,
        fileId,
        blob._id!
      );

      res.send({
        hasError: false,
        blobId: blob._id,
      });
    } catch (err) {
      logger.error(err as Error);

      await dispatch.blobService.markBlobAsErroneous(
        bucketId,
        fileId,
        blob._id!
      );

      throw err;
    }
  } catch (ex) {
    prepareAndSendCustomApiErrorResponse(ex, req, res);
  }
};
