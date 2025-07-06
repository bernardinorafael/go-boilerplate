package code

import (
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
	"github.com/nanorand/nanorand"
)

const (
	ttl         = time.Minute * 3
	codeSize    = 6
	maxAttempts = 3
)

type code struct {
	id        string
	userID    string
	code      string
	active    bool
	attempts  int64
	usedAt    *time.Time
	expiresAt time.Time
	createdAt time.Time
	updatedAt time.Time
}

func NewEntityFromModel(m model.Code) *code {
	return &code{
		id:        m.ID,
		userID:    m.UserID,
		code:      m.Code,
		active:    m.Active,
		attempts:  m.Attempts,
		usedAt:    m.UsedAt,
		expiresAt: m.ExpiresAt,
		createdAt: m.CreatedAt,
		updatedAt: m.UpdatedAt,
	}
}

func NewEntity(userID string) *code {
	now := time.Now()
	c, _ := nanorand.Gen(codeSize)

	code := code{
		id:        uid.New("code"),
		userID:    userID,
		code:      c,
		active:    true,
		attempts:  0,
		usedAt:    nil,
		expiresAt: now.Add(ttl),
		createdAt: now,
		updatedAt: now,
	}

	return &code
}

func (c *code) IsExpired() bool {
	return time.Now().After(c.expiresAt)
}

func (c *code) IsMaxAttempts() bool {
	return c.attempts >= maxAttempts
}

func (c *code) Deactivate() {
	c.active = false
	c.updatedAt = time.Now()
}

func (c *code) IncrementAttempt() {
	c.attempts++
	c.updatedAt = time.Now()
}

func (c *code) MarkAsUsed() {
	now := time.Now()
	c.usedAt = &now
	c.active = false
	c.updatedAt = now
}

func (c *code) Model() model.Code {
	return model.Code{
		ID:        c.id,
		UserID:    c.userID,
		Code:      c.code,
		Active:    c.active,
		Attempts:  c.attempts,
		UsedAt:    c.usedAt,
		ExpiresAt: c.expiresAt,
		CreatedAt: c.createdAt,
		UpdatedAt: c.updatedAt,
	}
}
