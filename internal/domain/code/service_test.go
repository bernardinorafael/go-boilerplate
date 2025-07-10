package code

import (
	"context"
	"os"
	"testing"
	"time"

	codemock "github.com/bernardinorafael/go-boilerplate/__mocks/domain/code"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_VerifyCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := codemock.NewMockRepository(ctrl)
	logger := log.NewWithOptions(os.Stdout, log.Options{})

	mockService := &service{log: logger, repo: mockRepository}

	t.Run("success case - valid code", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Minute * 5)

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  0,
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("error case - code not found", func(t *testing.T) {
		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(nil, fault.NewNotFound("code not found")).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.NotFound, fault.GetTag(err))
	})

	t.Run("error case - repository error", func(t *testing.T) {
		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(nil, fault.NewBadRequest("database error")).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})

	t.Run("error case - expired code", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(-time.Minute) // código expirado há 1 minuto

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  0,
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.Forbidden, fault.GetTag(err))
	})

	t.Run("error case - max attempts reached", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Minute * 5)

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  3, // max attempts
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.Forbidden, fault.GetTag(err))
	})

	t.Run("error case - invalid code", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Minute * 5)

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  0,
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "654321") // código incorreto

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("error case - update fails after successful verification", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Minute * 5)

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  0,
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(fault.NewBadRequest("update failed")).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})

	t.Run("error case - update fails after expired code", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(-time.Minute) // código expirado

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  0,
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(fault.NewBadRequest("update failed")).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})

	t.Run("error case - update fails after max attempts", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Minute * 5)

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  3, // max attempts
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(fault.NewBadRequest("update failed")).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "123456")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})

	t.Run("error case - update fails after invalid code", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Minute * 5)

		codeModel := model.Code{
			ID:        "code_id",
			UserID:    "user_id",
			Code:      "123456",
			Active:    true,
			Attempts:  0,
			UsedAt:    nil,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepository.
			EXPECT().
			GetByUserID(gomock.Any(), "user_id").
			Return(&codeModel, nil).
			Times(1)

		mockRepository.
			EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(fault.NewBadRequest("update failed")).
			Times(1)

		result, err := mockService.VerifyCode(context.Background(), "user_id", "654321")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})
}

func TestService_CreateCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := codemock.NewMockRepository(ctrl)
	logger := log.NewWithOptions(os.Stdout, log.Options{})

	mockService := &service{log: logger, repo: mockRepository}

	t.Run("success case", func(t *testing.T) {
		mockRepository.
			EXPECT().
			InactivateAll(gomock.Any(), "user_id").
			Return(nil).
			Times(1)

		mockRepository.
			EXPECT().
			Insert(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		err := mockService.CreateCode(context.Background(), "user_id")

		assert.NoError(t, err)
	})

	t.Run("error case - inactivate all fails", func(t *testing.T) {
		mockRepository.
			EXPECT().
			InactivateAll(gomock.Any(), "user_id").
			Return(fault.NewBadRequest("database error")).
			Times(1)

		err := mockService.CreateCode(context.Background(), "user_id")

		assert.Error(t, err)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})

	t.Run("error case - insert fails", func(t *testing.T) {
		mockRepository.
			EXPECT().
			InactivateAll(gomock.Any(), "user_id").
			Return(nil).
			Times(1)

		mockRepository.
			EXPECT().
			Insert(gomock.Any(), gomock.Any()).
			Return(fault.NewBadRequest("insert failed")).
			Times(1)

		err := mockService.CreateCode(context.Background(), "user_id")

		assert.Error(t, err)
		assert.Equal(t, fault.BadRequest, fault.GetTag(err))
	})
}
