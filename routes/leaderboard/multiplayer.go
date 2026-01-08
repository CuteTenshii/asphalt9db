package leaderboard

import (
	"asphalt9db/models"
	"asphalt9db/utils"
	"encoding/json"
	"net/http"
	"strings"

	"gorm.io/gorm/clause"
)

func Multiplayer(w http.ResponseWriter, r *http.Request) {
	db := models.DB()
	lb, err := utils.GetMultiplayerTopLeaderboard(100, 0)
	if err != nil {
		panic(err)
	}

	utils.RenderHtml(w, "leaderboard/multiplayer.html", map[string]interface{}{
		"data": lb.Data,
	})

	// Process data
	for _, v := range lb.Data {
		var foxyJson struct {
			Tier     int `json:"_tier"`
			ClubData struct {
				ID         string         `json:"id"`
				Timestamp  int64          `json:"timestamp"`
				Name       string         `json:"name"`
				Logo       utils.ClubLogo `json:"logo"`
				LastEditor struct {
					Credential     string   `json:"credential"`
					AllCredentials []string `json:"allCredentials"`
				} `json:"last_editor"`
			} `json:"_club_data"`
			Platform       models.Platform `json:"_platform"`
			AllCredentials []string        `json:"_allCredentials"`
			Alias          string          `json:"_alias"`
		}
		_ = json.Unmarshal([]byte(v.FoxyJson), &foxyJson)

		club := models.Club{
			ID:          foxyJson.ClubData.ID,
			DisplayName: foxyJson.ClubData.Name,
			Score:       0,
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&club).Error; err != nil {
			panic(err)
		}

		player := models.Player{
			ID:          strings.Split(v.Credential, ":")[1], // Split the "fed_id:<uuid>" part to keep only the UUID
			Name:        v.DisplayName,
			Alias:       foxyJson.Alias,
			ClubID:      &foxyJson.ClubData.ID,
			Platform:    foxyJson.Platform,
			Credentials: utils.MakeCredentialsMap(foxyJson.AllCredentials),
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&player).Error; err != nil {
			panic(err)
		}
	}
}
