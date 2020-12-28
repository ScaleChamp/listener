package controllers

import (
	"database/sql"
	"github.com/labstack/echo"
	"gitlab.com/scalablespace/listener/lib/components"
	"net/http"
)

type healthController struct {
	db *sql.DB
}

func (h *healthController) Health(c echo.Context) error {
	if err := h.db.Ping(); err != nil {
		return echo.ErrInternalServerError
	}
	return echo.NewHTTPError(http.StatusOK)
}

func NewHealthController(db *sql.DB) components.HealthController {
	return &healthController{db}
}
