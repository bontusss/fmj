package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmj/internal/email"
	"fmj/internal/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/oauth2/v2"
	"log"
)

type Service interface {
	Register(ctx context.Context, fullName, email, password string) error
	Login(email, password string) (*models.User, error)
	VerifyEmail(ctx context.Context, code string) error
	HandleGoogleLogin(ctx context.Context, googleUser *oauth2.Userinfo) (*models.User, error)
}

type service struct {
	repo  Repository
	email email.Service
}

type GoogleUser struct {
	ID       string
	Email    string
	Name     string
	Picture  string
	Verified bool
}

func (s *service) HandleGoogleLogin(ctx context.Context, googleUser *oauth2.Userinfo) (*models.User, error) {
	// Check if user exists by Google ID
	existingUser, err := s.repo.FindUserByGoogleID(ctx, googleUser.Id)
	if err == nil {
		return existingUser, nil
	}

	// Check if user exists by email
	existingUser, err = s.repo.FindUserByEmail(googleUser.Email)
	if err == nil {
		// Link Google account to existing user
		existingUser.GoogleID = googleUser.Id
		existingUser.Avatar = googleUser.Picture
		existingUser.Provider = "google"
		if err := s.repo.UpdateUser(ctx, existingUser); err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// Create new user
	user := &models.User{
		FullName: googleUser.Name,
		Email:    googleUser.Email,
		GoogleID: googleUser.Id,
		Avatar:   googleUser.Picture,
		Provider: "google",
		Verified: true, // Google accounts are pre-verified
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Send welcome email
	go func() {
		err := s.email.SendWelcomeEmail(user.Email, user.FullName)
		if err != nil {
			log.Fatal("email send welcome email error:", err)
		}
	}()

	return user, nil
}

func (s *service) Register(ctx context.Context, fullName, email, password string) error {
	// Check if user exists
	existing, _ := s.repo.FindUserByEmail(email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Generate verification code
	verificationCode := make([]byte, 32)
	if _, err := rand.Read(verificationCode); err != nil {
		return err
	}
	code := hex.EncodeToString(verificationCode)

	// Create user
	user := &models.User{
		FullName:         fullName,
		Email:            email,
		Password:         string(hashedPassword),
		Verified:         false,
		VerificationCode: code,
	}

	fmt.Printf("creating new user: %s\n", user.FullName)
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	// Send verification email
	fmt.Printf("sending verification email: %s\n", user.Email)
	return s.email.SendVerificationEmail(email, fullName, code)
}

func (s *service) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.Verified {
		return nil, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *service) VerifyEmail(ctx context.Context, code string) error {
	if err := s.repo.VerifyUser(ctx, code); err != nil {
		return errors.New("invalid verification code")
	}
	return nil
}

func NewService(repo Repository, emailSvc email.Service) Service {
	return &service{
		repo:  repo,
		email: emailSvc,
	}
}
