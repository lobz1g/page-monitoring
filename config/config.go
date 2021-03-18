package config

import (
	"encoding/json"
	"errors"
	"os"
)

const configFilename = "config/config.json"

type (
	cfg struct {
		Timeout   string   `json:"timeout"`
		Token     string   `json:"token"`
		Channel   string   `json:"channel"`
		DebugMode bool     `json:"debug"`
		Urls      []string `json:"url"`
	}
)

// readFile used for unit testing
var readFile = os.ReadFile

// unmarshal used for unit testing
var unmarshal = json.Unmarshal

func Get() (*cfg, error) {
	d, err := readFile(configFilename)
	if err != nil {
		return nil, err
	}

	var c *cfg
	err = unmarshal(d, &c)
	if err != nil {
		return nil, errors.New("cannot read config file")
	}
	return c, nil
}
