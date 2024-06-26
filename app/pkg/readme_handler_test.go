package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestReadmeNotOk(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/api/readme", func(c echo.Context) error {
		return NewReadmeHandler().Get(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name   string
		Query  string
		Status int
	}
	testCases := []TestCase{
		{
			Name:   "422",
			Query:  `?repo=repo&branch=branch&type=type`,
			Status: http.StatusUnprocessableEntity,
		},
		{
			Name:   "422",
			Query:  `?org=org&branch=branch&type=type`,
			Status: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		url := server.URL + "/api/readme" + tc.Query
		resp, err := http.Get(url)

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
	}
}
func TestReadmeOK(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/readme", func(c echo.Context) error {
		return NewReadmeHandler().Get(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name             string
		PreUploadRequest *UploadRequest
		Query            string
		Status           int
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
			},
			Query:  `?org=org&repo=repo&branch=branch&type=type`,
			Status: http.StatusOK,
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
			},
			Query:  `?org=org&repo=repo&branch=branch&type=type`,
			Status: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		c, err := NewUpload().Post(tc.PreUploadRequest)
		assert.NotNil(t, c)
		assert.NoError(t, err)
		url := server.URL + "/readme" + tc.Query
		resp, err := http.Get(url)

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
	}
}
