package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RobotsHandler struct {
	text string
}

func NewRobotsHandler() *RobotsHandler {
	return &RobotsHandler{
		text: ROBOTS_TXT,
	}
}

func (h *RobotsHandler) Get(c echo.Context) error {
	SetHeaderResponseText(c.Response().Header())
	return c.String(http.StatusOK, h.text)
}
