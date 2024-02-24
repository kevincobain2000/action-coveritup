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

func (e *Chart) GetType(name string) (*models.Type, error) {
	return e.typeModel.Get(name)
}

// func (e *Chart) GetInstaChartForPRCommits(req *ChartRequest, t *models.Type) ([]byte, error) {
// 	ret, err := e.coverageModel.GetLatestPRScoresForCommits(req.Org, req.Repo, req.PRNum, t.Name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(ret) > 0 {
// 		req.Branch = ret[0].BranchName
// 	}

// 	cReq := e.makeChartRequest(req, t)
// 	line := instachart.NewLineChart()
// 	xData := []string{}
// 	yData := []float64{}
// 	yyData := [][]float64{}
// 	names := []string{t.Name}

// 	for i := len(ret) - 1; i >= 0; i-- {
// 		r := ret[i]
// 		xData = append(xData, r.Commit)
// 		yData = append(yData, r.Score)
// 		yyData = append(yyData, yData)
// 	}
// 	if len(xData) == 0 {
// 		xData = append(xData, "0")
// 	}
// 	if len(yyData) == 0 {
// 		yyData = append(yyData, []float64{0})
// 	}

// 	return line.Get(xData, yyData, names, cReq)
// }
// func (e *Chart) GetInstaChartForPRUsers(req *ChartRequest, t *models.Type) ([]byte, error) {
// 	ret, err := e.coverageModel.GetLatestPRScoresForUsers(req.Org, req.Repo, req.PRNum, t.Name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pp.Println(ret)

// 	if len(ret) > 0 {
// 		req.Branch = fmt.Sprintf("%s #%d", ret[0].BranchName, req.PRNum)
// 	}

// 	cReq := e.makeChartRequest(req, t)
// 	bar := instachart.NewBarChart()
// 	xData := []string{}
// 	yData := []float64{}
// 	yyData := [][]float64{}
// 	names := []string{t.Name}

// 	for i := len(ret) - 1; i >= 0; i-- {
// 		r := ret[i]
// 		xData = append(xData, r.UserName)
// 		yData = append(yData, r.Score)
// 		yyData = append(yyData, yData)
// 	}
// 	if len(xData) == 0 {
// 		xData = append(xData, "0")
// 	}
// 	if len(yyData) == 0 {
// 		yyData = append(yyData, []float64{0})
// 	}

// 	return bar.GetVertical(xData, yyData, names, cReq)
// }

// Line chart with dates on x-axis and scores on y-axis
func (e *Chart) GetInstaChartForBranch(req *ChartRequest, t *models.Type) ([]byte, error) {
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
	if len(xData) == 0 {
		xData = append(xData, "0")
	}
	if len(yyData) == 0 {
		yyData = append(yyData, []float64{0})
	}

	return line.Get(xData, yyData, names, cReq)
}

// Line chart with dates on x-axis and scores on y-axis
func (e *Chart) GetInstaChartForUser(req *ChartRequest, t *models.Type) ([]byte, error) {
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
	if len(xData) == 0 {
		xData = append(xData, "0")
	}
	if len(yyData) == 0 {
		yyData = append(yyData, []float64{0})
	}

	return line.Get(xData, yyData, names, cReq)
}

// bar chart with branch names on x-axis and scores on y-axis
func (e *Chart) GetInstaChartForBranches(req *ChartRequest, t *models.Type) ([]byte, error) {
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

// bar chart with user names on x-axis and scores on y-axis
func (e *Chart) GetInstaChartForUsers(req *ChartRequest, t *models.Type) ([]byte, error) {
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

func (e *Chart) isMini(req *ChartRequest) bool {
	if req.Width <= 200 && req.Height <= 200 {
		return true
	}
	return false
}
func (e *Chart) makeChartRequest(req *ChartRequest, t *models.Type) *instachart.ChartRequest {
	title := req.Org + "/" + req.Repo
	subtitle := req.Branch
	if e.isMini(req) {
		title = t.Name
		subtitle = req.Branch
		req.Grid = "hide"
	}
	cReq := &instachart.ChartRequest{
		Output:        req.Output,
		Metric:        t.Metric,
		ChartTitle:    title,
		ChartSubtitle: subtitle,
		Theme:         req.Theme,
		Width:         req.Width,
		Height:        req.Height,
		Line:          req.Line,
		Grid:          req.Grid,
	}
	return cReq
}
