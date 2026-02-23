package xyz.nkrypt.android.data.crypto

import org.junit.Assert.assertArrayEquals
import org.junit.Assert.assertEquals
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner
import org.robolectric.annotation.Config

/**
 * Unit tests for CryptoUtils.
 * Uses Robolectric for Android API support on JVM.
 */
@RunWith(RobolectricTestRunner::class)
@Config(sdk = [33])
class CryptoUtilsTest {

    @Test
    fun encryptText_decryptText_roundtrip() {
        val plaintext = "secret password"
        val password = "master123"
        val payload = CryptoUtils.encryptText(plaintext, password)
        val decrypted = CryptoUtils.decryptText(payload, password)
        assertEquals(plaintext, decrypted)
    }

    @Test
    fun encrypt_decrypt_roundtrip() {
        val plaintext = "Hello, nkrypt!".toByteArray(Charsets.UTF_8)
        val salt = CryptoUtils.generateSalt()
        val key = CryptoUtils.createEncryptionKeyFromPassword("test-password", salt)
        val iv = CryptoUtils.generateIv()
        val encrypted = CryptoUtils.encrypt(key, iv, plaintext)
        val decrypted = CryptoUtils.decrypt(key, iv, encrypted)
        assertArrayEquals(plaintext, decrypted)
    }

    @Test
    fun buildCryptoHeader_unbuildCryptoHeader_roundtrip() {
        val iv = CryptoUtils.generateIv()
        val salt = CryptoUtils.generateSalt()
        val header = CryptoUtils.buildCryptoHeader(iv, salt)
        val (parsedIv, parsedSalt) = CryptoUtils.unbuildCryptoHeader(header)
        assertArrayEquals(iv, parsedIv)
        assertArrayEquals(salt, parsedSalt)
    }

    @Test
    fun wrongPassword_decryptText_throws() {
        val payload = CryptoUtils.encryptText("secret", "correct-password")
        try {
            CryptoUtils.decryptText(payload, "wrong-password")
            org.junit.Assert.fail("Expected exception for wrong password")
        } catch (_: Exception) {
            // Expected
        }
    }
}
