package session

import (
	"context"
	"fmt"
	"gulg/internal/modules/user"
	"gulg/pkg/fault"
	"gulg/pkg/logging"

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

func (s service) CreateSession(ctx context.Context, userId, ip, agent, refreshToken string) (sessionId string, err error) {
	userRecord, err := s.userService.GetUserByID(ctx, userId)
	if err != nil {
		return "", err // The error is already being handled in the user service
	}

	// Output: "Chrome em Windows"
	rawAgent := useragent.NewParser().Parse(agent)
	userAgent := fmt.Sprintf("%s em %s", rawAgent.GetBrowser(), rawAgent.GetOS())

	sess, err := New(userRecord.ID, ip, userAgent, refreshToken)
	if err != nil {
		s.log.Errorw(ctx, "failed to create session entity", logging.Err(err))
		return "", fault.NewUnprocessableEntity("failed to create session entity", err)
	}
	model := sess.ToModel()

	err = s.sessionRepo.Insert(ctx, model)
	if err != nil {
		s.log.Errorw(ctx, "failed to insert session entity", logging.Err(err))
		return "", fault.NewBadRequest("failed to insert session entity", err)
	}

	return sess.ID(), nil
}
