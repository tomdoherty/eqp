package eqp

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config stores the pattetns to match
type Config struct {
	URL       string `yaml:"url"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Insecure  string `yaml:"insecure"`
	Frequency string `yaml:"frequency"`

	Matches []struct {
		Name    string `yaml:"name"`
		Pattern string `yaml:"pattern"`
		Type    string `yaml:"type"`
		Seconds string `yaml:"seconds"`
		Index   string `yaml:"index"`
	} `yaml:"matches"`
}

func (c *Config) loadConfig(infile string) (err error) {
	file, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(data), &c)

	if err != nil {
		return err
	}
	return nil
}
