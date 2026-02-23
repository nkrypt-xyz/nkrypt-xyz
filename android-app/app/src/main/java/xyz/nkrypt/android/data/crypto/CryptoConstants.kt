package xyz.nkrypt.android.data.crypto

/**
 * Crypto constants matching nkrypt-xyz-core-nodejs for full compatibility with web client.
 */
object CryptoConstants {
    const val BUCKET_CRYPTO_SPEC = "NK001"
    const val IV_LENGTH = 12
    const val SALT_LENGTH = 16
    const val PBKDF2_ITERATIONS = 100_000
    const val PBKDF2_HASH_ALGORITHM = "SHA-256"
    const val KEY_LENGTH_BITS = 256
    const val GCM_TAG_LENGTH_BITS = 128
    const val ENCRYPTION_ALGORITHM = "AES/GCM/NoPadding"

    // Content hash spec (matches nkrypt-xyz-core-nodejs for sync compatibility)
    const val CONTENT_HASH_ALGORITHM = "SHA-256"
    const val CONTENT_HASH_SALT_LENGTH = 16
    const val CONTENT_HASH_META_KEY_HASH = "content_hash"
    const val CONTENT_HASH_META_KEY_SALT = "content_hash_salt"
}
