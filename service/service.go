package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
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
	// Could return '/resources' but tsuru doesn't send a trailing slash at the end
	// and with prefix empty route mapping doesn't work
	return "/"
}

// Middleware wraps all the requests around thesse middlewares
func (s *RiakService) Middleware(h http.Handler) http.Handler {
	h = BasicAuthHandler(h, s.Cfg.RiakAPIUsername, s.Cfg.RiakAPIPassword)
	return h
}

// JSONMiddleware wraps all the requests around these middlewares
func (s *RiakService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return j
}

// JSONEndpoints maps the routes with the endpoints
func (s *RiakService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	logrus.Debug("Registering endpoints...")

	return map[string]map[string]server.JSONEndpoint{

		"/resources/plans": map[string]server.JSONEndpoint{
			// Returs the available plans
			"GET": s.GetPlans,
		},

		"/resources": map[string]server.JSONEndpoint{
			// Creates a service instance
			"POST": s.CreateInstance,
		},

		"/resources/{name}/bind-app": map[string]server.JSONEndpoint{
			// Binds an instance with an application
			"POST": s.BindInstance,
			// Unbinds an instance from an application
			"DELETE": s.UnbindInstance,
		},

		// (Un)?Bind events to make custom stuff
		"/resources/{name}/bind": map[string]server.JSONEndpoint{
			// Bind app to instance event
			"POST": s.BindInstanceEvent,

			// Unbind app to instance event
			"DELETE": s.UnbindInstanceEvent,
		},

		"/resources/{name}": map[string]server.JSONEndpoint{
			// Removes the instance
			"DELETE": s.RemoveInstance,
		},

		"/resources/{name}/status": map[string]server.JSONEndpoint{
			// Checks the status of the instance
			"GET": s.CheckInstanceStatus,
		},
	}
}
