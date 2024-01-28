package pkg

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type PRHandler struct {
	PR *PR
}

func NewPRHandler() *PRHandler {
	return &PRHandler{
		PR: NewPR(),
	}
}

type PRRequest struct {
	Org        string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo       string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch     string `json:"branch" query:"branch" validate:"required,ascii" message:"ascii branch is required"`
	BaseBranch string `json:"base_branch" query:"base_branch" validate:"required,ascii" message:"ascii base_branch is required"`

	host   string
	scheme string
}

func (h *PRHandler) Get(c echo.Context) error {
	req := new(PRRequest)

	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)

	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}

	req.host = c.Request().Host
	req.scheme = c.Scheme()

	types, err := h.PR.TypesChangedSince(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(types) == 0 {
		return echo.NewHTTPError(http.StatusConflict, errors.New("no types changed"))
	}
	prComment, err := h.PR.Get(req, types)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusCreated, prComment)
}
