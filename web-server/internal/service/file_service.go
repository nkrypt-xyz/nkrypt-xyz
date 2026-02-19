package service

import (
	"context"
	"encoding/json"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/randstr"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type FileService struct {
	fileRepo *repository.FileRepository
}

func NewFileService(fileRepo *repository.FileRepository) *FileService {
	return &FileService{fileRepo: fileRepo}
}

func (s *FileService) FindFileByID(ctx context.Context, bucketID, fileID string) (*model.File, error) {
	return s.fileRepo.FindByID(ctx, bucketID, fileID)
}

func (s *FileService) FindFileByNameAndParent(ctx context.Context, name, bucketID, parentDirectoryID string) (*model.File, error) {
	return s.fileRepo.FindByNameAndParent(ctx, name, bucketID, parentDirectoryID)
}

func (s *FileService) CreateFile(ctx context.Context, name, bucketID string, metaData interface{}, encryptedMetaData, createdByUserID, parentDirectoryID string) (string, error) {
	existing, _ := s.fileRepo.FindByNameAndParent(ctx, name, bucketID, parentDirectoryID)
	if existing != nil {
		return "", apperror.NewUserError("DUPLICATE_FILE_NAME", "A file with this name already exists in the directory.")
	}
	metaBytes, err := json.Marshal(metaData)
	if err != nil {
		return "", apperror.NewDeveloperError("INVALID_METADATA", "Failed to serialize metadata.")
	}
	id, err := randstr.GenerateID(16)
	if err != nil {
		return "", apperror.NewDeveloperError("ID_GENERATION_FAILED", "Failed to generate file ID.")
	}
	f := &model.File{
		ID:                       id,
		BucketID:                 bucketID,
		ParentDirectoryID:        parentDirectoryID,
		Name:                     name,
		MetaData:                 metaBytes,
		EncryptedMetaData:        encryptedMetaData,
		SizeAfterEncryptionBytes: 0,
		CreatedByUserID:          createdByUserID,
	}
	if err := s.fileRepo.Create(ctx, f); err != nil {
		return "", err
	}
	return id, nil
}

func (s *FileService) RenameFile(ctx context.Context, bucketID, fileID, name string) error {
	return s.fileRepo.UpdateName(ctx, bucketID, fileID, name)
}

func (s *FileService) MoveFile(ctx context.Context, bucketID, fileID, newParentDirectoryID, newName string) error {
	existing, _ := s.fileRepo.FindByNameAndParent(ctx, newName, bucketID, newParentDirectoryID)
	if existing != nil && existing.ID != fileID {
		return apperror.NewUserError("DUPLICATE_FILE_NAME", "A file with this name already exists in the target directory.")
	}
	return s.fileRepo.Move(ctx, bucketID, fileID, newParentDirectoryID, newName)
}

func (s *FileService) SetMetaData(ctx context.Context, bucketID, fileID string, metaData interface{}) error {
	metaBytes, err := json.Marshal(metaData)
	if err != nil {
		return apperror.NewDeveloperError("INVALID_METADATA", "Failed to serialize metadata.")
	}
	return s.fileRepo.UpdateMetaData(ctx, bucketID, fileID, metaBytes)
}

func (s *FileService) SetEncryptedMetaData(ctx context.Context, bucketID, fileID, encryptedMetaData string) error {
	return s.fileRepo.UpdateEncryptedMetaData(ctx, bucketID, fileID, encryptedMetaData)
}

func (s *FileService) SetContentUpdatedAt(ctx context.Context, bucketID, fileID string) error {
	return s.fileRepo.UpdateContentUpdatedAt(ctx, bucketID, fileID)
}

func (s *FileService) UpdateSize(ctx context.Context, bucketID, fileID string, size int64) error {
	return s.fileRepo.UpdateSize(ctx, bucketID, fileID, size)
}

func (s *FileService) DeleteFile(ctx context.Context, bucketID, fileID string) error {
	return s.fileRepo.Delete(ctx, bucketID, fileID)
}

func (s *FileService) ListFilesInDirectory(ctx context.Context, bucketID, parentDirectoryID string) ([]model.File, error) {
	return s.fileRepo.ListByDirectory(ctx, bucketID, parentDirectoryID)
}
