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
}

// LoadRiakAPIConfigFromEnv loads the riakapi service configuration from the env
func (r *RiakAPI) LoadRiakAPIConfigFromEnv(riakCfg *Riak) {
	config.LoadEnvConfig(r)

	// Warn if salt is disabled
	if r.RiakAPISalt == "" {
		logrus.Warning("'RIAKAPI_SALT' not set, not salting the passwords")
	}

	// Warn if security is disabled
	if r.RiakAPIPassword == "" {
		logrus.Warning("'RIAKAPI_PASSWORD' not set, service security is disabled")
	}
}
