package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName         string             `bson:"first_name" json:"first_name"`
	LastName          string             `bson:"last_name" json:"last_name"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"password,omitempty"`
	Type              string             `bson:"type" json:"type"`
	Status            string             `bson:"status" json:"status"`
	ProfileImage      string             `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
	PasswordChangedAt time.Time          `bson:"password_changed_at,omitempty" json:"password_changed_at,omitempty"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}
