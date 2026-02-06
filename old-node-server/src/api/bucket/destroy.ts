import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import { requireBucketAuthorizationByBucketId } from "../../utility/access-control-utils.js";
import { throwOnFalsy, UserError } from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  // We want to make the user type in the name to ensure intention
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
      BucketPermission.DESTROY
    );

    let existingBucket = await dispatch.bucketService.findBucketByName(name);

    throwOnFalsy(
      UserError,
      existingBucket,
      "BUCKET_NOT_FOUND",
      `The requested bucket does not exist.`
    );

    throwOnFalsy(
      UserError,
      existingBucket.name === name,
      "BUCKET_NAME_MISMATCH",
      `You have incorrectly entered the bucket name.`
    );

    await dispatch.bucketService.removeBucket(bucketId);

    let directory = await dispatch.directoryService.findRootDirectoryByBucketId(bucketId);
    if (directory) {
      await dispatch.directoryService.deleteDirectory(bucketId, directory._id!);
      dispatch.directoryService.deleteDirectoryAndChildrenInTheBackground(bucketId, directory)
        .then(() => {
          logger.log("Deletion of directory and children finished in the background.");
        })
        .catch(ex => {
          logger.log("Deletion of directory and children failed in the background.");
          logger.error(ex);
        });
    }

    return {};
  }
}
