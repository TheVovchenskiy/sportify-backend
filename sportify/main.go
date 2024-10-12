package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

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

	port := ":8080"
	fmt.Printf("listen %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		panic(err)
	}
}
