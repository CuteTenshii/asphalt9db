package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/carlmjohnson/requests"
)

func GetAccessToken() (string, error) {
	cacheKey := "as9:access_token"
	rdb := GetRedis()
	ctx := context.Background()

	if val, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		return val, nil
	}

	anonymousFed := os.Getenv("GAMELOFT_ANONYMOUS_FED")
	if anonymousFed == "" {
		panic("No GAMELOFT_ANONYMOUS_FED environment variable provided")
	}
	anonymousPwd := os.Getenv("GAMELOFT_ANONYMOUS_PASSWORD")
	if anonymousPwd == "" {
		panic("No GAMELOFT_ANONYMOUS_PASSWORD environment variable provided")
	}

	deviceInformation := map[string]string{
		"id":         "2008971685286699808", // Any value is accepted, but you can get yours in C:\Users\<username>\Documents\Gameloft\Asphalt 9 Legends\device_launch_info
		"model":      "System Product Name",
		"country":    "US",
		"language":   "en",
		"firmware":   "15",
		"resolution": "1920.0x1080.0",
	}
	deviceInfoJson, _ := json.Marshal(deviceInformation)
	body := url.Values{}
	body.Set("client_id", "asphalt:6811:85636:48.0.5a:android:googleplay")
	body.Set("device_information", string(deviceInfoJson))
	body.Set("encrypt_tokens", "false")
	body.Set("password", anonymousPwd)
	body.Set("scope", "auth storage_ro lobby message transaction social_group_ro chat_subscribe leaderboard_ro")

	var resp AuthorizeResponse
	err := requests.
		URL(fmt.Sprintf("https://fed.gold0028.gameloft.com/v2/users/anonymous%%3A%s/authorize", anonymousFed)).
		BodyForm(body).
		UserAgent("JWeb").
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return "", err
	}
	rdb.Set(ctx, cacheKey, resp.AccessToken, 4*time.Hour)

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
