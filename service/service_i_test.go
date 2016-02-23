package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/NYTimes/gizmo/server"
	riak "github.com/basho/riak-go-client"

	"github.com/tsuru/riakapi/config"
	"github.com/tsuru/riakapi/service/client"
)

var (
	serviceITestCfg *config.ServiceConfig

	envVars = map[string]string{
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
)

func createIntegrationConfig() *config.ServiceConfig {
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
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})

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

func TestIntegrationInstanceBindingOk(t *testing.T) {

	// Prepare
	serviceTestClient := client.NewRiak(serviceITestCfg)
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})
	instance := "test-instance"
	plan := "tsuru-counter"
	appHost := "myapp.test.org"
	testKey := "MyTestAwesomeKey1234567890"
	uri := fmt.Sprintf("/resources?name=%s&plan=%s&team=myteam&user=username", instance, plan)

	// Create a new instance
	r, _ := http.NewRequest("POST", uri, nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("Coudn't prepare the isntance for the test")
	}

	// Bind our fresh created instance
	uri = fmt.Sprintf("/resources/%s/bind-app?app-host=%s", instance, appHost)
	wantCode := http.StatusCreated
	wantBody := map[string]string{
		"RIAK_HOSTS":       envVars["RIAK_HOSTS"],
		"RIAK_HTTP_PORT":   envVars["RIAK_HTTP_PORT"],
		"RIAK_PB_PORT":     envVars["RIAK_PB_PORT"],
		"RIAK_USER":        fmt.Sprintf("tsuru_%s", appHost),
		"RIAK_PASSWORD":    "",
		"RIAK_BUCKET_TYPE": plan,
		"RIAK_BUCKET":      instance,
	}
	r, _ = http.NewRequest("POST", uri, nil)
	w = httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != wantCode {
		t.Errorf("expected response code of %d; got %d", wantCode, w.Code)
	}

	var got map[string]string
	err := json.NewDecoder(w.Body).Decode(&got)

	if err != nil {
		t.Error("unable to JSON decode response body: ", err)
	}

	// save password for testing it and set blank the random password in order to compare
	pass := got["RIAK_PASSWORD"]
	got["RIAK_PASSWORD"] = ""

	// Check body json decoded
	if !reflect.DeepEqual(got, wantBody) {
		t.Errorf("expected response body of\n%#v;\ngot\n%#v", wantBody, got)
	}

	// Check getting and retrieving a key on the recent created bucket with the
	// username and password

	// Create riak connection with the new user
	u := got["RIAK_USER"]
	tc := &tls.Config{InsecureSkipVerify: true}
	a := &riak.AuthOptions{User: u, Password: pass, TlsConfig: tc}
	h := strings.Split(got["RIAK_HOSTS"], ":")[0]
	no := &riak.NodeOptions{
		RemoteAddress: fmt.Sprintf("%s:%s", h, got["RIAK_PB_PORT"]),
		AuthOptions:   a,
	}
	var n *riak.Node
	if n, err = riak.NewNode(no); err != nil {
		t.Errorf("Error creating node: %v", err)
	}
	co := &riak.ClusterOptions{Nodes: []*riak.Node{n}}
	var cluster *riak.Cluster
	if cluster, err = riak.NewCluster(co); err != nil {
		t.Errorf("Error connecting to riak: %v", err)
	}
	if err := cluster.Start(); err != nil {
		t.Errorf("Error connecting to riak: %v", err)
	}

	// set a key on the bucket
	value := rand.Intn(100)
	cmd, _ := riak.NewUpdateCounterCommandBuilder().
		WithBucketType(got["RIAK_BUCKET_TYPE"]).
		WithBucket(got["RIAK_BUCKET"]).
		WithKey(testKey).
		WithIncrement(int64(value)).
		Build()
	if err = cluster.Execute(cmd); err != nil {
		t.Errorf("Error setting test key/value: %v", err)
	}

	// Retrieve
	cmd, _ = riak.NewFetchCounterCommandBuilder().
		WithBucketType(got["RIAK_BUCKET_TYPE"]).
		WithBucket(got["RIAK_BUCKET"]).
		WithKey(testKey).
		Build()

	if err = cluster.Execute(cmd); err != nil {
		t.Errorf("Error fetching test key/value: %v", err)
	}

	fcc := cmd.(*riak.FetchCounterCommand)
	if int64(value) != fcc.Response.CounterValue {
		t.Errorf("Fetched counter error; expected: %d; got: %d", value, fcc.Response.CounterValue)
	}

}
