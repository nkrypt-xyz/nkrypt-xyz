# Models Reference

This page documents all request and response models used in the API.

## Table of Contents

- [AddUserRequest](#adduserrequest)
- [AddUserResponse](#adduserresponse)
- [AssertRequest](#assertrequest)
- [AssertResponse](#assertresponse)
- [BucketAuthorizationResponse](#bucketauthorizationresponse)
- [BucketListResponse](#bucketlistresponse)
- [BucketResponse](#bucketresponse)
- [CreateBlobResponse](#createblobresponse)
- [CreateBucketRequest](#createbucketrequest)
- [CreateBucketResponse](#createbucketresponse)
- [CreateDirectoryRequest](#createdirectoryrequest)
- [CreateDirectoryResponse](#createdirectoryresponse)
- [CreateFileRequest](#createfilerequest)
- [CreateFileResponse](#createfileresponse)
- [DeleteDirectoryRequest](#deletedirectoryrequest)
- [DeleteFileRequest](#deletefilerequest)
- [DestroyBucketRequest](#destroybucketrequest)
- [DirectoryResponse](#directoryresponse)
- [EmptySuccessResponse](#emptysuccessresponse)
- [FileResponse](#fileresponse)
- [FindUserRequest](#finduserrequest)
- [FindUserResponse](#finduserresponse)
- [GetDirectoryRequest](#getdirectoryrequest)
- [GetDirectoryResponse](#getdirectoryresponse)
- [GetFileRequest](#getfilerequest)
- [GetFileResponse](#getfileresponse)
- [LoginRequest](#loginrequest)
- [LoginResponse](#loginresponse)
- [LogoutAllSessionsRequest](#logoutallsessionsrequest)
- [LogoutRequest](#logoutrequest)
- [MetricsDiskResponse](#metricsdiskresponse)
- [MetricsGetSummaryResponse](#metricsgetsummaryresponse)
- [MoveDirectoryRequest](#movedirectoryrequest)
- [MoveFileRequest](#movefilerequest)
- [OverwriteUserPasswordRequest](#overwriteuserpasswordrequest)
- [RenameBucketRequest](#renamebucketrequest)
- [RenameDirectoryRequest](#renamedirectoryrequest)
- [RenameFileRequest](#renamefilerequest)
- [SessionListResponse](#sessionlistresponse)
- [SessionResponse](#sessionresponse)
- [SetBanningStatusRequest](#setbanningstatusrequest)
- [SetBucketAuthorizationRequest](#setbucketauthorizationrequest)
- [SetBucketMetaDataRequest](#setbucketmetadatarequest)
- [SetDirectoryEncryptedMetaDataRequest](#setdirectoryencryptedmetadatarequest)
- [SetDirectoryMetaDataRequest](#setdirectorymetadatarequest)
- [SetFileEncryptedMetaDataRequest](#setfileencryptedmetadatarequest)
- [SetFileMetaDataRequest](#setfilemetadatarequest)
- [SetGlobalPermissionsRequest](#setglobalpermissionsrequest)
- [UpdatePasswordRequest](#updatepasswordrequest)
- [UpdateProfileRequest](#updateprofilerequest)
- [UserListItemResponse](#userlistitemresponse)
- [UserListResponse](#userlistresponse)
- [UserResponse](#userresponse)
- [WriteQuantizedResponse](#writequantizedresponse)

---

## AddUserRequest

Admin requests

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `displayName` | string | **Yes** | Min: 4, Max: 128 |  |
| `userName` | string | **Yes** | Min: 4, Max: 32 |  |
| `password` | string | **Yes** | Min: 8, Max: 32 |  |


---

## AddUserResponse

AddUserResponse is the response for POST /api/admin/iam/add-user

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `userId` | string | No | - |  |


---

## AssertRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|


---

## AssertResponse

AssertResponse is the response for POST /api/user/assert

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `apiKey` | string | No | - |  |
| `user` | UserResponse | No | - |  |
| `session` | SessionResponse | No | - |  |


---

## BucketAuthorizationResponse

BucketAuthorizationResponse is one entry in bucketAuthorizations.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | No | - |  |
| `notes` | string | No | - |  |
| `permissions` | map[string]bool | No | - |  |


---

## BucketListResponse

BucketListResponse is the response for POST /api/bucket/list

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `bucketList` | []BucketResponse | No | - |  |


---

## BucketResponse

BucketResponse is the API representation of a bucket in list responses.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `_id` | string | No | - |  |
| `name` | string | No | - |  |
| `rootDirectoryId` | string | No | - |  |
| `cryptSpec` | string | No | - |  |
| `cryptData` | string | No | - |  |
| `metaData` | interface{} | No | - |  |
| `bucketAuthorizations` | []BucketAuthorizationResponse | No | - |  |
| `createdByUserIdentifier` | string | No | - |  |
| `createdAt` | int64 | No | - |  |
| `updatedAt` | int64 | No | - |  |


---

## CreateBlobResponse

CreateBlobResponse is the response for POST /api/blob/write

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `blobId` | string | No | - |  |


---

## CreateBucketRequest

Bucket requests

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 64 |  |
| `cryptSpec` | string | **Yes** | Min: 1, Max: 64 |  |
| `cryptData` | string | **Yes** | Min: 1, Max: 2048 |  |
| `metaData` | interface{} | **Yes** | - |  |


---

## CreateBucketResponse

CreateBucketResponse is the response for POST /api/bucket/create

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `bucketId` | string | No | - |  |
| `rootDirectoryId` | string | No | - |  |


---

## CreateDirectoryRequest

Directory requests

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `parentDirectoryId` | string | **Yes** | Length: 16, alphanum |  |
| `metaData` | interface{} | **Yes** | - |  |
| `encryptedMetaData` | string | **Yes** | Min: 1, Max: 1048576 |  |


---

## CreateDirectoryResponse

CreateDirectoryResponse is the response for POST /api/directory/create

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `directoryId` | string | No | - |  |


---

## CreateFileRequest

File requests

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `parentDirectoryId` | string | **Yes** | Length: 16, alphanum |  |
| `metaData` | interface{} | **Yes** | - |  |
| `encryptedMetaData` | string | **Yes** | Min: 1, Max: 1048576 |  |


---

## CreateFileResponse

CreateFileResponse is the response for POST /api/file/create

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `fileId` | string | No | - |  |


---

## DeleteDirectoryRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |


---

## DeleteFileRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |


---

## DestroyBucketRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `name` | string | **Yes** | Min: 1, Max: 64 |  |


---

## DirectoryResponse

DirectoryResponse is the API representation of a directory.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `_id` | string | No | - |  |
| `bucketId` | string | No | - |  |
| `parentDirectoryId` | string | No | - |  |
| `name` | string | No | - |  |
| `metaData` | interface{} | No | - |  |
| `encryptedMetaData` | string | No | - |  |
| `createdByUserIdentifier` | string | No | - |  |
| `createdAt` | int64 | No | - |  |
| `updatedAt` | int64 | No | - |  |


---

## EmptySuccessResponse

EmptySuccessResponse is used for endpoints that return only { hasError: false }

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |


---

## FileResponse

FileResponse is the API representation of a file.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `_id` | string | No | - |  |
| `bucketId` | string | No | - |  |
| `parentDirectoryId` | string | No | - |  |
| `name` | string | No | - |  |
| `metaData` | interface{} | No | - |  |
| `encryptedMetaData` | string | No | - |  |
| `sizeAfterEncryptionBytes` | int64 | No | - |  |
| `createdByUserIdentifier` | string | No | - |  |
| `createdAt` | int64 | No | - |  |
| `updatedAt` | int64 | No | - |  |
| `contentUpdatedAt` | int64 | No | - |  |


---

## FindUserRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `filters` | []FindUserFilter | **Yes** | dive |  |
| `includeGlobalPermissions` | bool | No | - |  |


---

## FindUserResponse

FindUserResponse is the response for POST /api/user/find

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `userList` | []UserResponse | No | - |  |


---

## GetDirectoryRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |


---

## GetDirectoryResponse

GetDirectoryResponse is the response for POST /api/directory/get

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `directory` | DirectoryResponse | No | - |  |
| `childDirectoryList` | []DirectoryResponse | No | - |  |
| `childFileList` | []FileResponse | No | - |  |


---

## GetFileRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |


---

## GetFileResponse

GetFileResponse is the response for POST /api/file/get

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `file` | FileResponse | No | - |  |


---

## LoginRequest

LoginRequest matches the API spec for /api/user/login.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userName` | string | **Yes** | Min: 4, Max: 32 |  |
| `password` | string | **Yes** | Min: 8, Max: 32 |  |


---

## LoginResponse

LoginResponse is the response for POST /api/user/login

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `apiKey` | string | No | - |  |
| `user` | UserResponse | No | - |  |
| `session` | SessionResponse | No | - |  |


---

## LogoutAllSessionsRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `message` | string | **Yes** | Min: 4, Max: 124 |  |


---

## LogoutRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `message` | string | **Yes** | Min: 4, Max: 124 |  |


---

## MetricsDiskResponse

MetricsDiskResponse is the disk section of the metrics summary.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `usedBytes` | int64 | No | - |  |
| `totalBytes` | int64 | No | - |  |


---

## MetricsGetSummaryResponse

MetricsGetSummaryResponse is the response for POST /api/metrics/get-summary

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `disk` | MetricsDiskResponse | No | - |  |


---

## MoveDirectoryRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |
| `newParentDirectoryId` | string | **Yes** | Length: 16, alphanum |  |
| `newName` | string | **Yes** | Min: 1, Max: 256 |  |


---

## MoveFileRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |
| `newParentDirectoryId` | string | **Yes** | Length: 16, alphanum |  |
| `newName` | string | **Yes** | Min: 1, Max: 256 |  |


---

## OverwriteUserPasswordRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | **Yes** | Length: 16, alphanum |  |
| `newPassword` | string | **Yes** | Min: 8, Max: 32 |  |


---

## RenameBucketRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `name` | string | **Yes** | Min: 1, Max: 64 |  |


---

## RenameDirectoryRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |


---

## RenameFileRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |


---

## SessionListResponse

SessionListResponse is the response for POST /api/user/list-all-sessions

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `sessionList` | []SessionListItem | No | - |  |


---

## SessionResponse

SessionResponse is the minimal session representation used in login/assert.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `_id` | string | No | - |  |


---

## SetBanningStatusRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | **Yes** | Length: 16, alphanum |  |
| `isBanned` | bool | No | - |  |


---

## SetBucketAuthorizationRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `targetUserId` | string | **Yes** | Length: 16, alphanum |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `permissionsToSet` | map[string]bool | **Yes** | - |  |


---

## SetBucketMetaDataRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `metaData` | interface{} | **Yes** | - |  |


---

## SetDirectoryEncryptedMetaDataRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `encryptedMetaData` | string | **Yes** | Min: 1, Max: 1048576 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |


---

## SetDirectoryMetaDataRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `metaData` | interface{} | **Yes** | - |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |


---

## SetFileEncryptedMetaDataRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `encryptedMetaData` | string | **Yes** | Min: 1, Max: 1048576 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |


---

## SetFileMetaDataRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `metaData` | interface{} | **Yes** | - |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |


---

## SetGlobalPermissionsRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | **Yes** | Length: 16, alphanum |  |
| `globalPermissions` | map[string]bool | **Yes** | - |  |


---

## UpdatePasswordRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `currentPassword` | string | **Yes** | Min: 8, Max: 32 |  |
| `newPassword` | string | **Yes** | Min: 8, Max: 32 |  |


---

## UpdateProfileRequest

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `displayName` | string | **Yes** | Min: 4, Max: 128 |  |


---

## UserListItemResponse

UserListItemResponse is a minimal user in list responses.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `_id` | string | No | - |  |
| `userName` | string | No | - |  |
| `displayName` | string | No | - |  |
| `isBanned` | bool | No | - |  |


---

## UserListResponse

UserListResponse is the response for POST /api/user/list

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `userList` | []UserListItemResponse | No | - |  |


---

## UserResponse

UserResponse is the JSON representation of a user in API responses.

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `_id` | string | No | - |  |
| `userName` | string | No | - |  |
| `displayName` | string | No | - |  |
| `isBanned` | bool | No | - |  |
| `globalPermissions` | map[string]bool | No | - |  |


---

## WriteQuantizedResponse

WriteQuantizedResponse is the response for POST /api/blob/write-quantized

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `blobId` | string | No | - |  |
| `bytesTransfered` | int64 | No | - |  |


---

