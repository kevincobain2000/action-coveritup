package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	url string
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{
		url: DOCS_URL,
	}
}

func (h *IndexHandler) Get(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, h.url)
}
