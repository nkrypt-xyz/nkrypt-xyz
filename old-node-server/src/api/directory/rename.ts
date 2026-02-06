import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import {
  ensureDirectoryBelongsToBucket,
  requireBucketAuthorizationByBucketId,
} from "../../utility/access-control-utils.js";
import {
  throwOnFalsy,
  throwOnTruthy,
  UserError,
} from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  name: string;
  bucketId: string;
  directoryId: string;
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
        directoryId: validators.id,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { name, bucketId, directoryId } = body;

    await ensureDirectoryBelongsToBucket(bucketId, directoryId);

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MANAGE_CONTENT
    );

    let existingDirectory = await dispatch.directoryService.findDirectoryById(
      bucketId,
      directoryId
    );

    // No need to do anything. It's the same name and same bucket
    if (
      existingDirectory &&
      existingDirectory._id === directoryId &&
      existingDirectory.name === name
    ) {
      return {};
    }

    throwOnFalsy(
      UserError,
      existingDirectory,
      "DIRECTORY_NOT_IN_BUCKET",
      `Given directory does not belong to the given bucket`
    );

    await dispatch.directoryService.setDirectoryName(
      bucketId,
      directoryId,
      name
    );

    return {};
  }
}
