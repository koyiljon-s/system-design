package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	Config *oauth2.Config
}

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleService() *GoogleService {
	return &GoogleService{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *GoogleService) LoginURL(state string) string {
	return g.Config.AuthCodeURL(state)
}

func (g *GoogleService) GetUser(code string) (*GoogleUser, error) {

	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := g.Config.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("google userinfo failed")
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
