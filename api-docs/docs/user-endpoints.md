# User Endpoints

This page documents all **user** related endpoints.

## Table of Contents

- [POST /user/login](#post--user-login)
- [POST /user/assert](#post--user-assert)
- [POST /user/logout](#post--user-logout)
- [POST /user/logout-all-sessions](#post--user-logout-all-sessions)
- [POST /user/list-all-sessions](#post--user-list-all-sessions)
- [POST /user/list](#post--user-list)
- [POST /user/find](#post--user-find)
- [POST /user/update-profile](#post--user-update-profile)
- [POST /user/update-password](#post--user-update-password)

---

## POST /user/login {#post--user-login}

🌐 **Public Endpoint** (No authentication required)

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userName` | string | **Yes** | Min: 4, Max: 32 |  |
| `password` | string | **Yes** | Min: 8, Max: 32 |  |

### Response

**Success (200):**

Response Model: [`LoginResponse`](./models.md#loginresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `apiKey` | string | No | - |  |
| `user` | UserResponse | No | - |  |
| `session` | SessionResponse | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /user/assert {#post--user-assert}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|

### Response

**Success (200):**

Response Model: [`AssertResponse`](./models.md#assertresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `apiKey` | string | No | - |  |
| `user` | UserResponse | No | - |  |
| `session` | SessionResponse | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /user/logout {#post--user-logout}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `message` | string | **Yes** | Min: 4, Max: 124 |  |

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

## POST /user/logout-all-sessions {#post--user-logout-all-sessions}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `message` | string | **Yes** | Min: 4, Max: 124 |  |

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

## POST /user/list-all-sessions {#post--user-list-all-sessions}

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

## POST /user/list {#post--user-list}

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

## POST /user/find {#post--user-find}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `filters` | []FindUserFilter | **Yes** | dive |  |
| `includeGlobalPermissions` | bool | No | - |  |

### Response

**Success (200):**

Response Model: [`FindUserResponse`](./models.md#finduserresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `userList` | []UserResponse | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /user/update-profile {#post--user-update-profile}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `displayName` | string | **Yes** | Min: 4, Max: 128 |  |

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

## POST /user/update-password {#post--user-update-password}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `currentPassword` | string | **Yes** | Min: 8, Max: 32 |  |
| `newPassword` | string | **Yes** | Min: 8, Max: 32 |  |

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

