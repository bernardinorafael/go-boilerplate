package code

import (
	"testing"
	"time"
	"unicode"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
	"github.com/stretchr/testify/assert"
)

func TestCodeEntity(t *testing.T) {
	t.Run("should successfully create entity", func(t *testing.T) {
		userID := uid.New("user")
		entity := NewEntity(userID)

		assert.NotNil(t, entity)
		assert.Equal(t, entity.userID, userID)
		assert.NotEmpty(t, entity.id)
		assert.NotEmpty(t, entity.code)
		assert.True(t, entity.active)
		assert.Equal(t, entity.attempts, int64(0))
		assert.Nil(t, entity.usedAt)
		assert.True(t, entity.expiresAt.After(time.Now()))
		assert.Equal(t, entity.createdAt, entity.updatedAt)
	})

	t.Run("should create entity from model", func(t *testing.T) {
		model := model.Code{
			ID:        "code_0314612f4c66a7a34a2ekj89",
			UserID:    "user_030d8d2f1f808a1brba9dc31",
			Code:      "123456",
			Active:    true,
			Attempts:  1,
			UsedAt:    nil,
			ExpiresAt: time.Now().Add(time.Minute * 15),
			CreatedAt: time.Now().Add(-time.Hour * 24 * 2),
			UpdatedAt: time.Now().Add(-time.Hour * 24 * 2),
		}

		entity := NewEntityFromModel(model)

		assert.NotNil(t, entity)
		assert.Equal(t, entity.id, model.ID)
		assert.Equal(t, entity.userID, model.UserID)
		assert.Equal(t, entity.code, model.Code)
		assert.Equal(t, entity.active, model.Active)
		assert.Equal(t, entity.attempts, model.Attempts)
		assert.Equal(t, entity.usedAt, model.UsedAt)
		assert.Equal(t, entity.expiresAt, model.ExpiresAt)
		assert.Equal(t, entity.createdAt, model.CreatedAt)
		assert.Equal(t, entity.updatedAt, model.UpdatedAt)
	})

	t.Run("should generate code correctly", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		assert.NotNil(t, entity)
		assert.Len(t, entity.code, 6)
		assert.NotContains(t, entity.code, unicode.Punct)  // Cannot contain any punctuation
		assert.NotContains(t, entity.code, unicode.Symbol) // Cannot contain any special characters
		assert.NotContains(t, entity.code, unicode.Space)  // Cannot contain any whitespace characters
	})

	t.Run("should correctly check if a code is expired", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		assert.NotNil(t, entity)
		assert.False(t, entity.IsExpired())

		entity.expiresAt = time.Now().Add(-time.Hour)
		assert.True(t, entity.IsExpired())
	})

	t.Run("should check if a code reaches the max attempts", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		// Code initially has 0 attempts
		assert.NotNil(t, entity)
		assert.Equal(t, entity.attempts, int64(0))
		assert.False(t, entity.IsMaxAttempts())

		entity.IncrementAttempt()
		assert.Equal(t, entity.attempts, int64(1))
		assert.False(t, entity.IsMaxAttempts())

		entity.IncrementAttempt()
		assert.Equal(t, entity.attempts, int64(2))
		assert.False(t, entity.IsMaxAttempts())

		entity.IncrementAttempt()
		assert.Equal(t, entity.attempts, int64(3))
		assert.True(t, entity.IsMaxAttempts())
	})

	t.Run("should deactivate a code", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		assert.NotNil(t, entity)
		assert.True(t, entity.active)

		// Capture the original updatedAt time
		originalUpdatedAt := entity.updatedAt

		// Small delay to ensure different timestamps
		time.Sleep(time.Millisecond)

		entity.Deactivate()

		assert.False(t, entity.active)
		assert.True(t, entity.updatedAt.After(originalUpdatedAt))
	})

	t.Run("should increment attempts", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		assert.NotNil(t, entity)
		assert.Equal(t, entity.attempts, int64(0))

		originalUpdatedAt := entity.updatedAt

		time.Sleep(time.Millisecond)

		entity.IncrementAttempt()

		assert.Equal(t, entity.attempts, int64(1))
		assert.True(t, entity.updatedAt.After(originalUpdatedAt))
	})

	t.Run("should mark a code as used", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		assert.NotNil(t, entity)
		assert.Nil(t, entity.usedAt)

		originalUpdatedAt := entity.updatedAt

		time.Sleep(time.Millisecond)

		entity.MarkAsUsed()

		assert.NotNil(t, entity.usedAt)
		assert.True(t, entity.usedAt.After(originalUpdatedAt))
		assert.True(t, entity.updatedAt.After(originalUpdatedAt))
	})

	t.Run("should convert entity to model", func(t *testing.T) {
		entity := NewEntity(uid.New("user"))

		assert.NotNil(t, entity)

		model := entity.Model()

		assert.Equal(t, model.ID, entity.id)
		assert.Equal(t, model.UserID, entity.userID)
		assert.Equal(t, model.Code, entity.code)
		assert.Equal(t, model.Active, entity.active)
		assert.Equal(t, model.Attempts, entity.attempts)
		assert.Equal(t, model.UsedAt, entity.usedAt)
		assert.Equal(t, model.ExpiresAt, entity.expiresAt)
		assert.Equal(t, model.CreatedAt, entity.createdAt)
		assert.Equal(t, model.UpdatedAt, entity.updatedAt)
	})
}
