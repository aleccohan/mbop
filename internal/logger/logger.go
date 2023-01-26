package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var Log logr.Logger

func Init() error {
	// TODO: add cloudwatch hook
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("who watches the watchmen (%v)?", err)
	}

	Log = zapr.NewLogger(zapLog)

	return nil
}
