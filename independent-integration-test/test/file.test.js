/* What do we test?

1. We create a bucket and dir A and dir B
2. We create file P in dir A
3. We rename P to PAlt
4. We move P to dir B
5. We set some metadata and encrypted metadata

T-1. We delete P

*/

import Joi from "joi";

import {
  callHappyPostJsonApi,
  callHappyPostJsonApiWithAuth,
  validateObject,
} from "./testlib/common-api-test-utils.js";

import {
  directorySchema,
  fileSchema,
  bucketListSchema,
  userAssertion,
  errorOfCode,
} from "./testlib/common-test-schema.js";

import { validators } from "../dist/validators.js";

const DEFAULT_USER_NAME = "admin";
const DEFAULT_PASSWORD = "PleaseChangeMe@YourEarliest2Day";

const TEST_BUCKET_NAME = "BuckX-" + Date.now();

const TEST_DIRECTORY_A_NAME = "DirA-" + Date.now();
const TEST_DIRECTORY_B_NAME = "DirB-" + Date.now();

const TEST_FILE_P_NAME = "FileP-" + Date.now();
const TEST_FILE_P_NAME_ALT = "FilePalt-" + Date.now();

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
  idOfFileP: null,
};

describe("File Suite", () => {
  test("(user/login): Preparational", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: DEFAULT_USER_NAME,
      password: DEFAULT_PASSWORD,
    });

    await validateObject(data, userAssertion);

    vars.apiKey = data.apiKey;
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

    vars.bucketId = data.bucketId;
    vars.rootDirectoryId = data.rootDirectoryId;
  });

  test("(directory/create): BuckXRoot/DirA", async () => {
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

  test("(directory/create): BuckXRoot/DirB", async () => {
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

  test("(directory/get): BuckXRoot/*", async () => {
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

  test("(file/create): BuckXRoot/DirA/FileP", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/file/create",
      {
        name: TEST_FILE_P_NAME,
        bucketId: vars.bucketId,
        parentDirectoryId: vars.idOfDirectoryA,
        encryptedMetaData: TEST_INITIAL_ENCRYPTED_METADATA,
        metaData: TEST_INITIAL_METADATA,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      fileId: validators.id,
    });

    vars.idOfFileP = data.fileId;
  });

  test("(file/get): BuckXRoot/DirA/FileP", async () => {
    const data = await callHappyPostJsonApiWithAuth(200, vars.apiKey, "/file/get", {
      bucketId: vars.bucketId,
      fileId: vars.idOfFileP,
    });

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      file: fileSchema,
    });
  });

  test("(file/rename)  BuckXRoot/DirA/FileP => BuckXRoot/DirA/FilePAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/file/rename",
      {
        name: TEST_FILE_P_NAME_ALT,
        bucketId: vars.bucketId,
        fileId: vars.idOfFileP,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(file/get): BuckXRoot/DirA/FilePAlt Ensure rename worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200, vars.apiKey, "/file/get", {
      bucketId: vars.bucketId,
      fileId: vars.idOfFileP,
    });

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      file: fileSchema,
    });

    expect(data.file.name).toBe(TEST_FILE_P_NAME_ALT);
  });

  test("(file/set-metadata) BuckXRoot/DirA/FilePAlt", async () => {
    let metaData = Object.assign(
      {},
      TEST_INITIAL_METADATA,
      TEST_NEW_METADATA_PART
    );
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/file/set-metadata",
      {
        metaData: metaData,
        bucketId: vars.bucketId,
        fileId: vars.idOfFileP,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(file/set-metadata) BuckXRoot/DirA/FilePAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/file/set-encrypted-metadata",
      {
        encryptedMetaData: TEST_NEW_ENCRYPTED_METADATA,
        bucketId: vars.bucketId,
        fileId: vars.idOfFileP,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(file/get): BuckXRoot/DirA/FilePAlt Ensure setting metaData and encryptedMetaData worked", async () => {
    const data = await callHappyPostJsonApiWithAuth(200, vars.apiKey, "/file/get", {
      bucketId: vars.bucketId,
      fileId: vars.idOfFileP,
    });

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      file: fileSchema,
    });

    let metaData = Object.assign(
      {},
      TEST_INITIAL_METADATA,
      TEST_NEW_METADATA_PART
    );
    expect(data.file.encryptedMetaData).toEqual(TEST_NEW_ENCRYPTED_METADATA);
    expect(data.file.metaData).toStrictEqual(metaData);
  });

  test("(file/move) BuckXRoot/DirA/FilePAlt => BuckXRoot/DirB/FilePAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200, vars.apiKey, "/file/move", {
      bucketId: vars.bucketId,
      fileId: vars.idOfFileP,
      newParentDirectoryId: vars.idOfDirectoryB,
      newName: TEST_FILE_P_NAME_ALT,
    });

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): BuckXRoot/DirB/* Ensure move worked", async () => {
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
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(0);
    expect(data.childFileList.length).toEqual(1);

    expect(
      data.childFileList.find((file) => file._id === vars.idOfFileP)
    ).toBeTruthy();
  });

  test("(file/delete) BuckXRoot/DirB/FilePAlt", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/file/delete",
      {
        bucketId: vars.bucketId,
        fileId: vars.idOfFileP,
      }
    );

    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
    });
  });

  test("(directory/get): BuckXRoot/DirB/* Ensure file delete worked", async () => {
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
      childDirectoryList: Joi.array().required().items(directorySchema),
      childFileList: Joi.array().required().items(fileSchema),
    });

    expect(data.childDirectoryList.length).toEqual(0);
    expect(data.childFileList.length).toEqual(0);
  });

  test("(file/get): BuckXRoot/DirB/FilePAlt Ensure file was deleted", async () => {
    const data = await callHappyPostJsonApiWithAuth(400, vars.apiKey, "/file/get", {
      bucketId: vars.bucketId,
      fileId: vars.idOfFileP,
    });

    await validateObject(data, {
      hasError: validators.hasErrorTruthy,
      error: errorOfCode("FILE_NOT_IN_BUCKET"),
    });
  });

  // eof
});
