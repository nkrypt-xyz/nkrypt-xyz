package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type BlobHandler struct {
	bucketSvc *service.BucketService
	fileSvc   *service.FileService
	blobSvc   *service.BlobService
}

func NewBlobHandler(bucketSvc *service.BucketService, fileSvc *service.FileService, blobSvc *service.BlobService) *BlobHandler {
	return &BlobHandler{bucketSvc: bucketSvc, fileSvc: fileSvc, blobSvc: blobSvc}
}

// Read handles POST /api/blob/read/:bucketId/:fileId
func (h *BlobHandler) Read(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	bucketID := chi.URLParam(r, "bucketId")
	fileID := chi.URLParam(r, "fileId")

	if len(bucketID) != 16 || len(fileID) != 16 {
		SendErrorResponse(w, apperror.NewUserError("INVALID_PATH_PARAMS", "Invalid bucket or file ID"))
		return
	}

	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, bucketID, "VIEW_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	file, err := h.fileSvc.FindFileByID(r.Context(), bucketID, fileID)
	if err != nil || file == nil {
		SendErrorResponse(w, apperror.NewUserError("FILE_NOT_IN_BUCKET", "The requested file could not be found in this bucket."))
		return
	}

	blob, err := h.blobSvc.FindLatestFinishedBlob(r.Context(), bucketID, fileID)
	if err != nil || blob == nil {
		SendErrorResponse(w, apperror.NewUserError("BLOB_NOT_FOUND", "No finished blob found for this file."))
		return
	}

	size, err := h.blobSvc.GetBlobSize(r.Context(), blob.ID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("nk-crypto-meta", blob.CryptoMetaHeaderContent)
	w.Header().Set("Access-Control-Expose-Headers", "nk-crypto-meta")
	w.WriteHeader(http.StatusOK)

	_, _ = h.blobSvc.StreamBlobToWriter(r.Context(), blob.ID, w)
}

// Write handles POST /api/blob/write/:bucketId/:fileId
func (h *BlobHandler) Write(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	bucketID := chi.URLParam(r, "bucketId")
	fileID := chi.URLParam(r, "fileId")

	if len(bucketID) != 16 || len(fileID) != 16 {
		SendErrorResponse(w, apperror.NewUserError("INVALID_PATH_PARAMS", "Invalid bucket or file ID"))
		return
	}

	cryptoMeta := r.Header.Get("nk-crypto-meta")
	if cryptoMeta == "" {
		SendErrorResponse(w, apperror.NewUserError("MISSING_CRYPTO_META", "Missing nk-crypto-meta header"))
		return
	}

	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, bucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	file, err := h.fileSvc.FindFileByID(r.Context(), bucketID, fileID)
	if err != nil || file == nil {
		SendErrorResponse(w, apperror.NewUserError("FILE_NOT_IN_BUCKET", "The requested file could not be found in this bucket."))
		return
	}

	blob, err := h.blobSvc.CreateInProgressBlob(r.Context(), bucketID, fileID, cryptoMeta, authData.UserID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Stream the request body to MinIO (size -1 means unknown, MinIO will handle streaming)
	if err := h.blobSvc.UploadBlobFromReader(r.Context(), blob.ID, r.Body, -1); err != nil {
		_ = h.blobSvc.MarkBlobErroneous(r.Context(), blob.ID)
		SendErrorResponse(w, err)
		return
	}

	// Get the uploaded size
	size, err := h.blobSvc.GetBlobSize(r.Context(), blob.ID)
	if err != nil {
		_ = h.blobSvc.MarkBlobErroneous(r.Context(), blob.ID)
		SendErrorResponse(w, err)
		return
	}

	// Mark blob as finished
	if err := h.blobSvc.MarkBlobFinished(r.Context(), blob.ID); err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Update file metadata
	if err := h.fileSvc.UpdateSize(r.Context(), bucketID, fileID, size); err != nil {
		SendErrorResponse(w, err)
		return
	}
	if err := h.fileSvc.SetContentUpdatedAt(r.Context(), bucketID, fileID); err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Delete old blobs
	_ = h.blobSvc.RemoveAllOtherBlobs(r.Context(), bucketID, fileID, blob.ID)

	SendSuccess(w, map[string]interface{}{
		"blobId": blob.ID,
	})
}

// WriteQuantized handles POST /api/blob/write-quantized/:bucketId/:fileId/:blobId/:offset/:shouldEnd
// For chunked uploads
func (h *BlobHandler) WriteQuantized(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	bucketID := chi.URLParam(r, "bucketId")
	fileID := chi.URLParam(r, "fileId")
	blobIDParam := chi.URLParam(r, "blobId")
	offsetParam := chi.URLParam(r, "offset")
	shouldEndParam := chi.URLParam(r, "shouldEnd")

	// Validate IDs
	if len(bucketID) != 16 || len(fileID) != 16 {
		SendErrorResponse(w, apperror.NewUserError("INVALID_PATH_PARAMS", "Invalid bucket or file ID"))
		return
	}

	// Parse offset
	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil {
		SendErrorResponse(w, apperror.NewUserError("INVALID_PATH_PARAMS", "Invalid offset parameter"))
		return
	}

	// Parse shouldEnd
	shouldEnd := shouldEndParam == "true"

	// Check bucket permissions
	if err := service.RequireBucketPermission(r.Context(), h.bucketSvc, authData.UserID, bucketID, "MANAGE_CONTENT"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Check file exists
	file, err := h.fileSvc.FindFileByID(r.Context(), bucketID, fileID)
	if err != nil || file == nil {
		SendErrorResponse(w, apperror.NewUserError("FILE_NOT_IN_BUCKET", "The requested file could not be found in this bucket."))
		return
	}

	var blobID string

	// If blobId is "null", generate new blob
	if blobIDParam == "null" || blobIDParam == "" {
		cryptoMeta := r.Header.Get("nk-crypto-meta")
		blob, err := h.blobSvc.CreateInProgressBlob(r.Context(), bucketID, fileID, cryptoMeta, authData.UserID)
		if err != nil {
			SendErrorResponse(w, err)
			return
		}
		blobID = blob.ID
	} else {
		blobID = blobIDParam
		// Verify blob exists and is in progress
		_, err := h.blobSvc.GetInProgressBlob(r.Context(), bucketID, fileID, blobID)
		if err != nil {
			SendErrorResponse(w, err)
			return
		}
	}

	// Upload chunk
	bytesWritten, err := h.blobSvc.AppendChunkToBlob(r.Context(), blobID, offset, r.Body)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	// If this is the final chunk, finalize the blob
	if shouldEnd {
		// Compose all chunks into final blob
		if err := h.blobSvc.FinalizeChunkedBlob(r.Context(), blobID); err != nil {
			SendErrorResponse(w, err)
			return
		}

		// Mark blob as finished
		if err := h.blobSvc.MarkBlobFinished(r.Context(), blobID); err != nil {
			SendErrorResponse(w, err)
			return
		}

		// Update file size
		blobSize, err := h.blobSvc.GetBlobSize(r.Context(), blobID)
		if err == nil {
			_ = h.fileSvc.UpdateSize(r.Context(), bucketID, fileID, blobSize)
		}

		// Remove old blobs
		_ = h.blobSvc.RemoveAllOtherBlobs(r.Context(), bucketID, fileID, blobID)
	}

	SendSuccess(w, map[string]interface{}{
		"blobId":          blobID,
		"bytesTransfered": bytesWritten,
	})
}
