package product

import (
	"net/http"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	minNameLength = 3
	maxNameLength = 100
	minPrice      = 1
)

type product struct {
	id      string
	name    string
	price   int64
	created time.Time
	updated time.Time
}

func NewEntity(name string, price int64) (*product, error) {
	p := product{
		id:      uid.New("prod"),
		name:    name,
		price:   price * 100, // Convert to cents
		created: time.Now(),
		updated: time.Now(),
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

func NewFromModel(model model.Product) *product {
	return &product{
		id:      model.ID,
		name:    model.Name,
		price:   model.Price,
		created: model.Created,
		updated: model.Updated,
	}
}

func (p *product) ChangeName(name string) {
	p.name = name
	p.updated = time.Now()
}

func (p *product) ChangePrice(price int64) {
	p.price = price
	p.updated = time.Now()
}

func (p product) validate() error {
	err := v.ValidateStruct(
		&p,
		v.Field(
			&p.name,
			v.Required.Error("this field is required"),
			v.Length(minNameLength, maxNameLength).Error("name must be between 3 and 100 characters"),
		),
		v.Field(
			&p.price,
			v.Required.Error("this field is required"),
			v.Min(minPrice).Error("price must be greater than 0"),
		),
	)
	if err != nil {
		return fault.New(
			"failed to validate product",
			fault.WithHTTPCode(http.StatusUnprocessableEntity),
			fault.WithTag(fault.ValidationError),
			fault.WithValidationError(err),
		)
	}

	return nil
}

func (p *product) Model() model.Product {
	return model.Product{
		ID:      p.id,
		Name:    p.name,
		Price:   p.price,
		Created: p.created,
		Updated: p.updated,
	}
}
