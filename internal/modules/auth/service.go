package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/mail"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/session"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/crypto"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
)

const (
	accessTokenDuration  = time.Minute * 15    // 15 minutes
	refreshTokenDuration = time.Hour * 24 * 30 // 30 days
)

// TODO: Remove userService dependency and user only the userRepo
// TODO: Remove sessionService dependency and user only the sessionRepo
type service struct {
	log            logging.Logger
	userService    user.Service
	userRepo       user.Repository
	sessionService session.Service
	sessionRepo    session.Repository
	mailer         *mail.Mail
	cache          *cache.Cache
	secretKey      string
}

func NewService(
	log logging.Logger,
	userService user.Service,
	userRepo user.Repository,
	sessionService session.Service,
	sessionRepo session.Repository,
	mailer *mail.Mail,
	cache *cache.Cache,
	secretKey string,
) Service {
	return &service{
		log:            log,
		userService:    userService,
		userRepo:       userRepo,
		sessionService: sessionService,
		sessionRepo:    sessionRepo,
		mailer:         mailer,
		cache:          cache,
		secretKey:      secretKey,
	}
}

func (s service) Logout(ctx context.Context) error {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error(ctx, "context does not contain auth key")
		return fault.NewUnauthorized("access token no provided")
	}

	sessRecord, err := s.sessionRepo.GetActiveByUserID(ctx, c.UserID)
	if err != nil {
		return fault.NewBadRequest("failed to retrieve active session")
	} else if sessRecord == nil {
		return fault.NewNotFound("active session not found")
	}

	session := session.NewFromModel(*sessRecord)
	session.Deactivate()

	err = s.sessionRepo.Update(ctx, session.ToModel())
	if err != nil {
		return fault.NewBadRequest("failed to deactivate session")
	}

	return nil
}

func (s service) Activate(ctx context.Context, userId string) error {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
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

	user := user.NewFromModel(*userRecord)
	user.Enable()

	err = s.userRepo.Update(ctx, user.ToModel())
	if err != nil {
		s.log.Errorw(ctx, "failed to update user", logging.Err(err))
		return fault.NewBadRequest("failed to update user")
	}

	return nil
}

func (s service) GetSignedUser(ctx context.Context) (*dto.UserResponse, error) {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error(ctx, "context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}

	var user *dto.UserResponse

	err := s.cache.GetStruct(ctx, c.UserID, &user)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CACHE_MISS:
			s.log.Infow(ctx, "user not found in cache", logging.Err(err))
		default:
			s.log.Errorw(ctx, "failed to get user from cache", logging.Err(err))
		}
	}

	// If the user is found in the cache, return it
	if user != nil {
		s.log.Info(ctx, "returning user from cache")
		return user, nil
	}

	userRecord, err := s.userRepo.GetByID(ctx, c.UserID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user = &dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Username:  userRecord.Username,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		Locked:    userRecord.Locked,
		Created:   userRecord.Created,
		Updated:   userRecord.Updated,
	}

	go func() {
		// Setting new context to avoid context deadline exceeded
		err = s.cache.SetStruct(context.Background(), user.ID, user, time.Minute*30)
		if err != nil {
			s.log.Errorw(ctx, "failed to set user in cache", logging.Err(err))
		}
	}()

	return user, nil
}

func (s service) Register(ctx context.Context, input dto.CreateUser) error {
	user, err := s.userService.CreateUser(ctx, input)
	if err != nil {
		s.log.Errorw(ctx, "failed to create user", logging.Err(err))
		return err // The error is already being handled in the user service
	}

	err = s.cache.SetStruct(ctx, user.ID, user, time.Minute*30)
	if err != nil {
		s.log.Errorw(ctx, "failed to set user in cache", logging.Err(err))
	}

	// TODO: Send to a queue
	// go func() {
	// 	params := mail.SendParams{
	// 		From:    <your-notification-sender>,
	// 		To:      input.Email,
	// 		Subject: "Activate your account",
	// 		File:    "activate_user.html",
	// 		Data: map[string]any{
	// 			"ActivationLink": <your-activation-link>,
	// 			"Name":           input.Name,
	// 		},
	// 	}
	// 	err := s.mailService.Send(params)
	// 	if err != nil {
	// 		s.log.Errorw(ctx, "failed to send email", logging.Err(err))
	// 	}
	// }()

	return nil
}

func (s service) Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to get user by id")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

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

	err = s.sessionRepo.DeactivateAll(ctx, userRecord.ID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to deactivate user sessions")
	}

	// Access token with 15 minutes expiration
	accessToken, _, err := token.Gen(s.secretKey, userRecord.ID, accessTokenDuration)
	if err != nil {
		s.log.Errorw(ctx, "failed to generate access token", logging.Err(err))
		return nil, fault.NewUnauthorized(err.Error())
	}
	// Refresh token with 30 days expiration
	refreshToken, _, err := token.Gen(s.secretKey, userRecord.ID, refreshTokenDuration)
	if err != nil {
		s.log.Errorw(ctx, "failed to generate refresh token", logging.Err(err))
		return nil, fault.NewUnauthorized(err.Error())
	}

	sessionParams := dto.CreateSession{
		UserID:       userRecord.ID,
		IP:           ip,
		Agent:        agent,
		RefreshToken: refreshToken,
	}
	sessionId, err := s.sessionService.CreateSession(ctx, sessionParams)
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
