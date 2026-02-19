package service

import (
	"context"
	"encoding/json"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/randstr"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type BucketService struct {
	bucketRepo    *repository.BucketRepository
	directoryRepo *repository.DirectoryRepository
}

func NewBucketService(bucketRepo *repository.BucketRepository, directoryRepo *repository.DirectoryRepository) *BucketService {
	return &BucketService{bucketRepo: bucketRepo, directoryRepo: directoryRepo}
}

func (s *BucketService) FindBucketByID(ctx context.Context, id string) (*model.Bucket, error) {
	return s.bucketRepo.FindByID(ctx, id)
}

func (s *BucketService) FindBucketByName(ctx context.Context, name string) (*model.Bucket, error) {
	return s.bucketRepo.FindByName(ctx, name)
}

// CreateBucketWithRootID creates a bucket, its root directory, and creator permission; returns bucket, rootDirectoryID, error.
func (s *BucketService) CreateBucketWithRootID(ctx context.Context, name, cryptSpec, cryptData string, metaData interface{}, createdByUserID string) (*model.Bucket, string, error) {
	existing, _ := s.bucketRepo.FindByName(ctx, name)
	if existing != nil {
		return nil, "", apperror.NewUserError("DUPLICATE_BUCKET_NAME", "A bucket with this name already exists.")
	}
	metaBytes, err := json.Marshal(metaData)
	if err != nil {
		return nil, "", apperror.NewDeveloperError("INVALID_METADATA", "Failed to serialize metadata.")
	}
	bucketID, err := randstr.GenerateID(16)
	if err != nil {
		return nil, "", apperror.NewDeveloperError("ID_GENERATION_FAILED", "Failed to generate bucket ID.")
	}
	bucket := &model.Bucket{
		ID:              bucketID,
		Name:            name,
		CryptSpec:       cryptSpec,
		CryptData:       cryptData,
		MetaData:        metaBytes,
		CreatedByUserID: createdByUserID,
	}
	if err := s.bucketRepo.Create(ctx, bucket); err != nil {
		return nil, "", err
	}
	rootDirID, err := randstr.GenerateID(16)
	if err != nil {
		return nil, "", apperror.NewDeveloperError("ID_GENERATION_FAILED", "Failed to generate root directory ID.")
	}
	rootDir := &model.Directory{
		ID:                rootDirID,
		BucketID:          bucketID,
		ParentDirectoryID: nil,
		Name:              name,
		MetaData:          []byte("{}"),
		EncryptedMetaData: "",
		CreatedByUserID:   createdByUserID,
	}
	if err := s.directoryRepo.Create(ctx, rootDir); err != nil {
		return nil, "", err
	}
	perm := &model.BucketPermission{
		BucketID:               bucketID,
		UserID:                 createdByUserID,
		Notes:                  "Created this bucket",
		PermModify:             true,
		PermManageAuthorization: true,
		PermDestroy:            true,
		PermViewContent:        true,
		PermManageContent:      true,
	}
	if err := s.bucketRepo.CreatePermission(ctx, perm); err != nil {
		return nil, "", err
	}
	return bucket, rootDirID, nil
}

func (s *BucketService) ListBucketsForUser(ctx context.Context, userID string) ([]model.BucketListItem, error) {
	bucketIDs, err := s.bucketRepo.ListBucketIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(bucketIDs) == 0 {
		return nil, nil
	}
	roots, err := s.directoryRepo.ListRootsByBucketIDs(ctx, bucketIDs)
	if err != nil {
		return nil, err
	}
	rootByBucket := make(map[string]string)
	for _, d := range roots {
		rootByBucket[d.BucketID] = d.ID
	}
	var result []model.BucketListItem
	for _, bucketID := range bucketIDs {
		b, err := s.bucketRepo.FindByID(ctx, bucketID)
		if err != nil || b == nil {
			continue
		}
		rootID := rootByBucket[bucketID]
		perms, err := s.bucketRepo.ListPermissionsByBucketID(ctx, bucketID)
		if err != nil {
			return nil, err
		}
		auths := make([]model.BucketAuthorizationItem, 0, len(perms))
		for _, p := range perms {
			auths = append(auths, model.BucketAuthorizationItem{
				UserID:      p.UserID,
				Notes:       p.Notes,
				Permissions: bucketPermissionToMap(&p),
			})
		}
		result = append(result, model.BucketListItem{
			ID:                   b.ID,
			Name:                 b.Name,
			RootDirectoryID:      rootID,
			CryptSpec:            b.CryptSpec,
			CryptData:            b.CryptData,
			MetaData:             b.MetaData,
			CreatedByUserID:      b.CreatedByUserID,
			CreatedAt:            b.CreatedAt,
			UpdatedAt:            b.UpdatedAt,
			BucketAuthorizations: auths,
		})
	}
	return result, nil
}

func bucketPermissionToMap(p *model.BucketPermission) map[string]bool {
	return map[string]bool{
		"MODIFY":               p.PermModify,
		"MANAGE_AUTHORIZATION":  p.PermManageAuthorization,
		"DESTROY":               p.PermDestroy,
		"VIEW_CONTENT":         p.PermViewContent,
		"MANAGE_CONTENT":        p.PermManageContent,
	}
}

func (s *BucketService) RenameBucket(ctx context.Context, bucketID, name string) error {
	existing, _ := s.bucketRepo.FindByName(ctx, name)
	if existing != nil && existing.ID != bucketID {
		return apperror.NewUserError("DUPLICATE_BUCKET_NAME", "A bucket with this name already exists.")
	}
	return s.bucketRepo.UpdateName(ctx, bucketID, name)
}

func (s *BucketService) SetBucketMetaData(ctx context.Context, bucketID string, metaData interface{}) error {
	metaBytes, err := json.Marshal(metaData)
	if err != nil {
		return apperror.NewDeveloperError("INVALID_METADATA", "Failed to serialize metadata.")
	}
	return s.bucketRepo.UpdateMetaData(ctx, bucketID, metaBytes)
}

func (s *BucketService) DestroyBucket(ctx context.Context, bucketID string) error {
	return s.bucketRepo.Delete(ctx, bucketID)
}

func (s *BucketService) SetBucketAuthorization(ctx context.Context, bucketID, targetUserID string, permissionsToSet map[string]bool, authorizingUserName string) error {
	p, err := s.bucketRepo.FindPermission(ctx, bucketID, targetUserID)
	if err != nil || p == nil {
		p = &model.BucketPermission{
			BucketID:               bucketID,
			UserID:                 targetUserID,
			Notes:                  "Authorized by @" + authorizingUserName,
			PermModify:             false,
			PermManageAuthorization: false,
			PermDestroy:            false,
			PermViewContent:        false,
			PermManageContent:      false,
		}
		if err := s.bucketRepo.CreatePermission(ctx, p); err != nil {
			return err
		}
	}
	if v, ok := permissionsToSet["MODIFY"]; ok {
		p.PermModify = v
	}
	if v, ok := permissionsToSet["MANAGE_AUTHORIZATION"]; ok {
		p.PermManageAuthorization = v
	}
	if v, ok := permissionsToSet["DESTROY"]; ok {
		p.PermDestroy = v
	}
	if v, ok := permissionsToSet["VIEW_CONTENT"]; ok {
		p.PermViewContent = v
	}
	if v, ok := permissionsToSet["MANAGE_CONTENT"]; ok {
		p.PermManageContent = v
	}
	return s.bucketRepo.UpdatePermission(ctx, p)
}

func (s *BucketService) GetUserBucketPermissions(ctx context.Context, bucketID, userID string) (*model.BucketPermission, error) {
	return s.bucketRepo.FindPermission(ctx, bucketID, userID)
}
