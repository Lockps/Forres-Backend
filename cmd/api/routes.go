package main

import (
	"net/http"

	"github.com/Lockps/Forres-release-version/cmd/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.FirstPage)
	mux.Post("/getdata", app.postDataHandler)
	mux.Get("/read", app.ReadFile)
	mux.Post("/post", database.CreateUsers)
	mux.Post("/test", database.ValidateUserHandler)
	mux.Get("/test", database.WithJWTAuth(database.Test))

	mux.Post("/testpostbooking", database.AddBookingToDB)
	mux.Get("/testgetbooking", database.GetUnAvaliableSeat)

	return mux
}
