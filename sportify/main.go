package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RunTgHandler(handler Handler) error {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Post("/message", handler.TryCreateEvent)

	port := ":8090"
	fmt.Printf("listen bot input %s\n", port)

	return http.ListenAndServe(port, r)
}

func main() {
	simpleEventStorage, err := NewSimpleEventStorage()
	if err != nil {
		panic(err)
	}

	handler := Handler{Storage: simpleEventStorage}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get("/events", handler.GetEvents)
	r.Get("/event/{id}", handler.GetEvent)
	r.Put("/event/sub/{id}", handler.SubscribeEvent)

	r.Get("/img/*", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/img/" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		fs := http.StripPrefix("/img/", http.FileServer(http.Dir("./photos")))

		fs.ServeHTTP(w, r)
	})

	go func() {
		if err := RunTgHandler(handler); err != nil {
			panic(err)
		}
	}()

	port := ":8080"
	fmt.Printf("listen %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		panic(err)
	}
}
