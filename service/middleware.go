package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

// BasicAuthHandler checks if the request is authorized
func BasicAuthHandler(h http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if password != "" { // check authentication disabled
			// Check access
			if reqUser, reqPass, ok := r.BasicAuth(); ok {
				// Wrong password and/or user
				if reqUser != username || reqPass != password {
					logrus.Error("Not authorized access")
					http.Error(w, "Login Required", http.StatusUnauthorized)
					return
				}
			}
		}
		// all good
		h.ServeHTTP(w, r)
	})
}
