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

		// Check len because api returns different order of the slice each time
		if got != test.wantBody {
			t.Errorf("expected response body of\n%#v;\ngot\n%#v", test.wantBody, got)
		}

		// Check state on dummy client is correct
		if !reflect.DeepEqual(test.wantDummyBuckets, test.givenClient.Buckets) {
			t.Errorf("expected dummy buckets %v; \ngot: %v", test.wantDummyBuckets, test.givenClient.Buckets)
		}

	}

}
