package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/coocood/freecache"
	"github.com/kevincobain2000/action-coveritup/db"
	"github.com/sirupsen/logrus"
)

type Github struct {
	log   *logrus.Logger
	api   string
	cache *freecache.Cache
}

func NewGithub() *Github {
	return &Github{
		log:   Logger(),
		api:   os.Getenv("GITHUB_API"),
		cache: db.Cache(),
	}
}

// VerifyGithubToken verifies the github token
// /repos/:owner/:repo/statuses/commit
func (g *Github) VerifyGithubToken(token, orgName, repoName, commit string) error {
	if token == "" {
		err := errors.New("token is empty")
		g.log.Error(err)
		return err
	}

	// look up cache, don't call Github API if it's already authenticated
	cacheKey := []byte(MD5(fmt.Sprintf("%s/%s/%s", orgName, repoName, token)))
	ret, err := g.cache.Get(cacheKey)
	if err == nil && string(ret) == "true" {
		return nil
	}

	url := g.getEndpoint(g.api, orgName, repoName, commit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		g.log.Error(err)
		return err
	}
	g.setHeader(req, token)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		g.log.Error(err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			g.log.Info(err)
		}
	}()

	if err != nil {
		g.log.Error(err)
		return err
	}

	if resp.StatusCode == 422 {
		err := fmt.Errorf("github auth response code is unprocessable entity %d", resp.StatusCode)
		g.log.Error(err)
		return err
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		err := fmt.Errorf("github auth response code is unauthorized %d", resp.StatusCode)
		g.log.Error(err)
		return err
	}
	if resp.StatusCode >= 500 {
		err := fmt.Errorf("github auth response code is server error %d", resp.StatusCode)
		g.log.Error(err)
		return err
	}
	g.log.Info("github auth response code is " + resp.Status)
	err = g.cache.Set(cacheKey, []byte("true"), 60*60*24*7)
	if err != nil {
		g.log.Error(err)
	}

	return err
}

// /repos/{owner}/{repo}/commits/{ref}
func (g *Github) getEndpoint(api, orgName, repoName, commit string) string {
	return fmt.Sprintf("%s/repos/%s/%s/commits/%s", api, orgName, repoName, commit)
}

func (g *Github) setHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")
}
