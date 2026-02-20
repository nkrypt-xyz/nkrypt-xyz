# Blob Endpoints

This page documents all **blob** related endpoints.

## Table of Contents

- [POST /blob/read/{bucketId}/{fileId}](#post--blob-read-bucketid-fileid)
- [POST /blob/write/{bucketId}/{fileId}](#post--blob-write-bucketid-fileid)
- [POST /blob/write-quantized/{bucketId}/{fileId}/{blobId}/{offset}/{shouldEnd}](#post--blob-write-quantized-bucketid-fileid-blobid-offset-shouldend)

---

## POST /blob/read/{bucketId}/{fileId} {#post--blob-read-bucketid-fileid}

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

## POST /blob/write/{bucketId}/{fileId} {#post--blob-write-bucketid-fileid}

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

## POST /blob/write-quantized/{bucketId}/{fileId}/{blobId}/{offset}/{shouldEnd} {#post--blob-write-quantized-bucketid-fileid-blobid-offset-shouldend}

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

