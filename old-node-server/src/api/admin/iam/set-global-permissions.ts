import Joi from "joi";
import { GlobalPermission } from "../../../constant/global-permission.js";
import { Generic } from "../../../global.js";
import { AbstractApi } from "../../../lib/abstract-api.js";
import { requireGlobalPermission } from "../../../utility/access-control-utils.js";
import { throwOnTruthy, UserError } from "../../../utility/coded-error.js";
import { validators } from "../../../validators.js";

type CurrentRequest = {
  userId: string,
  globalPermissions: Record<string, boolean>,
}

export class Api extends AbstractApi {
  get isEnabled(): boolean {
    return true;
  }

  get requiresAuthentication() {
    return true;
  }

  get requestSchema() {
    return Joi.object().required().keys({
      userId: validators.id,
      globalPermissions: Joi.object().required()
    });
  }

  async handle(body: CurrentRequest) {
    let { userId, globalPermissions } = body;

    await requireGlobalPermission(
      this.interimData.user,
      GlobalPermission.MANAGE_ALL_USER
    );

    let user = await dispatch.userService.findUserByIdOrFail(userId);

    Object.keys(GlobalPermission).forEach(key => {
      if (globalPermissions.hasOwnProperty(key)) {
        user.globalPermissions[key] = !!globalPermissions[key];
      }
    });

    await dispatch.adminService.setGlobalPermission(
      userId,
      user.globalPermissions
    );

    return {};
  }
}
