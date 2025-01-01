package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmj/config"
	"fmj/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/oauth2/v2"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

type Handler struct {
	service      Service
	config       *config.Config
	oauth2Config *goauth2.Config
}

func NewHandler(service Service, cfg *config.Config) *Handler {
	oauth2Config := &goauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleCallbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	return &Handler{service: service, config: cfg, oauth2Config: oauth2Config}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.GET("/login", h.ShowLogin)
		auth.POST("/login", h.Login)
		auth.GET("/register", h.ShowRegister)
		auth.POST("/register", h.Register)
		auth.GET("/verify", h.VerifyEmail)
		auth.GET("/logout", h.Logout)
		auth.GET("/google/login", h.GoogleLogin)
		auth.GET("/google/callback", h.GoogleCallback)
	}
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	// Generate random state
	state := make([]byte, 16)
	rand.Read(state)

	// Store state in session
	session := sessions.Default(c)
	session.Set("oauth_state", hex.EncodeToString(state))
	session.Save()

	// Redirect to Google
	url := h.oauth2Config.AuthCodeURL(hex.EncodeToString(state), goauth2.AccessTypeOffline, goauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	var data map[string]interface{}
	toastPage := filepath.Join("templates", "partials", "toast.html")
	indexPage := filepath.Join("templates", "pages", "index.html")

	// Verify state
	session := sessions.Default(c)
	expectedState := session.Get("oauth_state")
	if expectedState != c.Query("state") {
		data = map[string]interface{}{
			"Error": "An error occurred, try again",
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Google callback state mismatch", slog.String("state", c.Query("state")))
		return
	}
	session.Delete("oauth_state")

	// Exchange code for token
	code := c.Query("code")
	token, err := h.oauth2Config.Exchange(c, code)
	if err != nil {
		data = map[string]interface{}{
			"Error": "An error occurred, try again",
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Failed to exchange code", slog.String("error", err.Error()))
		return
	}

	// copied now
	// Check token expiration
	if token.Expiry.Before(time.Now()) {
		data = map[string]interface{}{
			"Error": "Session expired. Please try logging in again.",
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Google callback error", slog.String("error", "Token expired"))
		return
	}

	// Refresh token if available (optional)
	if token.RefreshToken != "" && token.Expiry.Before(time.Now().Add(-5*time.Minute)) {
		tokenSource := h.oauth2Config.TokenSource(context.Background(), token)
		newToken, refreshErr := tokenSource.Token()
		if refreshErr != nil {
			data = map[string]interface{}{
				"Error": "Failed to refresh token. Please try again.",
			}
			utils.Render(c, indexPage, nil)
			utils.Render(c, toastPage, data)
			slog.Error("Google callback error", slog.String("error", refreshErr.Error()))
			return
		}
		token = newToken
	}
	//end copy

	// Verify ID token
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		data = map[string]interface{}{
			"Error": "ID token missing in response. Please try again.",
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Google callback error", slog.String("error", "ID token missing"))
		return
	}

	payload, err := idtoken.Validate(context.Background(), idToken, h.config.GoogleClientID)
	if err != nil {
		// Log additional debugging information
		slog.Error("Failed to validate ID token",
			slog.String("error", err.Error()),
			slog.String("idToken", idToken),
			slog.String("clientID", h.config.GoogleClientID),
		)
		data = map[string]interface{}{
			"Error": "Failed to validate ID token. Please try again.",
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Google callback error", slog.String("error", err.Error()))
		return
	}

	slog.Info("ID token successfully validated", slog.Any("payload", payload))

	// Extract user info from claims
	emailVerified := payload.Claims["email_verified"].(bool)
	userinfo := &oauth2.Userinfo{
		Id:            payload.Claims["sub"].(string),
		Email:         payload.Claims["email"].(string),
		Name:          payload.Claims["name"].(string),
		Picture:       payload.Claims["picture"].(string),
		VerifiedEmail: &emailVerified,
	}

	// Handle user login/registration
	user, err := h.service.HandleGoogleLogin(c, userinfo)
	if err != nil {
		data = map[string]interface{}{
			"Error": "An error occurred, try again",
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Google callback error", slog.String("error", err.Error()))
		return
	}

	// Set session
	session.Set("user_id", user.ID.Hex())
	session.Save()

	data = map[string]interface{}{
		"Success": "Login in successful",
	}
	utils.Render(c, indexPage, nil)
	utils.Render(c, toastPage, data)
	slog.Error("Google callback error", slog.String("error", err.Error()))
}

func (h *Handler) ShowLogin(c *gin.Context) {
	// Define paths to the user templates.
	loginPage := filepath.Join("templates", "auth", "login.html")
	utils.Render(c, loginPage, nil)
}

func (h *Handler) Login(c *gin.Context) {
	var data map[string]interface{}
	toastPage := filepath.Join("templates", "partials", "toast.html")

	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.service.Login(email, password)
	if err != nil {
		// Prepare error data if login fails.
		data = map[string]interface{}{
			"Error": err.Error(),
		}
		slog.Error("Error logging a user in database", slog.String("email", email), slog.String("password", password), slog.String("error", err.Error()))
		utils.Render(c, toastPage, data)
		return
	}

	// Login succeeded, set session.
	session := sessions.Default(c)
	session.Set("user_id", user.ID.Hex())
	if err := session.Save(); err != nil {
		data = map[string]interface{}{
			"Error": "An error occurred while starting your session.",
		}
		slog.Error("An error occurred while saving the session", "error", err)
		utils.Render(c, toastPage, data)
		return
	}

	// Redirect to dashboard on successful login.
	c.Header("HX-Redirect", "/dashboard")
}

func (h *Handler) ShowRegister(c *gin.Context) {
	registerPage := filepath.Join("templates", "auth", "register.html")
	utils.Render(c, registerPage, nil)
}

func (h *Handler) Register(c *gin.Context) {
	var data map[string]interface{}
	toastPage := filepath.Join("templates", "partials", "toast.html")
	fullName := c.PostForm("full_name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if err := h.service.Register(c, fullName, email, password); err != nil {
		data = map[string]interface{}{
			"Error": err.Error(),
		}
		utils.Render(c, toastPage, data)
		slog.Error("Error registering user in database", slog.String("email", email), slog.String("full_name", fullName), slog.String("email", email), slog.String("error", err.Error()))
		return
	}

	data = map[string]interface{}{
		"Success": "Registration successful! Please check your email to verify your account.",
	}
	utils.Render(c, toastPage, data)
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var data map[string]interface{}
	toastPage := filepath.Join("templates", "partials", "toast.html")
	indexPage := filepath.Join("templates", "pages", "index.html")

	code := c.Query("code")

	if err := h.service.VerifyEmail(c, code); err != nil {
		data = map[string]interface{}{
			"Error": err.Error(),
		}
		utils.Render(c, indexPage, nil)
		utils.Render(c, toastPage, data)
		slog.Error("Error verifying user email", slog.String("error", err.Error()))
		return
	}

	data = map[string]interface{}{
		"Success": "Email verified successfully! You can now login.",
	}
	utils.Render(c, indexPage, nil)
	utils.Render(c, toastPage, data)
}

func (h *Handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		return
	}
	c.Redirect(http.StatusSeeOther, "/auth/login")
}
