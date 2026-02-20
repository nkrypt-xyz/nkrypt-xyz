package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type DirectoryHandler struct {
	bucketSvc    *service.BucketService
	directorySvc *service.DirectoryService
}

func NewDirectoryHandler(bucketSvc *service.BucketService, directorySvc *service.DirectoryService) *DirectoryHandler {
	return &DirectoryHandler{bucketSvc: bucketSvc, directorySvc: directorySvc}
}

func directoryToResponse(d *model.Directory) model.DirectoryResponse {
	var metaData interface{}
	if len(d.MetaData) > 0 {
		_ = json.Unmarshal(d.MetaData, &metaData)
	} else {
		metaData = map[string]interface{}{}
	}
	resp := model.DirectoryResponse{
		ID:                      d.ID,
		BucketID:                d.BucketID,
		ParentDirectoryID:       d.ParentDirectoryID,
		Name:                    d.Name,
		MetaData:                metaData,
		EncryptedMetaData:       d.EncryptedMetaData,
		CreatedByUserIdentifier: d.CreatedByUserID + "@.",
		CreatedAt:               d.CreatedAt.UnixMilli(),
		UpdatedAt:               d.UpdatedAt.UnixMilli(),
	}
	return resp
}

func fileToResponse(f *model.File) model.FileResponse {
	var metaData interface{}
	if len(f.MetaData) > 0 {
		_ = json.Unmarshal(f.MetaData, &metaData)
	} else {
		metaData = map[string]interface{}{}
	}
	return model.FileResponse{
		ID:                       f.ID,
		BucketID:                 f.BucketID,
		ParentDirectoryID:        f.ParentDirectoryID,
		Name:                     f.Name,
		MetaData:                 metaData,
		EncryptedMetaData:        f.EncryptedMetaData,
		SizeAfterEncryptionBytes: f.SizeAfterEncryptionBytes,
		CreatedByUserIdentifier:  f.CreatedByUserID + "@.",
		CreatedAt:                f.CreatedAt.UnixMilli(),
		UpdatedAt:                f.UpdatedAt.UnixMilli(),
		ContentUpdatedAt:         f.ContentUpdatedAt.UnixMilli(),
	}
}

// Create handles POST /api/directory/create
func (h *DirectoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.CreateDirectoryRequest
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
	dirID, err := h.directorySvc.CreateDirectory(r.Context(), req.Name, req.BucketID, req.MetaData, req.EncryptedMetaData, authData.UserID, req.ParentDirectoryID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.CreateDirectoryResponse{
		HasError:    false,
		DirectoryID: dirID,
	})
}

// Get handles POST /api/directory/get
func (h *DirectoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.GetDirectoryRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "VIEW_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	dir, childDirs, childFiles, err := h.directorySvc.GetDirectoryContents(r.Context(), req.BucketID, req.DirectoryID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	childDirList := make([]model.DirectoryResponse, 0, len(childDirs))
	for i := range childDirs {
		childDirList = append(childDirList, directoryToResponse(&childDirs[i]))
	}
	childFileList := make([]model.FileResponse, 0, len(childFiles))
	for i := range childFiles {
		childFileList = append(childFileList, fileToResponse(&childFiles[i]))
	}
	SendSuccess(w, &model.GetDirectoryResponse{
		HasError:           false,
		Directory:          directoryToResponse(dir),
		ChildDirectoryList: childDirList,
		ChildFileList:      childFileList,
	})
}

// Rename handles POST /api/directory/rename
func (h *DirectoryHandler) Rename(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.RenameDirectoryRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.DirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.directorySvc.RenameDirectory(r.Context(), req.BucketID, req.DirectoryID, req.Name); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// Move handles POST /api/directory/move
func (h *DirectoryHandler) Move(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.MoveDirectoryRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.DirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.NewParentDirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.directorySvc.MoveDirectory(r.Context(), req.BucketID, req.DirectoryID, req.NewParentDirectoryID, req.NewName); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// Delete handles POST /api/directory/delete
func (h *DirectoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.DeleteDirectoryRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.DirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.directorySvc.DeleteDirectory(r.Context(), req.BucketID, req.DirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetMetaData handles POST /api/directory/set-metadata
func (h *DirectoryHandler) SetMetaData(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.SetDirectoryMetaDataRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.DirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.directorySvc.SetMetaData(r.Context(), req.BucketID, req.DirectoryID, req.MetaData); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetEncryptedMetaData handles POST /api/directory/set-encrypted-metadata
func (h *DirectoryHandler) SetEncryptedMetaData(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.SetDirectoryEncryptedMetaDataRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.EnsureDirectoryBelongsToBucket(r.Context(), h.directorySvc, req.BucketID, req.DirectoryID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.directorySvc.SetEncryptedMetaData(r.Context(), req.BucketID, req.DirectoryID, req.EncryptedMetaData); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}
