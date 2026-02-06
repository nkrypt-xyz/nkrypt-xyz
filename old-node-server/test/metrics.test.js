import Joi from "joi";

import {
  callRawPostApi,
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

import { generateRandomBase64String } from "../dist/utility/string-utils.js";
import { Readable } from "stream";
import { writeFile } from "fs/promises";
import { createReadStream, readFileSync, writeFileSync } from "fs";

import crypto from "crypto";

import assert from "assert";
import compareStream from "stream-equal";
import pathlib from "path";
import { isABrowserReadableStream, isANodejsReadableStream } from "./testlib/common-stream-test-utils.js";
import constants from "../dist/constant/common-constants.js"

const TIMEOUT_FOR_LONG_RUNNING_TASKS = 2 * 60 * 1_000;

const DEFAULT_USER_NAME = "admin";
const DEFAULT_PASSWORD = "PleaseChangeMe@YourEarliest2Day";

const TEST_BUCKET_NAME = "BuckX-" + Date.now();
const TEST_FILE_P_NAME = "FileP-" + Date.now();

const TEST_INITIAL_METADATA = { createdFromApp: "Integration testing" };
const TEST_INITIAL_ENCRYPTED_METADATA = "PLACEHOLDER";

const TEST_STRING = generateRandomBase64String(1024 * 1024);
const TEST_STRING_2 = generateRandomBase64String(1024 * 1024);

const TEST_DATA_LOCATION = "./nkrypt-xyz-local-data/";
const TEST_BINARY_FILE_1_LOCATION = "./nkrypt-xyz-local-data/test-candidate-1.dat";
const TEST_BINARY_FILE_1_SIZE = 50_000_000;

const TEST_FAKE_CRYPTO_META_HEADER = {
  [constants.webServer.BLOB_API_CRYPTO_META_HEADER_NAME]: "None"
};

let vars = {
  apiKey: null,
  bucketId: null,
  rootDirectoryId: null,
  idOfDirectoryA: null,
  idOfDirectoryB: null,
  idOfFileP: null,
  testLocalRandomFile1Path: null
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
