package app

import (
	"context"
	"log"

	"github.com/katyafirstova/auth_service/internal/api/user"
	"github.com/katyafirstova/auth_service/internal/client/db"
	"github.com/katyafirstova/auth_service/internal/client/db/pg"
	"github.com/katyafirstova/auth_service/internal/closer"
	"github.com/katyafirstova/auth_service/internal/config"
	"github.com/katyafirstova/auth_service/internal/config/env"
	"github.com/katyafirstova/auth_service/internal/repository"
	userRepository "github.com/katyafirstova/auth_service/internal/repository/user"
	"github.com/katyafirstova/auth_service/internal/service"
	userService "github.com/katyafirstova/auth_service/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	dbClient   db.Client

	userRepo    repository.UserRepository
	userService service.UserService
	userImpl    *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to create db client: %s", err.Error())
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
