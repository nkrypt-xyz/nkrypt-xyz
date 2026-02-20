package handler

import (
	"net/http"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/crypto"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type UserHandler struct {
	userSvc    *service.UserService
	sessionSvc *service.SessionService
	authSvc    *service.AuthService
	cfg        *config.Config
}

func NewUserHandler(userSvc *service.UserService, sessionSvc *service.SessionService, authSvc *service.AuthService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
		authSvc:    authSvc,
		cfg:        cfg,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	ctx := r.Context()

	user, err := h.userSvc.FindUserByUserName(ctx, req.UserName)
	if err != nil || user == nil {
		SendErrorResponse(w, apperror.NewUserError("USER_NOT_FOUND", "User not found"))
		return
	}
	if user.IsBanned {
		SendErrorResponse(w, apperror.NewUserError("USER_BANNED", "User is banned"))
		return
	}

	ok, err := crypto.VerifyPassword(
		req.Password,
		user.PasswordHash,
		user.PasswordSalt,
		h.cfg.Crypto.Argon2Memory,
		h.cfg.Crypto.Argon2Iterations,
		h.cfg.Crypto.Argon2KeyLength,
		h.cfg.Crypto.Argon2Parallelism,
	)
	if err != nil || !ok {
		SendErrorResponse(w, apperror.NewUserError("PASSWORD_INVALID", "Invalid password"))
		return
	}

	session, apiKey, err := h.sessionSvc.CreateNewUniqueSession(ctx, user)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	userResp := model.UserResponse{
		ID:                user.ID,
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		IsBanned:          user.IsBanned,
		GlobalPermissions: buildGlobalPermissions(user),
	}

	sessionResp := model.SessionResponse{ID: session.ID}

	SendSuccess(w, &model.LoginResponse{
		HasError: false,
		APIKey:   apiKey,
		User:     userResp,
		Session:  sessionResp,
	})
}

// Assert returns the current authenticated user and session info.
func (h *UserHandler) Assert(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	user := authData.User
	userResp := model.UserResponse{
		ID:                user.ID,
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		IsBanned:          user.IsBanned,
		GlobalPermissions: buildGlobalPermissions(user),
	}

	sessionResp := model.SessionResponse{ID: authData.SessionID}

	SendSuccess(w, &model.AssertResponse{
		HasError: false,
		APIKey:   authData.ApiKey,
		User:     userResp,
		Session:  sessionResp,
	})
}

// Logout expires the current session.
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	var req model.LogoutRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.sessionSvc.ExpireSessionByID(r.Context(), authData.SessionID, authData.ApiKey, authData.UserID, req.Message); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// LogoutAllSessions expires all sessions for the current user.
func (h *UserHandler) LogoutAllSessions(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	var req model.LogoutAllSessionsRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.sessionSvc.ExpireAllSessionsByUserID(r.Context(), authData.UserID, req.Message); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// ListAllSessions returns up to 20 sessions for the current user.
func (h *UserHandler) ListAllSessions(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	sessions, err := h.sessionSvc.ListSessionsByUserID(r.Context(), authData.UserID, authData.SessionID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.SessionListResponse{
		HasError:    false,
		SessionList: sessions,
	})
}

// List returns all users (basic info only).
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Auth middleware already enforced; just list.
	users, err := h.userSvc.ListAllUsers(r.Context())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Map to API response format.
	resp := make([]model.UserListItemResponse, len(users))
	for i, u := range users {
		resp[i] = model.UserListItemResponse{
			ID:          u.ID,
			UserName:    u.UserName,
			DisplayName: u.DisplayName,
			IsBanned:    u.IsBanned,
		}
	}

	SendSuccess(w, &model.UserListResponse{
		HasError: false,
		UserList: resp,
	})
}

// Find searches users by userId and/or userName filters.
func (h *UserHandler) Find(w http.ResponseWriter, r *http.Request) {
	var req model.FindUserRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	var ids []string
	var names []string
	for _, f := range req.Filters {
		switch f.By {
		case "userId":
			if f.UserID != "" {
				ids = append(ids, f.UserID)
			}
		case "userName":
			if f.UserName != "" {
				names = append(names, f.UserName)
			}
		}
	}

	users, err := h.userSvc.QueryUsers(r.Context(), ids, names)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	resp := make([]model.UserResponse, len(users))
	for i, u := range users {
		ur := model.UserResponse{
			ID:          u.ID,
			UserName:    u.UserName,
			DisplayName: u.DisplayName,
			IsBanned:    u.IsBanned,
		}
		if req.IncludeGlobalPermissions {
			ur.GlobalPermissions = buildGlobalPermissions(&u)
		}
		resp[i] = ur
	}

	SendSuccess(w, &model.FindUserResponse{
		HasError: false,
		UserList: resp,
	})
}

// UpdateProfile updates the authenticated user's display name.
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	var req model.UpdateProfileRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.userSvc.UpdateDisplayName(r.Context(), authData.UserID, req.DisplayName); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// UpdatePassword changes the authenticated user's password and expires all sessions.
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	var req model.UpdatePasswordRequest
	if err := ParseAndValidateBody(r, &req); err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Verify current password
	user, err := h.userSvc.FindUserByIDOrFail(r.Context(), authData.UserID)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	ok, err := crypto.VerifyPassword(
		req.CurrentPassword,
		user.PasswordHash,
		user.PasswordSalt,
		h.cfg.Crypto.Argon2Memory,
		h.cfg.Crypto.Argon2Iterations,
		h.cfg.Crypto.Argon2KeyLength,
		h.cfg.Crypto.Argon2Parallelism,
	)
	if err != nil || !ok {
		SendErrorResponse(w, apperror.NewUserError("PASSWORD_INVALID", "Invalid password"))
		return
	}

	// Hash new password
	newHash, newSalt, err := crypto.HashPassword(
		req.NewPassword,
		h.cfg.Crypto.Argon2Memory,
		h.cfg.Crypto.Argon2Iterations,
		h.cfg.Crypto.Argon2KeyLength,
		h.cfg.Crypto.Argon2Parallelism,
	)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	if err := h.userSvc.UpdatePassword(r.Context(), authData.UserID, newHash, newSalt); err != nil {
		SendErrorResponse(w, err)
		return
	}

	// Expire all sessions for this user.
	if err := h.sessionSvc.ExpireAllSessionsByUserID(r.Context(), authData.UserID, "Password updated by user"); err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.EmptySuccessResponse{HasError: false})
}

// helper to build globalPermissions map from user flags.
func buildGlobalPermissions(u *model.User) map[string]bool {
	return map[string]bool{
		"MANAGE_ALL_USER": u.PermManageAllUser,
		"CREATE_USER":     u.PermCreateUser,
		"CREATE_BUCKET":   u.PermCreateBucket,
	}
}

