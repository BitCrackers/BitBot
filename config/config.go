package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Filter struct {
	Response string
	Delete   bool
	RegExp   []string
}

type Config struct {
	GuildID              string
	Debug                bool
	MuteRoleID           string
	ModLogChannelId      string
	Filters              []Filter
	Moderators           []string
	JanitorCycleDuration int
}

func Load() (Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return Config{}, fmt.Errorf("error while getting executable directory: %v", err)
	}
	ex = filepath.ToSlash(ex) // For the operating systems which use backslash...
	configPath := path.Join(path.Dir(ex), "config.toml")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {

		b, err := toml.Marshal(Config{
			Filters: []Filter{{}},
		})
		if err != nil {
			return Config{}, fmt.Errorf("error while marshalling empty config: %v", err)
		}

		err = ioutil.WriteFile(configPath, b, 0644)
		if err != nil {
			return Config{}, fmt.Errorf("error while writing empty config to file: %v", err)
		}

		logrus.Info("Empty config file created")
		os.Exit(0)
	}

	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file %v", err)
	}

	var cfg Config
	if err = toml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while decoding config: %v", err)
	}

	return cfg, nil
}
