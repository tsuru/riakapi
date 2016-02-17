package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// ServiceConfig holds all the configuration of the application
type ServiceConfig struct {
	*Riak
	*SSH
	*RiakAPI

	*config.Server
}

// NewServiceConfig creates a service config instance
func NewServiceConfig() *ServiceConfig {
	cfg := &ServiceConfig{
		Server:  &config.Server{},
		Riak:    &Riak{},
		SSH:     &SSH{},
		RiakAPI: &RiakAPI{},
	}

	cfg.LoadConfiguration()

	return cfg
}

// LoadConfiguration loads application configuration
func (s *ServiceConfig) LoadConfiguration() {
	config.LoadEnvConfig(s)
	config.LoadEnvConfig(s.Server)
	s.Riak.LoadRiakConfigFromEnv()
	s.SSH.LoadSSHConfigFromEnv(s.Riak)
	s.RiakAPI.LoadRiakAPIConfigFromEnv(s.Riak)
	logrus.Info("Service configuration loaded")
}
