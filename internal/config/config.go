package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds all configuration for the barcodekit service.
type Config struct {
	Addr   string
	NoAuth bool
}

// Defaults
const (
	DefaultAddr   = ":8080"
	DefaultNoAuth = false
)

// Load reads configuration from defaults, env vars, and flags.
// Priority: defaults < env vars < flags.
func Load() *Config {
	cfg := &Config{
		Addr:   DefaultAddr,
		NoAuth: DefaultNoAuth,
	}

	// Env vars
	if v := os.Getenv("BARCODEKIT_ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("BARCODEKIT_NO_AUTH"); v == "true" || v == "1" {
		cfg.NoAuth = true
	}

	// Flags
	flag.StringVar(&cfg.Addr, "addr", cfg.Addr, "listen address")
	flag.BoolVar(&cfg.NoAuth, "no-auth", cfg.NoAuth, "disable auth (dev mode)")
	flag.Parse()

	return cfg
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s no-auth=%v", c.Addr, c.NoAuth)
}
