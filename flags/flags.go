package flags

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

var (
	flagsClient            *Client
	errClientNotConfigured = errors.New("client not configured")
)

const (
	defaultSDKKeyEnvironmentVariable = "LAUNCHDARKLY_SDK_KEY"
)

func Configure(opts ...ConfigOption) error {
	c, err := NewClient(opts...)
	if err != nil {
		return fmt.Errorf("configure client: %w", err)
	}

	flagsClient = c
	return nil
}

func Connect() error {
	if flagsClient == nil {
		return errClientNotConfigured
	}

	return flagsClient.Connect()
}

func GetDefaultClient() (*Client, error) {
	if flagsClient == nil {
		return nil, errClientNotConfigured
	}

	return flagsClient, nil
}

type ConfigOption func(c *Client)

func WithSDKKey(key string) ConfigOption {
	return func(c *Client) {
		c.sdkKey = key
	}
}

func WithInitWait(t time.Duration) ConfigOption {
	return func(c *Client) {
		c.initWait = t
	}
}

func WithRelayProxy(proxyURL *url.URL) ConfigOption {
	return func(c *Client) {
		c.relayProxyURL = proxyURL.String()
	}
}

type FlagName string

type Client struct {
	sdkKey        string
	initWait      time.Duration
	relayProxyURL string
	wrappedClient *ld.LDClient
}

func NewClient(opts ...ConfigOption) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}

	if c.sdkKey == "" {
		defaultSDKKey, ok := os.LookupEnv(defaultSDKKeyEnvironmentVariable)
		if !ok {
			return nil, errors.New("LaunchDarkly SDK key not supplied via config option and the LAUNCHDARKLY_SDK_KEY environment variable does not exist")
		}
		c.sdkKey = defaultSDKKey
	}

	return c, nil
}

func (c *Client) Connect() error {
	ldConfig := ld.Config{}

	if c.relayProxyURL != "" {
		ldConfig.DataSource = ldcomponents.StreamingDataSource().BaseURI(c.relayProxyURL)
	}

	wrappedClient, err := ld.MakeCustomClient(c.sdkKey, ldConfig, c.initWait)
	if err != nil {
		return fmt.Errorf("create LaunchDarkly client: %w", err)
	}

	flagsClient.wrappedClient = wrappedClient

	return nil
}

func (c *Client) QueryBool(key FlagName, user User, defaultValue bool) (bool, error) {
	return c.wrappedClient.BoolVariation(string(key), user.ldUser, defaultValue)
}

func (c *Client) QueryString(key FlagName, user User, defaultValue string) (string, error) {
	return c.wrappedClient.StringVariation(string(key), user.ldUser, defaultValue)
}

func (c *Client) QueryInt(key FlagName, user User, defaultValue int) (int, error) {
	return c.wrappedClient.IntVariation(string(key), user.ldUser, defaultValue)
}

func (c *Client) Shutdown() error {
	return c.wrappedClient.Close()
}
