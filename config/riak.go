package config

import "fmt"

// Riak holds riak configuration
type Riak struct {
	Host string `envconfig:"RIAK_HOST"`
	Port int    `envconfig:"RIAK_PORT"`
}

func (r *Riak) String() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
