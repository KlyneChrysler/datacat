// Package config owns environment configuration: the shape here, loading
// and validation in load.go (twelve-factor III).
package config

type Config struct {
	TargetURL string
}
