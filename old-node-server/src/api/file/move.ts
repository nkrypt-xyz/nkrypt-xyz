import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import {
  ensureDirectoryBelongsToBucket,
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
  bucketId: string;
  fileId: string;
  newParentDirectoryId: string;
  newName: string;
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
        fileId: validators.id,
        newParentDirectoryId: validators.id,
        newName: validators.bucketName,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { newName, bucketId, fileId, newParentDirectoryId } = body;

    await ensureFileBelongsToBucket(bucketId, fileId);

    await ensureDirectoryBelongsToBucket(bucketId, newParentDirectoryId);

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MANAGE_CONTENT
    );

    let existingFile = await dispatch.fileService.findFileByNameAndParent(
      newName,
      bucketId,
      newParentDirectoryId
    );

    throwOnTruthy(
      UserError,
      existingFile,
      "FILE_NAME_ALREADY_IN_USE",
      `A file with the provided name "${newName}" already exists in the new parent file.`
    );

    await dispatch.fileService.moveFile(
      bucketId,
      fileId,
      newParentDirectoryId,
      newName
    );

    return {};
  }
}
