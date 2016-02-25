package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/tsuru/riakapi/utils"
)

const (
	// MissingParamsMsg message when parameters are required
	MissingParamsMsg = "Missing parameters"
	// BucketCreationFailMsg message when bucket creation fails
	BucketCreationFailMsg = "Error declaring bucket type"
	// ErrorBucketStatusMsg message when bucket status is in error state
	ErrorBucketStatusMsg = "Bucket error"
	// UserGrantingFailMsg message when granting access to users fails
	UserGrantingFailMsg = "Error granting user"
	// UserRevokingFailMsg message when revoking access to users fails
	UserRevokingFailMsg = "Error revoking user"
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

	bucketName, _ := mux.Vars(r)["name"]
	userWord := r.URL.Query().Get("app-host")
	if userWord == "" {
		logrus.Errorf("Could not bind the instance: %s", MissingParamsMsg)
		return http.StatusInternalServerError, MissingParamsMsg, nil
	}

	// Create the user and pass (if not present already from previous instances)
	user, pass, err := s.Client.EnsureUserPresent(userWord)

	if err != nil {
		logrus.Errorf("Could not Bind the instance: %s", err)
		return http.StatusInternalServerError, UserGrantingFailMsg, nil

	}

	// Grant access on bucket
	err = s.Client.GrantUserAccess(user, bucketName)
	if err != nil {
		logrus.Errorf("Could not Bind the instance: %s", err)
		return http.StatusInternalServerError, UserGrantingFailMsg, nil
	}

	rHosts, err := json.Marshal(s.Cfg.RiakClusterHosts)
	if err != nil {
		logrus.Errorf("Could not Bind the instance: %s", err)
		return http.StatusInternalServerError, UserGrantingFailMsg, nil
	}

	// The required env vars
	envVars := map[string]string{
		"RIAK_HOSTS":       string(rHosts),
		"RIAK_HTTP_PORT":   strconv.Itoa(s.Cfg.RiakHTTPPort),
		"RIAK_PB_PORT":     strconv.Itoa(s.Cfg.RiakPBPort),
		"RIAK_USER":        user,
		"RIAK_PASSWORD":    pass,
		"RIAK_BUCKET_TYPE": s.Client.GetBucketType(bucketName),
		"RIAK_BUCKET":      bucketName,
	}

	// If there is a certificate then set
	if s.Cfg.RiakRootCaCert != "" {
		envVars["RIAK_ROOT_CA_CERT"] = s.Cfg.RiakRootCaCert
	}

	logrus.Infof("Instace '%s' binded to '%s'", bucketName, userWord)
	return http.StatusCreated, envVars, nil
}

// UnbindInstance Unbinds the instance from the app on Tsuru, this translates to
// remove credentials from the desired bucket
func (s *RiakService) UnbindInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'UnbindInstance' endpoint")

	bucketName, _ := mux.Vars(r)["name"]
	userWord := r.URL.Query().Get("app-host")
	if userWord == "" {
		logrus.Errorf("Could not unbind the instance: %s", MissingParamsMsg)
		return http.StatusInternalServerError, MissingParamsMsg, nil
	}
	// Revoke access to the user
	username := utils.GenerateUsername(userWord)
	err := s.Client.RevokeUserAccess(username, bucketName)

	// TODO: Delete user
	// NOTE: Keep track of users instances and delete on last one

	if err != nil {
		logrus.Errorf("Could not unbind the instance: %s", err)
		return http.StatusInternalServerError, UserRevokingFailMsg, nil
	}

	logrus.Infof("Instace '%s' unbinded from '%s'", bucketName, userWord)
	return http.StatusOK, "", nil
}

// BindInstanceEvent Processes the event from tsuru when an app is binded to a service instance
func (s *RiakService) BindInstanceEvent(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'BindInstanceEvent' endpoint (no need to implement)")
	return http.StatusCreated, "", nil
}

// UnbindInstanceEvent Processes the event from tsuru when an app is unbinded from a service instance
func (s *RiakService) UnbindInstanceEvent(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'UnbindInstanceEvent' endpoint (no need to implement)")
	return http.StatusOK, "", nil
}

// RemoveInstance Remove instance Removes the instance from tsuru. Translated to riak,  delete
// all the keys from the bucket (causing bucket deletion) -> not a good choice, not deleting bucket
// Bucket will persist 'forever'
func (s *RiakService) RemoveInstance(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'RemoveInstance' endpoint")
	return http.StatusOK, "", nil
}

// CheckInstanceStatus Checks the status of an instance on tsuru. Translated to riak,
// Checks the status of the bucket
func (s *RiakService) CheckInstanceStatus(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Executing 'CheckInstanceStatus' endpoint")

	bucketName, _ := mux.Vars(r)["name"]
	ok, err := s.Client.IsAlive(bucketName)
	logrus.Infof("Instace '%s' status ok: %t", bucketName, ok)
	if ok {
		return http.StatusNoContent, nil, nil
	}
	logrus.Errorf("Bucket error: %v", err)
	return http.StatusInternalServerError, ErrorBucketStatusMsg, nil
}
