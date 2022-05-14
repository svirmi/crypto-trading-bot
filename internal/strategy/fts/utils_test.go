package fts

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Initialize(false, logrus.TraceLevel)
	code := m.Run()
	os.Exit(code)
}
