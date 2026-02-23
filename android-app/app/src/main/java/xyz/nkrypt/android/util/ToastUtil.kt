package xyz.nkrypt.android.util

import android.content.Context
import android.graphics.Color
import android.os.Handler
import android.os.Looper
import android.view.Gravity
import android.widget.FrameLayout
import android.widget.TextView
import android.widget.Toast

/**
 * Shows a short-lived green success toast. Safe to call from any thread.
 */
fun showSuccessToast(context: Context, message: String = "Download completed.") {
    Handler(Looper.getMainLooper()).post {
        val toast = Toast(context)
        toast.duration = Toast.LENGTH_SHORT
        toast.setGravity(Gravity.BOTTOM, 0, 96)
        val container = FrameLayout(context).apply {
            setBackgroundColor(Color.parseColor("#4CAF50"))
            val textView = TextView(context).apply {
                text = message
                setTextColor(Color.WHITE)
                setPadding(48, 24, 48, 24)
            }
            addView(textView)
        }
        toast.view = container
        toast.show()
    }
}
