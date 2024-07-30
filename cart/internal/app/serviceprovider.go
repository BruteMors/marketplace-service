package app

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/config"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart"
	cartRepository "github.com/BruteMors/marketplace-service/cart/internal/repository/inmemory/cart"
	cartService "github.com/BruteMors/marketplace-service/cart/internal/service/cart"
	"github.com/BruteMors/marketplace-service/cart/pkg/httpclient"
	"github.com/BruteMors/marketplace-service/cart/pkg/lomsservice"
	"github.com/BruteMors/marketplace-service/cart/pkg/productservice"
	"github.com/go-playground/validator/v10"
)

type serviceProvider struct {
	httpConfig                *config.HTTPServerConfig
	httpClient                *httpclient.HttpClient
	productService            *productservice.ProductService
	cartHttpApi               *cart.HttpApi
	cartService               *cartService.Service
	cartRepository            *cartRepository.Repository
	cartRepositoryWithMetrics *cartRepository.RepositoryWithMetrics
	lomsService               *lomsservice.Client
	validator                 *validator.Validate
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) HTTPServerConfig() *config.HTTPServerConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPServerConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) HTTPClient(_ context.Context) *httpclient.HttpClient {
	if s.httpClient == nil {
		cfg, err := config.NewHTTPClientConfig()
		if err != nil {
			log.Fatalf("failed to get http client config: %s", err.Error())
		}

		timeout, err := time.ParseDuration(cfg.Timeout)
		if err != nil {
			log.Fatalf("failed to parse http client timeout: %s", err.Error())
		}

		retries, err := strconv.Atoi(cfg.Retries)
		if err != nil {
			log.Fatalf("failed to parse http client retries: %s", err.Error())
		}

		retryDelay, err := time.ParseDuration(cfg.RetryDelay)
		if err != nil {
			log.Fatalf("failed to parse http client retry delay: %s", err.Error())
		}

		s.httpClient = httpclient.New(timeout, retries, retryDelay)
	}

	return s.httpClient
}

func (s *serviceProvider) ProductService(_ context.Context) *productservice.ProductService {
	if s.productService == nil {
		cfg, err := config.NewProductServiceConfig()
		if err != nil {
			log.Fatalf("failed to get product service config: %s", err.Error())
		}

		s.productService = productservice.New(
			s.HTTPClient(context.Background()),
			cfg.Address,
			cfg.Token,
			cfg.GetProductRPSLimit,
		)
	}
	return s.productService
}

func (s *serviceProvider) Validator(_ context.Context) *validator.Validate {
	if s.validator == nil {
		s.validator = validator.New(validator.WithRequiredStructEnabled())
	}

	return s.validator
}

func (s *serviceProvider) CartRepository(_ context.Context) *cartRepository.Repository {
	if s.cartRepository == nil {
		s.cartRepository = cartRepository.NewRepository()
	}

	return s.cartRepository
}

func (s *serviceProvider) CartRepositoryWithMetrics(ctx context.Context) *cartRepository.RepositoryWithMetrics {
	if s.cartRepositoryWithMetrics == nil {
		s.cartRepositoryWithMetrics = cartRepository.NewRepositoryWithMetrics(s.CartRepository(ctx))
	}

	return s.cartRepositoryWithMetrics
}

func (s *serviceProvider) CartService(ctx context.Context) *cartService.Service {
	if s.cartService == nil {
		s.cartService = cartService.NewCartService(
			s.ProductService(ctx),
			s.CartRepositoryWithMetrics(ctx),
			s.LomsService(ctx),
		)
	}

	return s.cartService
}

func (s *serviceProvider) CartHttpApi(ctx context.Context) *cart.HttpApi {
	if s.cartHttpApi == nil {
		s.cartHttpApi = cart.NewCartHttpApi(
			s.CartService(ctx),
			s.Validator(ctx),
		)
	}

	return s.cartHttpApi
}

func (s *serviceProvider) LomsService(_ context.Context) *lomsservice.Client {
	if s.lomsService == nil {
		client, err := lomsservice.NewClient()
		if err != nil {
			log.Fatalf("failed to get loms service client: %s", err.Error())
		}
		s.lomsService = client
	}

	return s.lomsService
}
