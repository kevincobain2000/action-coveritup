package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/coocood/freecache"
	"github.com/kevincobain2000/action-coveritup/db"
	"github.com/sirupsen/logrus"
)

type Github struct {
	body   *bytes.Buffer
	log    *logrus.Logger
	client *http.Client
	api    string
	cache  *freecache.Cache
}

func NewGithub() *Github {
	return &Github{
		body:   bytes.NewBuffer([]byte(`{"state":"success", "context": "coveritup", "description": "authenticated"}`)),
		log:    Logger(),
		client: &http.Client{},
		api:    os.Getenv("GITHUB_API"),
		cache:  db.Cache(),
	}
}

// VerifyGithubToken verifies the github token
// /repos/:owner/:repo/statuses/commit_hash
func (g *Github) VerifyGithubToken(token, orgName, repoName, commitHash string) error {
	if token == "" {
		err := errors.New("token is empty")
		g.log.Error(err)
		return err
	}

	// look up cache, don't call Github API if it's already authenticated
	cacheKey := []byte(fmt.Sprintf("%s/%s/%s", orgName, repoName, token))
	ret, err := g.cache.Get(cacheKey)
	if err == nil && string(ret) == "true" {
		return nil
	}

	url := g.getEndpoint(g.api, orgName, repoName, commitHash)

	req, err := http.NewRequest("POST", url, g.body)
	if err != nil {
		g.log.Error(err)
		return err
	}
	g.setHeader(req, token)

	resp, err := g.request(req)
	if err != nil {
		g.log.Error(err)
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		err := fmt.Errorf("github auth response code is %d", resp.StatusCode)
		g.log.Error(err)
		return err
	}
	err = g.cache.Set(cacheKey, []byte("true"), 60*60*24*7)

	return err
}

func (g *Github) getEndpoint(api, orgName, repoName, commitHash string) string {
	return fmt.Sprintf("%s/repos/%s/%s/statuses/%s", api, orgName, repoName, commitHash)
}

func (g *Github) setHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")
}

func (g *Github) request(req *http.Request) (*http.Response, error) {
	resp, err := g.client.Do(req)
	if err != nil {
		g.log.Error(err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			g.log.Info(err)
		}
	}()
	return resp, nil
}
