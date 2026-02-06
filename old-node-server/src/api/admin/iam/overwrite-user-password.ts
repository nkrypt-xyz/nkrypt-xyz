import Joi from "joi";
import { GlobalPermission } from "../../../constant/global-permission.js";
import { Generic } from "../../../global.js";
import { AbstractApi } from "../../../lib/abstract-api.js";
import { requireGlobalPermission } from "../../../utility/access-control-utils.js";
import { throwOnTruthy, UserError } from "../../../utility/coded-error.js";
import { calculateHashOfString } from "../../../utility/security-utils.js";
import { validators } from "../../../validators.js";

type CurrentRequest = {
  userId: string,
  newPassword: string;
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
      newPassword: validators.password,
    });
  }

  async handle(body: CurrentRequest) {
    let { userId, newPassword } = body;

    await requireGlobalPermission(
      this.interimData.user,
      GlobalPermission.MANAGE_ALL_USER
    );

    await dispatch.userService.findUserByIdOrFail(userId);

    let newPasswordBlock = calculateHashOfString(newPassword);

    await dispatch.userService.updateUserPassword(
      userId,
      newPasswordBlock
    );

    await dispatch.sessionService.expireAllSessionByUserId(
      userId,
      `All sessions expired due to password being overwritten by admin ${this.interimData.user.userName}.`
    );
    return {};
  }
}
