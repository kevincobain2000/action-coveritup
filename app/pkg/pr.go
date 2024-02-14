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
	chUrls := []string{}
	mdText.H4("CoverItUp Report")
	mdText.PlainText("")

	for _, t := range types {
		y := make([]float64, 2)
		row := make([]string, 3)
		row[0] = "*" + t.Name + "*"

		sb, err := p.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.BaseBranch, t.Name)
		if err != nil {
			y[0] = 0
			row[1] = ""
		} else {
			y[0] = sb.Score
			row[1] = F64NumberToK(&sb.Score) + "" + t.Metric
		}

		s, err := p.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.Branch, t.Name)
		if err != nil {
			y[1] = 0
			row[2] = ""
		} else {
			ud := p.UpOrDown(&sb.Score, &s.Score)
			y[1] = s.Score
			row[2] = F64NumberToK(&s.Score) + "" + t.Metric + ud
			if ud != "" && ud != "-" {
				row[2] = "**" + row[2] + "**"
			}
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

		u := fmt.Sprintf("%s://%s%sbar?title=%s&metric=%s&width=%s&height=%s&grid=hide&output=%s&theme=%s&grid=%s",
			req.scheme, req.host, os.Getenv("BASE_URL"), req.Org+"/"+req.Repo, t.Metric, "385", "320", "svg", req.Theme, "hide")

		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		u = u + "&data=" + string(jsonData)
		urls = append(urls, u)

		cu := fmt.Sprintf("%s://%s%schart?org=%s&repo=%s&pr_num=%d&type=%s&theme=%s&height=%d&line=%s", req.scheme, req.host, os.Getenv("BASE_URL"), req.Org, req.Repo, req.PRNum, t.Name, req.Theme, 320, "fill")
		chUrls = append(chUrls, cu)
	}
	mdText.Table(mdTable)
	images := ""
	for _, u := range urls {
		images += md.Image("chart", u)
	}
	mdText.PlainText(images)

	cImages := ""
	for _, u := range chUrls {
		cImages += fmt.Sprintf("<img src='%s' alt='commit history' />", u)
	}

	mdText.Details(fmt.Sprintf("Commit history for this PR %d", req.PRNum), cImages)

	mdText.PlainText("")
	readmeLink := fmt.Sprintf("%s://%s%sreadme?org=%s&repo=%s&branch=%s",
		req.scheme, req.host, os.Getenv("BASE_URL"), req.Org, req.Repo, req.Branch)
	mdText.PlainTextf(md.Link("Add Badges and Charts to Readme", readmeLink))

	return mdText.String(), nil
}

func (p *PR) UpOrDown(baseScore *float64, branchScore *float64) string {
	if *baseScore > *branchScore {
		return "-"
	}
	if *baseScore < *branchScore {
		return "+"
	}
	return ""
}

func (p *PR) TypesChangedSince(req *PRRequest) ([]models.Type, error) {
	typesChanged := []models.Type{}
	types, err := p.typeModel.GetBranchTypesFor(req.Org, req.Repo, []string{req.BaseBranch, req.Branch})
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
		if F64NumberToK(&sbs[0].Score) != F64NumberToK(&sbs[1].Score) {
			typesChanged = append(typesChanged, t)
			continue
		}
	}
	return typesChanged, nil
}
