package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	v1 "github.com/priyawadhwa/prbot/pkg/config/v1"
	"gopkg.in/yaml.v3"
)

// Get returns the config from a file
func Get(filename string) (*v1.Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "reading %s", filename)
	}
	var c *v1.Config
	if err := yaml.Unmarshal(contents, &c); err != nil {
		return nil, errors.Wrap(err, "unmarshalling config")
	}
	return c, nil
}
