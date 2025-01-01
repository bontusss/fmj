package auth

import (
	"context"
	"fmj/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	VerifyUser(ctx context.Context, code string) error
	FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
}

type repository struct {
	db  *mongo.Database
	ctx context.Context
}

func (r repository) FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (r repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(r.ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r repository) VerifyUser(ctx context.Context, code string) error {
	_, err := r.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"verification_code": code},
		bson.M{"$set": bson.M{"verified": true, "verification_code": ""}},
	)
	return err
}

func NewRepository(db *mongo.Database, ctx context.Context) Repository {
	return &repository{db, ctx}
}
