package main

import (
	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
	"gitlab.qdqmedia.com/shared-projects/riakapi/service"
	"gitlab.qdqmedia.com/shared-projects/riakapi/service/client"
)

func main() {
	// Load configuration
	cfg := config.NewServiceConfig()

	logrus.Info("Starting Riak API service...")

	server.Init("riak-api", cfg.Server)

	// Create the client
	client := client.NewNilClient()
	rkSrv := service.NewRiakService(cfg, client)
	err := server.Register(rkSrv)

	if err != nil {
		logrus.Fatalf("Unable to register service: %v", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
