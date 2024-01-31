package utils

import (
	"fmt"
	"log"
)

var (
	// Log is the logger.
	Log *Logger
)

// Log is a logger.
type Logger struct {
	// Logger is the logger.
	Logger *log.Logger
}