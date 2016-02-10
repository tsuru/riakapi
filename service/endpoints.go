package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

// GetPlans returns a json with the available plans on tsuru. Translated to riak,
// this are the bucket types
func (s *RiakService) GetPlans(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'GetPlans' endpoint")
}

// CreateInstance Creates a new instance on Tsuru, this translates to a new
// bucket of the desired bucket type on Riak
func (s *RiakService) CreateInstance(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'CreateInstance' endpoint")

}

//BindInstance Binds an app to an instance on Tsuru, this translates to a new
// authentication credentias and authorization for teh desired bucket
func (s *RiakService) BindInstance(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'BindInstance' endpoint")

}

// UnbindInstance Unbinds the instance from the app on Tsuru, this translates to
// remove credentials from the desired bucket
func (s *RiakService) UnbindInstance(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'UnbindInstance' endpoint")
}

// BindInstanceEvent Processes the event from tsuru when an app is binded to a service instance
func (s *RiakService) BindInstanceEvent(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'BindInstanceEvent' endpoint")
}

// UnbindInstanceEvent Processes the event from tsuru when an app is unbinded from a service instance
func (s *RiakService) UnbindInstanceEvent(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'UnbindInstanceEvent' endpoint")
}

// RemoveInstance Remove instance Removes the instance from tsuru. Translated to riak,  delete
// all the keys from the bucket (causing bucket deletion)
func (s *RiakService) RemoveInstance(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'RemoveInstance' endpoint")
}

// CheckInstanceStatus Checks the status of an instance on tsuru. TRanslated to riak,
// Checks teh status of the bucket
func (s *RiakService) CheckInstanceStatus(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Executing 'CheckInstanceStatus' endpoint")
}
