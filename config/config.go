// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"
	"time"
)

type Config struct {
	Period       time.Duration `config:"period"`
	Repositories []Repository  `config:"repositories"`
}

type Repository struct {
	Name      string  `config:"name"`
	Namespace *string `config:"namespace"`
}

var DefaultConfig = Config{
	Period:       1 * time.Second,
	Repositories: []Repository{},
}

func (r *Repository) FullName() string {
	if r.Namespace == nil {
		return fmt.Sprintf("library/%s", r.Name)
	} else {
		return fmt.Sprintf("%s/%s", *r.Namespace, r.Name)
	}
}
