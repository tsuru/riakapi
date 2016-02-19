package config

import (
	"fmt"

	"github.com/NYTimes/gizmo/config"
)

// Riak holds riak configuration
type Riak struct {
	// RiakHost is the host where app will connect to execute riak actions
	RiakHost string `envconfig:"RIAK_HOST"`
	// RiakPort is the port where app will connect to execute riak actions
	RiakPort int `envconfig:"RIAK_PORT"`
	// RiakUser is the user which app will connect to execute riak actions
	RiakUser string `envconfig:"RIAK_USER"`
	// RiakPass is the password which app will connect to execute riak actions
	RiakPass string `envconfig:"RIAK_PASSWORD"`
	// RiakCaCert is the CA certificate used to authenticate the TLS connection with riak server
	RiakCaCert string `envconfig:"RIAK_CA_PATH"`
	// RiakInsecureTLS makes an insecure TLS connection (we should be sure of the riak server and no MitM attacks are possible)
	RiakInsecureTLS int `envconfig:"RIAK_INSECURE_TLS"`
	// RiakServerName is the riaks server name for the certificates authentication
	RiakServerName string `envconfig:"RIAK_SERVER_NAME"`
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
