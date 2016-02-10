package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
	"gitlab.qdqmedia.com/shared-projects/riakapi/service/client"
)

// RiakService expose tsuru api for riak service
type RiakService struct {

	// Application configuration
	Cfg *config.ServiceConfig

	// resources client (normally riak)
	Client client.Client
}

// NewRiakService creates a new services ready to register on the server
func NewRiakService(c *config.ServiceConfig, client client.Client) *RiakService {
	logrus.Debug("New riak service created")
	return &RiakService{
		Cfg:    c,
		Client: client,
	}
}

// Prefix returns the url prefix for all the endpoints of this service
func (s *RiakService) Prefix() string {
	return "/resources"
}

// Middleware wraps all the requests around this middlewares
func (s *RiakService) Middleware(h http.Handler) http.Handler {
	return h
}

// Endpoints maps the routes with the endpoints
func (s *RiakService) Endpoints() map[string]map[string]http.HandlerFunc {
	logrus.Debug("Registering endpoints...")

	return map[string]map[string]http.HandlerFunc{

		"/plans": map[string]http.HandlerFunc{
			// Returs the available plans
			"GET": s.GetPlans,
		},

		"/": map[string]http.HandlerFunc{
			// Creates a service instance
			"POST": s.CreateInstance,
		},

		"/{name}/bind-app": map[string]http.HandlerFunc{
			// Binds an instance with an application
			"POST": s.BindInstance,
			// Unbinds an instance from an application
			"DELETE": s.UnbindInstance,
		},

		// (Un)?Bind events to make custom stuff
		"/{name}/bind": map[string]http.HandlerFunc{
			// Bind app to instance event
			"PUT": s.BindInstanceEvent,

			// Unbind app to instance event
			"DELETE": s.UnbindInstanceEvent,
		},

		"/{name}": map[string]http.HandlerFunc{
			// Removes the instance
			"DELETE": s.RemoveInstance,
		},

		"/{name}/status": map[string]http.HandlerFunc{
			// Checks the status of the instance
			"GET": s.CheckInstanceStatus,
		},
	}
}
