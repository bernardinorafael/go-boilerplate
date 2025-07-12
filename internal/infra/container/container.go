package container

import (
	"context"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/config"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/category"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/code"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/product"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/user"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/pg"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/mail"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/charmbracelet/log"
)

type Container struct {
	Config  *config.Config
	Logger  *log.Logger
	Metrics *metric.Metric

	Cache *cache.Cache
	DB    *pg.Database
	Mail  *mail.Mail

	ProductRepo  product.Repository
	UserRepo     user.Repository
	CodeRepo     code.Repository
	CategoryRepo category.Repository

	ProductService  product.Service
	UserService     user.Service
	CodeService     code.Service
	CategoryService category.Service
}

// NewContainer creates a new container instance.
//
// It initializes the infrastructure, repositories, and services.
//
// Parameters:
//   - ctx: The context for the container.
//   - cfg: The configuration for the container.
//   - logger: The logger for the container.
func NewContainer(ctx context.Context, cfg *config.Config, logger *log.Logger) (*Container, error) {
	container := &Container{
		Config:  cfg,
		Logger:  logger,
		Metrics: metric.New(),
	}

	err := container.initInfra(ctx)
	if err != nil {
		return nil, err
	}

	container.initRepositories()
	container.initServices()

	return container, nil
}

func (c *Container) initInfra(ctx context.Context) error {
	cache, err := cache.New(ctx, c.Config.RedisHost, c.Config.RedisPort, c.Config.RedisPassword)
	if err != nil {
		c.Logger.Error("failed to connect to cache", "error", err)
		return err
	}
	c.Cache = cache

	db, err := pg.NewConnection(c.Config.PostgresDSN)
	if err != nil {
		c.Logger.Error("failed to connect database", "error", err)
		return err
	}
	c.DB = db

	mailService := mail.New(ctx, c.Logger, c.Config.ResendKey, time.Second*5)
	c.Mail = mailService

	return nil
}

func (c *Container) initRepositories() {
	timeout := time.Second * 2

	c.ProductRepo = product.NewRepo(c.DB.DB(), timeout)
	c.UserRepo = user.NewRepo(c.DB.DB(), timeout)
	c.CodeRepo = code.NewRepo(c.DB.DB(), timeout)
	c.CategoryRepo = category.NewRepo(c.DB.DB(), timeout)
}

func (c *Container) initServices() {
	c.CategoryService = category.NewService(category.ServiceConfig{
		Log:          c.Logger,
		CategoryRepo: c.CategoryRepo,
	})

	c.CodeService = code.NewService(code.ServiceConfig{
		Log:      c.Logger,
		CodeRepo: c.CodeRepo,
		Metrics:  c.Metrics,
		Cache:    c.Cache,
		Mail:     c.Mail,
	})

	c.ProductService = product.NewService(product.ServiceConfig{
		Log:         c.Logger,
		ProductRepo: c.ProductRepo,
		Metrics:     c.Metrics,
		Cache:       c.Cache,
	})

	c.UserService = user.NewService(user.ServiceConfig{
		Log:     c.Logger,
		Metrics: c.Metrics,
		Cache:   c.Cache,
		Mail:    c.Mail,

		UserRepo:    c.UserRepo,
		CodeService: c.CodeService,

		AccessTokenDuration: c.Config.JWTAccessTokenDuration,
		SecretKey:           c.Config.JWTSecretKey,
	})
}

func (c *Container) Close() error {
	if c.Cache != nil {
		c.Cache.Close()
	}

	if c.DB != nil {
		c.DB.Close()
	}

	return nil
}
