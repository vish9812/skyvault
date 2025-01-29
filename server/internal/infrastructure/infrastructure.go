package infrastructure

import (
	"context"
	"fmt"
	"skyvault/internal/infrastructure/internal/authinfra"
	"skyvault/internal/infrastructure/internal/repository"
	"skyvault/internal/infrastructure/internal/storage"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"time"
)

type Infrastructure struct {
	app        *appconfig.App
	Repository *repository.Repository
	Storage    *storage.Storage
	Auth       *authinfra.AuthInfra
}

// NewInfrastructure initializes all infrastructure components
func NewInfrastructure(app *appconfig.App) *Infrastructure {
	instance := &Infrastructure{
		app: app,
	}

	instance.Repository = repository.NewRepository(app)
	instance.Storage = storage.NewStorage(app)
	instance.Auth = authinfra.NewAuthInfra(app)

	return instance
}

// Cleanup performs cleanup of all infrastructure components
func (i *Infrastructure) Cleanup(ctx context.Context) error {
	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var finalErr error

	// Run cleanup in parallel
	errChan := make(chan error, 2)
	go func() {
		if err := i.Repository.Cleanup(); err != nil {
			errChan <- apperror.NewAppError(err, "i.Cleanup:repository.Cleanup")
			return
		}
		errChan <- nil
	}()

	go func() {
		if err := i.Storage.Cleanup(); err != nil {
			errChan <- apperror.NewAppError(err, "i.Cleanup:storage.Cleanup")
			return
		}
		errChan <- nil
	}()

	// Collect errors
	for i := 0; i < len(errChan); i++ {
		select {
		case err := <-errChan:
			if err != nil {
				if finalErr == nil {
					finalErr = err
				} else {
					finalErr = fmt.Errorf("%w: %w", finalErr, err)
				}
			}
		case <-ctx.Done():
			return apperror.NewAppError(ctx.Err(), "i.Cleanup:Timeout")
		}
	}

	return finalErr
}

// Health checks the health of all infrastructure components
func (i *Infrastructure) Health(ctx context.Context) error {
	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var finalErr error

	// Run health checks in parallel
	errChan := make(chan error, 2)
	go func() {
		if err := i.Repository.Health(ctx); err != nil {
			errChan <- apperror.NewAppError(err, "Infrastructure.Health:repository.Health")
			return
		}
		errChan <- nil
	}()

	go func() {
		if err := i.Storage.Health(ctx); err != nil {
			errChan <- apperror.NewAppError(err, "Infrastructure.Health:storage.Health")
			return
		}
		errChan <- nil
	}()

	// Collect errors
	for i := 0; i < len(errChan); i++ {
		select {
		case err := <-errChan:
			if err != nil {
				if finalErr == nil {
					finalErr = err
				} else {
					finalErr = fmt.Errorf("%w: %w", finalErr, err)
				}
			}
		case <-ctx.Done():
			return apperror.NewAppError(ctx.Err(), "Infrastructure.Health:Timeout")
		}
	}

	return finalErr
}
