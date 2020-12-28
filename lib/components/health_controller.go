package components

import "github.com/labstack/echo"

type HealthController interface {
	Health(echo.Context) error
}
