package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NYTimes/gizmo/server"

	"github.com/tsuru/riakapi/config"
	"github.com/tsuru/riakapi/service/client"
)

var (
	serviceITestCfg *config.ServiceConfig
)

func createIntegrationConfig() *config.ServiceConfig {

	envVars := map[string]string{
		"RIAK_HOSTS":        os.Getenv("RIAK_PORT_8087_TCP_ADDR"),
		"RIAK_HTTP_PORT":    os.Getenv("RIAK_PORT_8098_TCP_PORT"),
		"RIAK_PB_PORT":      os.Getenv("RIAK_PORT_8087_TCP_PORT"),
		"RIAK_USER":         "riakapi",
		"RIAK_PASSWORD":     "riakapi",
		"RIAK_INSECURE_TLS": "1",
		"SSH_HOST":          os.Getenv("RIAK_PORT_22_TCP_HOST"),
		"SSH_PORT":          os.Getenv("RIAK_PORT_22_TCP_PORT"),
		"SSH_USER":          "riakapi",
		"SSH_PASSWORD":      "riakapi",
	}

	// Set env vars to load configuration
	for k, v := range envVars {
		os.Setenv(k, v)
	}

	return config.NewServiceConfig()
}

func TestMain(m *testing.M) {
	// Create configuration
	serviceITestCfg = createIntegrationConfig()

	// Run tests
	os.Exit(m.Run())
}

// TestIntegrationInstanceCreationOk Creates a new bucket on a bucket type. we
// check if the bucket & the bucket type are present
func TestIntegrationInstanceCreationOk(t *testing.T) {
	serviceTestClient := client.NewRiak(serviceITestCfg)

	uri := "/resources?name=test-bucket&plan=tsuru-counter&team=myteam&user=username"
	wantBody := ""
	wantCode := http.StatusOK

	// Create our dummy server (with config & client)
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceTestCfg, Client: serviceTestClient})

	// Create the request
	r, _ := http.NewRequest("POST", uri, nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != wantCode {
		t.Errorf("expected response code of %d; got %d", wantCode, w.Code)
	}

	var got interface{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("unable to JSON decode response body: ", err)
	}

	// Check body
	if got != wantBody {
		t.Errorf("expected response body of\n%#v;\ngot\n%#v", wantBody, got)
	}

	// Check correct bucket type
	if serviceTestClient.GetBucketType("test-bucket") != "tsuru-counter" {
		t.Error("Bucket not created correctly")
	}
}
