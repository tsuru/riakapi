package config

import (
	"fmt"

	"github.com/NYTimes/gizmo/config"
)

// Riak holds riak configuration
type Riak struct {
	RiakHost string `envconfig:"RIAK_HOST"`
	RiakPort int    `envconfig:"RIAK_PORT"`
}

func (r *Riak) String() string {
	return fmt.Sprintf("%s:%d", r.RiakHost, r.RiakPort)
}

//LoadRiakConfigFromEnv Loads riak configuration from enviroment setting defaults if not preset
func (r *Riak) LoadRiakConfigFromEnv() {
	config.LoadEnvConfig(r)
	if r.RiakHost == "" {
		r.RiakHost = "127.0.0.1"
	}

	if r.RiakPort == 0 {
		r.RiakPort = 8087
	}
}
