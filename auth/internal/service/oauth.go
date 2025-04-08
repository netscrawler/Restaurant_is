package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	pb "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

type OAuthService struct {
	oauthRepo  repository.OAuthRepository
	clientRepo repository.ClientRepository
	config     *oauth2.Config
}

func NewOAuthService(oauthRepo repository.OAuthRepository, clientRepo repository.ClientRepository, clientID, clientSecret, redirectURL string) *OAuthService {
	return &OAuthService{
		oauthRepo:  oauthRepo,
		clientRepo: clientRepo,
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     yandex.Endpoint,
			Scopes:       []string{"login:email", "login:info"},
		},
	}
}

func (s *OAuthService) HandleCallback(
	ctx context.Context,
	provider string,
	code string,
) (*pb.LoginResponse, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://login.yandex.ru/info")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	// Process user info and create or update user in the system
	// ...

	return &pb.LoginResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil
}

func (s *OAuthService) GetAuthURL(provider string) (string, error) {
	if provider != "yandex" {
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	authURL := s.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return authURL, nil
}
