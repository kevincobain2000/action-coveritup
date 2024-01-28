package pkg

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDestroyErrors(t *testing.T) {

	BeforeEach()
	defer AfterEach()
	e := echo.New()
	e.POST("/destroy", func(c echo.Context) error {
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
		url := server.URL + "/destroy"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(tc.Body))

		assert.NoError(t, err)
		assert.Equal(t, tc.Status, resp.StatusCode)
	}
}
