package utils

import (
	"asphalt9db/models"
	"strings"
)

type Player struct {
	Credential     string          `json:"credential"`
	Name           string          `json:"name"`
	Alias          string          `json:"alias"`
	Platform       models.Platform `json:"platform"`
	AllCredentials []string        `json:"allCredentials"`
}

func MakeCredentialsMap(credentials []string) map[string]string {
	f := make(map[string]string)
	for _, credential := range credentials {
		split := strings.Split(credential, ":")
		if len(split) != 2 {
			continue
		}

		key := split[0]
		value := split[1]
		f[key] = value
	}

	return f
}
