package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type ReadmeHandler struct {
	Readme *Readme
}

func NewReadmeHandler() *ReadmeHandler {
	return &ReadmeHandler{
		Readme: NewReadme(),
	}
}

type ReadmeRequest struct {
	Org    string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo   string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch string `json:"branch" query:"branch" validate:"ascii,excludes=/" message:"ascii branch is required"`
	User   string `json:"user" query:"user" validate:"ascii,excludes=/" message:"ascii user is required"`

	host   string
	scheme string
}

func (h *ReadmeHandler) Get(c echo.Context) error {
	req := new(ReadmeRequest)
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
	types, err := h.Readme.GetTypes(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(types) == 0 {
		return c.String(http.StatusNotFound, "no types found for this org/repo")
	}
	str, err := h.Readme.Get(req, types)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, str)

}
