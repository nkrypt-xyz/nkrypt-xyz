# System Endpoints

This page documents all **system** related endpoints.

## Table of Contents

- [GET /healthz](#get--healthz)
- [GET /readyz](#get--readyz)
- [GET /metrics](#get--metrics)
- [POST /metrics/get-summary](#post--metrics-get-summary)

---

## GET /healthz {#get--healthz}

🌐 **Public Endpoint** (No authentication required)

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

## GET /readyz {#get--readyz}

🌐 **Public Endpoint** (No authentication required)

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

## GET /metrics {#get--metrics}

🌐 **Public Endpoint** (No authentication required)

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

## POST /metrics/get-summary {#post--metrics-get-summary}

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

