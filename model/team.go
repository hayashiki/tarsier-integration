package model

import (
	"time"
)

const TeamKind = "teams"

type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" datastore:"name"`
	// secret managerでの永続化版も作成しておく
	Token     string    `json:"token" datastore:"token"`
	CreatedAt time.Time `datastore:"CreatedAt,noindex"`
	UpdatedAt time.Time `datastore:"UpdatedAt,noindex"`
}
