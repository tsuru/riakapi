package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	gizmoConfig "github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
	"gitlab.qdqmedia.com/shared-projects/riakapi/service/client"
)

var serviceTestCfg = &config.ServiceConfig{Riak: &config.Riak{}, Server: &gizmoConfig.Server{}}

//func setUp() {
//
//}
//
//func TestMain(m *testing.M) {
//	setUp()
//	os.Exit(m.Run())
//}

func TestGetPlans(t *testing.T) {

	serviceTestClient := client.NewDummy()

	correctPlans := []interface{}{
		map[string]interface{}{"name": client.BucketTypeCounter, "description": "Bucket type of counter data type"},
		map[string]interface{}{"name": client.BucketTypeSet, "description": "Bucket type of set data type"},
		map[string]interface{}{"name": client.BucketTypeMap, "description": "Bucket type of map data type"},
	}

	tests := []struct {
		givenURI    string
		givenClient *client.Dummy
		givenConfig *config.ServiceConfig
		givenMethod string

		wantCode int
		wantBody interface{}
	}{
		{
			givenURI:    "/resources/plans",
			givenClient: serviceTestClient,
			givenConfig: serviceTestCfg,
			givenMethod: "GET",

			wantCode: http.StatusOK,
			wantBody: correctPlans,
		},
	}

	for _, test := range tests {
		// Create our dummy server (with config & client)
		srvr := server.NewSimpleServer(nil)
		srvr.Register(&RiakService{Cfg: test.givenConfig, Client: test.givenClient})

		// Create the request
		r, _ := http.NewRequest(test.givenMethod, test.givenURI, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("expected response code of %d; got %d", test.wantCode, w.Code)
		}

		var got interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error("unable to JSON decode response body: ", err)
		}

		// Check len because api returns different order of the slice each time
		if len(got.([]interface{})) != len(test.wantBody.([]interface{})) {
			t.Errorf("expected response body of\n%#v;\ngot\n%#v", test.wantBody, got)
		}

	}
}

func TestInstanceCreation(t *testing.T) {
	serviceTestClient := client.NewDummy()

	tests := []struct {
		givenURI          string
		givenClient       *client.Dummy
		givenConfig       *config.ServiceConfig
		givenMethod       string
		givenDummyBuckets map[string]string

		wantCode         int
		wantBody         interface{}
		wantDummyBuckets map[string]string
	}{
		{
			givenURI:          "/resources",
			givenClient:       serviceTestClient,
			givenConfig:       serviceTestCfg,
			givenMethod:       "POST",
			givenDummyBuckets: map[string]string{},

			wantCode:         http.StatusInternalServerError,
			wantBody:         MissingParamsMsg,
			wantDummyBuckets: map[string]string{},
		},
		{
			givenURI:          "/resources?name=test-bucket&plan=wrong&team=myteam&user=username",
			givenClient:       serviceTestClient,
			givenConfig:       serviceTestCfg,
			givenMethod:       "POST",
			givenDummyBuckets: map[string]string{},

			wantCode:         http.StatusInternalServerError,
			wantBody:         BucketCreationFailMsg,
			wantDummyBuckets: map[string]string{},
		},
		{
			givenURI:          "/resources?name=test-bucket&plan=tsuru-counter&team=myteam&user=username",
			givenClient:       serviceTestClient,
			givenConfig:       serviceTestCfg,
			givenMethod:       "POST",
			givenDummyBuckets: map[string]string{},

			wantCode:         http.StatusOK,
			wantBody:         "",
			wantDummyBuckets: map[string]string{"test-bucket": "tsuru-counter"},
		},
		{ // Same test as previous one, will conflict the name
			givenURI:          "/resources?name=test-bucket&plan=tsuru-counter&team=myteam&user=username",
			givenClient:       serviceTestClient,
			givenConfig:       serviceTestCfg,
			givenMethod:       "POST",
			givenDummyBuckets: map[string]string{"test-bucket": "tsuru-counter"},

			wantCode:         http.StatusInternalServerError,
			wantBody:         BucketCreationFailMsg,
			wantDummyBuckets: map[string]string{"test-bucket": "tsuru-counter"},
		},
	}

	for _, test := range tests {

		// Set our initial state of the database
		test.givenClient.Buckets = test.givenDummyBuckets

		// Create our dummy server (with config & client)
		srvr := server.NewSimpleServer(nil)
		srvr.Register(&RiakService{Cfg: test.givenConfig, Client: test.givenClient})

		// Create the request
		r, _ := http.NewRequest(test.givenMethod, test.givenURI, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("expected response code of %d; got %d", test.wantCode, w.Code)
		}

		var got interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error("unable to JSON decode response body: ", err)
		}

		// Check body
		if got != test.wantBody {
			t.Errorf("expected response body of\n%#v;\ngot\n%#v", test.wantBody, got)
		}

		// Check state on dummy client is correct
		if !reflect.DeepEqual(test.wantDummyBuckets, test.givenClient.Buckets) {
			t.Errorf("expected dummy buckets %v; \ngot: %v", test.wantDummyBuckets, test.givenClient.Buckets)
		}

	}

}

func TestInstanceBindingOK(t *testing.T) {
	serviceTestClient := client.NewDummy()

	tests := []struct {
		givenURI          string
		givenClient       *client.Dummy
		givenConfig       *config.ServiceConfig
		givenMethod       string
		givenDummyUsers   map[string]*client.UserProps
		givenDummyBuckets map[string]string

		wantCode       int
		wantBody       map[string]string
		wantDummyUsers map[string]*client.UserProps
	}{
		{ // Check new app binding
			givenURI:          "/resources/testinstance/bind-app?app-host=myapp.tsuru.io",
			givenClient:       serviceTestClient,
			givenConfig:       serviceTestCfg,
			givenMethod:       "POST",
			givenDummyUsers:   map[string]*client.UserProps{},
			givenDummyBuckets: map[string]string{"testinstance": "testbuckettype"},

			wantCode: http.StatusCreated,
			wantBody: map[string]string{
				"RIAK_BUCKET":      "testinstance",
				"RIAK_BUCKET_TYPE": "testbuckettype",
				"RIAK_HOST":        "",
				"RIAK_HTTP_PORT":   "8098",
				"RIAK_PASSWORD":    "myapp.tsuru.io",
				"RIAK_PB_PORT":     "8087",
				"RIAK_USER":        "tsuru_myapp.tsuru.io",
			},
			wantDummyUsers: map[string]*client.UserProps{
				"tsuru_myapp.tsuru.io": &client.UserProps{
					Username: "tsuru_myapp.tsuru.io",
					Password: "myapp.tsuru.io",
					ACL:      []string{"testinstance"},
				},
			},
		},
		{ // Check when already binded
			givenURI:    "/resources/testinstance/bind-app?app-host=myapp.tsuru.io",
			givenClient: serviceTestClient,
			givenConfig: serviceTestCfg,
			givenMethod: "POST",
			givenDummyUsers: map[string]*client.UserProps{
				"tsuru_myapp.tsuru.io": &client.UserProps{
					Username: "tsuru_myapp.tsuru.io",
					Password: "myapp.tsuru.io",
					ACL:      []string{"testinstance"},
				},
			},
			givenDummyBuckets: map[string]string{"testinstance": "testbuckettype"},

			wantCode: http.StatusCreated,
			wantBody: map[string]string{
				"RIAK_BUCKET":      "testinstance",
				"RIAK_BUCKET_TYPE": "testbuckettype",
				"RIAK_HOST":        "",
				"RIAK_HTTP_PORT":   "8098",
				"RIAK_PASSWORD":    "myapp.tsuru.io",
				"RIAK_PB_PORT":     "8087",
				"RIAK_USER":        "tsuru_myapp.tsuru.io",
			},
			wantDummyUsers: map[string]*client.UserProps{
				"tsuru_myapp.tsuru.io": &client.UserProps{
					Username: "tsuru_myapp.tsuru.io",
					Password: "myapp.tsuru.io",
					ACL:      []string{"testinstance"},
				},
			},
		},
	}

	for _, test := range tests {

		// Set our initial state of the database
		test.givenClient.Users = test.givenDummyUsers
		test.givenClient.Buckets = test.givenDummyBuckets

		// Create our dummy server (with config & client)
		srvr := server.NewSimpleServer(nil)
		srvr.Register(&RiakService{Cfg: test.givenConfig, Client: test.givenClient})

		// Create the request
		r, _ := http.NewRequest(test.givenMethod, test.givenURI, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("expected response code of %d; got %d", test.wantCode, w.Code)
		}

		var got map[string]string
		err := json.NewDecoder(w.Body).Decode(&got)

		// Dont fail the test if the decode fails when the unmarshalling is for
		// invalid API calls
		if err != nil {
			t.Error("unable to JSON decode response body: ", err)
		}

		// Check body json decoded
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("expected response body of\n%#v;\ngot\n%#v", test.wantBody, got)
		}

		// Check state on dummy client is correct
		if len(test.wantDummyUsers) != len(test.givenClient.Users) {
			t.Errorf("expected dummy users %v; \ngot: %v", test.wantDummyUsers, test.givenClient.Users)
		}
		for k, v := range test.wantDummyUsers {
			u := test.givenClient.Users[k]
			if v.Username != u.Username || v.Password != u.Password || len(v.ACL) != len(u.ACL) {
				t.Errorf("expected dummy user %v; \ngot: %v", *v, *u)
			}
		}

	}
}

func TestInstanceBindingWrong(t *testing.T) {
	serviceTestClient := client.NewDummy()

	tests := []struct {
		givenURI    string
		givenClient *client.Dummy
		givenConfig *config.ServiceConfig
		givenMethod string

		wantCode int
		wantBody string
	}{
		{
			givenURI:    "/resources/testinstance/bind-app",
			givenClient: serviceTestClient,
			givenConfig: serviceTestCfg,
			givenMethod: "POST",

			wantCode: http.StatusInternalServerError,
			wantBody: MissingParamsMsg,
		},
	}

	for _, test := range tests {
		srvr := server.NewSimpleServer(nil)
		srvr.Register(&RiakService{Cfg: test.givenConfig, Client: test.givenClient})

		// Create the request
		r, _ := http.NewRequest(test.givenMethod, test.givenURI, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("expected response code of %d; got %d", test.wantCode, w.Code)
		}

		var got string
		json.NewDecoder(w.Body).Decode(&got)

		if got != test.wantBody {
			t.Errorf("Expected body: %s ; got: %s", test.wantBody, got)
		}
	}
}
