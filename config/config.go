package config

import (
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Bubble `yaml:"bubble"`
	}

	Bubble struct {
		Quit     []string `yaml:"quit" env-default:"q"`
		Check    []string `yaml:"check" env-default:"space"`
		SwapUp   []string `yaml:"swapup" env-default:"shift+up"`
		SwapDown []string `yaml:"swapdown" env-default:"shift+down"`
		Insert   []string `yaml:"insert" env-default:"enter"`
		Remove   []string `yaml:"remove" env-default:"backspace"`
		Edit     []string `yaml:"edit" env-default:"right"`
		EditExit []string `yaml:"editexit" env-default:"esc"`
		Help     []string `yaml:"help" env-default:"?"`
		Undo     []string `yaml:"undo" env-default:"ctrl+z"`
		Up       []string `yaml:"up" env-default:"up,j"`
		Down     []string `yaml:"down" env-default:"down,k"`
		Bin      []string `yaml:"bin" env-default:"ctrl+b"`
		Restore  []string `yaml:"restore" env-default:"left"`
		EmptyBin []string `yaml:"emptybin" env-default:"c"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(os.Getenv("HOME")+"/.config/todo/config.yml", cfg)
	if err != nil && strings.HasPrefix(err.Error(), "config file") {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
