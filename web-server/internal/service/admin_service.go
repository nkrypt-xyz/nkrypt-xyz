package service

import (
	"context"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/crypto"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/randstr"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type AdminService struct {
	userRepo *repository.UserRepository
	sessSvc  *SessionService
	cfg      *config.Config
}

func NewAdminService(userRepo *repository.UserRepository, sessSvc *SessionService, cfg *config.Config) *AdminService {
	return &AdminService{
		userRepo: userRepo,
		sessSvc:  sessSvc,
		cfg:      cfg,
	}
}

// CreateDefaultAdminIfNotExists seeds the default admin user if missing.
func (a *AdminService) CreateDefaultAdminIfNotExists(ctx context.Context) error {
	u, err := a.userRepo.FindUserByUserName(ctx, a.cfg.IAM.DefaultAdminUsername)
	if err == nil && u != nil {
		return nil
	}

	hash, salt, err := crypto.HashPassword(
		a.cfg.IAM.DefaultAdminPassword,
		a.cfg.Crypto.Argon2Memory,
		a.cfg.Crypto.Argon2Iterations,
		a.cfg.Crypto.Argon2KeyLength,
		a.cfg.Crypto.Argon2Parallelism,
	)
	if err != nil {
		return err
	}

	id, err := randstr.GenerateID(16)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := &model.User{
		ID:                id,
		DisplayName:       a.cfg.IAM.DefaultAdminDisplayName,
		UserName:          a.cfg.IAM.DefaultAdminUsername,
		PasswordHash:      hash,
		PasswordSalt:      salt,
		IsBanned:          false,
		PermManageAllUser: true,
		PermCreateUser:    true,
		PermCreateBucket:  true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	return a.userRepo.CreateUser(ctx, admin)
}

// AddUser creates a new user with default permissions.
func (a *AdminService) AddUser(ctx context.Context, displayName, userName, password string) (string, error) {
	// Check username uniqueness
	if existing, _ := a.userRepo.FindUserByUserName(ctx, userName); existing != nil {
		return "", apperror.NewUserError("DUPLICATE_USERNAME", "User name is already taken")
	}

	hash, salt, err := crypto.HashPassword(
		password,
		a.cfg.Crypto.Argon2Memory,
		a.cfg.Crypto.Argon2Iterations,
		a.cfg.Crypto.Argon2KeyLength,
		a.cfg.Crypto.Argon2Parallelism,
	)
	if err != nil {
		return "", err
	}

	id, err := randstr.GenerateID(16)
	if err != nil {
		return "", err
	}

	now := time.Now()
	user := &model.User{
		ID:                id,
		DisplayName:       displayName,
		UserName:          userName,
		PasswordHash:      hash,
		PasswordSalt:      salt,
		IsBanned:          false,
		PermManageAllUser: false,
		PermCreateUser:    false,
		PermCreateBucket:  true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := a.userRepo.CreateUser(ctx, user); err != nil {
		return "", err
	}
	return id, nil
}

func (a *AdminService) SetGlobalPermissions(ctx context.Context, userID string, perms map[string]bool) error {
	// Ensure user exists
	if _, err := a.userRepo.FindUserByID(ctx, userID); err != nil {
		return apperror.NewUserError("USER_NOT_FOUND", "The requested user could not be found.")
	}

	manageAll := perms["MANAGE_ALL_USER"]
	createUser := perms["CREATE_USER"]
	createBucket := perms["CREATE_BUCKET"]

	return a.userRepo.UpdateUserGlobalPermissions(ctx, userID, manageAll, createUser, createBucket)
}

func (a *AdminService) SetBanningStatus(ctx context.Context, userID string, isBanned bool) error {
	return a.userRepo.UpdateUserBanStatus(ctx, userID, isBanned)
}

func (a *AdminService) OverwriteUserPassword(ctx context.Context, userID, newPassword string) error {
	hash, salt, err := crypto.HashPassword(
		newPassword,
		a.cfg.Crypto.Argon2Memory,
		a.cfg.Crypto.Argon2Iterations,
		a.cfg.Crypto.Argon2KeyLength,
		a.cfg.Crypto.Argon2Parallelism,
	)
	if err != nil {
		return err
	}

	if err := a.userRepo.UpdateUserPassword(ctx, userID, hash, salt); err != nil {
		return err
	}

	return a.sessSvc.ExpireAllSessionsByUserID(ctx, userID, "Password overwritten by admin")
}


