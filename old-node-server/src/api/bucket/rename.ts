import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import { requireBucketAuthorizationByBucketId } from "../../utility/access-control-utils.js";
import { throwOnTruthy, UserError } from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  name: string;
  bucketId: string;
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
        name: validators.bucketName,
        bucketId: validators.id,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { name, bucketId } = body;

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MODIFY
    );

    let existingBucket = await dispatch.bucketService.findBucketByName(name);

    // No need to do anything. It's the same name and same bucket
    if (
      existingBucket &&
      existingBucket._id === bucketId &&
      existingBucket.name === name
    ) {
      return {};
    }

    throwOnTruthy(
      UserError,
      existingBucket,
      "BUCKET_NAME_ALREADY_IN_USE",
      `A bucket with the provided name ${name} already exists.`
    );

    await dispatch.bucketService.setBucketName(bucketId, name);

    return {};
  }
}
