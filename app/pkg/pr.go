package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/kevincobain2000/action-coveritup/models"
	md "github.com/nao1215/markdown"
	"github.com/sirupsen/logrus"
)

const (
	DOWN_SYMBOL      = "📉"
	UP_SYMBOL        = "📈"
	NO_CHANGE_SYMBOL = ""
)

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
	totalPlus := 0
	totalMinus := 0
	verdict := ""
	commitBranch := ""
	commitBaseBranch := ""
	isFirstPR := p.coverageModel.IsFirstPR(req.Org, req.Repo, req.PRNum)

	mdText := md.NewMarkdown(os.Stdout)
	mdTable := md.TableSet{
		Rows: [][]string{},
	}
	baseAndBranchImgUrls := []string{} // stores urls for bar charts for comparison of base and branch
	commitHistoryImgUrls := []string{} // stores urls for commit history trends (line charts)
	userHistoryImgUrls := []string{}   // stores urls for user history trends (line charts)
	mdText.H4("CoverItUp Report")

	for _, t := range types {
		y := make([]float64, 2)
		tableRow := make([]string, 4)
		tableRow[0] = "*" + t.Name + "*"

		scoreBaseBranch, err := p.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.BaseBranch, t.Name)
		if commitBaseBranch == "" {
			commitBaseBranch = scoreBaseBranch.Commit
		}

		if err != nil {
			y[0] = 0
			tableRow[1] = ""
		} else {
			y[0] = scoreBaseBranch.Score
			tableRow[1] = F64NumberToK(&scoreBaseBranch.Score) + "" + t.Metric
		}

		scoreBranch, err := p.coverageModel.GetLatestBranchScoreByPR(req.Org, req.Repo, req.Branch, t.Name, req.PRNum)
		if commitBranch == "" {
			commitBranch = scoreBranch.Commit
		}
		if err != nil {
			y[1] = 0
			tableRow[2] = ""
			tableRow[3] = ""
		} else {
			symbol := p.UpOrDown(&scoreBaseBranch.Score, &scoreBranch.Score)
			y[1] = scoreBranch.Score
			tableRow[2] = F64NumberToK(&scoreBranch.Score) + "" + t.Metric
			tableRow[3] = symbol
			if symbol == UP_SYMBOL || symbol == DOWN_SYMBOL {
				tableRow[2] = "**" + tableRow[2] + "**"
			}
			if symbol == UP_SYMBOL {
				totalPlus++
			}
			if symbol == DOWN_SYMBOL {
				totalMinus++
			}
		}
		mdTable.Rows = append(mdTable.Rows, tableRow)

		bbData := struct {
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
		bbData.Y = append(bbData.Y, y)
		bbData.Z = append(bbData.Z, y)

		// this is a static bar chart, that doesn't change once commented
		bbUrl := fmt.Sprintf("%s://%s%sbar?title=%s&metric=%s&width=%d&height=%d&grid=hide&output=%s&theme=%s&grid=%s",
			req.scheme,
			req.host,
			os.Getenv("BASE_URL"),
			req.Org+"/"+req.Repo,
			t.Metric,
			385,
			320,
			"svg",
			req.Theme,
			"hide")

		jsonData, err := json.Marshal(bbData)
		if err != nil {
			return "", err
		}
		bbUrl = bbUrl + "&data=" + string(jsonData)
		baseAndBranchImgUrls = append(baseAndBranchImgUrls, bbUrl)

		// this is dynamic bar chart, that changes upon new commits, so previous commits comments also change dynamically
		chUrl := fmt.Sprintf("%s://%s%schart?org=%s&repo=%s&pr_num=%d&type=%s&theme=%s&height=%d&line=%s&pr_history=%s&last_commit=%s",
			req.scheme,
			req.host,
			os.Getenv("BASE_URL"),
			req.Org,
			req.Repo,
			req.PRNum,
			t.Name,
			req.Theme,
			320,
			"fill",
			"commits",
			commitBranch)
		commitHistoryImgUrls = append(commitHistoryImgUrls, chUrl)

		uhUrl := fmt.Sprintf("%s://%s%schart?org=%s&repo=%s&pr_num=%d&type=%s&theme=%s&height=%d&line=%s&pr_history=%s&last_commit=%s",
			req.scheme,
			req.host,
			os.Getenv("BASE_URL"),
			req.Org,
			req.Repo,
			req.PRNum,
			t.Name,
			req.Theme,
			320,
			"fill",
			"users",
			commitBranch)
		userHistoryImgUrls = append(userHistoryImgUrls, uhUrl)
	}

	basicTable, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Type",
			fmt.Sprintf("`%s`", req.BaseBranch),
			fmt.Sprintf("`%s`", req.Branch),
			fmt.Sprintf("[%s](/%s/%s/commit/%s) from [%s](/%s/%s/pull/%d/commits/%s)",
				commitBaseBranch,
				req.Org,
				req.Repo,
				commitBaseBranch,
				commitBranch,
				req.Org,
				req.Repo,
				req.PRNum,
				commitBranch)).
		Format(mdTable.Rows)

	if err != nil {
		return "", err
	}

	if totalPlus > totalMinus {
		verdict = UP_SYMBOL
	} else if totalPlus < totalMinus {
		verdict = DOWN_SYMBOL
	} else if totalPlus == totalMinus && totalPlus > 0 {
		verdict = UP_SYMBOL + " and " + DOWN_SYMBOL
	}
	mdText.PlainText("")
	if isFirstPR {
		mdText.PlainText(basicTable)
	} else {
		mdText.Details(fmt.Sprintf("Comparison Table - <b>%d</b> Types %s",
			len(types),
			verdict),
			"\n"+basicTable+"\n")
	}

	images := ""
	for _, u := range baseAndBranchImgUrls {
		images += fmt.Sprintf("<img loading='eager' src='%s' alt='base vs branch' />", u)
	}

	mdText.PlainText("")
	mdText.Details(fmt.Sprintf("Comparisons Chart - <code>%s</code> from <code>%s</code>", req.BaseBranch, req.Branch), "\n"+images+"\n")

	cImages := ""
	for _, u := range commitHistoryImgUrls {
		cImages += fmt.Sprintf("<img loading='eager' src='%s' alt='commit history' />", u)
	}

	mdText.PlainText("")
	uptoCommitsText := fmt.Sprintf("Upto <code>%s</code> for <b>#%d</b>", commitBranch, req.PRNum)
	mdText.Details("Commits History", "\n"+uptoCommitsText+"\n"+cImages+"\n")

	uImages := ""
	for _, u := range userHistoryImgUrls {
		uImages += fmt.Sprintf("<img loading='eager' src='%s' alt='user history' />", u)
	}
	mdText.PlainText("")
	mdText.Details("Users History", "\n"+uptoCommitsText+"\n"+uImages+"\n")

	mdText.PlainText("")
	readmeLink := fmt.Sprintf("%s://%s%sreadme?org=%s&repo=%s&branch=%s",
		req.scheme,
		req.host,
		os.Getenv("BASE_URL"),
		req.Org,
		req.Repo,
		req.Branch)
	mdText.PlainTextf(md.Link("Embed README.md", readmeLink))

	return mdText.String(), nil
}

func (p *PR) UpOrDown(baseScore *float64, branchScore *float64) string {
	if *baseScore > *branchScore {
		return DOWN_SYMBOL
	}
	if *baseScore < *branchScore {
		return UP_SYMBOL
	}
	return NO_CHANGE_SYMBOL
}

func (p *PR) TypesToReport(req *PRRequest) ([]models.Type, error) {
	types, err := p.typeModel.GetBranchTypesFor(req.Org, req.Repo, []string{req.BaseBranch, req.Branch}, req.Types)
	return types, err
}
