package infrastructure

import (
	"context"
	"skyvault/internal/infrastructure/internal/store_db"
	"skyvault/internal/infrastructure/internal/store_file"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

type Infrastructure struct {
	DBStore   *store_db.Store
	FileStore *store_file.Store
	Repo      *store_db.Repo
	app       *appconfig.App
}

// NewInfrastructure initializes all infrastructure components
func NewInfrastructure(ctx context.Context, app *appconfig.App) *Infrastructure {
	instance := &Infrastructure{
		app: app,
	}

	// Initialize database
	instance.DBStore = store_db.NewStore(app, app.Config.DB.DSN)
	instance.Repo = store_db.NewRepo(instance.DBStore)

	// Initialize file storage
	instance.FileStore = store_file.NewStore(app)

	return instance
}

// Cleanup performs cleanup of all infrastructure components
func (i *Infrastructure) Cleanup(ctx context.Context) error {
	var err error

	// Cleanup database connection
	if err := i.DBStore.Cleanup(); err != nil {
		err = apperror.NewAppError(err, "i.Cleanup:DBStore.Cleanup")
	}

	// Cleanup file storage
	if err := i.FileStore.Cleanup(); err != nil {
		err = apperror.NewAppError(err, "i.Cleanup:FileStore.Cleanup")
	}

	return err
}

// Health checks the health of all infrastructure components
func (i *Infrastructure) Health(ctx context.Context) error {
	var err error

	if err := i.DBStore.Health(ctx); err != nil {
		err = apperror.NewAppError(err, "i.Health:DBStore.Health")
	}

	if err := i.FileStore.Health(ctx); err != nil {
		err = apperror.NewAppError(err, "i.Health:FileStore.Health")
	}

	return err
}
