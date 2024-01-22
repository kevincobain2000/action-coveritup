package pkg

import (
	"github.com/kevincobain2000/action-coveritup/models"
)

type Upload struct {
	coverage  *models.Coverage
	repo      *models.Repo
	org       *models.Org
	user      *models.User
	typeModel *models.Type
}

func NewUpload() *Upload {
	return &Upload{}
}

func (c *Upload) Post(req *UploadRequest) (*models.Coverage, error) {
	t, err := c.typeModel.Get(req.Type)
	if err != nil {
		return nil, err
	}

	if t.ID == 0 {
		t, err = c.typeModel.Create(req.Type, req.Metric)
		if err != nil {
			return nil, err
		}
	}

	o, err := c.org.Get(req.Org)
	if err != nil {
		return nil, err
	}
	if o.ID == 0 {
		o, err = c.org.Create(req.Org)
		if err != nil {
			return nil, err
		}
	}

	r, err := c.repo.Get(o.ID, req.Repo)
	if err != nil {
		return nil, err
	}
	if r.ID == 0 {
		r, err = c.repo.Create(o.ID, req.Repo)
		if err != nil {
			return nil, err
		}
	}

	u, err := c.user.Get(req.User)
	if err != nil {
		return nil, err
	}
	if u.ID == 0 {
		u, err = c.user.Create(req.User)
		if err != nil {
			return nil, err
		}
	}

	req.Branches += " " + req.Branch
	err = c.coverage.SoftDeleteCoverages(o.ID, r.ID, req.Branches)
	if err != nil {
		return nil, err
	}
	nc, err := c.coverage.Create(o.ID,
		r.ID,
		u.ID,
		t.ID,
		req.Branch,
		StringToInt(req.PRNum),
		TakeFirst(req.Commit, 7),
		SToF32(req.Score))
	if err != nil {
		return nil, err
	}

	return nc, nil
}
