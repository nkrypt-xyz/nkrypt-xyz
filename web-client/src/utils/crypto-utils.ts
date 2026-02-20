import {
  ENCRYPTION_ALGORITHM,
  ENCRYPTION_TAGLENGTH_IN_BITS,
  IV_LENGTH,
  PASSPHRASE_DERIVEKEY_ALGORITHM,
  PASSPHRASE_DERIVEKEY_GENERATED_KEYLENGTH,
  PASSPHRASE_DERIVEKEY_HASH_ALGORITHM,
  PASSPHRASE_DERIVEKEY_ITERATIONS,
  PASSPHRASE_IMPORTKEY_ALGORITHM,
  SALT_LENGTH,
  BUCKET_CRYPTO_SPEC,
} from "constants/crypto-specs";
import { testConstants } from "constants/test-constants";
import { convertSmallBufferToString, convertSmallStringToBuffer, convertSmallUint8ArrayToString } from "./buffer-utils";

export const makeRandomIv = async () => {
  if (testConstants.WEAKEN_CRYPTO_FOR_TESTING) {
    console.warn("WARNING! using predefined IV for testing ONLY. This significantly reduces the strength of the cryptography and must NEVER be used in production.");
    return { iv: testConstants.TEST_IV };
  }
  const iv = window.crypto.getRandomValues(new Uint8Array(IV_LENGTH));
  return { iv };
};

export const generateIv = () => {
  return window.crypto.getRandomValues(new Uint8Array(IV_LENGTH));
};

export const makeRandomSalt = async () => {
  if (testConstants.WEAKEN_CRYPTO_FOR_TESTING) {
    console.warn("WARNING! using predefined SALT for testing ONLY. This significantly reduces the strength of the cryptography and must NEVER be used in production.");
    return { salt: testConstants.TEST_SALT };
  }
  const salt = window.crypto.getRandomValues(new Uint8Array(SALT_LENGTH));
  return { salt };
};

export const createEncryptionKeyFromPassword = async (encryptionPassword: string, salt: Uint8Array) => {
  if (!encryptionPassword) {
    throw new Error("encryptionPassword is required");
  }

  if (!salt) {
    throw new Error("salt is required");
  }

  const encodedPassphrase = new TextEncoder().encode(encryptionPassword);

  const keyMaterial = await window.crypto.subtle.importKey("raw", encodedPassphrase, PASSPHRASE_IMPORTKEY_ALGORITHM, false, ["deriveBits", "deriveKey"]);

  const key = await window.crypto.subtle.deriveKey(
    {
      name: PASSPHRASE_DERIVEKEY_ALGORITHM,
      salt: salt as BufferSource,
      iterations: PASSPHRASE_DERIVEKEY_ITERATIONS,
      hash: PASSPHRASE_DERIVEKEY_HASH_ALGORITHM,
    },
    keyMaterial,
    {
      name: ENCRYPTION_ALGORITHM,
      length: PASSPHRASE_DERIVEKEY_GENERATED_KEYLENGTH,
    },
    true,
    ["encrypt", "decrypt"]
  );

  return { key, salt };
};

export const createEncryptionKeyForBucket = async (bucketPassword: string) => {
  const encodedPassphrase = new TextEncoder().encode(bucketPassword);
  const keyMaterial = await window.crypto.subtle.importKey("raw", encodedPassphrase, "PBKDF2", false, ["deriveBits", "deriveKey"]);

  const salt = new Uint8Array(16);

  const key = await window.crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      salt,
      iterations: PASSPHRASE_DERIVEKEY_ITERATIONS,
      hash: "SHA-256",
    },
    keyMaterial,
    {
      name: "AES-GCM",
      length: 256,
    },
    true,
    ["encrypt", "decrypt"]
  );

  return key;
};

export const encryptText = async (text: string, encryptionPassword: string) => {
  const encoder = new TextEncoder();
  const encodedData = encoder.encode(text);

  const { salt } = await makeRandomSalt();
  const { key } = await createEncryptionKeyFromPassword(encryptionPassword, salt);

  const { iv } = await makeRandomIv();

  const cipher = await window.crypto.subtle.encrypt(
    {
      name: ENCRYPTION_ALGORITHM,
      iv: iv as BufferSource,
      tagLength: ENCRYPTION_TAGLENGTH_IN_BITS,
    },
    key,
    encodedData
  );

  return {
    cipher: convertSmallBufferToString(cipher),
    iv: convertSmallUint8ArrayToString(iv),
    salt: convertSmallUint8ArrayToString(salt),
  };
};

export const decryptText = async ({ cipher, iv, salt }: { cipher: string; iv: string; salt: string }, encryptionPassword: string): Promise<string> => {
  const cipherBuffer = convertSmallStringToBuffer(cipher);
  const ivBuffer = convertSmallStringToBuffer(iv);
  const saltBuffer = convertSmallStringToBuffer(salt);

  const { key } = await createEncryptionKeyFromPassword(encryptionPassword, new Uint8Array(saltBuffer));

  const encodedData = await window.crypto.subtle.decrypt(
    {
      name: ENCRYPTION_ALGORITHM,
      iv: new Uint8Array(ivBuffer),
      tagLength: ENCRYPTION_TAGLENGTH_IN_BITS,
    },
    key,
    cipherBuffer
  );

  const data = new TextDecoder().decode(encodedData);
  return data;
};

export const encryptObject = async (object: Record<string, any>, encryptionPassword: string): Promise<string> => {
  const text = JSON.stringify(object);
  const encrypted = await encryptText(text, encryptionPassword);
  return JSON.stringify(encrypted);
};

export const decryptToObject = async (encryptedText: string, encryptionPassword: string): Promise<any> => {
  const encrypted = JSON.parse(encryptedText);
  const decrypted = await decryptText(encrypted, encryptionPassword);
  return JSON.parse(decrypted);
};

export const encryptBuffer = async ({ iv, key }: { iv: Uint8Array; key: CryptoKey }, buffer: ArrayBuffer): Promise<ArrayBuffer> => {
  const encryptedBuffer = await window.crypto.subtle.encrypt(
    {
      name: ENCRYPTION_ALGORITHM,
      iv: iv as BufferSource,
      tagLength: ENCRYPTION_TAGLENGTH_IN_BITS,
    },
    key,
    buffer
  );
  return encryptedBuffer;
};

export const decryptBuffer = async ({ iv, key }: { iv: Uint8Array; key: CryptoKey }, buffer: ArrayBuffer): Promise<ArrayBuffer> => {
  try {
    const decryptedBuffer = await window.crypto.subtle.decrypt(
      {
        name: ENCRYPTION_ALGORITHM,
        iv: iv as BufferSource,
        tagLength: ENCRYPTION_TAGLENGTH_IN_BITS,
      },
      key,
      buffer
    );
    return decryptedBuffer;
  } catch (ex) {
    console.error("Decryption error:", ex);
    throw ex;
  }
};

export const generateCryptoSpec = async () => {
  const { iv } = await makeRandomIv();
  const { salt } = await makeRandomSalt();

  return {
    spec: BUCKET_CRYPTO_SPEC,
    iv: convertSmallUint8ArrayToString(iv),
    salt: convertSmallUint8ArrayToString(salt),
  };
};

export const encryptCryptoData = async (password: string, cryptSpec: any): Promise<string> => {
  const { iv, salt } = cryptSpec;
  const ivBuffer = convertSmallStringToBuffer(iv);
  const saltBuffer = convertSmallStringToBuffer(salt);

  const { key } = await createEncryptionKeyFromPassword(password, new Uint8Array(saltBuffer));

  const testData = JSON.stringify({ test: "valid" });
  const encrypted = await encryptText(testData, password);

  return JSON.stringify(encrypted);
};

export const decryptCryptoData = async (password: string, cryptData: string): Promise<boolean> => {
  try {
    const encrypted = JSON.parse(cryptData);
    await decryptText(encrypted, password);
    return true;
  } catch {
    return false;
  }
};
