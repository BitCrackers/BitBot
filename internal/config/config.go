package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pelletier/go-toml"
)

type Filters struct {
	Response string
	Delete   bool
	RegExp   []string
}
type Config struct {
	GuildID string
	Debug   bool

	Filters            []Filters
	AutoReplyWithBuild []string
}

var emptyConfig = Config{
	GuildID: "",
	Debug:   false,
	Filters: []Filters{
		{
			Delete:   false,
			Response: "",
			RegExp:   []string{""},
		},
	},
	AutoReplyWithBuild: []string{},
}

var C Config

func Load() error {
	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	if _, err = os.Stat(path.Join(wd, "config.toml")); os.IsNotExist(err) {

		b, err := toml.Marshal(emptyConfig)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path.Join(wd, "config.toml"), b, 0644)
		if err != nil {
			return err
		}

		return errors.New("no config was found one has been created")
	}

	b, err := ioutil.ReadFile(path.Join(wd, "config.toml"))

	if err != nil {
		return err
	}

	err = toml.Unmarshal(b, &C)

	for _, f := range C.Filters {
		fmt.Printf("> Filter Loaded: %s\n", f.RegExp)
	}

	if err != nil {
		return err
	}

	return nil
}
