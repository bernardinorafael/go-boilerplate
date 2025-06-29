package code

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/mail"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/charmbracelet/log"
)

type ServiceConfig struct {
	CodeRepo Repository
	Log      *log.Logger
	Metrics  *metric.Metric
	Cache    *cache.Cache
	Mail     *mail.Mail
}

type service struct {
	log     *log.Logger
	metrics *metric.Metric
	cache   *cache.Cache
	mail    *mail.Mail
	repo    Repository
}

func NewService(c ServiceConfig) *service {
	return &service{
		log:     c.Log,
		repo:    c.CodeRepo,
		metrics: c.Metrics,
		cache:   c.Cache,
		mail:    c.Mail,
	}
}

func (s service) VerifyCode(ctx context.Context, userID, code string) (bool, error) {
	model, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			return false, fault.NewNotFound("no code found for user")
		default:
			s.log.Error("failed to retrieve code", "err", err)
			return false, fault.NewBadRequest("failed to verify code")
		}
	}

	entity := NewEntityFromModel(*model)
	entity.IncrementAttempt()

	if entity.IsExpired() {
		s.log.Debug("otp code has expired", "userId", userID)

		// If the code has expired, we increment the attempts and return an error
		if err := s.repo.Update(ctx, entity.Model()); err != nil {
			s.log.Error("failed to increment code attempt", "err", err)
			return false, fault.NewBadRequest("failed to verify code")
		}

		return false, fault.NewForbidden("code has expired")
	}

	if entity.IsMaxAttempts() {
		s.log.Info("max code attempts reached")

		// If the code has reached the max attempts, we increment the attempts and return an error
		if err := s.repo.Update(ctx, entity.Model()); err != nil {
			s.log.Error("failed to increment code attempt", "err", err)
			return false, fault.NewBadRequest("failed to verify code")
		}

		return false, fault.NewForbidden("max code attempts reached for this code")
	}

	if entity.code != code {
		s.log.Debug("otp code verification failed", "userId", userID)

		// If the code does not match, we increment the attempts and return an error
		if err := s.repo.Update(ctx, entity.Model()); err != nil {
			s.log.Error("failed to increment code attempt", "err", err)
			return false, fault.NewBadRequest("failed to verify code")
		}

		return false, nil
	}

	entity.MarkAsUsed()

	err = s.repo.Update(ctx, entity.Model())
	if err != nil {
		s.log.Error("failed to update code", "err", err, "userId", userID)
		return false, fault.NewBadRequest("failed to verify code")
	}

	s.log.Debug("otp code verification succeeded", "userId", userID)

	return true, nil
}

func (s service) CreateCode(ctx context.Context, userID string) error {
	s.log.Debug("trying to create code", "userId", userID)

	err := s.repo.InactivateAll(ctx, userID)
	if err != nil {
		s.log.Error("failed to inactivate all codes", "err", err)
		return fault.NewBadRequest("failed to create code")
	}

	entity := NewEntity(userID)

	err = s.repo.Insert(ctx, entity.Model())
	if err != nil {
		s.log.Error("failed to insert code", "err", err)
		return fault.NewBadRequest("failed to create otp code")
	}

	return nil
}
