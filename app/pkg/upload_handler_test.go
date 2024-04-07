package pkg

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevincobain2000/action-coveritup/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostUploadErrors(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.POST("/upload", func(c echo.Context) error {
		return NewUploadHandler().Post(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name   string
		Body   []byte
		Status int
	}
	testCases := []TestCase{
		{
			Name:   "422",
			Body:   []byte(`{"repo":"repo","branch":"branch","type":"type"}`),
			Status: http.StatusUnprocessableEntity,
		},
		{
			Name:   "422",
			Body:   []byte(`{"org":"org", "repo":"repo","branch":"branch","type":"type"}`),
			Status: http.StatusUnprocessableEntity,
		},
		{
			Name:   "422",
			Body:   []byte(`{"org":"org","repo":"repo","user": "user", "branch":"branch","type":"type"}`),
			Status: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		url := server.URL + "/upload"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(tc.Body))

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
	}
}

func TestPostUploadOK(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.POST("/upload", func(c echo.Context) error {
		return NewUploadHandler().Post(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name   string
		Body   []byte
		Status int
	}
	testCases := []TestCase{
		{
			Name:   "200",
			Body:   []byte(`{"org":"org","repo":"repo","user": "user", "branch":"branch","type":"type","commit":"commit"}`),
			Status: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		url := server.URL + "/upload"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(tc.Body))

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
	}
}
func TestPostCRUD(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.POST("/upload", func(c echo.Context) error {
		return NewUploadHandler().Post(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name   string
		Body   []byte
		Status int
	}
	testCases := []TestCase{
		{
			Name:   "200",
			Body:   []byte(`{"org":"org","repo":"repo","user": "user", "branch":"branch1","type":"type","commit":"commit", "branches": "master develop", "pr_num": "1234"}`),
			Status: http.StatusOK,
		},
		{
			Name:   "200",
			Body:   []byte(`{"org":"org","repo":"repo","user": "user", "branch":"branch2","type":"type","commit":"commit", "branches": "branch1 master develop", "pr_num": "1234"}`),
			Status: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		url := server.URL + "/upload"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(tc.Body))

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
	}

	cm := models.Coverage{}
	ret, err := cm.GetAllBranches("org", "repo", "type")
	assert.NoError(t, err)
	assert.Equal(t, []string{"branch1", "branch2"}, ret)
}
