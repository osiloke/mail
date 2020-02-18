package main

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"os"
)

func main() {
	uri := os.Getenv("redisURI")
	if uri == "" {
		uri = "redis://127.0.0.1:6379"
	}
	println("setting up recoverable machinery with broker " + uri)
	cnf := &config.Config{
		Broker:        uri,
		DefaultQueue:  "machinery_tasks",
		ResultBackend: uri,
	}
	server, err := machinery.NewServer(cnf)
	if err != nil {
		panic(err)
	}
	sig := &tasks.Signature{
		Name: "notification",
		// RetryCount: 3,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: "addonConfig",
			},
			{
				Type:  "string",
				Value: "params",
			},
			{
				Type:  "string",
				Value: "data",
			},
			{
				Type:  "string",
				Value: "trace",
			},
		},
	}
	_, err = server.SendTask(sig)
}
