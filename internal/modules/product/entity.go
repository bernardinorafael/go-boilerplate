package product

import (
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
)

const (
	minNameLength = 3
	minPrice      = 1
)

type product struct {
	id      string
	name    string
	price   int64
	created time.Time
	updated time.Time
}

func New(name string, price int64) (*product, error) {
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

func (p *product) validate() error {
	if p.name == "" {
		return fault.NewUnprocessableEntity("name is required")
	}

	if len(p.name) < minNameLength {
		return fault.NewUnprocessableEntity("name must be at least 3 characters")
	}

	if p.price < minPrice {
		return fault.NewUnprocessableEntity("price must be greater than 0")
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
