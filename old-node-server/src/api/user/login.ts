import Joi from "joi";
import { AbstractApi } from "../../lib/abstract-api.js";
import { throwOnFalsy, throwOnTruthy, UserError } from "../../utility/coded-error.js";
import { extract } from "../../utility/misc-utils.js";
import { compareHashWithString } from "../../utility/security-utils.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  userName: string;
  password: string;
};

export class Api extends AbstractApi {
  get isEnabled(): boolean {
    return true;
  }

  get requiresAuthentication() {
    return false;
  }

  get requestSchema() {
    return Joi.object()
      .keys({
        userName: validators.userName,
        password: validators.password,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { userName, password } = body;

    let user = await dispatch.userService.findUserOrFail(userName);

    throwOnTruthy(UserError,
      user.isBanned,
      "USER_BANNED",
      "You are currently banned from logging in."
    );

    let isPasswordCorrect = compareHashWithString(
      password,
      user.password.salt,
      user.password.hash
    );
    throwOnFalsy(
      UserError,
      isPasswordCorrect,
      "INCORRECT_PASSWORD",
      "The password you have used is not correct."
    );

    let { session, apiKey } =
      await dispatch.sessionService.createNewUniqueSession(user);

    return {
      apiKey,
      user: extract(user, ["_id", "userName", "displayName", "globalPermissions"]),
      session: extract(session, ["_id"]),
    };
  }
}
