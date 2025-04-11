package session

import (
	"context"
	"fmt"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"

	"github.com/bernardinorafael/gogem/pkg/fault"
	"github.com/medama-io/go-useragent"
)

type service struct {
	log         logging.Logger
	sessionRepo Repository
	userService user.Service
	secretKey   string
}

func NewService(log logging.Logger, sessionRepo Repository, userService user.Service, secretKey string) Service {
	return &service{
		log:         log,
		sessionRepo: sessionRepo,
		userService: userService,
		secretKey:   secretKey,
	}
}

func (s service) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewAccessToken, error) {
	claims, err := token.Verify(s.secretKey, refreshToken)
	if err != nil {
		return nil, fault.NewUnauthorized("invalid refresh token")
	}

	sessRecord, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve session")
	}
	session := NewFromModel(*sessRecord)

	if session.IsExpired() {
		return nil, fault.NewBadRequest("session has expired")
	}

	if session.UserID() != claims.UserID {
		return nil, fault.NewUnauthorized("unauthorized user")
	}

	newAccessToken, _, err := token.Gen(s.secretKey, claims.UserID, time.Minute*15)
	if err != nil {
		return nil, fault.NewBadRequest("failed to generate access token")
	}

	return &dto.RenewAccessToken{
		AccessToken:        newAccessToken,
		AccessTokenExpires: time.Now().Add(time.Minute * 15),
	}, nil
}

func (s service) GetAllSessions(ctx context.Context) ([]dto.SessionResponse, error) {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error(ctx, "context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}

	records, err := s.sessionRepo.GetAllByUserID(ctx, c.UserID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve sessions")
	}

	// Pre-allocate the slice to avoid reallocations
	// This is more efficient than appending to the slice
	sessions := make([]dto.SessionResponse, len(records))
	for i, s := range records {
		sessions[i] = dto.SessionResponse{
			ID:      s.ID,
			Agent:   s.Agent,
			IP:      s.IP,
			Active:  s.Active,
			Created: s.Created,
			Updated: s.Updated,
		}
	}

	return sessions, nil
}

func (s service) CreateSession(ctx context.Context, input dto.CreateSession) (sessionId string, err error) {
	userRecord, err := s.userService.GetUserByID(ctx, input.UserID)
	if err != nil {
		return "", err // The error is already being handled in the user service
	}

	rawAgent := useragent.NewParser().Parse(input.Agent)
	agent := fmt.Sprintf("%s em %s", rawAgent.Browser(), rawAgent.OS())
	// In case the user agent is not a browser, we will use "unknown agent"
	// Likely to happen in mobile devices or in CLI/Postman and similar tools
	// Output: "<browser> em <os>"
	if rawAgent.Browser() == "" || rawAgent.OS() == "" {
		agent = "unknown agent"
	}

	sess, err := New(userRecord.ID, input.IP, agent, input.RefreshToken)
	if err != nil {
		s.log.Errorw(ctx, "failed to create session entity", logging.Err(err))
		return "", fault.NewUnprocessableEntity("failed to create session entity")
	}
	model := sess.ToModel()

	err = s.sessionRepo.Insert(ctx, model)
	if err != nil {
		s.log.Errorw(ctx, "failed to insert session entity", logging.Err(err))
		return "", fault.NewBadRequest("failed to insert session entity")
	}

	return sess.ID(), nil
}
