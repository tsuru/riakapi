package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/NYTimes/gizmo/server"
	riak "github.com/basho/riak-go-client"

	"github.com/tsuru/riakapi/config"
	"github.com/tsuru/riakapi/service/client"
)

var (
	serviceITestCfg *config.ServiceConfig

	envVars = map[string]string{
		"RIAK_HOSTS":        fmt.Sprintf(`[{"host":"%[1]s","server_name":"%[1]s"}]`, os.Getenv("RIAK_PORT_8087_TCP_ADDR")),
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

// ---------------------------- Setup & teardown ------------------------------
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

// ---------------------------- Helper functions ------------------------------
func incCounter(c *riak.Cluster, bucketType, bucket, key string, value int) error {
	cmd, _ := riak.NewUpdateCounterCommandBuilder().
		WithBucketType(bucketType).
		WithBucket(bucket).
		WithKey(key).
		WithIncrement(int64(value)).
		Build()
	if err := c.Execute(cmd); err != nil {
		return err
	}
	return nil
}

func getCounter(c *riak.Cluster, bucketType, bucket, key string) (int64, error) {
	cmd, _ := riak.NewFetchCounterCommandBuilder().
		WithBucketType(bucketType).
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err := c.Execute(cmd); err != nil {
		return 0, err
	}

	fcc := cmd.(*riak.FetchCounterCommand)

	return fcc.Response.CounterValue, nil
}

// ---------------------------- Integration tests ------------------------------

// TestIntegrationInstanceCreationOk Creates a new bucket on a bucket type. we
// check if the bucket & the bucket type are present
func TestIntegrationInstanceCreationOk(t *testing.T) {
	serviceTestClient := client.NewRiak(serviceITestCfg)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	plan := "tsuru-counter"
	instance := fmt.Sprintf("test-instance-%d", rnd.Int())
	uri := fmt.Sprintf("/resources?name=%s&plan=%s&team=myteam&user=username", instance, plan)
	wantBody := ""
	wantCode := http.StatusOK

	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})

	// Create the instance
	r, _ := http.NewRequest("POST", uri, nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	// Check response
	if w.Code != wantCode {
		t.Errorf("expected response code of %d; got %d", wantCode, w.Code)
	}

	var got interface{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("unable to JSON decode response body: ", err)
	}

	if got != wantBody {
		t.Errorf("expected response body of\n%#v;\ngot\n%#v", wantBody, got)
	}

	// Check correct bucket type
	if serviceTestClient.GetBucketType(instance) != plan {
		t.Error("Bucket not created correctly")
	}
}

func TestIntegrationInstanceBindingOk(t *testing.T) {

	// Prepare
	serviceTestClient := client.NewRiak(serviceITestCfg)
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instance := fmt.Sprintf("test-instance-%d", rnd.Int())
	plan := "tsuru-counter"
	appHost := "myapp.test.org"
	testKey := fmt.Sprintf("MyTestAwesomeKey_%d", rnd.Int()) // Random key always to avoid collisions between runs
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

	// Check response
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
	got["RIAK_PASSWORD"] = pass

	// Check getting and retrieving a key on the recent created bucket with the
	// username and password

	cfg := createIntegrationConfig()
	cfg.RiakUser = got["RIAK_USER"]
	cfg.RiakPass = got["RIAK_PASSWORD"]
	cluster, err := client.NewRiakCluster(cfg)
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Errorf("Error connecting to riak with new credentials: %v", err)
	}

	// Set a key on the bucket
	value := rnd.Intn(100)
	if err = incCounter(cluster, got["RIAK_BUCKET_TYPE"], got["RIAK_BUCKET"], testKey, value); err != nil {
		t.Errorf("Error setting test key/value: %v", err)
	}

	// Retrieve
	var res int64
	if res, err = getCounter(cluster, got["RIAK_BUCKET_TYPE"], got["RIAK_BUCKET"], testKey); err != nil {
		t.Errorf("Error fetching test key/value: %v", err)
	}

	if int64(value) != res {
		t.Errorf("Fetched counter error; expected: %d; got: %d", value, res)
	}
}

func TestIntegrationInstanceUnbindingOk(t *testing.T) {
	// Prepare
	serviceTestClient := client.NewRiak(serviceITestCfg)
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instance := fmt.Sprintf("test-instance-%d", rnd.Int())
	plan := "tsuru-counter"
	appHost := "myapp.test.org"
	testKey := fmt.Sprintf("MyTestAwesomeKey_%d", rnd.Int()) // Random key always to avoid collisions between runs
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

	cfg := createIntegrationConfig()
	cfg.RiakUser = got["RIAK_USER"]
	cfg.RiakPass = got["RIAK_PASSWORD"]
	cluster, err := client.NewRiakCluster(cfg)
	if err != nil {
		t.Error(err)
	}

	// Set a key on the bucket (this should be allowed)
	value := rnd.Intn(100)
	if err = incCounter(cluster, got["RIAK_BUCKET_TYPE"], got["RIAK_BUCKET"], testKey, value); err != nil {
		t.Errorf("Error setting test key/value: %v", err)
	}
	// unbind our instance
	r, _ = http.NewRequest("DELETE", uri, nil)
	w = httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("Coudn't prepare the isntance for the test")
	}

	// This shouldn't be allowed
	// Set a key on the bucket
	value = rnd.Intn(100)
	if err = incCounter(cluster, got["RIAK_BUCKET_TYPE"], got["RIAK_BUCKET"], testKey, value); err == nil {
		t.Errorf("Should have raised permission error, it didn't")
	}

	checkError := fmt.Sprintf("Permission denied: User '%s' does not have 'riak_kv.put' on %s/%s", got["RIAK_USER"], got["RIAK_BUCKET_TYPE"], got["RIAK_BUCKET"])

	if !strings.Contains(err.Error(), checkError) {
		t.Errorf("Expected error: %s\n; got: %s", checkError, err.Error())
	}
}

func TestIntegrationInstanceStatusOk(t *testing.T) {
	// Prepare
	serviceTestClient := client.NewRiak(serviceITestCfg)
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instance := fmt.Sprintf("test-instance-%d", rnd.Int())
	plan := "tsuru-counter"
	uri := fmt.Sprintf("/resources?name=%s&plan=%s&team=myteam&user=username", instance, plan)

	// Create a new instance ok
	r, _ := http.NewRequest("POST", uri, nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("Coudn't prepare the instance for the test")
	}

	// Check instance status
	uri = fmt.Sprintf("/resources/%s/status", instance)
	wantCode := http.StatusNoContent

	r, _ = http.NewRequest("GET", uri, nil)
	w = httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != wantCode {
		t.Errorf("Error checking status; expect: %d\ngot: %d", wantCode, w.Code)
	}
}
func TestIntegrationInstanceStatusWrong(t *testing.T) {
	serviceTestClient := client.NewRiak(serviceITestCfg)
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&RiakService{Cfg: serviceITestCfg, Client: serviceTestClient})

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instance := fmt.Sprintf("test-instance-%d", rnd.Int())
	uri := fmt.Sprintf("/resources/wrong_%s/status", instance)
	wantCode := http.StatusInternalServerError

	// Check not existent instance status
	r, _ := http.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != wantCode {
		t.Errorf("Error checking status; expect: %d\ngot: %d", wantCode, w.Code)
	}
}
