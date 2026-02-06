import Joi from "joi";
import { getDefaultGlobalPermissionsForNewStandardUser, GlobalPermission } from "../../../constant/global-permission.js";
import { Generic } from "../../../global.js";
import { AbstractApi } from "../../../lib/abstract-api.js";
import { requireGlobalPermission } from "../../../utility/access-control-utils.js";
import { throwOnTruthy, UserError } from "../../../utility/coded-error.js";
import { validators } from "../../../validators.js";

type CurrentRequest = {
  displayName: string;
  userName: string;
  password: string;
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
        displayName: validators.displayName,
        userName: validators.userName,
        password: validators.password,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { displayName, userName, password } = body;

    await requireGlobalPermission(
      this.interimData.user,
      GlobalPermission.CREATE_USER
    );

    let exists = await dispatch.userService.findUserByUserName(userName);
    throwOnTruthy(
      UserError,
      exists,
      "USER_NAME_ALREADY_IN_USE",
      `The provided UserName ${userName} is already in use.`
    );

    let user: Generic = await dispatch.adminService.addUser(
      displayName,
      userName,
      password,
      getDefaultGlobalPermissionsForNewStandardUser()
    );

    return { userId: user._id };
  }
}
