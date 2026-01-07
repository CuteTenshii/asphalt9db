package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/carlmjohnson/requests"
)

func GetLeaderboard[T any](uuid string, limit, offset int) (*LeaderboardResponse[T], error) {
	rdb := GetRedis()
	ctx := context.Background()
	if val, err := rdb.Get(ctx, fmt.Sprintf("leaderboard:%s:%d", uuid, limit)).Result(); err == nil {
		var j LeaderboardResponse[T]
		err := json.Unmarshal([]byte(val), &j)
		if err != nil {
			return nil, err
		}
		return &j, nil
	}

	var resp *LeaderboardResponse[T]
	err := requests.
		URL(fmt.Sprintf("https://fed.gold0028.gameloft.com/v1/leaderboards/desc/Leaderboard_%s", uuid)).
		Param("limit", strconv.Itoa(limit)).
		Param("offset", strconv.Itoa(offset)).
		UserAgent("JWeb").
		Header("Access-Token", MustGetAccessToken()).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	j, _ := json.Marshal(resp)
	rdb.Set(ctx, fmt.Sprintf("leaderboard:%s:%d", uuid, limit), j, 1*time.Hour)

	return resp, nil
}

func GetClubsLeaderboard(limit, offset int) (*LeaderboardResponse[LeaderboardClubResponse], error) {
	return GetLeaderboard[LeaderboardClubResponse]("8321aa3c-d5f7-11f0-b8b4-b8ca3a634708", limit, offset)
}

type LeaderboardResponse[T any] struct {
	Data         []T    `json:"data"`
	ID           string `json:"id"`
	Created      string `json:"created"`
	TotalEntries int64  `json:"total_entries"`
}

type LeaderboardClubResponse struct {
	ID          string   `json:"id"`
	LastEditor  Player   `json:"_lastEditor"`
	Credential  string   `json:"credential"`
	DisplayName string   `json:"display_name"`
	Alias       string   `json:"alias"`
	ClanSize    string   `json:"clan_size"`
	ClanID      string   `json:"clan_id"`
	Logo        ClubLogo `json:"_logo"`
	LastUpdate  string   `json:"last_update"`
	Score       float64  `json:"score"`
	Rank        float64  `json:"rank"`
}
