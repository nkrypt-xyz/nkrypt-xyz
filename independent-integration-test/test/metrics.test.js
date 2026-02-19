import Joi from "joi";

import {
  callHappyPostJsonApi,
  callHappyPostJsonApiWithAuth,
  validateObject,
} from "./testlib/common-api-test-utils.js";

import {
  userAssertion,
} from "./testlib/common-test-schema.js";

import { validators } from "../dist/validators.js";

const DEFAULT_USER_NAME = "admin";
const DEFAULT_PASSWORD = "PleaseChangeMe@YourEarliest2Day";

let vars = {
  apiKey: null,
};

describe("Metrics Suite", () => {
  test("(user/login): Preparational", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: DEFAULT_USER_NAME,
      password: DEFAULT_PASSWORD,
    });

    await validateObject(data, userAssertion);

    vars.apiKey = data.apiKey;
  });

  test("(metrics/get-summary): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/metrics/get-summary", {}
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      disk: Joi.object().required().keys({
        usedBytes: Joi.number().required().min(-1),
        totalBytes: Joi.number().required().min(-1)
      })
    });

  });

  // eof
});
