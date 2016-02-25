package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// RiakHost is a helper struct for decoding json configuration
type RiakHost struct {
	Host       string `json:"host"`
	ServerName string `json:"server_name,omitempty"`
}

// Riak holds riak configuration
type Riak struct {
	// RiakHosts is a json  array of host and server_name strings hash.
	//  server_name can be blan, if blank then host will be used
	// Example:
	//	[
	//		{
	//		  "host": "c1.test.org",
	//		  "server_name": "c1"
	//		},
	//		{
	//		  "host": "c2.test.org"
	//		}
	//	]
	RiakHosts string `envconfig:"RIAK_HOSTS"`
	// RiakHTTPPort is the port where http listens on riak cluster
	RiakHTTPPort int `envconfig:"RIAK_HTTP_PORT"`
	// RiakPBPort is the port where protobuffer listens on riak cluster (required)
	RiakPBPort int `envconfig:"RIAK_PB_PORT"`
	// RiakUser is the user which app will connect to execute riak actions
	RiakUser string `envconfig:"RIAK_USER"`
	// RiakPass is the password which app will connect to execute riak actions
	RiakPass string `envconfig:"RIAK_PASSWORD"`
	// RiakRootCaCertPath path to the root CA certificate used to authenticate the TLS connection with riak server
	RiakRootCaCertPath string `envconfig:"RIAK_ROOT_CA_PATH"`
	// RiakRootCaCert Root CA cert content file (alternative to RIAK_ROOT_CA_PATH)
	RiakRootCaCert string `envconfig:"RIAK_ROOT_CA"`
	// RiakInsecureTLS makes an insecure TLS connection (we should be sure of the riak server and no MitM attacks are possible)
	RiakInsecureTLS int `envconfig:"RIAK_INSECURE_TLS",default:"0"`

	// RiakClusterHosts is a custom attr to set the riak cluster hosts in the correct way
	RiakClusterHosts []*RiakHost
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

	if r.RiakRootCaCertPath != "" {
		var pemData []byte
		var err error
		if pemData, err = ioutil.ReadFile(r.RiakRootCaCertPath); err != nil {
			logrus.Fatalf("Error reading ca cert: %v", err)
		}
		r.RiakRootCaCert = string(pemData)
	}

	// Convert json to our clister objects
	var cluster []*RiakHost
	err := json.Unmarshal([]byte(r.RiakHosts), &cluster)

	if err != nil {
		logrus.Fatalf("Wrong RIAK_HOSTS format")
	}

	// if no servername then set the same as host
	for _, n := range cluster {
		if n.ServerName == "" {
			n.ServerName = n.Host
		}
	}
	r.RiakClusterHosts = cluster

}
