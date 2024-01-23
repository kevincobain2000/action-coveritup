package pkg

import (
	"fmt"
	"os"

	md "github.com/go-spectest/markdown"
	"github.com/kevincobain2000/action-coveritup/models"
)

type Readme struct {
	typeModel *models.Type
}

func NewReadme() *Readme {
	return &Readme{}
}

func (r *Readme) GetTypes(req *ReadmeRequest) ([]models.Type, error) {
	types, err := r.typeModel.GetTypesFor(req.Org, req.Repo)
	if err != nil {
		return nil, err
	}
	return types, nil
}
func (r *Readme) Get(req *ReadmeRequest, types []models.Type) (string, error) {
	mdText := md.NewMarkdown(os.Stdout)

	mdText.H1("CoverItUp Report").
		PlainText("").
		H2("Badges").
		PlainText("")

	for _, t := range types {
		u := fmt.Sprintf("%s://%s%s/badge?org=%s&repo=%s&type=%s&branch=%s",
			req.scheme, req.host, os.Getenv("BASE_URL"), req.Org, req.Repo, t.Name, req.Branch)
		mdText.PlainTextf(md.Image(t.Name, u))
	}

	mdText.PlainText("").
		H2("Charts").
		PlainText("")
	mdText.H3("Branch").
		PlainText("")
	if req.Branch != "" {

		mdText.H4("Single").
			PlainText("")

		for _, t := range types {
			u := fmt.Sprintf("%s://%s%s/chart?org=%s&repo=%s&type=%s&branch=%s",
				req.scheme, req.host, os.Getenv("BASE_URL"), req.Org, req.Repo, t.Name, req.Branch)
			mdText.PlainTextf(md.Image(t.Name, u))
		}
		mdText.PlainText("")
	}

	mdText.H4("All").PlainText("")
	for _, t := range types {
		u := fmt.Sprintf("%s://%s/chart?org=%s&repo=%s&type=%s&branches=all",
			req.scheme, req.host, req.Org, req.Repo, t.Name)
		mdText.PlainTextf(md.Image(t.Name, u))
	}
	mdText.PlainText("")

	mdText.H3("User").PlainText("")

	if req.User != "" {
		mdText.H3("User").
			PlainText("")
		mdText.H4("Single").
			PlainText("")

		for _, t := range types {
			u := fmt.Sprintf("%s://%s%s/chart?org=%s&repo=%s&type=%s&user=%s",
				req.scheme, req.host, os.Getenv("BASE_URL"), req.Org, req.Repo, t.Name, req.User)
			mdText.PlainTextf(md.Image(t.Name, u))
		}

		mdText.PlainText("")
	}

	mdText.H3("All").
		PlainText("")
	for _, t := range types {
		u := fmt.Sprintf("%s://%s%s/chart?org=%s&repo=%s&type=%s&users=all",
			req.scheme, req.host, os.Getenv("BASE_URL"), req.Org, req.Repo, t.Name)
		mdText.PlainTextf(md.Image(t.Name, u))
	}
	return mdText.String(), nil
}
