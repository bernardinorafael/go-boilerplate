package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/session"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/crypto"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
)

const (
	accessTokenDuration  = time.Minute * 15    // 15 minutes
	refreshTokenDuration = time.Hour * 24 * 30 // 30 days
)

type service struct {
	log            logging.Logger
	userService    user.Service
	sessionService session.Service
	secretKey      string
}

func NewService(
	log logging.Logger,
	userService user.Service,
	sessionService session.Service,
	secretKey string,
) Service {
	return &service{
		log:            log,
		userService:    userService,
		sessionService: sessionService,
		secretKey:      secretKey,
	}
}

func (s service) GetSigned(ctx context.Context, userId string) (*model.User, error) {
	userRecord, err := s.userService.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err // The error is already being handled in the user service
	}

	return userRecord, nil
}

func (s service) Register(ctx context.Context, input dto.CreateUser) error {
	err := s.userService.CreateUser(ctx, input)
	if err != nil {
		s.log.Errorw(ctx, "failed to create user", logging.Err(err))
		return err // The error is already being handled in the user service
	}

	s.log.Infow(ctx, "user created", logging.String("user_email", input.Email))
	return nil
}

func (s service) Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Errorw(ctx, "failed to get user by email", logging.Err(err))
		return nil, err // The error is already being handled in the user service
	}

	if !crypto.PasswordMatches(password, user.Password) {
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	if user.Locked {
		return nil, fault.New(
			"user is locked",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.LOCKED_USER),
		)
	}

	if !user.Enabled {
		return nil, fault.New(
			"user must enable account to login",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.DISABLED_USER),
		)
	}

	// Access token with 15 minutes expiration
	accessToken, _, err := token.Gen(s.secretKey, user.ID, accessTokenDuration)
	if err != nil {
		s.log.Errorw(ctx, "failed to generate access token", logging.Err(err))
		return nil, fault.NewUnauthorized(err.Error())
	}
	// Refresh token with 30 days expiration
	refreshToken, _, err := token.Gen(s.secretKey, user.ID, refreshTokenDuration)
	if err != nil {
		s.log.Errorw(ctx, "failed to generate refresh token", logging.Err(err))
		return nil, fault.NewUnauthorized(err.Error())
	}

	sessionId, err := s.sessionService.CreateSession(ctx, user.ID, ip, agent, refreshToken)
	if err != nil {
		return nil, err // The error is already being handled in the user service
	}

	response := dto.LoginResponse{
		SessionID:    sessionId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}
