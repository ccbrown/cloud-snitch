package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewId(t *testing.T) {
	id := NewId("x")
	parts := strings.Split(string(id), "-")
	assert.Len(t, parts, 2)
	assert.Equal(t, "x", parts[0])
	assert.Len(t, parts[1], 22)
}
