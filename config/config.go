// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period     time.Duration `config:"period"`
	Repository string        `config:"repository"`
}

var DefaultConfig = Config{
	Period:     1 * time.Second,
	Repository: "library/nginx",
}
