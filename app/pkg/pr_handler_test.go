package pkg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPRErrors(t *testing.T) {

	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/pr", func(c echo.Context) error {
		return NewPRHandler().Get(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name             string
		Query            string
		Status           int
		ContentType      string
		ResponseContains string
	}
	testCases := []TestCase{
		{
			Name:        "409",
			Query:       `?org=org&repo=repo&branch=branch&base_branch=master&pr_num=2&type=type`,
			ContentType: "application/json; charset=UTF-8",
			Status:      http.StatusConflict,
		},
		{
			Name:        "422",
			Query:       `?org=org&branch=branch&type=type`,
			ContentType: "application/json; charset=UTF-8",
			Status:      http.StatusUnprocessableEntity,
		},
		{
			Name:        "422",
			Query:       `?org=org&repo=repo&type=type`,
			ContentType: "application/json; charset=UTF-8",
			Status:      http.StatusUnprocessableEntity,
		},
		{
			Name:        "422",
			Query:       `?org=org&repo=repo&branch=branch`,
			ContentType: "application/json; charset=UTF-8",
			Status:      http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		url := server.URL + "/pr" + tc.Query
		resp, err := http.Get(url)

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
		assert.Equal(t, tc.ContentType, resp.Header.Get("Content-Type"))
		if tc.ResponseContains != "" {
			data, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			assert.Contains(t, string(data), tc.ResponseContains)
		}
	}
}

func TestPROK(t *testing.T) {

	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/pr", func(c echo.Context) error {
		return NewPRHandler().Get(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name             string
		PreUploadRequest *UploadRequest
		Query            string
		Status           int
		ContentType      string
		ResponseContains string
	}
	testCases := []TestCase{
		{
			Name: "200",
			PreUploadRequest: &UploadRequest{
				Org:    "org",
				Repo:   "repo",
				Type:   "type",
				Branch: "branch",
				Commit: "commit",
				PRNum:  "2",
			},
			Query:       `?org=org&repo=repo&branch=branch&base_branch=master&pr_num=2&type=type`,
			Status:      http.StatusCreated,
			ContentType: "text/plain; charset=UTF-8",
		},
		{
			Name: "200",
			PreUploadRequest: &UploadRequest{
				Org:    "org",
				Repo:   "repo",
				Type:   "type",
				Branch: "branch",
				Commit: "commit",
				Score:  "1999",
				PRNum:  "2",
			},
			Query:       `?org=org&repo=repo&branch=branch&base_branch=master&pr_num=2&types=type`,
			Status:      http.StatusCreated,
			ContentType: "text/plain; charset=UTF-8",
		},
	}

	for _, tc := range testCases {
		c, err := NewUpload().Post(tc.PreUploadRequest)
		assert.NotNil(t, c)
		assert.NoError(t, err)
		url := server.URL + "/pr" + tc.Query
		resp, err := http.Get(url)

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
		assert.Equal(t, tc.ContentType, resp.Header.Get("Content-Type"))
		if tc.ResponseContains != "" {
			data, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			assert.Contains(t, string(data), tc.ResponseContains)
		}
	}
}
