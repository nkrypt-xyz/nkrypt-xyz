import Joi from "joi";
import { AbstractApi } from "../../lib/abstract-api.js";

type CurrentRequest = {
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
      .keys({})
      .required();
  }

  async handle(body: CurrentRequest) {
    let disk = await dispatch.metricsService.getDiskUsage();

    return { disk };
  }
}
