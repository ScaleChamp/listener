package config

import (
	"github.com/labstack/echo"
	"gitlab.com/scalablespace/listener/lib/components"
)

func routes(hc components.HealthController) *echo.Echo {
	e := echo.New()
	e.GET("/health", hc.Health)
	return e
}
