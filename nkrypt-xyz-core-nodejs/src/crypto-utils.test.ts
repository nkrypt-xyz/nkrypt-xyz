import { describe, it, expect, beforeEach, afterEach } from "vitest";
import {
  encryptText,
  decryptText,
  encryptObject,
  decryptToObject,
  encryptCryptoData,
  decryptCryptoData,
  createEncryptionKeyFromPassword,
  encryptBuffer,
  decryptBuffer,
  makeRandomIv,
  makeRandomSalt,
  buildCryptoHeader,
  unbuildCryptoHeader,
  convertSmallUint8ArrayToString,
} from "./index.js";
import { testConstants } from "./test-constants.js";

describe("crypto-utils", () => {
  beforeEach(() => {
    testConstants.WEAKEN_CRYPTO_FOR_TESTING = true;
  });

  afterEach(() => {
    testConstants.WEAKEN_CRYPTO_FOR_TESTING = false;
  });

  describe("encryptText / decryptText", () => {
    it("round-trips plain text", async () => {
      const plain = "Hello, World!";
      const password = "test-password-123";
      const encrypted = await encryptText(plain, password);
      expect(encrypted).toHaveProperty("cipher");
      expect(encrypted).toHaveProperty("iv");
      expect(encrypted).toHaveProperty("salt");

      const decrypted = await decryptText(encrypted, password);
      expect(decrypted).toBe(plain);
    });

    it("produces different ciphertext each time (with random IV/salt)", async () => {
      testConstants.WEAKEN_CRYPTO_FOR_TESTING = false;
      const plain = "same content";
      const password = "password";
      const enc1 = await encryptText(plain, password);
      const enc2 = await encryptText(plain, password);
      expect(enc1.cipher).not.toBe(enc2.cipher);
      expect(await decryptText(enc1, password)).toBe(plain);
      expect(await decryptText(enc2, password)).toBe(plain);
      testConstants.WEAKEN_CRYPTO_FOR_TESTING = true;
    });

    it("fails with wrong password", async () => {
      const plain = "secret";
      const encrypted = await encryptText(plain, "correct-password");
      await expect(decryptText(encrypted, "wrong-password")).rejects.toThrow();
    });
  });

  describe("encryptObject / decryptToObject", () => {
    it("round-trips object", async () => {
      const obj = { foo: "bar", count: 42, nested: { a: 1 } };
      const password = "obj-password";
      const encrypted = await encryptObject(obj as Record<string, unknown>, password);
      expect(typeof encrypted).toBe("string");

      const decrypted = await decryptToObject(encrypted, password);
      expect(decrypted).toEqual(obj);
    });
  });

  describe("encryptCryptoData / decryptCryptoData", () => {
    it("validates correct password", async () => {
      const ivResult = await makeRandomIv();
      const saltResult = await makeRandomSalt();
      const cryptSpec = {
        iv: convertSmallUint8ArrayToString(ivResult.iv),
        salt: convertSmallUint8ArrayToString(saltResult.salt),
      };
      const password = "bucket-password";
      const cryptData = await encryptCryptoData(password, cryptSpec);
      const isValid = await decryptCryptoData(password, cryptData);
      expect(isValid).toBe(true);
    });

    it("rejects wrong password", async () => {
      const cryptSpec = {
        iv: convertSmallUint8ArrayToString(testConstants.TEST_IV),
        salt: convertSmallUint8ArrayToString(testConstants.TEST_SALT),
      };
      const cryptData = await encryptCryptoData("correct", cryptSpec);
      const isValid = await decryptCryptoData("wrong", cryptData);
      expect(isValid).toBe(false);
    });
  });

  describe("encryptBuffer / decryptBuffer", () => {
    it("round-trips ArrayBuffer", async () => {
      const plain = new TextEncoder().encode("binary data here").buffer;
      const password = "buffer-password";
      const salt = (await makeRandomSalt()).salt;
      const { key } = await createEncryptionKeyFromPassword(password, salt);
      const iv = (await makeRandomIv()).iv;

      const encrypted = await encryptBuffer({ iv, key }, plain);
      expect(encrypted.byteLength).toBeGreaterThan(plain.byteLength);

      const decrypted = await decryptBuffer({ iv, key }, encrypted);
      expect(new Uint8Array(decrypted)).toEqual(new Uint8Array(plain));
    });
  });

  describe("buildCryptoHeader / unbuildCryptoHeader integration", () => {
    it("header format matches crypto-utils output", async () => {
      const ivResult = await makeRandomIv();
      const saltResult = await makeRandomSalt();
      const ivStr = convertSmallUint8ArrayToString(ivResult.iv);
      const saltStr = convertSmallUint8ArrayToString(saltResult.salt);

      const header = buildCryptoHeader(ivStr, saltStr);
      const [decodedIv, decodedSalt] = unbuildCryptoHeader(header);
      expect(new Uint8Array(decodedIv)).toEqual(ivResult.iv);
      expect(new Uint8Array(decodedSalt)).toEqual(saltResult.salt);
    });
  });
});
