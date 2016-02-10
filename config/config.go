package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// ServiceConfig holds all the configuration of the application
type ServiceConfig struct {
	*Riak

	*config.Server
}

// NewServiceConfig creates a service config instance
func NewServiceConfig() *ServiceConfig {
	cfg := &ServiceConfig{
		Server: &config.Server{},
		Riak:   &Riak{},
	}

	cfg.LoadConfiguration()

	return cfg
}

// LoadConfiguration loads application configuration
func (s *ServiceConfig) LoadConfiguration() {
	config.LoadEnvConfig(s)
	config.LoadEnvConfig(s.Server)
	config.LoadEnvConfig(s.Riak)

	logrus.Info("Service configuration loaded")
}
