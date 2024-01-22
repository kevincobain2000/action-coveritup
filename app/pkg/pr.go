package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	md "github.com/go-spectest/markdown"
	"github.com/kevincobain2000/action-coveritup/models"
	"github.com/sirupsen/logrus"
)

const ()

type PR struct {
	coverageModel *models.Coverage
	typeModel     *models.Type
	log           *logrus.Logger
}

func NewPR() *PR {
	return &PR{
		log: Logger(),
	}
}

func (p *PR) Get(req *PRRequest, types []models.Type) (string, error) {
	mdText := md.NewMarkdown(os.Stdout)
	mdTable := md.TableSet{
		Header: []string{"Type", req.BaseBranch, req.Branch},
		Rows:   [][]string{},
	}
	urls := []string{}
	mdText.H4("CoverItUp Report")
	mdText.PlainText("")

	for _, t := range types {
		y := make([]float64, 2)
		row := make([]string, 3)
		row[0] = t.Name

		sb, err := p.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.BaseBranch, t.Name)
		if err != nil {
			y[0] = 0
			row[1] = ""
		} else {
			y[0] = sb.Score
			row[1] = F64NumberToK(&sb.Score) + " " + t.Metric
		}

		s, err := p.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.Branch, t.Name)
		if err != nil {
			y[1] = 0
			row[2] = ""
		} else {
			y[1] = s.Score
			row[2] = F64NumberToK(&s.Score) + " " + t.Metric + p.UpOrDown(&sb.Score, &s.Score)
		}
		mdTable.Rows = append(mdTable.Rows, row)

		data := struct {
			X []string    `json:"x"`
			Y [][]float64 `json:"y"`
			Z [][]float64 `json:"z"`
			N []string    `json:"names"`
		}{
			X: []string{req.BaseBranch, req.Branch},
			Y: [][]float64{},
			Z: [][]float64{},
			N: []string{t.Name},
		}
		data.Y = append(data.Y, y)
		data.Z = append(data.Z, y)

		u := fmt.Sprintf("%s://%s/bar?title=%s&metric=%s&width=%s&height=%s&output=%s&theme=%s",
			req.scheme, req.host, req.Org+"/"+req.Repo, t.Metric, "385", "320", "svg", "dark")

		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		u = u + "&data=" + string(jsonData)
		urls = append(urls, u)
	}
	mdText.Table(mdTable)
	for _, u := range urls {
		mdText.PlainTextf(md.Image("chart", u))
	}

	mdText.PlainText("")
	readmeLink := fmt.Sprintf("%s://%s/readme?org=%s&repo=%s&branch=%s",
		req.scheme, req.host, req.Org, req.Repo, req.Branch)
	mdText.PlainTextf(md.Link("Add to Readme", readmeLink))

	return mdText.String(), nil
}

func (p *PR) UpOrDown(baseScore *float64, branchScore *float64) string {
	if *baseScore > *branchScore {
		return "+"
	}
	if *baseScore < *branchScore {
		return "-"
	}
	return ""
}

func (p *PR) TypesChangedSince(req *PRRequest) ([]models.Type, error) {
	typesChanged := []models.Type{}
	types, err := p.typeModel.GetTypesFor(req.Org, req.Repo)
	if err != nil {
		return typesChanged, err
	}
	for _, t := range types {
		sbs, err := p.coverageModel.GetLatestBranchScoresWithPR(req.Org, req.Repo, req.Branch, t.Name, 2)
		if err != nil {
			p.log.Error(err)
			typesChanged = append(typesChanged, t)
			continue
		}
		if len(sbs) < 2 {
			typesChanged = append(typesChanged, t)
			continue
		}

		if float64(sbs[0].PRNum) != float64(sbs[1].PRNum) {
			typesChanged = append(typesChanged, t)
			continue
		}
		if sbs[0].Score != sbs[1].Score {
			typesChanged = append(typesChanged, t)
			continue
		}
	}
	return typesChanged, nil
}
