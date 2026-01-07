package main

import (
	"asphalt9db/models"
	"asphalt9db/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	r := chi.NewRouter()
	db := models.DB()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderHtml(w, "home.html", nil)
	})
	r.Get("/leaderboard/clubs", func(w http.ResponseWriter, r *http.Request) {
		lb, err := utils.GetClubsLeaderboard(100, 0)
		if err != nil {
			panic(err)
		}

		utils.RenderHtml(w, "leaderboard/clubs.html", map[string]interface{}{
			"data": lb.Data,
		})

		// Process data
		for _, v := range lb.Data {
			club := models.Club{
				ID:          v.ID,
				DisplayName: v.DisplayName,
				Score:       int64(v.Score),
			}
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&club).Error; err != nil {
				panic(err)
			}

			player := models.Player{
				ID:          strings.Split(v.LastEditor.Credential, ":")[1], // Split the "fed_id:<uuid>" part to keep only the UUID
				Name:        v.LastEditor.Name,
				Alias:       v.LastEditor.Alias,
				ClubID:      &v.ID,
				Platform:    v.LastEditor.Platform,
				Credentials: utils.MakeCredentialsMap(v.LastEditor.AllCredentials),
			}
			if err := db.Create(&player).Error; err != nil {
				panic(err)
			}
		}
	})

	r.Get("/players/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var player models.Player
		if err := db.Select("id", "name", "club_id", "credentials").Where("id", id).First(&player).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				utils.RenderHtml(w, "errors/404.html", nil)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderHtml(w, "errors/500.html", nil)
			return
		}

		utils.RenderHtml(w, "players/player.html", map[string]interface{}{
			"data": player,
		})
	})

	http.ListenAndServe(":3000", r)
}
