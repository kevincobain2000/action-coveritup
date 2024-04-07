package pkg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetBadgeErrors(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/badge", func(c echo.Context) error {
		return NewBadgeHandler().Get(c)
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
			Name:             "404",
			Query:            `?org=org&repo=repo&branch=branch&type=type`,
			Status:           http.StatusOK,
			ContentType:      "image/svg+xml",
			ResponseContains: `fill="white">404</text>`,
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
		url := server.URL + "/badge" + tc.Query
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
func TestGetBadgeOK(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.GET("/badge", func(c echo.Context) error {
		return NewBadgeHandler().Get(c)
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
			Query:            `?org=org&repo=repo&branch=branch&type=type`,
			Status:           http.StatusOK,
			ContentType:      "image/svg+xml",
			ResponseContains: `<text x="115" y="15" fill="#fff">0</text>`,
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
			Query:            `?org=org&repo=repo&branch=branch&type=type`,
			Status:           http.StatusOK,
			ContentType:      "image/svg+xml",
			ResponseContains: `<text id="rlink" x="123.5" y="14">2.0k</text>`,
		},
	}

	for _, tc := range testCases {
		c, err := NewUpload().Post(tc.PreUploadRequest)
		assert.NotNil(t, c)
		assert.NoError(t, err)
		url := server.URL + "/badge" + tc.Query
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
