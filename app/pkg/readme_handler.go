package pkg

import (
	"net/http"

	"github.com/kevincobain2000/action-coveritup/models"
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

type ReadmeResponse struct {
	Types    []models.Type `json:"types"`
	Branches []string      `json:"branches"`
}

type ReadmeRequest struct {
	Org    string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo   string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch string `json:"branch" query:"branch" validate:"ascii" message:"ascii branch is required"`
	User   string `json:"user" query:"user" validate:"ascii,excludes=/" message:"ascii user is required"`
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

	types, err := h.Readme.GetTypes(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	branches, err := h.Readme.GetBranches(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := &ReadmeResponse{
		Types:    types,
		Branches: branches,
	}

	return c.JSON(http.StatusOK, res)
}
