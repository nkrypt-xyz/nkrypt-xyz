package service

import (
	"context"
	"strings"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
)

type AuthService struct {
	sessionSvc *SessionService
	userSvc    *UserService
	cfg        *config.Config
}

func NewAuthService(sessionSvc *SessionService, userSvc *UserService, cfg *config.Config) *AuthService {
	return &AuthService{
		sessionSvc: sessionSvc,
		userSvc:    userSvc,
		cfg:        cfg,
	}
}

// Authenticate validates the Authorization header and returns auth data for the authenticated user.
func (a *AuthService) Authenticate(ctx context.Context, authorizationHeader string) (*model.AuthData, error) {
	if strings.TrimSpace(authorizationHeader) == "" {
		return nil, apperror.NewUserError("AUTHORIZATION_HEADER_MISSING", "Authorization header is missing")
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return nil, apperror.NewUserError("AUTHORIZATION_HEADER_MALFORMATTED", "Authorization header is malformatted")
	}
	apiKey := strings.TrimSpace(parts[1])
	if len(apiKey) != a.cfg.IAM.APIKeyLength {
		return nil, apperror.NewUserError("AUTHORIZATION_HEADER_MALFORMATTED", "Authorization header is malformatted")
	}

	redisSession, err := a.sessionSvc.GetSessionByApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	if redisSession == nil {
		return nil, apperror.NewUserError("API_KEY_EXPIRED", "Your session has expired. Login again.")
	}

	// Defense-in-depth expiration check using createdAt and configured validity duration.
	createdAt := time.UnixMilli(redisSession.CreatedAt)
	if time.Since(createdAt) > a.cfg.IAM.SessionValidityDuration {
		return nil, apperror.NewUserError("API_KEY_EXPIRED", "Your session has expired. Login again.")
	}

	user, err := a.userSvc.FindUserByIDOrFail(ctx, redisSession.UserID)
	if err != nil {
		return nil, err
	}

	return &model.AuthData{
		ApiKey:    apiKey,
		UserID:    redisSession.UserID,
		SessionID: redisSession.SessionID,
		User:      user,
	}, nil
}

