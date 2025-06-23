package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"
	"github.com/charmbracelet/log"
)

type ServiceConfig struct {
	UserRepo Repository
	Log      *log.Logger
	Metrics  *metric.Metric
	Cache    *cache.Cache

	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	SecretKey            string
}

type service struct {
	log     *log.Logger
	repo    Repository
	metrics *metric.Metric
	cache   *cache.Cache

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	secretKey            string
}

func NewService(c ServiceConfig) *service {
	return &service{
		log:     c.Log,
		repo:    c.UserRepo,
		metrics: c.Metrics,
		cache:   c.Cache,

		accessTokenDuration:  c.AccessTokenDuration,
		refreshTokenDuration: c.RefreshTokenDuration,
		secretKey:            c.SecretKey,
	}
}

func (s service) GetSignedUser(ctx context.Context) (*dto.UserResponse, error) {
	s.log.Debug("trying to retrieve signed user")

	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error("context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}
	userID := c.UserID

	var cachedUser *dto.UserResponse
	err := s.cache.GetStruct(ctx, fmt.Sprintf("user:%s", userID), &cachedUser)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CacheMiss:
			s.log.Debug("cache miss for user", "id", userID)
			s.metrics.RecordCacheMiss("user")
		default:
			s.log.Error("failed to query user from cache", "err", err)
		}
	}

	if cachedUser != nil {
		s.log.Debug("cache hit for user", "id", userID)
		s.metrics.RecordCacheHit("product")
		return cachedUser, nil
	}

	userRecord, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to retrieve user", "err", err)
		s.metrics.RecordError("auth", "get-by-id")
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		s.log.Debug("user not found", "id", userID)
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:      userRecord.ID,
		Name:    userRecord.Name,
		Email:   userRecord.Email,
		Created: userRecord.Created,
		Updated: userRecord.Updated,
	}

	cacheKey := fmt.Sprintf("user:%s", userID)
	err = s.cache.SetStruct(ctx, cacheKey, user, time.Minute*30)
	if err != nil {
		s.log.Error("failed to caching user", "err", err)
		s.metrics.RecordError("users", "cache-user")
	}
	s.log.Debug("user stored in cache", "cacheKey", cacheKey)

	return &user, nil
}

func (s service) Register(ctx context.Context, input dto.CreateUser) (*dto.AuthResponse, error) {
	s.log.Debug(
		"trying to register a new user with",
		"name", input.Name,
		"email", input.Email,
	)

	u, err := NewEntity(input.Name, input.Email)
	if err != nil {
		s.log.Error("failed to create a user", "err", err)
		return nil, err // Error is already handled by the entity
	}

	err = s.repo.Insert(ctx, u.Model())
	if err != nil {
		if err = dbutil.VerifyDuplicatedConstraintKey(err); err != nil {
			s.log.Error("duplicated user", "email", input.Email)
			s.metrics.RecordError("user", "duplicated-user")
			return nil, err // Error is already handled by the helper
		}
		s.log.Error("failed to insert user", "err", err)
		s.metrics.RecordError("user", "insert-user")
		return nil, fault.NewBadRequest("failed to insert user")
	}

	accessToken, _, err := token.Gen(s.secretKey, u.id, s.accessTokenDuration)
	if err != nil {
		s.log.Error("failed to generate access token", "err", err)
		return nil, fault.NewBadRequest("failed to generate access token")
	}

	refreshToken, _, err := token.Gen(s.secretKey, u.id, s.refreshTokenDuration)
	if err != nil {
		s.log.Error("failed to generate refresh token", "err", err)
		return nil, fault.NewBadRequest("failed to generate access token")
	}

	res := dto.AuthResponse{
		UserID:       u.id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.log.Debug(
		"user created successfully",
		"details", strings.Join(
			[]string{
				fmt.Sprintf("id: %s", u.id),
				fmt.Sprintf("name: %s", u.name),
				fmt.Sprintf("email: %s", u.email),
			},
			"\n",
		),
	)

	return &res, nil
}
