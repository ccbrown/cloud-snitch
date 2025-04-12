package app_test

import (
	"testing"

	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	m.Run()
}
