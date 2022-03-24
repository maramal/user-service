package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID `db:"id,omitempty"`
	Email        string             `db:"email" json:"email"`
	RefreshToken string             `db:"refresh_token" json:"refresh_token"`
	UserAgent    string             `db:"user_agent" json:"user_agent"`
	ClientIP     string             `db:"client_ip" json:"client_ip"`
	IsBlocked    bool               `db:"is_blocked" json:"is_blocked"`
	CreatedAt    time.Time          `db:"created_at" json:"created_at"`
	ExpiresAt    time.Time          `db:"expires_at" json:"expires_at"`
}
