import Joi from "joi";
import { GlobalPermission } from "../../constant/global-permission.js";
import { Generic } from "../../global.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import { requireGlobalPermission } from "../../utility/access-control-utils.js";
import { throwOnTruthy, UserError } from "../../utility/coded-error.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  name: string;
  cryptSpec: string;
  cryptData: string;
  metaData: Generic;
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
        cryptSpec: validators.cryptSpec,
        cryptData: validators.cryptData,
        metaData: validators.metaData,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { name, cryptData, cryptSpec, metaData } = body;

    await requireGlobalPermission(
      this.interimData.user,
      GlobalPermission.CREATE_BUCKET
    );

    let exists = await dispatch.bucketService.findBucketByName(name);
    throwOnTruthy(
      UserError,
      exists,
      "BUCKET_NAME_ALREADY_IN_USE",
      `A bucket with the provided name ${name} already exists.`
    );

    let bucket: Generic = await dispatch.bucketService.createBucket(
      name,
      cryptSpec,
      cryptData,
      metaData,
      this.interimData.userId as string
    );

    let directory: Generic = await dispatch.directoryService.createDirectory(
      `${name} Root`,
      bucket._id,
      metaData,
      "{}",
      this.interimData.userId as string,
      null
    );

    return { bucketId: bucket._id, rootDirectoryId: directory._id };
  }
}
