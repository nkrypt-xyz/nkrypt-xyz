# Bucket Endpoints

This page documents all **bucket** related endpoints.

## Table of Contents

- [POST /bucket/create](#post--bucket-create)
- [POST /bucket/list](#post--bucket-list)
- [POST /bucket/rename](#post--bucket-rename)
- [POST /bucket/set-metadata](#post--bucket-set-metadata)
- [POST /bucket/set-authorization](#post--bucket-set-authorization)
- [POST /bucket/destroy](#post--bucket-destroy)

---

## POST /bucket/create {#post--bucket-create}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | **Yes** | Min: 1, Max: 64 |  |
| `cryptSpec` | string | **Yes** | Min: 1, Max: 64 |  |
| `cryptData` | string | **Yes** | Min: 1, Max: 2048 |  |
| `metaData` | interface{} | **Yes** | - |  |

### Response

**Success (200):**

Response Model: [`CreateBucketResponse`](./models.md#createbucketresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `bucketId` | string | No | - |  |
| `rootDirectoryId` | string | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /bucket/list {#post--bucket-list}

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

## POST /bucket/rename {#post--bucket-rename}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `name` | string | **Yes** | Min: 1, Max: 64 |  |

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

## POST /bucket/set-metadata {#post--bucket-set-metadata}

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

## POST /bucket/set-authorization {#post--bucket-set-authorization}

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

## POST /bucket/destroy {#post--bucket-destroy}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `bucketId` | string | **Yes** | Length: 16, alphanum |  |
| `name` | string | **Yes** | Min: 1, Max: 64 |  |

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

