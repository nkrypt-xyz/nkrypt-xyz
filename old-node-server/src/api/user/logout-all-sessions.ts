import Joi from "joi";
import { AbstractApi } from "../../lib/abstract-api.js";
import { validators } from "../../validators.js";

type CurrentRequest = {
  message: string;
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
        message: validators.logoutMessage,
      })
      .required();
  }

  async handle(body: CurrentRequest) {
    let { message } = body;

    await dispatch.sessionService.expireAllSessionByUserId(
      this.interimData.userId as string,
      message
    );

    return {};
  }
}
