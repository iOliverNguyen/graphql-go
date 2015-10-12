package debug

import (
	"log"
	"os"
)

func New(prefix string) *log.Logger {
	logger := log.New(os.Stderr, "[LOG] "+prefix+": ", log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}
