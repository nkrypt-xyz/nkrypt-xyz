import Joi from "joi";
import { AbstractApi } from "../../lib/abstract-api.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  displayName: string;
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
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { displayName } = body;

    await dispatch.userService.updateOwnCommonProperties(
      this.interimData.userId as string,
      displayName
    );

    return {};
  }
}
