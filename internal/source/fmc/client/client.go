package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
)

type FMCClient struct {
	HTTPClient     *http.Client
	BaseURL        string
	Username       string
	Password       string
	AccessToken    string
	RefreshToken   string
	DefaultTimeout time.Duration
	Logger         *logger.Logger
	Ctx            context.Context
}

// NewFMCClient creates a new FMC client with the given parameters.
// It authenticates to the FMC API and stores the access and refresh tokens.
func NewFMCClient(
	context context.Context,
	username string,
	password string,
	httpScheme string,
	hostname string,
	port int,
	httpClient *http.Client,
	logger *logger.Logger,
) (*FMCClient, error) {
	c := &FMCClient{
		HTTPClient:     httpClient,
		BaseURL:        fmt.Sprintf("%s://%s:%d/api", httpScheme, hostname, port),
		Username:       username,
		Password:       password,
		DefaultTimeout: time.Second * constants.DefaultAPITimeout,
		Logger:         logger,
		Ctx:            context,
	}

	aToken, rToken, err := c.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("authentication: %w", err)
	}

	c.AccessToken = aToken
	c.RefreshToken = rToken

	return c, nil
}
