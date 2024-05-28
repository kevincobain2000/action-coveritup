package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type ProgressHandler struct {
	Progress *Progress
}

func NewProgressHandler() *ProgressHandler {
	return &ProgressHandler{
		Progress: NewProgress(),
	}
}

type ProgressRequest struct {
	Org    string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo   string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch string `json:"branch" query:"branch" validate:"required,ascii" message:"ascii branch is required"`
	Type   string `json:"type" query:"type" validate:"ascii,required,excludes=/" message:"ascii type is required"`
	Theme  string `json:"theme" query:"theme" default:"light" validate:"oneof=light dark" message:"theme must be light or dark"`
}

func (h *ProgressHandler) Get(c echo.Context) error {
	req := new(ProgressRequest)
	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)

	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}

	t, err := h.Progress.GetType(req.Type)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if t.Metric != "%" {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "metric must be %")
	}

	res, err := h.Progress.Get(req, t)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ResponseSVG(c, res)
}
