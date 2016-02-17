package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// RiakAPI holds the configuration for the riak api service configuration
type RiakAPI struct {
	// RiakAPIUsername is the user used to authenticate against the API service
	RiakAPIUsername string `envconfig:"RIAKAPI_USERNAME"`

	// RiakAPIPassword is the password used to authenticate against the API service
	RiakAPIPassword string `envconfig:"RIAKAPI_PASSWORD"`

	// RiakAPISalt is the salt used for the password creation
	RiakAPISalt string `envconfig:"RIAKAPI_SALT"`

	// RiakAPIRiakHTTPPort is the port where http listens on riak cluster
	RiakAPIRiakHTTPPort int `envconfig:"RIAKAPI_RIAK_HTTP_PORT"`

	// RiakAPIRiakPBPort is the port where protobuffer listens on riak cluster
	RiakAPIRiakPBPort int `envconfig:"RIAKAPI_RIAK_PB_PORT"`
}

// LoadRiakAPIConfigFromEnv loads the riakapi service configuration from the env
func (r *RiakAPI) LoadRiakAPIConfigFromEnv(riakCfg *Riak) {
	config.LoadEnvConfig(r)
	//default protobuffer port
	if r.RiakAPIRiakPBPort == 0 {
		r.RiakAPIRiakPBPort = 8087
	}

	//default http port
	if r.RiakAPIRiakHTTPPort == 0 {
		r.RiakAPIRiakHTTPPort = 8098
	}

	// Warn if salt is disabled
	if r.RiakAPISalt == "" {
		logrus.Warning("'RIAKAPI_SALT' not set, not salting the passwords")
	}

	// Warn if security is disabled
	if r.RiakAPIPassword == "" {
		logrus.Warning("'RIAKAPI_PASSWORD' not set, service security is disabled")
	}
}
