# File Endpoints

This page documents all **file** related endpoints.

## Table of Contents

- [POST /file/create](#post--file-create)
- [POST /file/get](#post--file-get)
- [POST /file/rename](#post--file-rename)
- [POST /file/move](#post--file-move)
- [POST /file/delete](#post--file-delete)
- [POST /file/set-metadata](#post--file-set-metadata)
- [POST /file/set-encrypted-metadata](#post--file-set-encrypted-metadata)

---

## POST /file/create {#post--file-create}

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

Response Model: [`CreateFileResponse`](./models.md#createfileresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `fileId` | string | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /file/get {#post--file-get}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |

### Response

**Success (200):**

Response Model: [`GetFileResponse`](./models.md#getfileresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `file` | FileResponse | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /file/rename {#post--file-rename}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 256 |  |
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |

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

## POST /file/move {#post--file-move}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |
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

## POST /file/delete {#post--file-delete}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `fileId` | string | **Yes** | Length: 16, alphanum |  |

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

## POST /file/set-metadata {#post--file-set-metadata}

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

## POST /file/set-encrypted-metadata {#post--file-set-encrypted-metadata}

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

