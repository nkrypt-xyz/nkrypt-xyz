import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import {
  ensureFileBelongsToBucket,
  requireBucketAuthorizationByBucketId,
} from "../../utility/access-control-utils.js";
import { throwOnFalsy, UserError } from "../../utility/coded-error.js";
import { strip } from "../../utility/misc-utils.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
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
        bucketId: validators.id,
        fileId: validators.id,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { bucketId, fileId } = body;

    await ensureFileBelongsToBucket(bucketId, fileId);

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.VIEW_CONTENT
    );

    let file = await dispatch.fileService.findFileById(bucketId, fileId);

    throwOnFalsy(
      UserError,
      file,
      "FILE_NOT_IN_BUCKET",
      `Given file does not belong to the given bucket`
    );

    strip(file, ["collection"]);
    return { file };
  }
}
