package session

import (
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
)

const (
	// ttl is the time to live for a session
	ttl = time.Hour * 24 * 30 // 30 days
)

type session struct {
	id           string
	userId       string
	ip           string
	agent        string
	refreshToken string
	active       bool
	expires      time.Time
	created      time.Time
	updated      time.Time
}

func New(userId, ip, agent, refresh string) (*session, error) {
	s := session{
		id:           uid.New("sess"),
		userId:       userId,
		ip:           ip,
		agent:        agent,
		refreshToken: refresh,
		active:       true,
		expires:      time.Now().Add(ttl),
		created:      time.Now(),
		updated:      time.Now(),
	}

	if err := s.validate(); err != nil {
		return nil, fault.New(
			"failed to create session entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &s, nil
}

func (s *session) NewFromModel(m model.Session) *session {
	return &session{
		id:           m.ID,
		userId:       m.UserID,
		ip:           m.IP,
		agent:        m.Agent,
		refreshToken: m.RefreshToken,
		active:       m.Active,
		expires:      m.Expires,
		created:      m.Created,
		updated:      m.Updated,
	}
}

func (s *session) validate() error {
	if s.userId == "" {
		return fault.New("user id is required")
	}
	if s.ip == "" {
		return fault.New("ip is required")
	}
	if s.agent == "" {
		return fault.New("agent is required")
	}
	if s.refreshToken == "" {
		return fault.New("refresh token is required")
	}

	return nil
}

func (s *session) ToModel() model.Session {
	return model.Session{
		ID:           s.id,
		UserID:       s.userId,
		IP:           s.ip,
		Agent:        s.agent,
		RefreshToken: s.refreshToken,
		Active:       s.active,
		Expires:      s.expires,
		Created:      s.created,
		Updated:      s.updated,
	}
}

func (s *session) IsExpired() bool {
	return s.expires.Before(time.Now())
}

func (s *session) Activate() {
	s.active = true
	s.updated = time.Now()
}

func (s *session) Deactivate() {
	s.active = false
	s.updated = time.Now()
}

func (s *session) ID() string           { return s.id }
func (s *session) UserID() string       { return s.userId }
func (s *session) IP() string           { return s.ip }
func (s *session) Agent() string        { return s.agent }
func (s *session) RefreshToken() string { return s.refreshToken }
func (s *session) Active() bool         { return s.active }
func (s *session) Expires() time.Time   { return s.expires }
func (s *session) Created() time.Time   { return s.created }
func (s *session) Updated() time.Time   { return s.updated }
