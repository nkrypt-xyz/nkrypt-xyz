import Joi from "joi";

import {
  callHappyPostJsonApi,
  callHappyPostJsonApiWithAuth,
  validateObject,
} from "./testlib/common-api-test-utils.js";

import {
  directorySchema,
  bucketListSchema,
  userAssertion,
  fileSchema,
} from "./testlib/common-test-schema.js";

import { validators } from "../dist/validators.js";

const DEFAULT_USER_NAME = "admin";
const DEFAULT_PASSWORD = "PleaseChangeMe@YourEarliest2Day";

const TEST_BUCKET_NAME = "testBucket1-" + Date.now();
const TEST_BUCKET_NEW_NAME = "testBucket1Renamed-" + Date.now();
const TEST_DIRECTORY_A_NAME = "testDirA-" + Date.now();
const TEST_DIRECTORY_A_A_NAME = "testDirAA-" + Date.now();
const TEST_DIRECTORY_B_NAME = "testDirB-" + Date.now();
const TEST_DIRECTORY_B_NAME_ALT = "testDirBAlt-" + Date.now();

const TEST_USER_USER_NAME = "testuser1-" + Date.now();
const TEST_USER_DISPLAY_NAME = "Test User 1";
const TEST_USER_PASSWORD = "ExamplePassword";

const TEST_INITIAL_METADATA = { createdFromApp: "Integration testing" };
const TEST_INITIAL_ENCRYPTED_METADATA = "PLACEHOLDER";
const TEST_NEW_METADATA_PART = { example: "value" };
const TEST_NEW_ENCRYPTED_METADATA = "PLACEHOLDER2";

let vars = {
  apiKey: null,
  bucketId: null,
  rootDirectoryId: null,
  idOfDirectoryA: null,
  idOfDirectoryB: null,
  idOfDirectoryAA: null,
  testUserId: null,
  testUserApiKey: null,
};

describe("Bucket and Directory Suite", () => {
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

    vars.testUserId = data.userId;
  });

  test("(user/login): Preparational", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: TEST_USER_USER_NAME,
      password: TEST_USER_PASSWORD,
    });

    await validateObject(data, userAssertion);

    vars.testUserApiKey = data.apiKey;
  });

  test("(bucket/create): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/create",
      {
        name: TEST_BUCKET_NAME,
        cryptSpec: "V1:AES256",
        cryptData: "PLACEHOLDER",
        metaData: {},
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      bucketId: validators.id,
      rootDirectoryId: validators.id,
    });
  });

  test("(bucket/list): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/list",
      {}
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      bucketList: bucketListSchema,
    });

    let bucket = data.bucketList.find(
      (bucket) => bucket.name === TEST_BUCKET_NAME
    );
    expect(bucket).not.toBeFalsy();

    vars.bucketId = bucket._id;
    vars.rootDirectoryId = bucket.rootDirectoryId;
  });

  test("(directory/create): Bucket1Root/testDirA", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/create",
      {
        name: TEST_DIRECTORY_A_NAME,
        bucketId: vars.bucketId,
        parentDirectoryId: vars.rootDirectoryId,
        encryptedMetaData: TEST_INITIAL_ENCRYPTED_METADATA,
        metaData: TEST_INITIAL_METADATA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directoryId: validators.id,
    });

    vars.idOfDirectoryA = data.directoryId;
  });

  test("(directory/create): Bucket1Root/testDirB", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/create",
      {
        name: TEST_DIRECTORY_B_NAME,
        bucketId: vars.bucketId,
        parentDirectoryId: vars.rootDirectoryId,
        encryptedMetaData: TEST_INITIAL_ENCRYPTED_METADATA,
        metaData: TEST_INITIAL_METADATA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directoryId: validators.id,
    });

    vars.idOfDirectoryB = data.directoryId;
  });

  test("(directory/create) Bucket1Root/testDirA/testDirAA", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/create",
      {
        name: TEST_DIRECTORY_A_A_NAME,
        bucketId: vars.bucketId,
        parentDirectoryId: vars.idOfDirectoryA,
        encryptedMetaData: TEST_INITIAL_ENCRYPTED_METADATA,
        metaData: TEST_INITIAL_METADATA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directoryId: validators.id,
    });
  });

  test("(directory/get): Bucket1Root/*", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.rootDirectoryId,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(2);
  });

  test("(directory/get): Bucket1Root/testDirA/*", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(1);
  });

  test("(directory/get): Bucket1Root/testDirA/testDirAA/*", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().optional().max(0),
      childFileList: Joi.array().optional().max(0),
    });

    expect(data.childDirectoryList.length).toEqual(0);
  });

  test("(directory/rename) Bucket1Root/testDirB => Bucket1Root/testDirBAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/rename",
      {
        name: TEST_DIRECTORY_B_NAME_ALT,
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): Bucket1Root/testDirB/* Ensure rename worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required(),
      childFileList: Joi.array().optional().max(0),
    });

    expect(data.directory.name).toBe(TEST_DIRECTORY_B_NAME_ALT);
  });

  test("(directory/set-metadata) Bucket1Root/testDirBAlt", async () => {
    let metaData = Object.assign(
      {},
      TEST_INITIAL_METADATA,
      TEST_NEW_METADATA_PART
    );
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/set-metadata",
      {
        metaData: metaData,
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/set-metadata) Bucket1Root/testDirBAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/set-encrypted-metadata",
      {
        encryptedMetaData: TEST_NEW_ENCRYPTED_METADATA,
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): Bucket1Root/testDirBAlt/* Ensure metaData and encryptedMetaData worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required(),
      childFileList: Joi.array().required(),
    });

    let metaData = Object.assign(
      {},
      TEST_INITIAL_METADATA,
      TEST_NEW_METADATA_PART
    );
    expect(data.directory.encryptedMetaData).toEqual(
      TEST_NEW_ENCRYPTED_METADATA
    );
    expect(data.directory.metaData).toStrictEqual(metaData);
  });

  test("(directory/move) Bucket1Root/testDirBAlt => Bucket1Root/testDirA/testDirBAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/move",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
        newParentDirectoryId: vars.idOfDirectoryA,
        newName: TEST_DIRECTORY_B_NAME_ALT,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): Bucket1Root/testDirA/* Ensure move worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(2);

    expect(
      data.childDirectoryList.find(
        (directory) => directory._id === vars.idOfDirectoryB
      )
    ).toBeTruthy();
  });

  test("(directory/delete) Bucket1Root/testDirA/testDirBAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/delete",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryB,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): Bucket1Root/testDirA/* Ensure delete worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.idOfDirectoryA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(1);

    expect(
      data.childDirectoryList.find(
        (directory) => directory._id === vars.idOfDirectoryB
      )
    ).toBeFalsy();
  });

  test("(bucket/rename): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/rename",
      {
        name: TEST_BUCKET_NEW_NAME,
        bucketId: vars.bucketId,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(bucket/list): Ensure rename worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/list",
      {}
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      bucketList: bucketListSchema,
    });

    let bucket1 = data.bucketList.find(
      (bucket) => bucket.name === TEST_BUCKET_NAME
    );
    expect(bucket1).toBeFalsy();

    let bucket2 = data.bucketList.find(
      (bucket) =>
        bucket.name === TEST_BUCKET_NEW_NAME && bucket._id == vars.bucketId
    );
    expect(bucket2).not.toBeFalsy();
  });

  test("(bucket/set-metadata): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/set-metadata",
      {
        metaData: {
          exampleMetaData: "exampleValue",
        },
        bucketId: vars.bucketId,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(bucket/list): Ensure setting metaData worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/list",
      {}
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      bucketList: bucketListSchema,
    });

    let bucket = data.bucketList.find((bucket) => bucket._id == vars.bucketId);
    expect(bucket).not.toBeFalsy();
    expect(bucket.metaData).toHaveProperty("exampleMetaData");
    expect(bucket.metaData.exampleMetaData).toEqual("exampleValue");
  });

  test("(bucket/set-authorization): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/set-authorization",
      {
        targetUserId: vars.testUserId,
        bucketId: vars.bucketId,
        permissionsToSet: {
          VIEW_CONTENT: true,
        },
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): Bucket1Root/* as TestUser to ensure bucket authorization worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.testUserApiKey,
      "/directory/get",
      {
        bucketId: vars.bucketId,
        directoryId: vars.rootDirectoryId,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      directory: directorySchema,
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(1);
  });

  test("(bucket/destroy): Affirmative", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/destroy",
      {
        name: TEST_BUCKET_NEW_NAME,
        bucketId: vars.bucketId,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(bucket/list): Ensure destroy worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/bucket/list",
      {}
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      bucketList: bucketListSchema,
    });

    let bucket1 = data.bucketList.find((bucket) => bucket._id == vars.bucketId);
    expect(bucket1).toBeFalsy();
  });

  // eof
});
