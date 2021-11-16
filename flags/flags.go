package flags

import (
	"errors"
	"fmt"
	"os"
	"time"

	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

var flagsClient Client

const (
	defaultSDKKeyEnvironmentVariable = "LAUNCHDARKLY_SDK_KEY"
)

type Config struct {
	sdkKey string
}

type ConfigOption func(c *Config)

func WithSDKKey(key string) ConfigOption {
	return func(c *Config) {
		c.sdkKey = key
	}
}

func NewConfig(opts ...ConfigOption) (*Config, error) {
	c := &Config{}
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

func Start(c *Config) error {
	ldConfig := ld.Config{}

	wrappedClient, err := ld.MakeCustomClient(c.sdkKey, ldConfig, 2*time.Second)
	if err != nil {
		return fmt.Errorf("create LaunchDarkly client: %w", err)
	}

	flagsClient = &ldClient{
		wrappedClient: wrappedClient,
	}

	return nil
}

func GetClient() (*Client, error) {
	if flagsClient == nil {
		return nil, errors.New("client not started")
	}

	return &flagsClient, nil
}

type FlagName string

type Client interface {
	Shutdown() error
}

type ldClient struct {
	wrappedClient *ld.LDClient
}

func (c *ldClient) Shutdown() error {
	return c.wrappedClient.Close()
}
