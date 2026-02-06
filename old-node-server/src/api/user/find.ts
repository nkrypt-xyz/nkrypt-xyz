import Joi from "joi";
import { AbstractApi } from "../../lib/abstract-api.js";
import { User } from "../../model/core-db-entities.js";
import { validators } from "../../validators.js";

type CurrentRequest = { filters: { by: string, userName: string, userId: string }[], includeGlobalPermissions: boolean };

export class Api extends AbstractApi {
  get isEnabled(): boolean {
    return true;
  }

  get requiresAuthentication() {
    return true;
  }

  get requestSchema() {
    return Joi.object().keys({
      filters: Joi.array().items(
        Joi.object().keys({
          by: Joi.string().valid("userName", 'userId'),
          userName: validators.userName.optional().allow(null),
          userId: validators.id.optional().allow(null)
        })
      ),
      includeGlobalPermissions: Joi.boolean().required()
    });
  }

  async handle(body: CurrentRequest) {
    let userIdList = [];
    let userNameList = [];
    for (let filter of body.filters) {
      if (filter.by === "userId" && filter.userId !== null) {
        userIdList.push(filter.userId);
      } else if (filter.by === "userName" && filter.userName !== null) {
        userNameList.push(filter.userName)
      }
    }

    let userList = (await dispatch.userService.queryUsers(userIdList, userNameList)).map((user: User) => {
      let _user = {
        _id: user._id,
        userName: user.userName,
        displayName: user.displayName,
        isBanned: user.isBanned
      }
      if (body.includeGlobalPermissions) {
        (<any>_user).globalPermissions = user.globalPermissions
      }

      return _user;
    });

    return { userList };
  }
}
