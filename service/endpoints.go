package service

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	MissingParamsMsg      = "Missing parameters"
	BucketCreationFailMsg = "Error declaring bucket type"
)

// GetPlans returns a json with the available plans on tsuru. Translated to riak,
// this are the bucket types
func (s *RiakService) GetPlans(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'GetPlans' endpoint")

	plans, err := s.Client.GetBucketTypes()

	if err != nil {
		return http.StatusInternalServerError, map[string]error{"error": err}, err
	}

	return http.StatusOK, &plans, nil
}

// CreateInstance Creates a new instance on Tsuru, this translates to a new
// bucket of the desired bucket type on Riak
func (s *RiakService) CreateInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'CreateInstance' endpoint")

	bucketName := r.URL.Query().Get("name")
	bucketType := r.URL.Query().Get("plan")

	if bucketName == "" || bucketType == "" {
		logrus.Errorf("Could not create the instance: %s", MissingParamsMsg)
		return http.StatusInternalServerError, MissingParamsMsg, nil
	}

	err := s.Client.CreateBucket(bucketName, bucketType)

	if err != nil {
		logrus.Errorf("Could not create the instance: %s", err)
		return http.StatusInternalServerError, BucketCreationFailMsg, nil
	}

	logrus.Infof("Instace '%s' created", bucketName)
	return http.StatusOK, "", nil
}

//BindInstance Binds an app to an instance on Tsuru, this translates to a new
// authentication credentias and authorization for teh desired bucket
func (s *RiakService) BindInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'BindInstance' endpoint")

	bucketName, ok := mux.Vars(r)["name"]
	if !ok {
		logrus.Errorf("Could not bind the instance: %s", MissingParamsMsg)
		return http.StatusInternalServerError, MissingParamsMsg, nil
	}

	userWord := r.URL.Query().Get("app-host")

	if userWord == "" {
		logrus.Errorf("Could not bind the instance: %s", MissingParamsMsg)
		return http.StatusInternalServerError, MissingParamsMsg, nil
	}

	// Create the user and pass (if not present already from previous instances)
	user, pass, err := s.Client.EnsureUserPresent(userWord)

	if err != nil {
		logrus.Errorf("Could not Bind the instance: %s", err)
		return http.StatusInternalServerError, BucketCreationFailMsg, nil

	}

	// Grant access on bucket
	err = s.Client.GrantUserAccess(user, bucketName)
	if err != nil {
		logrus.Errorf("Could not Bind the instance: %s", err)
		return http.StatusInternalServerError, BucketCreationFailMsg, nil
	}

	// The required env vars
	envVars := map[string]string{
		"RIAK_HOST":        s.Cfg.RiakHost,
		"RIAK_HTTP_PORT":   strconv.Itoa(8098),
		"RIAK_PB_PORT":     strconv.Itoa(8087),
		"RIAK_USER":        user,
		"RIAK_PASSWORD":    pass,
		"RIAK_BUCKET_TYPE": s.Client.GetBucketType(bucketName),
		"RIAK_BUCKET":      bucketName,
	}

	return http.StatusCreated, envVars, nil
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
