package leaderboard

import (
	"asphalt9db/models"
	"asphalt9db/utils"
	"net/http"
	"strings"

	"gorm.io/gorm/clause"
)

func Clubs(w http.ResponseWriter, r *http.Request) {
	db := models.DB()
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
}
