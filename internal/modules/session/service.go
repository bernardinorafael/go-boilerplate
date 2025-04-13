package session

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/token"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"

	"github.com/medama-io/go-useragent"
)

type ServiceConfig struct {
	SessionRepo Repository
	UserService user.Service

	Cache     *cache.Cache
	Metrics   *metric.Metric
	SecretKey string
}

type service struct {
	sessionRepo Repository
	userService user.Service
	cache       *cache.Cache
	metrics     *metric.Metric
	secretKey   string
}

func NewService(c ServiceConfig) Service {
	return &service{
		sessionRepo: c.SessionRepo,
		userService: c.UserService,
		cache:       c.Cache,
		metrics:     c.Metrics,
		secretKey:   c.SecretKey,
	}
}

func (s service) GetSessionByUserID(ctx context.Context, userID string) (*dto.SessionResponse, error) {
	var cachedSession *model.Session

	err := s.cache.GetStruct(ctx, fmt.Sprintf("sess:%s", userID), &cachedSession)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CACHE_MISS:
			s.metrics.RecordCacheMiss("session-service")
			slog.Info("cache:miss session not found")
		default:
			slog.Error("failed to query session from cache")
		}
	}

	if cachedSession != nil {
		if time.Now().After(cachedSession.Expires) {
			return nil, fault.NewBadRequest("session has expired")
		}

		s.metrics.RecordCacheHit("session-service")

		return &dto.SessionResponse{
			ID:      cachedSession.ID,
			Agent:   cachedSession.Agent,
			IP:      cachedSession.IP,
			Active:  cachedSession.Active,
			Created: cachedSession.Created,
			Updated: cachedSession.Updated,
		}, nil
	}

	sessionRecord, err := s.sessionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		s.metrics.RecordError("sessions", "get-active-by-user-id")
		return nil, fault.NewBadRequest("failed to retrieve session")
	} else if sessionRecord == nil {
		return nil, fault.NewNotFound("session not found")
	}

	res := &dto.SessionResponse{
		ID:      sessionRecord.ID,
		Agent:   sessionRecord.Agent,
		IP:      sessionRecord.IP,
		Active:  sessionRecord.Active,
		Created: sessionRecord.Created,
		Updated: sessionRecord.Updated,
	}

	cacheKey := fmt.Sprintf("sess:%s", userID)
	err = s.cache.SetStruct(ctx, cacheKey, sessionRecord, time.Minute*30)
	if err != nil {
		slog.Error("failed to caching session", "error", err)
	}

	return res, nil
}

func (s service) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewAccessToken, error) {
	claims, err := token.Verify(s.secretKey, refreshToken)
	if err != nil {
		return nil, fault.NewUnauthorized("invalid refresh token")
	}

	sessRecord, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		s.metrics.RecordError("sessions", "get-by-refresh-token")
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
		slog.Error("context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}

	records, err := s.sessionRepo.GetAllByUserID(ctx, c.UserID)
	if err != nil {
		s.metrics.RecordError("sessions", "get-all-by-user-id")
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

func (s service) CreateSession(ctx context.Context, input dto.CreateSession) (*dto.SessionResponse, error) {
	userRecord, err := s.userService.GetUserByID(ctx, input.UserID)
	if err != nil {
		s.metrics.RecordError("sessions", "get-user-by-id")
		return nil, err // The error is already being handled in the user service
	}
	userID := userRecord.ID

	rawAgent := useragent.NewParser().Parse(input.Agent)
	agent := fmt.Sprintf("%s em %s", rawAgent.Browser(), rawAgent.OS())
	// In case the user agent is not a browser, we will use "unknown agent"
	// Likely to happen in mobile devices or in CLI/Postman and similar tools
	// Output: "<browser> em <os>"
	if rawAgent.Browser() == "" || rawAgent.OS() == "" {
		agent = "unknown agent"
	}

	sess, err := New(userID, input.IP, agent, input.RefreshToken)
	if err != nil {
		return nil, fault.NewUnprocessableEntity("failed to create session entity")
	}

	err = s.sessionRepo.Insert(ctx, sess.Model())
	if err != nil {
		s.metrics.RecordError("sessions", "insert-session")
		return nil, fault.NewBadRequest("failed to insert session entity")
	}

	res := dto.SessionResponse{
		ID:      sess.ID(),
		Agent:   sess.Agent(),
		IP:      sess.IP(),
		Active:  sess.Active(),
		Created: sess.Created(),
		Updated: sess.Updated(),
	}

	return &res, nil
}
