package log

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

func NewColorLog() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, nil))
}
