package pkg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetBadge(t *testing.T) {
	e := echo.New()

	BeforeEach()
	defer AfterEach()

	e.GET("/badge", func(c echo.Context) error {
		return NewBadgeHandler().Get(c)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	type TestCase struct {
		Name                     string
		Query                    string
		ExpectedStatus           int
		ExpectedResponseContains string
	}
	testCases := []TestCase{
		{
			Name:                     "404",
			Query:                    `?org=org&repo=repo&branch=branch&type=type`,
			ExpectedStatus:           http.StatusOK,
			ExpectedResponseContains: `<text x="96" y="14" fill="white">404</text>`,
		},
	}

	for _, tc := range testCases {
		url := server.URL + "/badge" + tc.Query
		resp, err := http.Get(url)
		data, _ := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Contains(t, string(data), tc.ExpectedResponseContains)
		assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)
		assert.Equal(t, "image/svg+xml", resp.Header.Get("Content-Type"))
	}
}
