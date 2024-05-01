package pkg

import (
	"github.com/kevincobain2000/action-coveritup/models"
	instachart "github.com/kevincobain2000/instachart/pkg"
)

type Table struct {
	typeModel     *models.Type
	coverageModel *models.Coverage
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) GetInstaTableForBranch(req *TableRequest) ([]byte, error) {
	cReq := t.makeChartRequest(req)
	table := instachart.NewTableChart()
	rData := [][]string{}

	types, err := t.typeModel.GetAllTypesFor(req.Org, req.Repo)
	if err != nil {
		return nil, err
	}

	names := []string{
		req.Branch,
		"Score",
	}
	for _, tt := range types {
		ret, err := t.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.Branch, tt.Name)
		if err != nil {
			return nil, err
		}
		score := F64NumberToK(&ret.Score) + " " + ret.Metric
		rData = append(rData, []string{tt.Name, score})
	}

	return table.Get(names, rData, cReq)
}

func (t *Table) makeChartRequest(req *TableRequest) *instachart.ChartRequest {
	cReq := &instachart.ChartRequest{
		Output: req.Output,
		Theme:  req.Theme,
		Width:  req.Width,
	}
	return cReq
}
