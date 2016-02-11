package service

import (
	"os"
	"testing"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
	"gitlab.qdqmedia.com/shared-projects/riakapi/service/client"
)

var (
	envVars                      = map[string]string{}
	serviceIntegrationTestConfig *config.ServiceConfig
	serviceIntegrationTestClient client.Client
)

// setUp will set the required settings on env vars for the configuration and
// create the integration client
func setUp() {
	for k, v := range envVars {
		os.Setenv(k, v)
	}

	serviceIntegrationTestConfig = config.NewServiceConfig()
	//serviceIntegrationTestClient = client.NewRiak(serviceIntegrationTestConfig.RiakHost, serviceIntegrationTestConfig.RiakPort)
}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
}
