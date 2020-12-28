package config

import (
	"github.com/caarlos0/env/v6"
	"gitlab.com/scalablespace/listener/app/models"
)

func newEnvironment() (models.Environment, error) {
	var e models.Environment
	return e, env.Parse(&e)
}
