import Joi from "joi";
import { AbstractApi } from "../../lib/abstract-api.js";

type CurrentRequest = Record<string, never>;

export class Api extends AbstractApi {
  get isEnabled(): boolean {
    return true;
  }

  get requiresAuthentication() {
    return true;
  }

  get requestSchema() {
    return Joi.object().required();
  }

  async handle(body: CurrentRequest) {
    let userList = (await dispatch.userService.listAllUsers()).map((user) => ({
      userName: user.userName,
      displayName: user.displayName,
      _id: user._id,
      isBanned: user.isBanned
    }));

    return { userList };
  }
}
