import Joi from "joi";
import { Generic } from "../../global.js";
import { AbstractApi } from "../../lib/abstract-api.js";
import { extract } from "../../utility/misc-utils.js";
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
    return Joi.object().required();
  }

  async handle(body: CurrentRequest) {
    let sessionList = await dispatch.sessionService.listSessionsByUserIdOrFail(
      this.interimData.userId as string,
    );

    sessionList.forEach(session => {
      (session as Generic).isCurrentSession = session._id === this.interimData.sessionId;
    });

    sessionList = sessionList.map(session => extract(session,
      ['isCurrentSession', 'hasExpired', 'expireReason', 'createdAt', 'expiredAt']));

    return { sessionList };
  }
}
