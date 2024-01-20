package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type UploadHandler struct {
	upload *Upload
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		upload: NewUpload(),
	}
}

type UploadRequest struct {
	Org      string `json:"org"  query:"org" validate:"required,ascii,excludes=/,max=255" message:"org is required"`
	Repo     string `json:"repo" query:"repo" validate:"required,ascii,excludes=/,max=255" message:"repo is required"`
	User     string `json:"user" query:"user" validate:"required,ascii,excludes=/,max=255" message:"user is required"`
	Type     string `json:"type" query:"type" validate:"required,ascii,required,excludes=/,max=32" message:"ascii type is required"`
	Metric   string `json:"metric" query:"metric" validate:"ascii,max=3" message:"ascii metric is required"`
	Branch   string `json:"branch" query:"branch" validate:"required,ascii,max=255" message:"ascii branch is required"`
	Commit   string `json:"commit" query:"commit" validate:"ascii,excludes=/,max=255" message:"ascii commit is required"`
	Score    string `json:"score" query:"score" validate:"ascii,excludes=/,max=12" message:"ascii score is required"`
	Branches string `json:"branches" query:"branches" validate:"ascii" message:"ascii branches is required"`
	PRNum    string `json:"pr_num" query:"pr_num" validate:"ascii,excludes=/,max=4" message:"ascii pr_num is required"`
}

func (h *UploadHandler) Post(c echo.Context) error {

	req := new(UploadRequest)
	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)
	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}

	// if err := h.github.VerifyGithubToken(c.Request().Header.Get("Authorization"), req.Org, req.Repo, req.Type); err != nil {
	// 	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	// }

	res, err := h.upload.Post(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	SetHeadersResponseJSON(c.Response().Header())
	return c.JSON(http.StatusOK, res)
}
