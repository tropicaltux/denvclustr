package denvclustr

import (
	"fmt"
	"log/slog"

	_ "github.com/tropicaltux/denvclustr/internal/logger"
)

func Run() {
	slog.Info("Project X started")
	fmt.Println("Hello, Project X!")
	slog.Info("Project X finished")
}
