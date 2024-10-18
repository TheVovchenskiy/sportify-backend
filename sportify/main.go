package main

import (
	"context"
	"flag"

	"github.com/TheVovchenskiy/sportify-backend/server"
)

func main() {
	configFile := flag.String(
		"configfile",
		"",
		"you should use file path to config file. NOT SECURE example: ./config.example.yaml",
	)
	flag.Parse()

	if configFile == nil {
		panic("flag --configfile is nil")
	}

	srv := server.Server{}
	baseCtx := context.Background()

	if err := srv.Run(baseCtx, *configFile); err != nil {
		panic(err)
	}
}
