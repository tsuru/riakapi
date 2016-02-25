package config

import (
	"golang.org/x/crypto/ssh"

	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// SSH holds SSH configuration
type SSH struct {
	// SSHHost is the host where ssh will connect to execute riak-admin commands
	SSHHost string `envconfig:"SSH_HOST"`
	// SSHPort is the port where ssh will connect to execute riak-admin commands
	SSHPort int `envconfig:"SSH_PORT"`
	// SSHUser is the user which ssh will connect to execute riak-admin commands (needs passwordless sudo for riak-admin)
	SSHUser string `envconfig:"SSH_USER"`
	// SSHPassword is the password which ssh will connect to execute riak-admin commands
	SSHPassword string `envconfig:"SSH_PASSWORD"`
	// SSHPrivateKey is the private key which ssh will connect to execute riak-admin commands
	SSHPrivateKey string `envconfig:"SSH_PRIVATE_KEY"`

	// SSHAuthMethods internal variable with the ssh auth methods prepared based on the settings provided
	SSHAuthMethods []ssh.AuthMethod
}

// getAuthMethods returns the correct auth methods based on the password or private key
func (s *SSH) getAuthMethods() []ssh.AuthMethod {
	var methods []ssh.AuthMethod
	if s.SSHPrivateKey != "" {
		key, err := ssh.ParsePrivateKey([]byte(s.SSHPrivateKey))
		if err != nil {
			logrus.Errorf("Error parsing private ssh key: %v", err)
			return nil
		}
		methods = append(methods, ssh.PublicKeys(key))
	}

	if s.SSHPassword != "" {
		methods = append(methods, ssh.Password(s.SSHPassword))
	}
	if len(methods) == 0 {
		logrus.Warning("No ssh authentication methods present")
	}
	logrus.Debugf("Processed '%d' ssh auth methods", len(methods))
	return methods
}

// LoadSSHConfigFromEnv Loads SSH env and populates to default ones
func (s *SSH) LoadSSHConfigFromEnv(riakCfg *Riak) {
	config.LoadEnvConfig(s)
	if s.SSHHost == "" {
		// if no host then use the first riak host, if not then localhost
		if len(riakCfg.RiakClusterHosts) > 0 {
			s.SSHHost = riakCfg.RiakClusterHosts[0].Host
		} else {
			s.SSHHost = "127.0.0.1"
		}
	}
	if s.SSHPort == 0 {
		s.SSHPort = 22
	}
	s.SSHAuthMethods = s.getAuthMethods()
}
