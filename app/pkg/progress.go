package pkg

import (
	"github.com/kevincobain2000/action-coveritup/models"
	gps "github.com/kevincobain2000/go-progress-svg"
)

type Progress struct {
	coverageModel *models.Coverage
	typeModel     *models.Type
}

func NewProgress() *Progress {
	return &Progress{}
}

func (b *Progress) GetType(name string) (*models.Type, error) {
	return b.typeModel.Get(name)
}

func (b *Progress) Get(req *ProgressRequest, t *models.Type) ([]byte, error) {
	ret, err := b.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.Branch, t.Name)

	if err != nil {
		return nil, err
	}

	var circleColor, progressColor, textColor, backgroundColor, captionColor string

	switch req.Theme {
	case "dark":
		circleColor = "#333333"
		textColor = "#FFFFFF"
		backgroundColor = "#000000"
		captionColor = "gray"
	case "light":
		circleColor = "#DDDDDD"
		textColor = "#36454F"
		backgroundColor = "#FFFFFF"
		captionColor = "gray"
	}

	switch {
	case ret.Score > 70:
		progressColor = "#77DD77" // pastel green
	case ret.Score >= 30:
		progressColor = "#FFB347" // pastel yellow
	default:
		progressColor = "#FF6961" // pastel red
	}

	circular, err := gps.NewCircular(func(o *gps.CircularOptions) error {
		o.Progress = int(ret.Score)
		o.Size = 120
		o.CircleWidth = 15
		o.ProgressWidth = 15
		o.CircleColor = circleColor
		o.ProgressColor = progressColor
		o.TextColor = textColor
		o.TextSize = 52
		o.ShowPercentage = true
		o.BackgroundColor = backgroundColor
		o.Caption = t.Name
		o.CaptionPos = "bottom"
		o.CaptionSize = 30
		o.CaptionColor = captionColor
		return nil
	})

	return []byte(circular.SVG()), err
}
