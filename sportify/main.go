package main

import (
	"context"

	"github.com/TheVovchenskiy/sportify-backend/server"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	srv := server.Server{}
	baseCtx := context.Background()

	viper.OnConfigChange(func(_ fsnotify.Event) {
		err := srv.ReRun(baseCtx)
		if err != nil {
			panic(err)
		}
	})
	viper.WatchConfig()

	if err := srv.Run(baseCtx); err != nil {
		panic(err)
	}
}
