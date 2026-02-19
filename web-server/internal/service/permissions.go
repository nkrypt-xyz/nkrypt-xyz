package service

import (
	"context"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
)

// RequireGlobalPermission checks that the given user has all listed global permissions.
func RequireGlobalPermission(user *model.User, permissions ...string) error {
	perms := map[string]bool{
		"MANAGE_ALL_USER": user.PermManageAllUser,
		"CREATE_USER":     user.PermCreateUser,
		"CREATE_BUCKET":   user.PermCreateBucket,
	}
	for _, p := range permissions {
		if !perms[p] {
			return apperror.NewUserError("INSUFFICIENT_GLOBAL_PERMISSION",
				"You do not have the required permissions. This action requires the \""+p+"\" permission.")
		}
	}
	return nil
}

// RequireBucketPermission ensures the bucket exists and the user has all listed bucket permissions.
func RequireBucketPermission(ctx context.Context, bucketSvc *BucketService, userID, bucketID string, permissions ...string) error {
	bucket, err := bucketSvc.FindBucketByID(ctx, bucketID)
	if err != nil || bucket == nil {
		return apperror.NewUserError("BUCKET_NOT_FOUND", "The requested bucket could not be found.")
	}
	perm, err := bucketSvc.GetUserBucketPermissions(ctx, bucketID, userID)
	if err != nil || perm == nil {
		return apperror.NewUserError("NO_AUTHORIZATION", "You do not have access to this bucket.")
	}
	perms := map[string]bool{
		"MODIFY":              perm.PermModify,
		"MANAGE_AUTHORIZATION": perm.PermManageAuthorization,
		"DESTROY":             perm.PermDestroy,
		"VIEW_CONTENT":       perm.PermViewContent,
		"MANAGE_CONTENT":     perm.PermManageContent,
	}
	for _, p := range permissions {
		if !perms[p] {
			return apperror.NewUserError("INSUFFICIENT_BUCKET_PERMISSION",
				"You do not have the required bucket permission: \""+p+"\".")
		}
	}
	return nil
}

// EnsureDirectoryBelongsToBucket returns an error if the directory is not in the bucket.
func EnsureDirectoryBelongsToBucket(ctx context.Context, dirSvc *DirectoryService, bucketID, directoryID string) error {
	dir, err := dirSvc.FindDirectoryByID(ctx, bucketID, directoryID)
	if err != nil || dir == nil {
		return apperror.NewUserError("DIRECTORY_NOT_IN_BUCKET", "The requested directory could not be found in this bucket.")
	}
	return nil
}

// EnsureFileBelongsToBucket returns an error if the file is not in the bucket.
func EnsureFileBelongsToBucket(ctx context.Context, fileSvc *FileService, bucketID, fileID string) error {
	file, err := fileSvc.FindFileByID(ctx, bucketID, fileID)
	if err != nil || file == nil {
		return apperror.NewUserError("FILE_NOT_IN_BUCKET", "The requested file could not be found in this bucket.")
	}
	return nil
}

