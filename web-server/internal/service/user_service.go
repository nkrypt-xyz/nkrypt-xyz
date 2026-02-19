package service

import (
	"context"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindUserByIDOrFail(ctx context.Context, id string) (*model.User, error) {
	u, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, apperror.NewUserError("USER_NOT_FOUND", "The requested user could not be found.")
	}
	return u, nil
}

func (s *UserService) FindUserByUserName(ctx context.Context, userName string) (*model.User, error) {
	return s.repo.FindUserByUserName(ctx, userName)
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]model.UserListItem, error) {
	return s.repo.ListAllUsers(ctx)
}

func (s *UserService) UpdateDisplayName(ctx context.Context, id, displayName string) error {
	return s.repo.UpdateUserDisplayName(ctx, id, displayName)
}

func (s *UserService) UpdatePassword(ctx context.Context, id, newHash, newSalt string) error {
	return s.repo.UpdateUserPassword(ctx, id, newHash, newSalt)
}

func (s *UserService) SetBanningStatus(ctx context.Context, id string, isBanned bool) error {
	return s.repo.UpdateUserBanStatus(ctx, id, isBanned)
}

// QueryUsers implements the FindUser behavior: filter by IDs and/or usernames.
func (s *UserService) QueryUsers(ctx context.Context, userIDs []string, userNames []string) ([]model.User, error) {
	var result []model.User

	if len(userIDs) > 0 {
		usersByID, err := s.repo.QueryUsersByIDs(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		result = append(result, usersByID...)
	}
	if len(userNames) > 0 {
		usersByName, err := s.repo.QueryUsersByUserNames(ctx, userNames)
		if err != nil {
			return nil, err
		}
		result = append(result, usersByName...)
	}

	// Deduplicate by ID if necessary.
	seen := make(map[string]bool)
	var out []model.User
	for _, u := range result {
		if seen[u.ID] {
			continue
		}
		seen[u.ID] = true
		out = append(out, u)
	}
	return out, nil
}


