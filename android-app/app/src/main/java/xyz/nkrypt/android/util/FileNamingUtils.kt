package xyz.nkrypt.android.util

import androidx.documentfile.provider.DocumentFile

/**
 * Returns a unique filename for the given parent directory.
 * On naming clash, inserts " (1)", " (2)", etc. before the extension.
 * E.g. "document.pdf" -> "document (1).pdf", not "document.pdf (1)".
 */
fun uniqueFileName(parent: DocumentFile, desiredName: String): String {
    val existing = parent.findFile(desiredName)
    if (existing == null) return desiredName

    val lastDot = desiredName.lastIndexOf('.')
    val (baseName, ext) = if (lastDot > 0) {
        desiredName.substring(0, lastDot) to desiredName.substring(lastDot)
    } else {
        desiredName to ""
    }

    var n = 1
    while (true) {
        val candidate = "$baseName ($n)$ext"
        if (parent.findFile(candidate) == null) return candidate
        n++
    }
}
