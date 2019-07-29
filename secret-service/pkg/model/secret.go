package model

import "time"

type Secret struct {
	Hash           string
	SecretText     string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	RemainingViews int
}
