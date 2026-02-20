# Directory Endpoints

This page documents all **directory** related endpoints.

## Table of Contents

- [POST /directory/create](#post--directory-create)
- [POST /directory/get](#post--directory-get)
- [POST /directory/rename](#post--directory-rename)
- [POST /directory/move](#post--directory-move)
- [POST /directory/delete](#post--directory-delete)
- [POST /directory/set-metadata](#post--directory-set-metadata)
- [POST /directory/set-encrypted-metadata](#post--directory-set-encrypted-metadata)

---

## POST /directory/create {#post--directory-create}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `parentDirectoryId` | string | **Yes** | Length: 16, alphanum |  |
| `metaData` | interface{} | **Yes** | - |  |
| `encryptedMetaData` | string | **Yes** | Min: 1, Max: 1048576 |  |

### Response

**Success (200):**

Response Model: [`CreateDirectoryResponse`](./models.md#createdirectoryresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `directoryId` | string | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /directory/get {#post--directory-get}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |

### Response

**Success (200):**

Response Model: [`GetDirectoryResponse`](./models.md#getdirectoryresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `directory` | DirectoryResponse | No | - |  |
| `childDirectoryList` | []DirectoryResponse | No | - |  |
| `childFileList` | []FileResponse | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /directory/rename {#post--directory-rename}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |

### Response

**Success (200):**

```json
{
  "hasError": false,
  ...
}
```

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /directory/move {#post--directory-move}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |
| `newParentDirectoryId` | string | **Yes** | Length: 16, alphanum |  |
| `newName` | string | **Yes** | Min: 1, Max: 256 |  |

### Response

**Success (200):**

```json
{
  "hasError": false,
  ...
}
```

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /directory/delete {#post--directory-delete}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `directoryId` | string | **Yes** | Length: 16, alphanum |  |

### Response

**Success (200):**

```json
{
  "hasError": false,
  ...
}
```

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /directory/set-metadata {#post--directory-set-metadata}

🔒 **Authentication Required**

### Request Body

No request body required.

### Response

**Success (200):**

```json
{
  "hasError": false,
  ...
}
```

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /directory/set-encrypted-metadata {#post--directory-set-encrypted-metadata}

🔒 **Authentication Required**

### Request Body

No request body required.

### Response

**Success (200):**

```json
{
  "hasError": false,
  ...
}
```

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

