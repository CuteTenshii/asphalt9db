package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/carlmjohnson/requests"
)

func GetAccessToken() (string, error) {
	cacheKey := "as9:access_token"
	rdb := GetRedis()
	ctx := context.Background()

	if val, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		return val, nil
	}

	steamId := os.Getenv("STEAM_ID")
	if steamId == "" {
		panic("No STEAM_ID environment variable provided")
	}

	deviceInformation := map[string]string{
		"id":         "",
		"model":      "System Product Name",
		"country":    "US",
		"language":   "en",
		"firmware":   "10.0.19045.5796",
		"resolution": "1920.0x1080.0",
	}
	deviceInfoJson, _ := json.Marshal(deviceInformation)
	body := url.Values{}
	body.Set("client_id", "asphalt:6811:85636:48.0.5a:windows:steam")
	body.Set("device_information", string(deviceInfoJson))
	body.Set("encrypt_tokens", "false")
	body.Set("password", "")
	body.Set("scope", "auth storage_ro lobby message transaction social_group_ro chat_subscribe leaderboard_ro")

	var resp AuthorizeResponse
	err := requests.
		URL(fmt.Sprintf("https://fed.gold0028.gameloft.com/v2/users/steam%%3A%s/authorize", steamId)).
		BodyForm(body).
		UserAgent("JWeb").
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return "", err
	}
	rdb.Set(ctx, cacheKey, resp.AccessToken, 0)

	return resp.AccessToken, nil
}

func MustGetAccessToken() string {
	auth, err := GetAccessToken()
	if err != nil {
		panic(err)
	}

	return auth
}

type AuthorizeResponse struct {
	AccessToken string `json:"access_token"`
}
