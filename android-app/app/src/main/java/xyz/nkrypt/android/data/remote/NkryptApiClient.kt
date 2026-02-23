package xyz.nkrypt.android.data.remote

import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import xyz.nkrypt.android.data.remote.api.NkryptApi
import java.util.concurrent.TimeUnit

object NkryptApiClient {

    /**
     * Normalizes server URL. Supports both HTTP and HTTPS.
     * - Explicit http:// or https:// preserved
     * - localhost / 127.0.0.1 default to http (for local dev)
     * - Other hosts default to https
     */
    private fun normalizeBaseUrl(url: String): String {
        var base = url.trim().removeSuffix("/")
        if (!base.startsWith("http://") && !base.startsWith("https://")) {
            val hostPart = base.substringBefore("/")
            val isLocal = hostPart.startsWith("localhost") ||
                hostPart.startsWith("127.0.0.1") ||
                hostPart.startsWith("10.") ||
                hostPart.startsWith("192.168.")
            base = if (isLocal) "http://$base" else "https://$base"
        }
        return "$base/"
    }

    fun create(baseUrl: String, apiKey: String? = null): NkryptApi {
        val clientBuilder = OkHttpClient.Builder()
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(60, TimeUnit.SECONDS)
            .writeTimeout(60, TimeUnit.SECONDS)
            .addInterceptor(HttpLoggingInterceptor().apply {
                level = HttpLoggingInterceptor.Level.NONE
            })
        if (apiKey != null) {
            clientBuilder.addInterceptor { chain ->
                val request = chain.request().newBuilder()
                    .header("Authorization", "Bearer $apiKey")
                    .build()
                chain.proceed(request)
            }
        }
        val client = clientBuilder.build()

        return Retrofit.Builder()
            .baseUrl(normalizeBaseUrl(baseUrl))
            .client(client)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(NkryptApi::class.java)
    }
}
