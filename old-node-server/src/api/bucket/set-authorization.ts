import Joi from "joi";
import { BucketPermission } from "../../constant/bucket-permission.js";
import { miscConstants } from "../../constant/misc-constants.js";
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
  targetUserId: string;
  bucketId: string;
  permissionsToSet: Record<string, boolean>;
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
        targetUserId: validators.id,
        bucketId: validators.id,
        permissionsToSet: validators.partialBucketPermissions,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { targetUserId, bucketId, permissionsToSet } = body;

    await requireBucketAuthorizationByBucketId(
      this.interimData.userId as string,
      bucketId,
      BucketPermission.MANAGE_AUTHORIZATION
    );

    let bucket = await dispatch.bucketService.findBucketById(bucketId);
    throwOnFalsy(
      UserError,
      bucket,
      "BUCKET_NOT_FOUND",
      `The requested bucket does not exist.`
    );

    let authorization = bucket.bucketAuthorizations.find(
      (authorization: Generic) => authorization.userId === targetUserId
    );

    if (!authorization) {
      await dispatch.bucketService.authorizeUserWithAllPermissionsForbidden(
        bucketId,
        targetUserId,
        miscConstants.BUCKET_NEW_DEFAULT_AUTHORIZATION_MESSAGE_fn(this.interimData.user.userName)
      );
      bucket = await dispatch.bucketService.findBucketById(bucketId);
      authorization = bucket.bucketAuthorizations.find(
        (authorization: Generic) => authorization.userId === targetUserId
      );
    }

    Object.keys(permissionsToSet).forEach((permission) => {
      authorization!.permissions[permission] = permissionsToSet[permission];
    });
    await dispatch.bucketService.setAuthorizationList(
      bucketId,
      bucket.bucketAuthorizations
    );

    return {};
  }
}
