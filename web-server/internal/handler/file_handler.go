package handler

import (
	"net/http"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type FileHandler struct {
	bucketSvc    *service.BucketService
	directorySvc *service.DirectoryService
	fileSvc      *service.FileService
	blobSvc      *service.BlobService
}

func NewFileHandler(bucketSvc *service.BucketService, directorySvc *service.DirectoryService, fileSvc *service.FileService, blobSvc *service.BlobService) *FileHandler {
	return &FileHandler{bucketSvc: bucketSvc, directorySvc: directorySvc, fileSvc: fileSvc, blobSvc: blobSvc}
}

// Create handles POST /api/file/create
func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.CreateFileRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.ParentDirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	fileID, err := h.fileSvc.CreateFile(r.Context(), req.Name, req.BucketID, req.MetaData, req.EncryptedMetaData, authData.UserID, req.ParentDirectoryID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.CreateFileResponse{
		HasError: false,
		FileID:   fileID,
	})
}

// Get handles POST /api/file/get
func (h *FileHandler) Get(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.GetFileRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "VIEW_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	file, err := h.fileSvc.FindFileByID(r.Context(), req.BucketID, req.FileID)
	if err != nil || file == nil {
		SendErrorResponse(w, apperror.NewUserError("FILE_NOT_IN_BUCKET", "The requested file could not be found in this bucket."))
		return
	}
	SendSuccess(w, &model.GetFileResponse{
		HasError: false,
		File:     fileToResponse(file),
	})
}

// Rename handles POST /api/file/rename
func (h *FileHandler) Rename(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.RenameFileRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureFileBelongsToBucket(r.Context(), h.fileSvc, req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.fileSvc.RenameFile(r.Context(), req.BucketID, req.FileID, req.Name); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// Move handles POST /api/file/move
func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.MoveFileRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureFileBelongsToBucket(r.Context(), h.fileSvc, req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.NewParentDirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.fileSvc.MoveFile(r.Context(), req.BucketID, req.FileID, req.NewParentDirectoryID, req.NewName); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// Delete handles POST /api/file/delete
func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.DeleteFileRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureFileBelongsToBucket(r.Context(), h.fileSvc, req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	// Delete all blobs for this file
	if err := h.blobSvc.RemoveAllBlobsOfFile(r.Context(), req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.fileSvc.DeleteFile(r.Context(), req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetMetaData handles POST /api/file/set-metadata
func (h *FileHandler) SetMetaData(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.SetFileMetaDataRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureFileBelongsToBucket(r.Context(), h.fileSvc, req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.fileSvc.SetMetaData(r.Context(), req.BucketID, req.FileID, req.MetaData); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetEncryptedMetaData handles POST /api/file/set-encrypted-metadata
func (h *FileHandler) SetEncryptedMetaData(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.SetFileEncryptedMetaDataRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureFileBelongsToBucket(r.Context(), h.fileSvc, req.BucketID, req.FileID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.fileSvc.SetEncryptedMetaData(r.Context(), req.BucketID, req.FileID, req.EncryptedMetaData); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}
