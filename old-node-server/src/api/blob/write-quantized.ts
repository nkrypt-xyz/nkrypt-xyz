import * as ExpressCore from "express-serve-static-core";
import { WriteStream } from "fs";
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

export const blobWriteQuantizedApiPath = "/api/blob/write-quantized/:bucketId/:fileId/:blobId/:offset/:shouldEnd";

let schema = Joi.object().required().keys({
  bucketId: validators.id,
  fileId: validators.id,
  blobId: Joi.string().allow(null, "null").length(16).required(),
  offset: Joi.number().min(0).required(),
  shouldEnd: Joi.boolean().required()
});

export const blobWriteQuantizedApiHandler = async (
  req: ExpressCore.Request,
  res: ExpressCore.Response
) => {
  try {
    let { apiKey, userId, sessionId, user } =
      await dispatch.authService.authenticate(req);

    let { bucketId, fileId, blobId, offset, shouldEnd } = await schema.validateAsync(req.params);
    blobId = (blobId === "null" ? null : blobId);

    if (offset > dispatch.config.blobStorage.maxFileSizeBytes) {
      throw new CodedError("BLOB_SIZE_EXCEEDS_LIMIT", "Rejected attempt to write file larger than allowed")
    }

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

    let blob, fileStream: WriteStream;
    if (blobId === null) {
      ({ blob, stream: fileStream } =
        await dispatch.blobService.createInProgressBlob(bucketId, fileId, cryptoMetaHeaderContent, userId));
    } else {
      ({ blob, stream: fileStream } =
        await dispatch.blobService.getInProgressBlob(bucketId, fileId, blobId, offset));
    }

    try {
      let byteCountWrapper = { byteCount: 0 };
      await pipeline(
        req,
        dispatch.blobService.createStandardSizeLimiter(offset, byteCountWrapper),
        fileStream
      );

      if (shouldEnd) {
        await dispatch.blobService.markBlobAsFinished(bucketId, fileId, blob._id!);

        await dispatch.blobService.removeAllOtherBlobs(
          bucketId,
          fileId,
          blob._id!
        );
      }

      res.send({
        hasError: false,
        blobId: blob._id,
        bytesTransfered: byteCountWrapper.byteCount
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
