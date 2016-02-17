package service

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizationMiddleware(t *testing.T) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))

	tests := []struct {
		givenUsername   string
		givenPassword   string
		givenAuthHeader string

		wantCode int
	}{
		{
			givenUsername:   "testuser",
			givenPassword:   "testpass",
			givenAuthHeader: authHeader,
			wantCode:        http.StatusOK,
		},
		{
			givenUsername:   "wronguser",
			givenPassword:   "testpass",
			givenAuthHeader: authHeader,
			wantCode:        http.StatusUnauthorized,
		},
		{
			givenUsername:   "testuser",
			givenPassword:   "wrongpass",
			givenAuthHeader: authHeader,
			wantCode:        http.StatusUnauthorized,
		},
	}

	for _, test := range tests {

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", test.givenAuthHeader)
		res := httptest.NewRecorder()

		BasicAuthHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), test.givenUsername, test.givenPassword).ServeHTTP(res, req)

		if test.wantCode != res.Code {
			t.Errorf("Expected code wrong, want: %d; got: %d", test.wantCode, res.Code)
		}
	}

}
