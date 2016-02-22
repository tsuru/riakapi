package main

import (
	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"

	"github.com/tsuru/riakapi/config"
	"github.com/tsuru/riakapi/service"
	"github.com/tsuru/riakapi/service/client"
)

func main() {
	// Load configuration
	cfg := config.NewServiceConfig()

	logrus.Info("Starting Riak API service...")

	server.Init("riak-api", cfg.Server)

	// Create the client
	//client := client.NewDummy()
	client := client.NewRiak(cfg)
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
