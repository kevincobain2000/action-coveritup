package pkg

import (
	"strings"

	"github.com/kevincobain2000/action-coveritup/models"
	"github.com/narqo/go-badge"
)

type Badge struct {
	coverageModel *models.Coverage
	typeModel     *models.Type
}

func NewBadge() *Badge {
	return &Badge{}
}

func (b *Badge) Get(req *BadgeRequest, t *models.Type) ([]byte, error) {
	ret, err := b.coverageModel.GetLatestBranchScore(req.Org, req.Repo, req.Branch, t.Name)

	if err != nil {
		return nil, err
	}
	scoreStr := F64NumberToK(&ret.Score)
	label := t.Name + " | " + req.Branch

	if t.Metric == "" {
		return badge.RenderBytesSocial(label, scoreStr, "#00000", "#00000", "#fff")
	}

	scoreStr += ret.Metric
	badgeColor, labelColor, color := b.getBadgeColors(t.Metric, ret.Score)
	return badge.RenderBytesFlat(label, scoreStr, badgeColor, labelColor, color)
}

func (b *Badge) Get404(req *BadgeRequest) ([]byte, error) {
	return badge.RenderBytesFlat(req.Branch+"|"+req.Type, "404", "#fff", "white", "red")
}

func (b *Badge) GetType(name string) (*models.Type, error) {
	if (name == models.TYPE_AVERAGE_PR_DAYS) || (name == models.TYPE_NUMBER_OF_CONTRIBUTORS) {
		return &models.Type{
			Name: name,
		}, nil
	}
	return b.typeModel.Get(name)
}

func (e *Badge) getBadgeColors(metric string, score float64) (badge.Color, badge.Color, badge.Color) {
	metric = strings.ToLower(metric)
	if metric == "%" {
		var color badge.Color
		var textColor badge.Color
		switch {
		case score >= 90:
			color = "green"
			textColor = "black"
		case score >= 75:
			color = "blue"
			textColor = "white"
		case score >= 50:
			color = "orange"
			textColor = "black"
		case score < 50:
			color = "red"
			textColor = "white"
		default:
			color = "white"
			textColor = "black"
		}
		return "#fff", textColor, color
	}
	if metric == "kb" || metric == "mb" || metric == "gb" {
		return "#fff", "gray", "yellow"
	}
	if metric == "ms" ||
		metric == "s" ||
		metric == "sec" ||
		metric == "min" ||
		metric == "days" ||
		metric == "d" ||
		metric == "h" ||
		metric == "hr" {
		return "#fff", "gray", "#FFFF00"
	}
	if metric == "alloc" {
		return "#fff", "gray", "#00FF00"
	}
	return "#bbb", "#fff", "gray"
}
