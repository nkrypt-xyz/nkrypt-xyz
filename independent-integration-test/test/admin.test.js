import Joi from "joi";

import {
  callHappyPostJsonApi,
  callHappyPostJsonApiWithAuth,
  callPostJsonApi,
  validateObject,
} from "./testlib/common-api-test-utils.js";

import { userAssertion, errorOfCode } from "./testlib/common-test-schema.js";

import { validators } from "../dist/validators.js";

const DEFAULT_USER_NAME = "admin";
const DEFAULT_PASSWORD = "PleaseChangeMe@YourEarliest2Day";

const TEST_USER_USER_NAME = "testuser1-" + Date.now();
const TEST_USER_DISPLAY_NAME = "Test User 1";
const TEST_USER_PASSWORD = "ExamplePassword";
const TEST_USER_OVERWRITTEN_PASSWORD = "ExamplePassword1";


let vars = {
  apiKey: null,
  newUserId: null
};

describe("Admin Suite", () => {
  test("(user/login): Preparational", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: DEFAULT_USER_NAME,
      password: DEFAULT_PASSWORD,
    });

    await validateObject(data, userAssertion);

    vars.apiKey = data.apiKey;
  });

  test("(admin/iam/add-user): Create new user", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/admin/iam/add-user",
      {
        displayName: TEST_USER_DISPLAY_NAME,
        userName: TEST_USER_USER_NAME,
        password: TEST_USER_PASSWORD,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      userId: validators.id,
    });

    vars.newUserId = data.userId;
  });

  test("(admin/iam/set-global-permissions): Set Global Permissions", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/admin/iam/set-global-permissions",
      {
        userId: vars.newUserId,
        globalPermissions: {
          CREATE_USER: true,
        }
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy
    });
  });

  test("(user/login): Ensure newly created user can log in and has the newly set permission", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: TEST_USER_USER_NAME,
      password: TEST_USER_PASSWORD,
    });

    await validateObject(data, userAssertion);

    expect(data.user.globalPermissions.CREATE_USER).toBe(true);
  });

  test("(admin/iam/set-banning-status): Set Banning Status", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/admin/iam/set-banning-status",
      {
        userId: vars.newUserId,
        isBanned: true
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy
    });
  });

  test("(user/login): Ensure user cannot log in", async () => {
    const response = await callPostJsonApi("/user/login", {
      userName: TEST_USER_USER_NAME,
      password: TEST_USER_PASSWORD,
    });
    expect(response.status).toEqual(403);
    let data = await response.json();

    await validateObject(data, {
      hasError: validators.hasErrorTruthy,
      error: errorOfCode("USER_BANNED"),
    });
  });

  test("(admin/iam/set-banning-status): Set Banning Status", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/admin/iam/set-banning-status",
      {
        userId: vars.newUserId,
        isBanned: false
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy
    });
  });

  test("(admin/iam/overwrite-user-password): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/admin/iam/overwrite-user-password",
      {
        userId: vars.newUserId,
        newPassword: TEST_USER_OVERWRITTEN_PASSWORD
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy
    });
  });


  test("(user/login): Ensure user can log in with overwritten password", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: TEST_USER_USER_NAME,
      password: TEST_USER_OVERWRITTEN_PASSWORD,
    });

    await validateObject(data, userAssertion);
  });
});
