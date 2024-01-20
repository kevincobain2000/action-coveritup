package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type BadgeHandler struct {
	Badge *Badge
}

func NewBadgeHandler() *BadgeHandler {
	return &BadgeHandler{
		Badge: NewBadge(),
	}
}

type BadgeRequest struct {
	Org    string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo   string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch string `json:"branch" query:"branch" validate:"required,ascii,excludes=/" message:"ascii branch is required"`
	Type   string `json:"type" query:"type" validate:"ascii,required,excludes=/" message:"ascii type is required"`
	Style  string `json:"style" query:"style" default:"badge" validate:"oneof=badge chart table" message:"style must be badge or chart"`
	Color  string `json:"color" query:"color" default:"blue" validate:"oneof=blue red green orange" message:"color must be blue, red, green or orange"`
}

func (h *BadgeHandler) Get(c echo.Context) error {
	req := new(BadgeRequest)
	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)

	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}

	t, err := h.Badge.GetType(req.Type)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if t.Name == "" {
		res, err := h.Badge.Get404(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseSVG(c, res)
	}

	res, err := h.Badge.Get(req, t)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ResponseSVG(c, res)
}
