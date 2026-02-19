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
import { createTestMinIOHelper } from "./testlib/minio-helper.js";

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

// Initialize MinIO helper for verification
const minioHelper = createTestMinIOHelper();

let vars = {
  apiKey: null,
  bucketId: null,
  rootDirectoryId: null,
  idOfDirectoryA: null,
  idOfDirectoryB: null,
  idOfFileP: null,
  testBlobId: null
};

describe("Blob Suite", () => {
  // Ensure MinIO bucket exists before tests
  beforeAll(async () => {
    await minioHelper.ensureBucket();
  });
  test("(user/login): Preparational", async () => {
    const data = await callHappyPostJsonApi(200, "/user/login", {
      userName: DEFAULT_USER_NAME,
      password: DEFAULT_PASSWORD,
    });

    await validateObject(data, userAssertion);

    vars.apiKey = data.apiKey;
  });

  test("(bucket/create): Preparational", async () => {
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

  test("(file/create): BuckXRoot/FileP", async () => {
    const data = await callHappyPostJsonApiWithAuth(200,
      vars.apiKey,
      "/file/create",
      {
        name: TEST_FILE_P_NAME,
        bucketId: vars.bucketId,
        parentDirectoryId: vars.rootDirectoryId,
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

  test("(blob/write) Into BuckXRoot/FileP", async () => {
    let endPoint = `/blob/write/${vars.bucketId}/${vars.idOfFileP}`;
    let data = await (
      await callRawPostApi(endPoint, vars.apiKey, TEST_STRING, TEST_FAKE_CRYPTO_META_HEADER)
    ).json();
    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      blobId: validators.id,
    });

    // Verify blob exists in MinIO
    const exists = await minioHelper.blobExists(data.blobId);
    expect(exists).toBe(true);

    vars.testBlobId = data.blobId;
  });

  test("(blob/write) Again Into BuckXRoot/FileP", async () => {
    let endPoint = `/blob/write/${vars.bucketId}/${vars.idOfFileP}`;
    let data = await (
      await callRawPostApi(endPoint, vars.apiKey, TEST_STRING_2, TEST_FAKE_CRYPTO_META_HEADER)
    ).json();
    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      blobId: validators.id,
    });

    // Verify blob exists in MinIO
    const exists = await minioHelper.blobExists(data.blobId);
    expect(exists).toBe(true);

    vars.testBlobId = data.blobId;
  });

  test("(blob/read) Read BuckXRoot/FileP", async () => {
    let endPoint = `/blob/read/${vars.bucketId}/${vars.idOfFileP}`;
    let data = await (
      await callRawPostApi(endPoint, vars.apiKey, "")
    ).text();
    expect(data).toEqual(TEST_STRING_2);
  });

  test("(blob/write) Again Stream Into BuckXRoot/FileP", async () => {
    let endPoint = `/blob/write/${vars.bucketId}/${vars.idOfFileP}`;

    const dumpData = crypto.randomBytes(TEST_BINARY_FILE_1_SIZE);
    writeFileSync(TEST_BINARY_FILE_1_LOCATION, dumpData);

    let stream = createReadStream(TEST_BINARY_FILE_1_LOCATION);

    let data = await (
      await callRawPostApi(endPoint, vars.apiKey, stream, TEST_FAKE_CRYPTO_META_HEADER)
    ).json();
    await validateObject(data, {
      hasError: validators.hasErrorFalsy,
      blobId: validators.id,
    });

    vars.testBlobId = data.blobId;

    // Verify blob exists in MinIO
    const exists = await minioHelper.blobExists(data.blobId);
    expect(exists).toBe(true);

    // Compare buffers: local file vs MinIO blob
    let f1 = readFileSync(TEST_BINARY_FILE_1_LOCATION);
    let f2 = await minioHelper.getBlob(data.blobId);
    expect(Buffer.compare(f1, f2)).toBe(0);
  });

  test("(external) Ensure stream comparison works with stream-equal vendor library", async () => {
    // Compare streams to make sure they are equal.
    let f1 = createReadStream(TEST_BINARY_FILE_1_LOCATION);
    let f2 = await minioHelper.getBlobStream(vars.testBlobId);
    let result = await compareStream(f1, f2)
    expect(result).toBe(true);
  }, TIMEOUT_FOR_LONG_RUNNING_TASKS);

  test("(blob/read) Read BuckXRoot/FileP", async () => {
    let endPoint = `/blob/read/${vars.bucketId}/${vars.idOfFileP}`;

    let data = (
      await callRawPostApi(endPoint, vars.apiKey, "")
    );

    // ensure we are not confusing the streams
    expect(isANodejsReadableStream(data.body)).toBe(true);
    expect(isABrowserReadableStream(data.body)).toBe(false);

    let fileStream = createReadStream(TEST_BINARY_FILE_1_LOCATION);

    // compare streams to make sure they are equal.
    let result = await compareStream(data.body, fileStream);
    expect(result).toBe(true);

  }, TIMEOUT_FOR_LONG_RUNNING_TASKS);

  test("(blob/write) Quantized Stream Into BuckXRoot/FileP", async () => {
    const dumpData = crypto.randomBytes(TEST_BINARY_FILE_1_SIZE);
    writeFileSync(TEST_BINARY_FILE_1_LOCATION, dumpData);

    const makeEndpoint = (blobId, offset, shouldEnd) => {
      return `/blob/write-quantized/${vars.bucketId}/${vars.idOfFileP}/${blobId}/${offset}/${shouldEnd}`;
    }

    let blobId = null;

    // 50000000  50_000_000
    const chunk1Start = 0;
    const chunk2Start = 20_000_000;
    const chunk3Start = 40_000_000;

    {
      let stream = createReadStream(TEST_BINARY_FILE_1_LOCATION, { start: chunk1Start, end: chunk2Start - 1 });
      let endPoint = makeEndpoint(blobId, chunk1Start, false);
      let data = await (
        await callRawPostApi(endPoint, vars.apiKey, stream, TEST_FAKE_CRYPTO_META_HEADER)
      ).json();
      await validateObject(data, {
        hasError: validators.hasErrorFalsy,
        blobId: validators.id,
        bytesTransfered: Joi.number().required()
      });
      expect(data.bytesTransfered).toBe(20_000_000);
      blobId = data.blobId;
      vars.testBlobId = data.blobId;
    }

    {
      let stream = createReadStream(TEST_BINARY_FILE_1_LOCATION, { start: chunk2Start, end: chunk3Start - 1 });
      let endPoint = makeEndpoint(blobId, chunk2Start, false);
      let data = await (
        await callRawPostApi(endPoint, vars.apiKey, stream, TEST_FAKE_CRYPTO_META_HEADER)
      ).json();
      await validateObject(data, {
        hasError: validators.hasErrorFalsy,
        blobId: validators.id,
        bytesTransfered: Joi.number().required()
      });
      expect(data.bytesTransfered).toBe(20_000_000);
      blobId = data.blobId;
    }

    {
      let stream = createReadStream(TEST_BINARY_FILE_1_LOCATION, { start: chunk3Start });
      let endPoint = makeEndpoint(blobId, chunk3Start, true);
      let data = await (
        await callRawPostApi(endPoint, vars.apiKey, stream, TEST_FAKE_CRYPTO_META_HEADER)
      ).json();
      await validateObject(data, {
        hasError: validators.hasErrorFalsy,
        blobId: validators.id,
        bytesTransfered: Joi.number().required()
      });
      expect(data.bytesTransfered).toBe(10_000_000);
      blobId = data.blobId;
    }

    // Compare buffers: local file vs MinIO blob
    let f1 = readFileSync(TEST_BINARY_FILE_1_LOCATION);
    let f2 = await minioHelper.getBlob(blobId);
    expect(Buffer.compare(f1, f2)).toBe(0);
  });

  // eof
});
