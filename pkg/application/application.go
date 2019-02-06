package application

import (
	"github.com/emilyzhang/goyagi/pkg/config"
	"github.com/emilyzhang/goyagi/pkg/database"
	"github.com/emilyzhang/goyagi/pkg/metrics"
	"github.com/emilyzhang/goyagi/pkg/sentry"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

// App contains necessary references that will be persisted throughout the
// application's lifecycle.
type App struct {
	Config  config.Config
	DB      *pg.DB
	Sentry  sentry.Sentry
	Metrics metrics.Metrics
}

// New creates a new instance of App with a Config and Metrics.
func New() (App, error) {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	sentry, err := sentry.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	m, err := metrics.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	return App{cfg, db, sentry, m}, nil
}
