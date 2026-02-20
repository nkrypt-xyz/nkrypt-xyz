package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type BucketHandler struct {
	bucketSvc *service.BucketService
}

func NewBucketHandler(bucketSvc *service.BucketService) *BucketHandler {
	return &BucketHandler{bucketSvc: bucketSvc}
}

// Create handles POST /api/bucket/create
func (h *BucketHandler) Create(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	if err := service.RequireGlobalPermission(authData.User, "CREATE_BUCKET"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	var req model.CreateBucketRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	bucket, rootDirID, err := h.bucketSvc.CreateBucketWithRootID(r.Context(), req.Name, req.CryptSpec, req.CryptData, req.MetaData, authData.UserID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.CreateBucketResponse{
		HasError:         false,
		BucketID:         bucket.ID,
		RootDirectoryID:  rootDirID,
	})
}

// List handles POST /api/bucket/list
func (h *BucketHandler) List(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	list, err := h.bucketSvc.ListBucketsForUser(r.Context(), authData.UserID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	bucketList := make([]model.BucketResponse, 0, len(list))
	for _, b := range list {
		var metaData interface{}
		if len(b.MetaData) > 0 {
			_ = json.Unmarshal(b.MetaData, &metaData)
		} else {
			metaData = map[string]interface{}{}
		}
		auths := make([]model.BucketAuthorizationResponse, 0, len(b.BucketAuthorizations))
		for _, a := range b.BucketAuthorizations {
			auths = append(auths, model.BucketAuthorizationResponse{
				UserID:      a.UserID,
				Notes:      a.Notes,
				Permissions: a.Permissions,
			})
		}
		bucketList = append(bucketList, model.BucketResponse{
			ID:                      b.ID,
			Name:                    b.Name,
			RootDirectoryID:         b.RootDirectoryID,
			CryptSpec:               b.CryptSpec,
			CryptData:               b.CryptData,
			MetaData:                metaData,
			BucketAuthorizations:    auths,
			CreatedByUserIdentifier: b.CreatedByUserID + "@.",
			CreatedAt:               b.CreatedAt.UnixMilli(),
			UpdatedAt:               b.UpdatedAt.UnixMilli(),
		})
	}
	SendSuccess(w, &model.BucketListResponse{
		HasError:   false,
		BucketList: bucketList,
	})
}

// Rename handles POST /api/bucket/rename
func (h *BucketHandler) Rename(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.RenameBucketRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MODIFY"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.bucketSvc.RenameBucket(r.Context(), req.BucketID, req.Name); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetMetaData handles POST /api/bucket/set-metadata
func (h *BucketHandler) SetMetaData(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.SetBucketMetaDataRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MODIFY"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.bucketSvc.SetBucketMetaData(r.Context(), req.BucketID, req.MetaData); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetAuthorization handles POST /api/bucket/set-authorization
func (h *BucketHandler) SetAuthorization(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.SetBucketAuthorizationRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "MANAGE_AUTHORIZATION"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.bucketSvc.SetBucketAuthorization(r.Context(), req.BucketID, req.TargetUserID, req.PermissionsToSet, authData.User.UserName); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// Destroy handles POST /api/bucket/destroy
func (h *BucketHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}
	var req model.DestroyBucketRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, req.BucketID, "DESTROY"); err != nil {
		SendErrorResponse(w, err)
		return
	}
	bucket, err := h.bucketSvc.FindBucketByID(r.Context(), req.BucketID)
	if err != nil || bucket == nil {
		SendErrorResponse(w, apperror.NewUserError("BUCKET_NOT_FOUND", "The requested bucket could not be found."))
		return
	}
	if bucket.Name != req.Name {
		SendErrorResponse(w, apperror.NewUserError("BUCKET_NAME_MISMATCH", "The bucket name does not match."))
		return
	}
	if err := h.bucketSvc.DestroyBucket(r.Context(), req.BucketID); err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}