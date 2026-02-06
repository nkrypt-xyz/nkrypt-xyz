import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { Generic } from "../../global.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import { requireBucketAuthorizationByBucketId } from "../../utility/access-control-utils.js";
import { throwOnTruthy, UserError } from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  name: string;
  bucketId: string;
  parentDirectoryId: string;
  metaData: Generic;
  encryptedMetaData: string;
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
        name: validators.directoryName,
        bucketId: validators.id,
        parentDirectoryId: validators.id,
        encryptedMetaData: validators.encryptedMetaData,
        metaData: validators.metaData,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { name, bucketId, parentDirectoryId, encryptedMetaData, metaData } =
      body;

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MANAGE_CONTENT
    );

    let exists = await dispatch.directoryService.findDirectoryByNameAndParent(
      name,
      bucketId,
      parentDirectoryId
    );
    throwOnTruthy(
      UserError,
      exists,
      "DIRECTORY_NAME_ALREADY_IN_USE",
      `A directory with the provided name ${name} already exists in the parent directory.`
    );

    let directory: Generic = await dispatch.directoryService.createDirectory(
      name,
      bucketId,
      metaData,
      encryptedMetaData,
      this.interimData.userId as string,
      parentDirectoryId
    );

    return { directoryId: directory._id };
  }
}
