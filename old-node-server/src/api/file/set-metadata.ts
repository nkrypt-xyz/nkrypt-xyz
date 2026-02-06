import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { Generic } from "../../global.js";
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
  metaData: Record<string, Generic>;
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
        metaData: validators.metaData,
        bucketId: validators.id,
        fileId: validators.id,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { metaData, bucketId, fileId } = body;

    await ensureFileBelongsToBucket(bucketId, fileId);

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MANAGE_CONTENT
    );

    await dispatch.fileService.setFileMetaData(bucketId, fileId, metaData);

    return {};
  }
}
