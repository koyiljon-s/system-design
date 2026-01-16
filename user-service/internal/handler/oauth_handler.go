package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"primejobs/user-service/internal/model"
	"primejobs/user-service/internal/repository"
	"primejobs/user-service/internal/service/jwt"
	"primejobs/user-service/internal/service/oauth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	userRepo *repository.UserRepository
	google   *oauth.GoogleService
}

func NewOAuthHandler(repo *repository.UserRepository) *OAuthHandler {
	return &OAuthHandler{
		userRepo: repo,
		google:   oauth.NewGoogleService(),
	}
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}

// GET /oauth/google
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "state error"})
		return
	}

	c.SetCookie(
		"oauth_state",
		state,
		int((5 * time.Minute).Seconds()),
		"/",
		"",
		false, // true in prod (HTTPS)
		true,
	)

	c.Redirect(http.StatusFound, h.google.LoginURL(state))
}

// GET /oauth/google/callback
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {

	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth response"})
		return
	}

	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != state {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "state mismatch"})
		return
	}

	googleUser, err := h.google.GetUser(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "google auth failed"})
		return
	}

	// Try find by provider
	user, err := h.userRepo.FindByProvider("google", googleUser.ID)

	if err != nil {

		if err == gorm.ErrRecordNotFound {

			var picture *string
			if googleUser.Picture != "" {
				picture = &googleUser.Picture
			}

			user = &model.User{
				Name:       googleUser.Name,
				Email:      googleUser.Email,
				Provider:   "google",
				ProviderID: googleUser.ID,
				PictureURL: picture,
			}

			if err := h.userRepo.Create(user); err != nil {
				c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
				return
			}

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  toUserResponse(user),
		"token": token,
	})
}
