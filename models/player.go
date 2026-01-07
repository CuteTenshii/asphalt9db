package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Platform int

const (
	PlatformSteam Platform = iota
	PlatformNintendoSwitch
	PlatformAndroid
	PlatformIOS
)

type CredentialsMap map[string]string

func (c *CredentialsMap) Scan(value interface{}) error {
	if value == nil {
		*c = make(map[string]string)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	return json.Unmarshal(bytes, c)
}

// Value implements driver.Value interface
func (c CredentialsMap) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

type Player struct {
	ID          string
	Name        string
	Alias       string
	ClubID      *string
	Club        *Club
	Platform    Platform
	Credentials CredentialsMap `gorm:"type:json"`
}
