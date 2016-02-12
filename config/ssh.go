package config

import (
	"golang.org/x/crypto/ssh"

	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// SSH holds SSH configuration
type SSH struct {
	SSHHost       string `envconfig:"SSH_HOST"`
	SSHPort       int    `envconfig:"SSH_PORT"`
	SSHUser       string `envconfig:"SSH_USER"`
	SSHPassword   string `envconfig:"SSH_PASSWORD"`
	SSHPrivateKey string `envconfig:"SSH_PRIVATE_KEY"`

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
	logrus.Debugf("Processed '%d' ssh auth methods", len(methods))
	return methods
}

// LoadSSHConfigFromEnv Loads SSH env and populates to default ones
func (s *SSH) LoadSSHConfigFromEnv(riakCfg *Riak) {
	config.LoadEnvConfig(s)
	if s.SSHHost == "" {
		// if no host then use the same as riak host, if not then localhost
		if riakCfg.RiakHost != "" {
			s.SSHHost = riakCfg.RiakHost
		} else {
			s.SSHHost = "127.0.0.1"
		}
	}
	if s.SSHPort == 0 {
		s.SSHPort = 22
	}
	s.SSHAuthMethods = s.getAuthMethods()
}
