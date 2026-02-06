import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import {
  ensureFileBelongsToBucket,
  requireBucketAuthorizationByBucketId,
} from "../../utility/access-control-utils.js";
import {
  throwOnFalsy,
  throwOnTruthy,
  UserError,
} from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  encryptedMetaData: string;
  bucketId: string;
  fileId: string;
};

export class Api extends AbstractApi {
  get isEnabled(): boolean {
    return true;
  }

  get requiresAuthentication() {
    return true;
  }

  get requestSchema() {
    return Joi.object()
      .keys({
        encryptedMetaData: validators.encryptedMetaData,
        bucketId: validators.id,
        fileId: validators.id,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { encryptedMetaData, bucketId, fileId } = body;

    await ensureFileBelongsToBucket(bucketId, fileId);

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MANAGE_CONTENT
    );

    await dispatch.fileService.setFileEncryptedMetaData(
      bucketId,
      fileId,
      encryptedMetaData
    );

    return {};
  }
}
