// pkg/application/application.go
package application

import (
	"github.com/emilyzhang/goyagi/pkg/config"
	"github.com/emilyzhang/goyagi/pkg/database"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

// App contains necessary references that will be persisted throughout the
// application's lifecycle.
type App struct {
	Config config.Config
	DB     *pg.DB
}

// New creates a new instance of App with a Config and Metrics.
func New() (App, error) {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	return App{cfg, db}, nil
}
