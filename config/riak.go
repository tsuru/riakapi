package config

import (
	"io/ioutil"
	"strings"

	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// Riak holds riak configuration
type Riak struct {
	// RiakHosts is an array of riak hosts separated by ";"
	RiakHosts string `envconfig:"RIAK_HOSTS"`
	// RiakHTTPPort is the port where http listens on riak cluster
	RiakHTTPPort int `envconfig:"RIAK_HTTP_PORT"`
	// RiakPBPort is the port where protobuffer listens on riak cluster (required)
	RiakPBPort int `envconfig:"RIAK_PB_PORT"`
	// RiakUser is the user which app will connect to execute riak actions
	RiakUser string `envconfig:"RIAK_USER"`
	// RiakPass is the password which app will connect to execute riak actions
	RiakPass string `envconfig:"RIAK_PASSWORD"`
	// RiakCaCert is the CA certificate used to authenticate the TLS connection with riak server
	RiakCaCertPath string `envconfig:"RIAK_CA_PATH"`
	// RiakInsecureTLS makes an insecure TLS connection (we should be sure of the riak server and no MitM attacks are possible)
	RiakInsecureTLS int `envconfig:"RIAK_INSECURE_TLS"`
	// RiakServerName is the riaks server name for the certificates authentication
	RiakServerName string `envconfig:"RIAK_SERVER_NAME"`

	// RiakClusterHosts is a custom attr to set the riak cluster hosts in the correct way
	RiakClusterHosts []string

	// RiakCaCert is a custo attr with the CA cert path content file
	RiakCaCert string
}

//LoadRiakConfigFromEnv Loads riak configuration from enviroment setting defaults if not preset
func (r *Riak) LoadRiakConfigFromEnv() {
	config.LoadEnvConfig(r)

	// Check required
	if r.RiakHosts == "" {
		logrus.Fatal("RIAK_HOSTS is required")
	}

	if r.RiakHTTPPort == 0 {
		r.RiakHTTPPort = 8098
	}

	if r.RiakPBPort == 0 {
		r.RiakPBPort = 8087
	}

	if r.RiakCaCertPath != "" {
		var pemData []byte
		var err error
		if pemData, err = ioutil.ReadFile(r.RiakCaCertPath); err != nil {
			logrus.Fatalf("Error reading ca cert: %v", err)
		}
		r.RiakCaCert = string(pemData)
	}

	// Split the hosts and create an slice with each one
	r.RiakClusterHosts = strings.Split(r.RiakHosts, ";")
}
