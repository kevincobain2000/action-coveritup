package pkg

import (
	"github.com/kevincobain2000/action-coveritup/models"
)

type Destroy struct {
	coverage *models.Coverage
}

func NewDestroy() *Destroy {
	return &Destroy{}
}

func (c *Destroy) Delete(req DestroyRequest) error {
	if req.Type != "" {
		return c.coverage.DeleteCoveragesByType(req.Org, req.Repo, req.Type)
	}

	return c.coverage.DeleteCoverages(req.Org, req.Repo)
}
