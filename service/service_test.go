package service

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gizmoConfig "github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
	"gitlab.qdqmedia.com/shared-projects/riakapi/service/client"
)

var dummyConfig = &config.ServiceConfig{Riak: &config.Riak{}, Server: &gizmoConfig.Server{}}
var dummyClient = client.NewNilClient()

//func setUp() {
//
//}
//
//func TestMain(m *testing.M) {
//	setUp()
//	os.Exit(m.Run())
//}

func TestGetPlansOK(t *testing.T) {
	tests := []struct {
		givenURI    string
		givenClient client.Client
		givenConfig *config.ServiceConfig
		givenMethod string
		givenBody   io.Reader

		wantCode int
		wantBody interface{}
	}{
		{
			givenURI:    "/resources/plans",
			givenClient: dummyClient,
			givenConfig: dummyConfig,
			givenMethod: "GET",
			givenBody:   strings.NewReader(""),

			wantCode: http.StatusOK,
			wantBody: "OK",
		},
	}

	for _, test := range tests {
		// Create our dummy server (with config & client)
		srvr := server.NewSimpleServer(test.givenConfig.Server)
		srvr.Register(&RiakService{Cfg: test.givenConfig, Client: test.givenClient})

		// Create the request
		r, _ := http.NewRequest(test.givenMethod, test.givenURI, test.givenBody)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("expected response code of %d; got %d", test.wantCode, w.Code)
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil || string(b) != test.wantBody {
			t.Errorf("expected response code of %s; got %s", test.wantBody, string(b))
		}

	}
}
