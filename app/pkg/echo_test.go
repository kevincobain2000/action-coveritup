package pkg

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEcho(t *testing.T) {
	var baseURL embed.FS
	var publicDir embed.FS
	e := NewEcho("", publicDir, baseURL)
	assert.NotNil(t, e)
}
