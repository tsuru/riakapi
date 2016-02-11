package config

import "fmt"

// Riak holds riak configuration
type Riak struct {
	RiakHost string `envconfig:"RIAK_HOST"`
	RiakPort int    `envconfig:"RIAK_PORT"`
}

func (r *Riak) String() string {
	return fmt.Sprintf("%s:%d", r.RiakHost, r.RiakPort)
}
