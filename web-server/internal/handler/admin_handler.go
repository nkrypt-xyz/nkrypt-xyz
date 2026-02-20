package handler

import (
	"net/http"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type AdminHandler struct {
	adminSvc *service.AdminService
	userSvc  *service.UserService
}

func NewAdminHandler(adminSvc *service.AdminService, userSvc *service.UserService) *AdminHandler {
	return &AdminHandler{
		adminSvc: adminSvc,
		userSvc:  userSvc,
	}
}

// AddUser handles POST /api/admin/iam/add-user
func (h *AdminHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	if err := service.RequireGlobalPermission(authData.User, "CREATE_USER"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	var req model.AddUserRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	userID, err := h.adminSvc.AddUser(r.Context(), req.DisplayName, req.UserName, req.Password)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.AddUserResponse{
		HasError: false,
		UserID:   userID,
	})
}

// SetGlobalPermissions handles POST /api/admin/iam/set-global-permissions
func (h *AdminHandler) SetGlobalPermissions(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	if err := service.RequireGlobalPermission(authData.User, "MANAGE_ALL_USER"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	var req model.SetGlobalPermissionsRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.adminSvc.SetGlobalPermissions(r.Context(), req.UserID, req.GlobalPermissions); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// SetBanningStatus handles POST /api/admin/iam/set-banning-status
func (h *AdminHandler) SetBanningStatus(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	if err := service.RequireGlobalPermission(authData.User, "MANAGE_ALL_USER"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	var req model.SetBanningStatusRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.adminSvc.SetBanningStatus(r.Context(), req.UserID, req.IsBanned); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// OverwriteUserPassword handles POST /api/admin/iam/overwrite-user-password
func (h *AdminHandler) OverwriteUserPassword(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	if err := service.RequireGlobalPermission(authData.User, "MANAGE_ALL_USER"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	var req model.OverwriteUserPasswordRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.adminSvc.OverwriteUserPassword(r.Context(), req.UserID, req.NewPassword); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

