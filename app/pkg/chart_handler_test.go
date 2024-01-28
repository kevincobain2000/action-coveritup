package pkg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetChartErrors(t *testing.T) {

	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/chart", func(c echo.Context) error {
		return NewChartHandler().Get(c)
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
			Name:        "404",
			Query:       `?org=org&repo=repo&branch=branch&type=type`,
			Status:      http.StatusNotFound,
			ContentType: "application/json; charset=UTF-8",
		},
		{
			Name:        "422",
			Query:       `?repo=repo&branch=branch&type=type`,
			ContentType: "application/json; charset=UTF-8",
			Status:      http.StatusUnprocessableEntity,
		},
		{
			Name:        "422",
			Query:       `?org=org&branch=branch&type=type`,
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
		url := server.URL + "/chart" + tc.Query
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

func TestGetChartOK(t *testing.T) {

	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/chart", func(c echo.Context) error {
		return NewChartHandler().Get(c)
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
			},
			Query:            `?org=org&repo=repo&branch=branch&type=type&output=svg`,
			Status:           http.StatusOK,
			ContentType:      "image/svg+xml",
			ResponseContains: `<svg xmlns=`,
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
			Query:       `?org=org&repo=repo&branch=branch&type=type`,
			Status:      http.StatusOK,
			ContentType: "image/png",
		},
	}

	for _, tc := range testCases {
		c, err := NewUpload().Post(tc.PreUploadRequest)
		assert.NotNil(t, c)
		assert.NoError(t, err)
		url := server.URL + "/chart" + tc.Query
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
