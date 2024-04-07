package pkg

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type DestroyHandler struct {
	Destroy *Destroy
	github  *Github
}

func NewDestroyHandler() *DestroyHandler {
	return &DestroyHandler{
		Destroy: NewDestroy(),
		github:  NewGithub(),
	}
}

type DestroyRequest struct {
	Org    string `json:"org"  form:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo   string `json:"repo" form:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Type   string `json:"type" form:"type" validate:"ascii,excludes=/" message:"ascii type is required"`
	Commit string `json:"commit" form:"commit" validate:"ascii,excludes=/" message:"ascii commit is required"`
}

func (h *DestroyHandler) Post(c echo.Context) error {
	req := new(DestroyRequest)
	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)

	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}

	if os.Getenv("GITHUB_API") != "" {
		if err := h.github.VerifyGithubToken(c.Request().Header.Get("Authorization"), req.Org, req.Repo, req.Commit); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
	}

	err = h.Destroy.Delete(*req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	SetHeadersResponseJSON(c.Response().Header())
	return c.JSON(http.StatusOK, `{"status": "destroyed"}`)
}
