// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Period     time.Duration `config:"period"`
	Repository string        `config:"repository"`
}

var DefaultConfig = Config{
	Period:     1 * time.Second,
	Repository: "nginx",
}

func (c *Config) RepositoryLong() string {
	if strings.Index(c.Repository, "/") == -1 {
		return fmt.Sprintf("library/%s", c.Repository)
	}
	return c.Repository
}
