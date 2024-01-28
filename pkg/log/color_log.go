package log
import (
	"log/slog"
	"github.com/lmittmann/tint"
	"os"
)

func NewColorLog() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, nil))
}