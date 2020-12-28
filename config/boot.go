package config

import (
	"context"
	"database/sql"
	"github.com/labstack/echo"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"log"
)

func boot(lc fx.Lifecycle, e *echo.Echo, db *sql.DB, env models.Environment, consumer components.Consumer, cs components.ConsumerStopper) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go consumer.Start()
			go func() {
				if err := e.Start(env.ServerAddress); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cs.Stop()
			return multierr.Combine(e.Shutdown(ctx), db.Close())
		},
	})
}
