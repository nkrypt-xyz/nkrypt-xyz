package xyz.nkrypt.android.util

import java.security.SecureRandom

private val secureRandom = SecureRandom()
private val CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789".toCharArray()

fun generateId16(): String {
    val sb = StringBuilder(16)
    repeat(16) {
        sb.append(CHARS[secureRandom.nextInt(CHARS.size)])
    }
    return sb.toString()
}
