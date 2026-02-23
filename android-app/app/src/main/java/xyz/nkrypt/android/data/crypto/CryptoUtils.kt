package xyz.nkrypt.android.data.crypto

import android.util.Base64
import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.SecretKey
import javax.crypto.SecretKeyFactory
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.PBEKeySpec
import javax.crypto.spec.SecretKeySpec

/**
 * Crypto utilities compatible with nkrypt-xyz web client.
 * Uses AES-256-GCM, PBKDF2 with SHA-256, 100k iterations.
 */
object CryptoUtils {

    private val secureRandom = SecureRandom()

    /**
     * Create encryption key from password and salt (PBKDF2, 100k iterations).
     */
    fun createEncryptionKeyFromPassword(password: String, salt: ByteArray): SecretKey {
        val spec = PBEKeySpec(
            password.toCharArray(),
            salt,
            CryptoConstants.PBKDF2_ITERATIONS,
            CryptoConstants.KEY_LENGTH_BITS
        )
        val factory = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256")
        val keyBytes = factory.generateSecret(spec).encoded
        return SecretKeySpec(keyBytes, "AES")
    }

    /**
     * Generate random IV (12 bytes for GCM).
     */
    fun generateIv(): ByteArray {
        val iv = ByteArray(CryptoConstants.IV_LENGTH)
        secureRandom.nextBytes(iv)
        return iv
    }

    /**
     * Generate random salt (16 bytes).
     */
    fun generateSalt(): ByteArray {
        val salt = ByteArray(CryptoConstants.SALT_LENGTH)
        secureRandom.nextBytes(salt)
        return salt
    }

    /**
     * Build crypto header: NK001|{base64(iv)}|{base64(salt)}
     */
    fun buildCryptoHeader(iv: ByteArray, salt: ByteArray): String {
        val ivB64 = Base64.encodeToString(iv, Base64.NO_WRAP)
        val saltB64 = Base64.encodeToString(salt, Base64.NO_WRAP)
        return "${CryptoConstants.BUCKET_CRYPTO_SPEC}|$ivB64|$saltB64"
    }

    /**
     * Parse crypto header into iv and salt.
     */
    fun unbuildCryptoHeader(header: String): Pair<ByteArray, ByteArray> {
        val parts = header.split("|")
        require(parts.size >= 3) { "Invalid crypto header format" }
        val iv = Base64.decode(parts[1], Base64.NO_WRAP)
        val salt = Base64.decode(parts[2], Base64.NO_WRAP)
        return iv to salt
    }

    /**
     * Encrypt plaintext with key and IV.
     */
    fun encrypt(key: SecretKey, iv: ByteArray, plaintext: ByteArray): ByteArray {
        val cipher = Cipher.getInstance(CryptoConstants.ENCRYPTION_ALGORITHM)
        val spec = GCMParameterSpec(CryptoConstants.GCM_TAG_LENGTH_BITS, iv)
        cipher.init(Cipher.ENCRYPT_MODE, key, spec)
        return cipher.doFinal(plaintext)
    }

    /**
     * Decrypt ciphertext with key and IV.
     */
    fun decrypt(key: SecretKey, iv: ByteArray, ciphertext: ByteArray): ByteArray {
        val cipher = Cipher.getInstance(CryptoConstants.ENCRYPTION_ALGORITHM)
        val spec = GCMParameterSpec(CryptoConstants.GCM_TAG_LENGTH_BITS, iv)
        cipher.init(Cipher.DECRYPT_MODE, key, spec)
        return cipher.doFinal(ciphertext)
    }

    /**
     * Encrypt text with password (generates IV and salt, returns JSON-encoded result).
     */
    fun encryptText(text: String, password: String): EncryptedPayload {
        val salt = generateSalt()
        val key = createEncryptionKeyFromPassword(password, salt)
        val iv = generateIv()
        val plaintext = text.toByteArray(Charsets.UTF_8)
        val ciphertext = encrypt(key, iv, plaintext)
        return EncryptedPayload(
            cipher = Base64.encodeToString(ciphertext, Base64.NO_WRAP),
            iv = Base64.encodeToString(iv, Base64.NO_WRAP),
            salt = Base64.encodeToString(salt, Base64.NO_WRAP)
        )
    }

    /**
     * Decrypt text with password.
     */
    fun decryptText(payload: EncryptedPayload, password: String): String {
        val iv = Base64.decode(payload.iv, Base64.NO_WRAP)
        val salt = Base64.decode(payload.salt, Base64.NO_WRAP)
        val ciphertext = Base64.decode(payload.cipher, Base64.NO_WRAP)
        val key = createEncryptionKeyFromPassword(password, salt)
        val plaintext = decrypt(key, iv, ciphertext)
        return String(plaintext, Charsets.UTF_8)
    }

    data class EncryptedPayload(
        val cipher: String,
        val iv: String,
        val salt: String
    )
}
