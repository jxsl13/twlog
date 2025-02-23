package config

import (
	"errors"
)

func NewSaidConfig() SaidConfig {
	return SaidConfig{
		Deduplicate: false,
	}
}

type SaidConfig struct {
	Deduplicate bool `koanf:"deduplicate" short:"D" description:"deduplicate objects based on all fields"`
	Extended    bool `koanf:"extended" short:"e" description:"add two additional fields, file and id to the output"`
	IPsOnly     bool `koanf:"ips.only" short:"i" description:"only print IP addresses"`
}

func (cfg *SaidConfig) Validate() error {

	if cfg.Extended && cfg.IPsOnly {
		return errors.New("extended and ips only flags are mutually exclusive")
	}

	return nil
}
