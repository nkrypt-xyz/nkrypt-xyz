export const BLOB_API_CRYPTO_META_HEADER_NAME = "nk-crypto-meta";

export const BUCKET_CRYPTO_SPEC = "NK001";

export const IV_LENGTH = 12;
export const SALT_LENGTH = 16;

export const PASSPHRASE_ENCODING = "utf-8";

export const PASSPHRASE_IMPORTKEY_ALGORITHM = "PBKDF2";
export const PASSPHRASE_DERIVEKEY_ALGORITHM = "PBKDF2";
export const PASSPHRASE_DERIVEKEY_ITERATIONS = 100000;
export const PASSPHRASE_DERIVEKEY_HASH_ALGORITHM = "SHA-256";
export const PASSPHRASE_DERIVEKEY_GENERATED_KEYLENGTH = 256;

export const ENCRYPTION_ALGORITHM = "AES-GCM";
export const ENCRYPTION_TAGLENGTH_IN_BITS = 128;

export const BLOB_CHUNK_SIZE_BYTES = 1024 * 1024 - 128 / 8;
export const BLOB_CHUNK_SIZE_INCLUDING_TAG_BYTES = 1024 * 1024;

/**
 * Content hash spec for file integrity and efficient sync.
 * Used in blob metadata (local buckets) and remote blob version metadata (web-client).
 */
export const CONTENT_HASH_ALGORITHM = "SHA-256";
export const CONTENT_HASH_SALT_LENGTH = 16;
export const CONTENT_HASH_OUTPUT_ENCODING = "hex" as const; // 64 chars for SHA-256

/** Metadata key for content hash in blob metadata JSON */
export const CONTENT_HASH_META_KEY_HASH = "content_hash";
/** Metadata key for content hash salt in blob metadata JSON */
export const CONTENT_HASH_META_KEY_SALT = "content_hash_salt";

export const CryptoConstant = {
  ALGO_AES_GCM_256: {
    CONTENT_FILE_HEADER: new Uint8Array([0x4e, 0x4b, 0x30, 0x30, 0x31]),
    IV_LENGTH: 12,
  },
};
