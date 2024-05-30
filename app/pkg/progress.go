package pkg

import (
	"fmt"

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
		textColor = "#36454F"
	case "light":
		circleColor = "#DDDDDD"
		textColor = "#36454F"
	}

	switch {
	case ret.Score > 50:
		progressColor = "#77DD77" // pastel green
		captionColor = "#00FF00"
		backgroundColor = "#D1E2C4"
	case ret.Score >= 30:
		progressColor = "#FFB347" // pastel yellow
		captionColor = "#FFA500"
		backgroundColor = "#FFD1A4"
	default:
		progressColor = "#FF6961" // pastel red
		captionColor = "#FF0000"
		backgroundColor = "#FFD1A4"
	}
	if req.Style == "bar" {
		bar, err := gps.NewBattery(func(o *gps.BatteryOptions) error {
			o.Progress = int(ret.Score)
			o.Width = 180
			o.Height = 40
			o.ProgressColor = progressColor
			o.TextColor = textColor
			o.TextSize = 14
			o.ProgressCaption = fmt.Sprintf("%d%%", int(ret.Score))
			o.Caption = t.Name
			o.CaptionSize = 14
			o.CaptionColor = captionColor
			o.BackgroundColor = backgroundColor
			o.CornerRadius = 10
			return nil
		})
		return []byte(bar.SVG()), err
	}

	circle, err := gps.NewCircular(func(o *gps.CircularOptions) error {
		o.Progress = int(ret.Score)
		o.Size = 100
		o.CircleWidth = 15
		o.ProgressWidth = 15
		o.CircleColor = circleColor
		o.ProgressColor = progressColor
		o.TextColor = textColor
		o.TextSize = 52
		o.ShowPercentage = true
		o.BackgroundColor = backgroundColor
		o.Caption = t.Name
		o.CaptionSize = 30
		o.CaptionColor = captionColor
		o.SegmentGap = 0
		return nil
	})

	return []byte(circle.SVG()), err
}
