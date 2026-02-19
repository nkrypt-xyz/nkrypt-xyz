package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/randstr"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type SessionService struct {
	redis  *redis.Client
	repo   *repository.SessionRepository
	config *config.Config
}

func NewSessionService(redis *redis.Client, repo *repository.SessionRepository, cfg *config.Config) *SessionService {
	return &SessionService{
		redis:  redis,
		repo:   repo,
		config: cfg,
	}
}

// CreateNewUniqueSession creates a new session with a unique API key.
func (s *SessionService) CreateNewUniqueSession(ctx context.Context, user *model.User) (*model.Session, string, error) {
	var apiKey string
	for i := 0; i < 99; i++ {
		key, err := randstr.GenerateAPIKey(s.config.IAM.APIKeyLength)
		if err != nil {
			return nil, "", err
		}
		exists, err := s.redis.Exists(ctx, "nk:session:"+key).Result()
		if err != nil {
			return nil, "", err
		}
		if exists == 0 {
			apiKey = key
			break
		}
		if i == 98 {
			return nil, "", apperror.NewDeveloperError("API_KEY_CREATION_FAILED", "Timed out generating unique API key")
		}
	}

	sessionID, err := randstr.GenerateID(16)
	if err != nil {
		return nil, "", err
	}
	now := time.Now()

	sessionData := model.RedisSessionData{
		SessionID: sessionID,
		UserID:    user.ID,
		CreatedAt: now.UnixMilli(),
	}
	jsonBytes, _ := json.Marshal(sessionData)

	pipe := s.redis.Pipeline()
	pipe.Set(ctx, "nk:session:"+apiKey, string(jsonBytes), s.config.IAM.SessionValidityDuration)
	pipe.SAdd(ctx, "nk:user_sessions:"+user.ID, apiKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, "", err
	}

	apiKeyHash := sha256Hex(apiKey)
	dbSession := &model.Session{
		ID:         sessionID,
		UserID:     user.ID,
		APIKeyHash: apiKeyHash,
		HasExpired: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.repo.CreateSession(ctx, dbSession); err != nil {
		return nil, "", err
	}

	return dbSession, apiKey, nil
}

func (s *SessionService) GetSessionByApiKey(ctx context.Context, apiKey string) (*model.RedisSessionData, error) {
	val, err := s.redis.Get(ctx, "nk:session:"+apiKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sess model.RedisSessionData
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *SessionService) GetSessionByIDOrFail(ctx context.Context, id string) (*model.Session, error) {
	sess, err := s.repo.FindSessionByID(ctx, id)
	if err != nil {
		return nil, apperror.NewUserError("SESSION_NOT_FOUND", "The requested session could not be found.")
	}
	return sess, nil
}

func (s *SessionService) ListSessionsByUserID(ctx context.Context, userID, currentSessionID string) ([]model.SessionListItem, error) {
	sessions, err := s.repo.ListSessionsByUserID(ctx, userID, 20)
	if err != nil {
		return nil, err
	}

	out := make([]model.SessionListItem, len(sessions))
	for i, sess := range sessions {
		item := model.SessionListItem{
			IsCurrentSession: sess.ID == currentSessionID,
			HasExpired:       sess.HasExpired,
			CreatedAt:        sess.CreatedAt.UnixMilli(),
			ExpireReason:     sess.ExpireReason,
		}
		if sess.ExpiredAt != nil {
			expiredAtMs := sess.ExpiredAt.UnixMilli()
			item.ExpiredAt = &expiredAtMs
		}
		out[i] = item
	}
	return out, nil
}

func (s *SessionService) ExpireSessionByID(ctx context.Context, sessionID, apiKey, userID, message string) error {
	pipe := s.redis.Pipeline()
	pipe.Del(ctx, "nk:session:"+apiKey)
	pipe.SRem(ctx, "nk:user_sessions:"+userID, apiKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return s.repo.ExpireSessionByID(ctx, sessionID, "Logout: "+message)
}

func (s *SessionService) ExpireAllSessionsByUserID(ctx context.Context, userID, message string) error {
	apiKeys, err := s.redis.SMembers(ctx, "nk:user_sessions:"+userID).Result()
	if err != nil {
		return err
	}
	if len(apiKeys) > 0 {
		keys := make([]string, len(apiKeys))
		for i, ak := range apiKeys {
			keys[i] = "nk:session:" + ak
		}
		pipe := s.redis.Pipeline()
		pipe.Del(ctx, keys...)
		pipe.Del(ctx, "nk:user_sessions:"+userID)
		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}
	}
	return s.repo.ExpireAllSessionsByUserID(ctx, userID, "ForceLogout: "+message)
}

func sha256Hex(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
