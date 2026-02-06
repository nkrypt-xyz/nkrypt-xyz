import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { Generic } from "../../global.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import { requireBucketAuthorizationByBucketId } from "../../utility/access-control-utils.js";
import {
  throwOnFalsy,
  throwOnTruthy,
  UserError,
} from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  bucketId: string;
  metaData: Record<string, Generic>;
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
        bucketId: validators.id,
        metaData: validators.metaData,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { bucketId, metaData } = body;

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MODIFY
    );

    let bucket = await dispatch.bucketService.findBucketById(bucketId);
    throwOnFalsy(
      UserError,
      bucket,
      "BUCKET_NOT_FOUND",
      `The requested bucket does not exist.`
    );

    await dispatch.bucketService.setBucketMetaData(bucketId, metaData);

    return {};
  }
}
