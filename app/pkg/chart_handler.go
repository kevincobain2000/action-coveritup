package pkg

import (
	"net/http"
	"strings"

	"github.com/kevincobain2000/action-coveritup/models"
	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type ChartHandler struct {
	Chart *Chart
}

func NewChartHandler() *ChartHandler {
	return &ChartHandler{
		Chart: NewChart(),
	}
}

type ChartRequest struct {
	Org        string `json:"org"  query:"org" validate:"required,ascii,excludes=/" message:"org is required"`
	Repo       string `json:"repo" query:"repo" validate:"required,ascii,excludes=/" message:"repo is required"`
	Branch     string `json:"branch" query:"branch" validate:"ascii" message:"ascii branch is required"`
	User       string `json:"user" query:"user" validate:"ascii,excludes=/" message:"ascii user is required"`
	BaseBranch string `json:"base_branch" query:"base_branch" validate:"ascii" message:"ascii base_branch is required"`
	Type       string `json:"type" query:"type" validate:"ascii" message:"ascii type is required"`
	Types      string `json:"types" query:"types" validate:"ascii" message:"ascii types are required"`
	Branches   string `json:"branches" query:"branches" validate:"ascii" message:"ascii branches is required"`
	Users      string `json:"users" query:"users" validate:"ascii" message:"ascii users is required"`
	PRNum      int    `json:"pr_num" query:"pr_num"`
	PRHistory  string `json:"pr_history" query:"pr_history" default:"commits" validate:"oneof=commits users" message:"pr_history must be commits or users"`
	LastCommit string `json:"last_commit" query:"last_commit" validate:"ascii" message:"ascii last_commit is required"`

	Output string `json:"output" query:"output" default:"png" validate:"oneof=svg png" message:"output must be svg or png"`
	Theme  string `json:"theme" query:"theme" default:"light" validate:"oneof=light dark" message:"theme must be light or dark"`
	Width  int    `json:"width" query:"width" default:"1024" validate:"min=1,max=2048" message:"width must be between 1 and 2048"`
	Height int    `json:"height" query:"height" default:"512" validate:"min=1,max=2048" message:"height must be between 1 and 2048"`
	Line   string `json:"line" query:"line" default:"nofill" validate:"oneof=nofill fill" message:"line must be fill or nofill"`
	Grid   string `json:"grid" query:"grid" default:"show" validate:"oneof=show hide" message:"grid must be show or hide"`
}

func (h *ChartHandler) Get(c echo.Context) error {
	req := new(ChartRequest)
	if err := BindRequest(c, req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	defaults.SetDefaults(req)
	TrimStringFields(req)

	msgs, err := ValidateRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msgs)
	}
	if req.Type == "" && req.Types == "" {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "type or types are required")
	}

	typesMany := strings.Split(req.Types, ",")

	if len(typesMany) > 1 {
		if req.Branch == "" {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "branch for types is required")
		}
		types := []*models.Type{}
		for _, t := range typesMany {
			tt, err := h.Chart.GetType(t)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
			}
			if tt.ID == 0 {
				return echo.NewHTTPError(http.StatusNotFound, "type from types not found")
			}
			types = append(types, tt)
		}
		for i := 1; i < len(types); i++ {
			if types[i].Metric != types[0].Metric {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, "types must have same metric")
			}
		}
		buf, err := h.Chart.GetInstaChartForTypes(req, req.Branch, types)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseMedia(c, buf, req.Output)
	}

	t, err := h.Chart.GetType(req.Type)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if t.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "type not found")
	}

	if req.Branches == "" && req.Users == "" && req.Branch == "" && req.User == "" && req.PRNum == 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "branches, users, pr, branch or user, branch or base_branch are required")
	}
	if req.PRNum > 0 {
		if req.PRHistory == "commits" {
			buf, err := h.Chart.GetInstaChartForPRCommits(req, t)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			return ResponseMedia(c, buf, req.Output)
		}
		if req.PRHistory == "users" {
			buf, err := h.Chart.GetInstaChartForPRUsers(req, t)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			return ResponseMedia(c, buf, req.Output)
		}
	}

	if req.Branch != "" && req.BaseBranch == "" {
		buf, err := h.Chart.GetInstaChartForBranch(req, t)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseMedia(c, buf, req.Output)
	}
	if req.Branch != "" && req.BaseBranch != "" {
		req.Branches = req.BaseBranch + "," + req.Branch
		buf, err := h.Chart.GetInstaChartForBranches(req, t)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseMedia(c, buf, req.Output)
	}
	if req.User != "" {
		buf, err := h.Chart.GetInstaChartForUser(req, t)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseMedia(c, buf, req.Output)
	}
	if req.Branches != "" {
		buf, err := h.Chart.GetInstaChartForBranches(req, t)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseMedia(c, buf, req.Output)
	}
	if req.Users != "" {
		buf, err := h.Chart.GetInstaChartForUsers(req, t)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ResponseMedia(c, buf, req.Output)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Client went wrong")
}
