/**
 * Permission constants aligned with backend (web-server/internal/service/permissions.go,
 * web-server/internal/handler/user_handler.go buildGlobalPermissions,
 * web-server/internal/service/bucket_service.go bucketPermissionToMap).
 */

/** Global permissions (user-level, from users.perm_* columns) */
export const GLOBAL_PERMISSION = {
  MANAGE_ALL_USER: "MANAGE_ALL_USER",
  CREATE_USER: "CREATE_USER",
  CREATE_BUCKET: "CREATE_BUCKET",
} as const;

export type GlobalPermission = (typeof GLOBAL_PERMISSION)[keyof typeof GLOBAL_PERMISSION];

/** Bucket permissions (per-bucket, from bucket_user_permissions) */
export const BUCKET_PERMISSION = {
  MODIFY: "MODIFY",
  MANAGE_AUTHORIZATION: "MANAGE_AUTHORIZATION",
  DESTROY: "DESTROY",
  VIEW_CONTENT: "VIEW_CONTENT",
  MANAGE_CONTENT: "MANAGE_CONTENT",
} as const;

export type BucketPermission = (typeof BUCKET_PERMISSION)[keyof typeof BUCKET_PERMISSION];
