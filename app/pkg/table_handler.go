package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type TableHandler struct {
	Table *Table
}

func NewTableHandler() *TableHandler {
	return &TableHandler{
		Table: NewTable(),
	}
}

type TableRequest struct {
	Org    string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo   string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch string `json:"branch" query:"branch" validate:"ascii" message:"ascii branch is required"`

	Output string `json:"output" query:"output" default:"png" validate:"oneof=svg png" message:"output must be svg or png"`
	Theme  string `json:"theme" query:"theme" default:"light" validate:"oneof=light dark" message:"theme must be light or dark"`
	Width  int    `json:"width" query:"width" default:"1024" validate:"min=1,max=2048" message:"width must be between 1 and 2048"`
}

func (h *TableHandler) Get(c echo.Context) error {
	req := new(TableRequest)
	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)
	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}

	res, err := h.Table.GetInstaTableForBranch(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ResponsePNG(c, res)
}
