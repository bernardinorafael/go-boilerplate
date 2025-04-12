package product

import (
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
)

type product struct {
	id      string
	name    string
	price   int64
	created time.Time
	updated time.Time
}

func New(name string, price int64) *product {
	p := product{
		id:      uid.New("prod"),
		name:    name,
		price:   price * 100, // Convert to cents
		created: time.Now(),
		updated: time.Now(),
	}

	return &p
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

func (p *product) Model() model.Product {
	return model.Product{
		ID:      p.id,
		Name:    p.name,
		Price:   p.price,
		Created: p.created,
		Updated: p.updated,
	}
}
