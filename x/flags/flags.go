package flags

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	lddynamodb "github.com/launchdarkly/go-server-sdk-dynamodb"
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

func WithProxyMode(proxyURL *url.URL) ConfigOption {
	return func(c *Client) {
		c.proxyModeConfig = &proxyModeConfig{
			relayProxyURL: proxyURL.String(),
		}
	}
}

func WithDaemonMode(dynamoTableName string, cacheTTL time.Duration) ConfigOption {
	return func(c *Client) {
		if c.daemonModeConfig == nil {
			c.daemonModeConfig = &daemonModeConfig{}
		}

		c.daemonModeConfig = &daemonModeConfig{
			dynamoTableName: dynamoTableName,
			cacheTTL:        cacheTTL,
		}
	}
}

func WithDynamoBaseURL(baseURL *url.URL) ConfigOption {
	return func(c *Client) {
		if c.daemonModeConfig == nil {
			c.daemonModeConfig = &daemonModeConfig{}
		}

		c.daemonModeConfig.dynamoBaseURL = baseURL.String()
	}
}

type FlagName string

type Client struct {
	sdkKey           string
	initWait         time.Duration
	proxyModeConfig  *proxyModeConfig
	daemonModeConfig *daemonModeConfig
	wrappedConfig    ld.Config
	wrappedClient    *ld.LDClient
}

type proxyModeConfig struct {
	relayProxyURL string
}

type daemonModeConfig struct {
	dynamoTableName string
	dynamoBaseURL   string
	cacheTTL        time.Duration
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

	if c.proxyModeConfig != nil && c.daemonModeConfig != nil {
		return nil, errors.New("cannot configure the SDK for Proxy and Daemon modes simultaneously")
	}

	if c.proxyModeConfig != nil {
		c.wrappedConfig = configForProxyMode(c.proxyModeConfig)
	} else if c.daemonModeConfig != nil {
		c.wrappedConfig = configForDaemonMode(c.daemonModeConfig)
	}

	return c, nil
}

func (c *Client) Connect() error {
	if c.wrappedClient != nil {
		return errors.New("attempted to call Connect on a connected client")
	}

	wrappedClient, err := ld.MakeCustomClient(c.sdkKey, c.wrappedConfig, c.initWait)
	if err != nil {
		return fmt.Errorf("create LaunchDarkly client: %w", err)
	}

	flagsClient.wrappedClient = wrappedClient

	return nil
}

func configForProxyMode(cfg *proxyModeConfig) ld.Config {
	return ld.Config{
		DataSource: ldcomponents.StreamingDataSource().BaseURI(cfg.relayProxyURL),
	}
}

func configForDaemonMode(cfg *daemonModeConfig) ld.Config {
	datastoreBuilder := lddynamodb.DataStore(cfg.dynamoTableName)

	if cfg.dynamoBaseURL != "" {
		datastoreBuilder.ClientConfig(aws.NewConfig().WithEndpoint(cfg.dynamoBaseURL))
	}

	return ld.Config{
		DataSource: ldcomponents.ExternalUpdatesOnly(),
		DataStore: ldcomponents.PersistentDataStore(
			datastoreBuilder,
		).CacheTime(cfg.cacheTTL),
	}
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

func (c *Client) RawClient() interface{} {
	return c.wrappedClient
}

func (c *Client) Shutdown() error {
	return c.wrappedClient.Close()
}
