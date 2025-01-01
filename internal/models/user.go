package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	FullName         string             `bson:"full_name"`
	Email            string             `bson:"email"`
	Password         string             `bson:"password"`
	Verified         bool               `bson:"verified"`
	VerificationCode string             `bson:"verification_code,omitempty"`
	GoogleID         string             `bson:"google_id,omitempty"`
	Avatar           string             `bson:"avatar,omitempty"`
	Provider         string             `bson:"provider,omitempty"` // "local" or "google"
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}
