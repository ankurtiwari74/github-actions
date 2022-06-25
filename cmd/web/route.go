package main

import (
	"bookings/internals/config"
	"bookings/internals/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	// return mux

	mux := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Use(middleware.Recoverer)
	mux.Use(writeToConsole)
	mux.Use(sessionLoad)
	mux.Use(NoSurf)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/luxury", handlers.Repo.Luxury)
	mux.Get("/general", handlers.Repo.General)
	mux.Get("/reservations", handlers.Repo.Reservations)
	mux.Post("/reservations", handlers.Repo.PostReservations)
	mux.Get("/search", handlers.Repo.Search)
	mux.Post("/search", handlers.Repo.PostSearch)
	mux.Post("/searchJSON", handlers.Repo.PostSearchJSON)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/search-summary", handlers.Repo.SearchSummary)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoomById)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
