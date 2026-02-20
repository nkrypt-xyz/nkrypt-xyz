# Admin Endpoints

This page documents all **admin** related endpoints.

## Table of Contents

- [POST /admin/iam/add-user](#post--admin-iam-add-user)
- [POST /admin/iam/set-global-permissions](#post--admin-iam-set-global-permissions)
- [POST /admin/iam/set-banning-status](#post--admin-iam-set-banning-status)
- [POST /admin/iam/overwrite-user-password](#post--admin-iam-overwrite-user-password)

---

## POST /admin/iam/add-user {#post--admin-iam-add-user}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `displayName` | string | **Yes** | Min: 4, Max: 128 |  |
| `userName` | string | **Yes** | Min: 4, Max: 32 |  |
| `password` | string | **Yes** | Min: 8, Max: 32 |  |

### Response

**Success (200):**

Response Model: [`AddUserResponse`](./models.md#adduserresponse)

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `hasError` | bool | No | - |  |
| `userId` | string | No | - |  |

**Error Responses:**

Common error codes:
- `ACCESS_DENIED` - Authentication required or insufficient permissions
- `VALIDATION_ERROR` - Request validation failed
- `NOT_FOUND` - Resource not found


---

## POST /admin/iam/set-global-permissions {#post--admin-iam-set-global-permissions}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | **Yes** | Length: 16, alphanum |  |
| `globalPermissions` | map[string]bool | **Yes** | - |  |

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

## POST /admin/iam/set-banning-status {#post--admin-iam-set-banning-status}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | **Yes** | Length: 16, alphanum |  |
| `isBanned` | bool | No | - |  |

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

## POST /admin/iam/overwrite-user-password {#post--admin-iam-overwrite-user-password}

🔒 **Authentication Required**

### Request Body

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `userId` | string | **Yes** | Length: 16, alphanum |  |
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

