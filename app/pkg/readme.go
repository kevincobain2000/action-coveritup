package pkg

import (
	"github.com/kevincobain2000/action-coveritup/models"
)

type Readme struct {
	typeModel     *models.Type
	coverageModel *models.Coverage
}

func NewReadme() *Readme {
	return &Readme{}
}

func (r *Readme) GetTypes(req *ReadmeRequest) ([]models.Type, error) {
	types, err := r.typeModel.GetAllTypesFor(req.Org, req.Repo)
	if err != nil {
		return nil, err
	}
	if len(types) == 0 {
		return []models.Type{}, nil
	}
	return types, nil
}

func (r *Readme) GetBranches(req *ReadmeRequest) ([]string, error) {
	branches, err := r.coverageModel.GetAllBranches(req.Org, req.Repo)
	if err != nil {
		return nil, err
	}
	if len(branches) == 0 {
		return []string{}, nil
	}
	return branches, nil
}
