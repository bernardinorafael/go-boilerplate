package category

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/pagination"
	"github.com/bernardinorafael/go-boilerplate/pkg/strutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
	"github.com/charmbracelet/log"
)

type ServiceConfig struct {
	Log          *log.Logger
	CategoryRepo Repository
}

type service struct {
	log  *log.Logger
	repo Repository
}

func NewService(c ServiceConfig) *service {
	return &service{
		log:  c.Log,
		repo: c.CategoryRepo,
	}
}

func (s service) FindAll(ctx context.Context, search dto.SearchParams) (*pagination.Paginated[dto.CategoryResponse], error) {
	s.log.Debug(
		"trying to retrieve categories with",
		"details", strings.Join(
			[]string{
				fmt.Sprintf("term: %s", search.Term),
				fmt.Sprintf("sort: %s", search.Sort),
				fmt.Sprintf("limit: %d", search.Limit),
				fmt.Sprintf("page: %d", search.Page),
			},
			"\n",
		),
	)

	records, totalItems, err := s.repo.FindAll(ctx, search)
	if err != nil {
		s.log.Error("failed to retrieve categories", "err", err)
		return nil, fault.NewBadRequest("failed to retrieve categories")
	}

	var categories = make([]dto.CategoryResponse, len(records))
	for i, c := range records {
		categories[i] = dto.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			Active:    c.Active,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	s.log.Debug("categories retrieved successfully", "totalItems", totalItems)

	paginatedResponse := pagination.New(categories, totalItems, search.Page, search.Limit)
	return &paginatedResponse, nil
}

func (s service) Create(ctx context.Context, input dto.CreateCategory) error {
	record, err := s.repo.FindByName(ctx, input.Name)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			// category not found, continue with creation
			s.log.Debug("category not found", "name", input.Name)
		default:
			s.log.Error("failed to retrieve category", "err", err)
			return fault.NewBadRequest("failed to retrieve category")
		}
	}

	if record != nil && !record.Active {
		s.log.Debug("category reactivated", "name", input.Name)

		record.Active = true
		record.UpdatedAt = time.Now()
		record.DeletedAt = nil

		err := s.repo.Update(ctx, *record)
		if err != nil {
			s.log.Error("failed to reactivate category", "err", err)
			return fault.NewBadRequest("failed to reactivate category")
		}

		return nil
	}

	newCategory := model.Category{
		ID:        uid.New("cat"),
		Name:      input.Name,
		Slug:      strutil.Slugify(input.Name),
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	err = s.repo.Insert(ctx, newCategory)
	if err != nil {
		if err := dbutil.VerifyDuplicatedConstraintKey(err); err != nil {
			s.log.Error("duplicated product", "name", input.Name, "err", err)
			return err // Error is already handled by the helper
		}
		s.log.Error("failed to create category", "err", err)
		return fault.NewBadRequest("failed to create category")
	}

	s.log.Debug(
		"product created successfully",
		"details", strings.Join(
			[]string{
				fmt.Sprintf("id: %s", newCategory.ID),
				fmt.Sprintf("name: %s", newCategory.Name),
			},
			"\n",
		),
	)

	return nil
}

func (s service) Delete(ctx context.Context, categoryID string) error {
	model, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			return fault.NewNotFound("category not found")
		default:
			s.log.Error("failed to retrieve category", "err", err)
			return fault.NewBadRequest("failed to delete category")
		}
	}
	if model.DeletedAt != nil {
		return fault.NewNotFound("category not found")
	}

	now := time.Now()
	model.Active = false
	model.UpdatedAt = now
	model.DeletedAt = &now

	err = s.repo.Update(ctx, *model)
	if err != nil {
		s.log.Error("failed to update category", "err", err)
		return fault.NewBadRequest("failed to delete category")
	}

	s.log.Debug("category deleted", "categoryId", model.ID)

	return nil
}

func (s service) GetByID(ctx context.Context, categoryID string) (*dto.CategoryResponse, error) {
	s.log.Debug("trying to retrieve category", "categoryId", categoryID)

	model, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			return nil, fault.NewNotFound("category not found")
		default:
			s.log.Error("failed to retrieve category", "err", err)
			return nil, fault.NewBadRequest("failed to retrieve category")
		}
	}

	return &dto.CategoryResponse{
		ID:        model.ID,
		Name:      model.Name,
		Slug:      model.Slug,
		Active:    model.Active,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}
