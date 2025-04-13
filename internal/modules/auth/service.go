package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/token"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/mail"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/session"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/crypto"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
)

const (
	accessTokenDuration  = time.Minute * 15    // 15 minutes
	refreshTokenDuration = time.Hour * 24 * 30 // 30 days
)

type ServiceConfig struct {
	SecretKey string

	UserService    user.Service
	UserRepo       user.Repository
	SessionService session.Service
	SessionRepo    session.Repository

	Mailer  *mail.Mail
	Cache   *cache.Cache
	Metrics *metric.Metric
}

// TODO: Remove userService dependency and user only the userRepo
// TODO: Remove sessionService dependency and user only the sessionRepo
type service struct {
	userService    user.Service
	userRepo       user.Repository
	sessionService session.Service
	sessionRepo    session.Repository
	mailer         *mail.Mail
	cache          *cache.Cache
	metrics        *metric.Metric
	secretKey      string
}

func NewService(c ServiceConfig) Service {
	return &service{
		userService:    c.UserService,
		userRepo:       c.UserRepo,
		sessionService: c.SessionService,
		sessionRepo:    c.SessionRepo,
		mailer:         c.Mailer,
		cache:          c.Cache,
		metrics:        c.Metrics,
		secretKey:      c.SecretKey,
	}
}

func (s service) Logout(ctx context.Context) error {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		slog.Error("context does not contain auth key")
		return fault.NewUnauthorized("access token no provided")
	}

	sessRecord, err := s.sessionRepo.GetActiveByUserID(ctx, c.UserID)
	if err != nil {
		s.metrics.RecordError("auth", "get-active-user-by-id")
		return fault.NewBadRequest("failed to retrieve active session")
	} else if sessRecord == nil {
		return fault.NewNotFound("active session not found")
	}

	sess := session.NewFromModel(*sessRecord)
	sess.Deactivate()

	err = s.sessionRepo.Update(ctx, sess.Model())
	if err != nil {
		return fault.NewBadRequest("failed to deactivate session")
	}

	err = s.cache.Delete(ctx, fmt.Sprintf("sess:%s", c.UserID))
	if err != nil {
		slog.Error("failed to delete session from cache", "error", err)
	}

	go func() {
		ctx := context.Background()

		cacheKey := fmt.Sprintf("sess:%s", c.UserID)
		has, err := s.cache.Has(ctx, cacheKey)
		if err != nil {
			slog.Error("failed to check if session is in cache", "error", err)
		} else if has {
			err = s.cache.Delete(ctx, c.UserID)
			if err != nil {
				slog.Error("failed to delete session from cache", "error", err)
			}
		}
	}()

	return nil
}

func (s service) Activate(ctx context.Context, userId string) error {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		s.metrics.RecordError("auth", "get-user-by-ud")
		return fault.NewBadRequest("failed to get user by id")
	} else if userRecord == nil {
		return fault.NewNotFound("user not found")
	}

	if userRecord.Enabled {
		return fault.New(
			"expired activation link",
			fault.WithHTTPCode(http.StatusBadRequest),
			fault.WithTag(fault.EXPIRED),
		)
	}

	u := user.NewFromModel(*userRecord)
	u.Enable()

	err = s.userRepo.Update(ctx, u.Model())
	if err != nil {
		s.metrics.RecordError("auth", "update-user")
		return fault.NewBadRequest("failed to update user")
	}

	return nil
}

func (s service) GetSignedUser(ctx context.Context) (*dto.UserResponse, error) {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		slog.Error("context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}

	userRecord, err := s.userRepo.GetByID(ctx, c.UserID)
	if err != nil {
		s.metrics.RecordError("auth", "get-user-by-id")
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := &dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Username:  userRecord.Username,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		Locked:    userRecord.Locked,
		Created:   userRecord.Created,
		Updated:   userRecord.Updated,
	}

	return user, nil
}

func (s service) Register(ctx context.Context, input dto.CreateUser) error {
	_, err := s.userService.CreateUser(ctx, input)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return err // The error is already being handled in the user service
	}

	// go func() {
	// 	params := mail.SendParams{
	// 		From:    <email-sender>,
	// 		To:      <user-email>,
	// 		Subject: "Activate your account",
	// 		File:    "activate_user.html",
	// 		Data: map[string]any{
	// 			// TODO: Change the activation link
	// 			"ActivationLink": <activation-link>,
	// 			"Name":           <user-name>,
	// 		},
	// 	}
	// 	err := s.mailer.Send(params)
	// 	if err != nil {
	// 		s.metrics.RecordError("mailer", "send")
	// 		slog.Error("failed to send email", "error", err)
	// 	}
	// }()

	return nil
}

func (s service) Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.metrics.RecordError("auth", "get-user-by-email")
		return nil, fault.NewBadRequest("failed to get user by email")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}
	userID := userRecord.ID

	if !crypto.PasswordMatches(password, userRecord.Password) {
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	if userRecord.Locked {
		return nil, fault.New(
			"user is locked",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.LOCKED_USER),
		)
	}

	if !userRecord.Enabled {
		return nil, fault.New(
			"user must enable account to login",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.DISABLED_USER),
		)
	}

	err = s.sessionRepo.DeactivateAll(ctx, userID)
	if err != nil {
		s.metrics.RecordError("auth", "deactivate-all-sessions")
		return nil, fault.NewBadRequest("failed to deactivate user sessions")
	}

	accessToken, _, err := token.Gen(s.secretKey, userID, accessTokenDuration)
	if err != nil {
		return nil, fault.NewUnauthorized(err.Error())
	}
	refreshToken, _, err := token.Gen(s.secretKey, userID, refreshTokenDuration)
	if err != nil {
		return nil, fault.NewUnauthorized(err.Error())
	}

	params := dto.CreateSession{
		IP:           ip,
		Agent:        agent,
		UserID:       userID,
		RefreshToken: refreshToken,
	}
	sess, err := s.sessionService.CreateSession(ctx, params)
	if err != nil {
		return nil, err // The error is already being handled in the user service
	}

	response := dto.LoginResponse{
		SessionID:    sess.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}
