package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

// GetPlans returns a json with the available plans on tsuru. Translated to riak,
// this are the data types
func (s *RiakService) GetPlans(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'GetPlans' endpoint")

	plans, err := s.Client.GetDataTypes()

	if err != nil {
		return http.StatusInternalServerError, map[string]error{"error": err}, err
	}

	return http.StatusOK, plans, nil
}

// CreateInstance Creates a new instance on Tsuru, this translates to a new
// bucket of the desired bucket type on Riak
func (s *RiakService) CreateInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'CreateInstance' endpoint")
	return http.StatusNotImplemented, nil, nil
}

//BindInstance Binds an app to an instance on Tsuru, this translates to a new
// authentication credentias and authorization for teh desired bucket
func (s *RiakService) BindInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'BindInstance' endpoint")
	return http.StatusNotImplemented, nil, nil
}

// UnbindInstance Unbinds the instance from the app on Tsuru, this translates to
// remove credentials from the desired bucket
func (s *RiakService) UnbindInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'UnbindInstance' endpoint")
	return http.StatusNotImplemented, nil, nil
}

// BindInstanceEvent Processes the event from tsuru when an app is binded to a service instance
func (s *RiakService) BindInstanceEvent(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'BindInstanceEvent' endpoint")
	return http.StatusNotImplemented, nil, nil
}

// UnbindInstanceEvent Processes the event from tsuru when an app is unbinded from a service instance
func (s *RiakService) UnbindInstanceEvent(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'UnbindInstanceEvent' endpoint")
	return http.StatusNotImplemented, nil, nil
}

// RemoveInstance Remove instance Removes the instance from tsuru. Translated to riak,  delete
// all the keys from the bucket (causing bucket deletion)
func (s *RiakService) RemoveInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'RemoveInstance' endpoint")
	return http.StatusNotImplemented, nil, nil
}

// CheckInstanceStatus Checks the status of an instance on tsuru. TRanslated to riak,
// Checks teh status of the bucket
func (s *RiakService) CheckInstanceStatus(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'CheckInstanceStatus' endpoint")
	return http.StatusNotImplemented, nil, nil
}
