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
	Branch     string `json:"branch" query:"branch" validate:"required,ascii" message:"branch is required"`
	BaseBranch string `json:"base_branch" query:"base_branch" validate:"required,ascii" message:"base_branch is required"`
	PRNum      int    `json:"pr_num" query:"pr_num" validate:"required,numeric" message:"pr_num is required"`
	Theme      string `json:"theme" query:"theme" default:"light" validate:"oneof=light dark" message:"theme must be light or dark"`
	Types      string `json:"types" query:"types" validate:"ascii" message:"ascii types are required"`
	DiffTypes  string `json:"diff_types" query:"diff_types" validate:"ascii" message:"ascii diff_types are required"`

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

	types, err := h.PR.TypesToReport(req)
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
