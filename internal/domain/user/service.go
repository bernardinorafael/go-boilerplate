package user

import (
	"context"
	"fmt"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"strings"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/code"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/mail"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"
	"github.com/charmbracelet/log"
)

type ServiceConfig struct {
	Log     *log.Logger
	Metrics *metric.Metric
	Cache   *cache.Cache
	Mail    *mail.Mail

	UserRepo    Repository
	CodeService code.Service

	AccessTokenDuration time.Duration
	SecretKey           string
}

type service struct {
	log     *log.Logger
	metrics *metric.Metric
	cache   *cache.Cache
	mail    *mail.Mail

	repo        Repository
	codeService code.Service

	accessTokenDuration time.Duration
	secretKey           string
}

func NewService(c ServiceConfig) *service {
	return &service{
		log:     c.Log,
		metrics: c.Metrics,
		cache:   c.Cache,
		mail:    c.Mail,

		repo:        c.UserRepo,
		codeService: c.CodeService,

		accessTokenDuration: c.AccessTokenDuration,
		secretKey:           c.SecretKey,
	}
}

func (s service) Verify(ctx context.Context, userID, code string) (*dto.AuthResponse, error) {
	s.log.Debug("trying to verify code", "code", code)

	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			s.log.Debug("user not found", "id", userID)
			return nil, fault.NewNotFound("user not found")
		default:
			s.log.Error("failed to retrieve user", "err", err)
			return nil, fault.NewBadRequest("failed to retrieve user")
		}
	}

	valid, err := s.codeService.VerifyCode(ctx, userID, code)
	if err != nil {
		return nil, err // error already handled
	}

	if !valid {
		return nil, fault.NewConflict("incorrect otp code")
	}

	accessToken, _, err := token.Gen(s.secretKey, userID, s.accessTokenDuration)
	if err != nil {
		return nil, fault.NewBadRequest("failed to generate access token")
	}

	return &dto.AuthResponse{
		UserID:      userID,
		AccessToken: accessToken,
	}, nil
}

func (s service) Login(ctx context.Context, email string) error {
	s.log.Debug("trying to login with", "email", email)

	record, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			s.log.Debug("user not found", "id", email)
			return fault.NewNotFound("user not found")
		default:
			s.log.Error("failed to retrieve user", "err", err)
			return fault.NewBadRequest("failed to retrieve user")
		}
	}

	err = s.codeService.CreateCode(ctx, record.ID)
	if err != nil {
		s.log.Error("failed to generate otp code", "err", err)
		return fault.NewBadRequest("failed to generate otp code")
	}

	s.log.Debug("successfully sent code to", "email", email)

	return nil
}

func (s service) GetSignedUser(ctx context.Context) (*dto.UserResponse, error) {
	s.log.Debug("trying to retrieve signed user")

	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error("context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}

	var cachedUser *dto.UserResponse
	err := s.cache.GetStruct(ctx, fmt.Sprintf("user:%s", c.UserID), &cachedUser)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CacheMiss:
			s.log.Debug("cache miss for user", "id", c.UserID)
			s.metrics.RecordCacheMiss("user")
		default:
			s.log.Error("failed to query user from cache", "err", err)
		}
	}

	if cachedUser != nil {
		s.log.Debug("cache hit for user", "id", c.UserID)
		s.metrics.RecordCacheHit("product")
		return cachedUser, nil
	}

	record, err := s.repo.GetByID(ctx, c.UserID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			s.log.Debug("user not found", "id", c.UserID)
			return nil, fault.NewNotFound("user not found")
		default:
			s.log.Error("failed to retrieve user", "err", err)
			s.metrics.RecordError("auth", "get-by-id")
			return nil, fault.NewBadRequest("failed to retrieve user")
		}
	}

	user := dto.UserResponse{
		ID:      record.ID,
		Name:    record.Name,
		Email:   record.Email,
		Created: record.Created,
		Updated: record.Updated,
	}

	cacheKey := fmt.Sprintf("user:%s", c.UserID)
	err = s.cache.SetStruct(ctx, cacheKey, user, time.Minute*30)
	if err != nil {
		s.log.Error("failed to caching user", "err", err)
		s.metrics.RecordError("users", "cache-user")
	}
	s.log.Debug("user stored in cache", "cacheKey", cacheKey)

	return &user, nil
}

func (s service) Register(ctx context.Context, input dto.CreateUser) error {
	s.log.Debug(
		"trying to register a new user with",
		"name", input.Name,
		"email", input.Email,
	)

	userEntity, err := NewEntity(input.Name, input.Email)
	if err != nil {
		s.log.Error("failed to create a user", "err", err)
		return err // Error is already handled by the entity
	}

	err = s.repo.Insert(ctx, userEntity.Model())
	if err != nil {
		if err = dbutil.VerifyDuplicatedConstraintKey(err); err != nil {
			s.log.Error("duplicated user", "email", input.Email)
			s.metrics.RecordError("user", "duplicated-user")
			return err // Error is already handled by the helper
		}
		s.log.Error("failed to insert user", "err", err)
		s.metrics.RecordError("user", "insert-user")
		return fault.NewBadRequest("failed to insert user")
	}

	s.log.Debug(
		"user created successfully",
		"details", strings.Join(
			[]string{
				fmt.Sprintf("id: %s", userEntity.id),
				fmt.Sprintf("name: %s", userEntity.name),
				fmt.Sprintf("email: %s", userEntity.email),
			},
			"\n",
		),
	)

	// go func() {
	// 	r := retry.New()
	// 	err := r.Do(func() error {
	// 		return s.mail.Send(mail.SendParams{
	// 			From:    mail.NotificationSender,
	// 			To:      "rafaelferreirab2@gmail.com",
	// 			Subject: "Seja bem-vindo!",
	// 			File:    mail.WelcomeTmpl,
	// 			Data: map[string]any{
	// 				"Name": "Rafael",
	// 			},
	// 		})
	// 	})
	// 	if err != nil {
	// 		s.log.Error("failed to send email with retry", "err", err)
	// 	}
	// }()

	return nil
}
