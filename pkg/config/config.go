package config

import (
	"fmt"
)

// Config is the application configuration.
type Config struct {
	// Server is the server configuration.
	Server ServerConfig
}

// ServerConfig is the server configuration.
type ServerConfig struct {
	// Port is the port to listen on.
	Port int
}

type LocalStoreConfig struct {
	// Path is the path to the database file.
	Path string
}