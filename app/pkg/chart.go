package pkg

import (
	"strings"

	"github.com/kevincobain2000/action-coveritup/models"
	instachart "github.com/kevincobain2000/instachart/pkg"
	"github.com/sirupsen/logrus"
)

const (
	STYLE_CHART = "chart"
	STYLE_TABLE = "table"
)

type Chart struct {
	coverageModel *models.Coverage
	typeModel     *models.Type
	log           *logrus.Logger
}

func NewChart() *Chart {
	return &Chart{
		log: Logger(),
	}
}

func (e *Chart) GetType(name string) (models.Type, error) {
	return e.typeModel.Get(name)
}

func (e *Chart) GetInstaChartForBranch(req *ChartRequest, t models.Type) ([]byte, error) {
	cReq := e.makeChartRequest(req, t)
	line := instachart.NewLineChart()
	xData := []string{}
	yData := []float64{}
	yyData := [][]float64{}

	names := []string{t.Name}
	ret, err := e.coverageModel.GetLatestBranchScores(req.Org, req.Repo, req.Branch, t.Name)
	if err != nil {
		return nil, err
	}

	for i := len(ret) - 1; i >= 0; i-- {
		r := ret[i]
		xData = append(xData, r.CreatedAt)
		yData = append(yData, r.Score)
		yyData = append(yyData, yData)
	}

	return line.Get(xData, yyData, names, cReq)
}
func (e *Chart) GetInstaChartForUser(req *ChartRequest, t models.Type) ([]byte, error) {
	cReq := e.makeChartRequest(req, t)
	line := instachart.NewLineChart()
	xData := []string{}
	yData := []float64{}
	yyData := [][]float64{}

	names := []string{t.Name}
	ret, err := e.coverageModel.GetLatestUserScores(req.Org, req.Repo, req.User, t.Name)
	if err != nil {
		return nil, err
	}

	for i := len(ret) - 1; i >= 0; i-- {
		r := ret[i]
		xData = append(xData, r.CreatedAt)
		yData = append(yData, r.Score)
		yyData = append(yyData, yData)
	}

	return line.Get(xData, yyData, names, cReq)
}
func (e *Chart) GetInstaChartForBranches(req *ChartRequest, t models.Type) ([]byte, error) {
	cReq := e.makeChartRequest(req, t)
	bar := instachart.NewBarChart()

	if req.Branches == "all" {
		bs, err := e.coverageModel.GetAllBranches(req.Org, req.Repo, t.Name)
		if err != nil {
			return nil, err
		}
		req.Branches = strings.Join(bs, ",")
	}

	xData := []string{}
	yData := []float64{}
	zData := []float64{}
	names := []string{t.Name}
	branches := strings.Split(req.Branches, ",")
	var hasErr error
	for _, branch := range branches {
		ret, err := e.coverageModel.GetLatestBranchScore(req.Org, req.Repo, branch, t.Name)
		if err != nil {
			hasErr = err
			break
		}
		xData = append(xData, ret.BranchName)
		yData = append(yData, ret.Score)
		zData = append(zData, ret.Score)
	}
	if hasErr != nil {
		return nil, hasErr
	}
	yyData := [][]float64{yData}
	zzData := [][]float64{zData}

	return bar.GetStacked(xData, yyData, zzData, names, cReq)
}
func (e *Chart) GetInstaChartForUsers(req *ChartRequest, t models.Type) ([]byte, error) {
	cReq := e.makeChartRequest(req, t)
	bar := instachart.NewBarChart()

	if req.Users == "all" {
		us, err := e.coverageModel.GetAllUsers(req.Org, req.Repo, t.Name)
		if err != nil {
			return nil, err
		}
		req.Users = strings.Join(us, ",")
	}

	xData := []string{}
	yData := []float64{}
	zData := []float64{}
	names := []string{t.Name}
	users := strings.Split(req.Users, ",")
	var hasErr error
	for _, user := range users {
		ret, err := e.coverageModel.GetLatestUserScore(req.Org, req.Repo, user, t.Name)
		if err != nil {
			hasErr = err
			break
		}
		xData = append(xData, ret.UserName)
		yData = append(yData, ret.Score)
		zData = append(zData, ret.Score)
	}
	if hasErr != nil {
		return nil, hasErr
	}
	yyData := [][]float64{yData}
	zzData := [][]float64{zData}

	return bar.GetStacked(xData, yyData, zzData, names, cReq)
}

func (e *Chart) makeChartRequest(req *ChartRequest, t models.Type) *instachart.ChartRequest {
	cReq := &instachart.ChartRequest{
		Output:        req.Output,
		Metric:        t.Metric,
		ChartTitle:    req.Org + "/" + req.Repo,
		ChartSubtitle: req.Branch,
		Theme:         req.Theme,
		Width:         req.Width,
		Height:        req.Height,
		Line:          req.Line,
	}
	return cReq
}
