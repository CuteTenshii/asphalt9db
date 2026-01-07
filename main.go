package main

import (
	"asphalt9db/routes"
	"asphalt9db/routes/leaderboard"
	"asphalt9db/routes/players"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	r.Get("/", routes.Home)
	r.Get("/leaderboard/clubs", leaderboard.Clubs)
	r.Get("/leaderboard/multiplayer", leaderboard.Multiplayer)
	r.Get("/players/{id}", players.Player)

	http.ListenAndServe(":3000", r)
}
