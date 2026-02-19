package service

import (
	"context"
	"encoding/json"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/randstr"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type DirectoryService struct {
	dirRepo  *repository.DirectoryRepository
	fileRepo *repository.FileRepository
}

func NewDirectoryService(dirRepo *repository.DirectoryRepository, fileRepo *repository.FileRepository) *DirectoryService {
	return &DirectoryService{dirRepo: dirRepo, fileRepo: fileRepo}
}

func (s *DirectoryService) FindDirectoryByID(ctx context.Context, bucketID, directoryID string) (*model.Directory, error) {
	return s.dirRepo.FindByID(ctx, bucketID, directoryID)
}

func (s *DirectoryService) FindDirectoryByNameAndParent(ctx context.Context, name, bucketID string, parentDirectoryID *string) (*model.Directory, error) {
	return s.dirRepo.FindByNameAndParent(ctx, name, bucketID, parentDirectoryID)
}

func (s *DirectoryService) CreateDirectory(ctx context.Context, name, bucketID string, metaData interface{}, encryptedMetaData string, createdByUserID string, parentDirectoryID string) (string, error) {
	parentID := &parentDirectoryID
	existing, _ := s.dirRepo.FindByNameAndParent(ctx, name, bucketID, parentID)
	if existing != nil {
		return "", apperror.NewUserError("DUPLICATE_DIRECTORY_NAME", "A directory with this name already exists in the parent.")
	}
	metaBytes, err := json.Marshal(metaData)
	if err != nil {
		return "", apperror.NewDeveloperError("INVALID_METADATA", "Failed to serialize metadata.")
	}
	id, err := randstr.GenerateID(16)
	if err != nil {
		return "", apperror.NewDeveloperError("ID_GENERATION_FAILED", "Failed to generate directory ID.")
	}
	d := &model.Directory{
		ID:                id,
		BucketID:          bucketID,
		ParentDirectoryID: parentID,
		Name:              name,
		MetaData:          metaBytes,
		EncryptedMetaData: encryptedMetaData,
		CreatedByUserID:   createdByUserID,
	}
	if err := s.dirRepo.Create(ctx, d); err != nil {
		return "", err
	}
	return id, nil
}

func (s *DirectoryService) GetDirectoryContents(ctx context.Context, bucketID, directoryID string) (*model.Directory, []model.Directory, []model.File, error) {
	dir, err := s.dirRepo.FindByID(ctx, bucketID, directoryID)
	if err != nil || dir == nil {
		return nil, nil, nil, apperror.NewUserError("DIRECTORY_NOT_IN_BUCKET", "The requested directory could not be found in this bucket.")
	}
	children, err := s.dirRepo.ListChildDirectories(ctx, bucketID, directoryID)
	if err != nil {
		return nil, nil, nil, err
	}
	files, err := s.fileRepo.ListByDirectory(ctx, bucketID, directoryID)
	if err != nil {
		return nil, nil, nil, err
	}
	return dir, children, files, nil
}

func (s *DirectoryService) RenameDirectory(ctx context.Context, bucketID, directoryID, name string) error {
	return s.dirRepo.UpdateName(ctx, bucketID, directoryID, name)
}

func (s *DirectoryService) SetMetaData(ctx context.Context, bucketID, directoryID string, metaData interface{}) error {
	metaBytes, err := json.Marshal(metaData)
	if err != nil {
		return apperror.NewDeveloperError("INVALID_METADATA", "Failed to serialize metadata.")
	}
	return s.dirRepo.UpdateMetaData(ctx, bucketID, directoryID, metaBytes)
}

func (s *DirectoryService) SetEncryptedMetaData(ctx context.Context, bucketID, directoryID, encryptedMetaData string) error {
	return s.dirRepo.UpdateEncryptedMetaData(ctx, bucketID, directoryID, encryptedMetaData)
}

// isDescendantOf walks up from dirID's ancestors; if we ever reach ancestorID, dirID is a descendant of ancestorID (cycle risk).
func (s *DirectoryService) isDescendantOf(ctx context.Context, bucketID, dirID, ancestorID string) (bool, error) {
	currentID := dirID
	for currentID != "" {
		if currentID == ancestorID {
			return true, nil
		}
		d, err := s.dirRepo.FindByID(ctx, bucketID, currentID)
		if err != nil || d == nil {
			return false, nil
		}
		if d.ParentDirectoryID == nil {
			return false, nil
		}
		currentID = *d.ParentDirectoryID
	}
	return false, nil
}

func (s *DirectoryService) MoveDirectory(ctx context.Context, bucketID, directoryID, newParentDirectoryID, newName string) error {
	descendant, err := s.isDescendantOf(ctx, bucketID, newParentDirectoryID, directoryID)
	if err != nil {
		return err
	}
	if descendant {
		return apperror.NewUserError("INVALID_MOVE", "Cannot move a directory into its own descendant.")
	}
	var newParentID *string
	if newParentDirectoryID != "" {
		newParentID = &newParentDirectoryID
	}
	return s.dirRepo.Move(ctx, bucketID, directoryID, newParentID, newName)
}

func (s *DirectoryService) DeleteDirectory(ctx context.Context, bucketID, directoryID string) error {
	return s.dirRepo.Delete(ctx, bucketID, directoryID)
}

func (s *DirectoryService) FindRootDirectoryByBucketID(ctx context.Context, bucketID string) (*model.Directory, error) {
	return s.dirRepo.FindRootByBucketID(ctx, bucketID)
}

func (s *DirectoryService) ListRootDirectoriesByBucketIDs(ctx context.Context, bucketIDs []string) ([]model.Directory, error) {
	return s.dirRepo.ListRootsByBucketIDs(ctx, bucketIDs)
}
