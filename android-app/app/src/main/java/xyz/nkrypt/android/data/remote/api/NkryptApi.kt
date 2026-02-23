package xyz.nkrypt.android.data.remote.api

import com.google.gson.annotations.SerializedName
import okhttp3.RequestBody
import okhttp3.ResponseBody
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Path
import retrofit2.http.Streaming

interface NkryptApi {

    @POST("api/user/login")
    suspend fun login(@Body body: LoginRequest): Response<LoginResponse>

    @POST("api/bucket/list")
    suspend fun listBuckets(@Body body: ApiKeyRequest): Response<BucketListResponse>

    @POST("api/directory/get")
    suspend fun getDirectory(@Body body: DirectoryGetRequest): Response<DirectoryGetResponse>

    @POST("api/directory/create")
    suspend fun createDirectory(@Body body: DirectoryCreateRequest): Response<DirectoryCreateResponse>

    @POST("api/directory/rename")
    suspend fun renameDirectory(@Body body: DirectoryRenameRequest): Response<ApiResponse>

    @POST("api/directory/move")
    suspend fun moveDirectory(@Body body: DirectoryMoveRequest): Response<ApiResponse>

    @POST("api/directory/delete")
    suspend fun deleteDirectory(@Body body: DirectoryDeleteRequest): Response<ApiResponse>

    @POST("api/file/create")
    suspend fun createFile(@Body body: FileCreateRequest): Response<FileCreateResponse>

    @POST("api/file/get")
    suspend fun getFile(@Body body: FileGetRequest): Response<FileGetResponse>

    @POST("api/file/rename")
    suspend fun renameFile(@Body body: FileRenameRequest): Response<ApiResponse>

    @POST("api/file/move")
    suspend fun moveFile(@Body body: FileMoveRequest): Response<ApiResponse>

    @POST("api/file/delete")
    suspend fun deleteFile(@Body body: FileDeleteRequest): Response<ApiResponse>

    @POST("api/blob/read/{bucketId}/{fileId}")
    @Streaming
    suspend fun readBlob(
        @Path("bucketId") bucketId: String,
        @Path("fileId") fileId: String,
        @Header("Authorization") authorization: String
    ): Response<ResponseBody>

    @POST("api/blob/write/{bucketId}/{fileId}")
    suspend fun writeBlob(
        @Path("bucketId") bucketId: String,
        @Path("fileId") fileId: String,
        @Header("Authorization") authorization: String,
        @Header("nk-crypto-meta") cryptoMeta: String,
        @Body body: RequestBody
    ): Response<BlobWriteResponse>
}

data class LoginRequest(val userName: String, val password: String)

data class LoginResponse(
    val hasError: Boolean,
    val user: UserDto?,
    val session: SessionDto?,
    val apiKey: String?,
    val error: ApiError?
)

data class UserDto(
    val _id: String,
    val userName: String,
    val displayName: String,
    val isBanned: Boolean,
    val globalPermissions: Map<String, Boolean>
)

data class SessionDto(val _id: String)

data class ApiError(val code: String, val message: String)

data class ApiKeyRequest(val apiKey: String)

data class BucketListResponse(
    val hasError: Boolean,
    val buckets: List<BucketDto>?,
    val error: ApiError?
)

data class BucketDto(
    val _id: String,
    val name: String,
    val rootDirectoryId: String?,
    val encryptionAlgo: String?,
    val metaData: Map<String, Any>?,
    val permissions: Map<String, Boolean>?,
    val createdAt: Long,
    val updatedAt: Long
)

data class ApiResponse(val hasError: Boolean, val error: ApiError?)

data class DirectoryRenameRequest(
    val apiKey: String,
    val bucketId: String,
    val directoryId: String,
    val name: String
)

data class DirectoryMoveRequest(
    val apiKey: String,
    val bucketId: String,
    val directoryId: String,
    val newParentDirectoryId: String,
    val newName: String
)

data class DirectoryDeleteRequest(
    val apiKey: String,
    val bucketId: String,
    val directoryId: String
)

data class FileRenameRequest(
    val apiKey: String,
    val bucketId: String,
    val fileId: String,
    val name: String
)

data class FileMoveRequest(
    val apiKey: String,
    val bucketId: String,
    val fileId: String,
    val newParentDirectoryId: String,
    val newName: String
)

data class FileDeleteRequest(
    val apiKey: String,
    val bucketId: String,
    val fileId: String
)

data class DirectoryGetRequest(
    val apiKey: String,
    val bucketId: String,
    val directoryId: String?
)

data class DirectoryGetResponse(
    val hasError: Boolean,
    val directory: DirectoryDto?,
    @SerializedName("childDirectoryList") val subDirectories: List<DirectoryDto>?,
    @SerializedName("childFileList") val files: List<FileDto>?,
    val error: ApiError?
)

data class DirectoryDto(
    val _id: String,
    val name: String,
    val metaData: Map<String, Any>?,
    val encryptedMetaData: String?,
    val createdAt: Long
)

data class DirectoryCreateRequest(
    val apiKey: String,
    val bucketId: String,
    val parentDirectoryId: String?,
    val name: String,
    val metaData: Map<String, Any>?,
    val encryptedMetaData: String
)

data class DirectoryCreateResponse(
    val hasError: Boolean,
    val directoryId: String?,
    val error: ApiError?
)

data class FileCreateRequest(
    val apiKey: String,
    val bucketId: String,
    val directoryId: String?,
    val name: String,
    val sizeInBytes: Long,
    val metaData: Map<String, Any>?,
    val encryptedMetaData: String
)

data class FileCreateResponse(
    val hasError: Boolean,
    val file: FileDto?,
    val error: ApiError?
)

data class FileGetRequest(
    val apiKey: String,
    val bucketId: String,
    val fileId: String
)

data class FileGetResponse(val hasError: Boolean, val file: FileDto?, val error: ApiError?)

data class FileDto(
    val _id: String,
    val name: String,
    val sizeInBytes: Long,
    val metaData: Map<String, Any>?,
    val encryptedMetaData: String?,
    val createdAt: Long
)

data class BlobWriteResponse(val hasError: Boolean, val blob: BlobDto?, val error: ApiError?)

data class BlobDto(val _id: String, val sizeInBytes: Long, val createdAt: Long)
