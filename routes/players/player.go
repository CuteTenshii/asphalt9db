package players

import (
	"asphalt9db/models"
	"asphalt9db/utils"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Player(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := models.DB()

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
}
