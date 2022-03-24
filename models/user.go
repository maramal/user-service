package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `db:"id,omitempty" json:"id,omitempty"`
	FirstName         string             `db:"first_name" json:"first_name"`
	LastName          string             `db:"last_name" json:"last_name"`
	Email             string             `db:"email" json:"email"`
	Password          string             `db:"password" json:"password,omitempty"`
	Type              string             `db:"type" json:"type"`
	Status            string             `db:"status" json:"status"`
	ProfileImage      string             `db:"profile_image,omitempty" json:"profile_image,omitempty"`
	PasswordChangedAt time.Time          `db:"password_changed_at,omitempty" json:"password_changed_at,omitempty"`
	CreatedAt         time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `db:"updated_at" json:"updated_at"`
}
