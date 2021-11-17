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

func WithRelayProxy(proxyURL *url.URL) ConfigOption {
	return func(c *Client) {
		c.relayProxyURL = proxyURL.String()
	}
}

func WithDaemonMode(dynamoTableName string, cacheTTL time.Duration) ConfigOption {
	return func(c *Client) {
		c.dynamoTableName = dynamoTableName
		c.cacheTTL = cacheTTL
	}
}

func WithDynamoBaseURL(baseURL *url.URL) ConfigOption {
	return func(c *Client) {
		c.dynamoBaseURL = baseURL.String()
	}
}

type FlagName string

type Client struct {
	sdkKey          string
	initWait        time.Duration
	relayProxyURL   string
	dynamoTableName string
	dynamoBaseURL   string
	cacheTTL        time.Duration
	wrappedClient   *ld.LDClient
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

	if c.relayProxyURL != "" && c.dynamoTableName != "" {
		return nil, errors.New("cannot supply both a Relay Proxy URL and a Dynamo table name")
	}

	return c, nil
}

func (c *Client) Connect() error {
	ldConfig := ld.Config{}

	if c.relayProxyURL != "" {
		ldConfig = configForProxyMode(c.relayProxyURL)
	} else if c.dynamoTableName != "" {
		ldConfig = configForDaemonMode(c.dynamoTableName, c.dynamoBaseURL, c.cacheTTL)
	}

	wrappedClient, err := ld.MakeCustomClient(c.sdkKey, ldConfig, c.initWait)
	if err != nil {
		return fmt.Errorf("create LaunchDarkly client: %w", err)
	}

	flagsClient.wrappedClient = wrappedClient

	return nil
}

func configForProxyMode(proxyURL string) ld.Config {
	return ld.Config{
		DataSource: ldcomponents.StreamingDataSource().BaseURI(proxyURL),
	}
}

func configForDaemonMode(dynamoTable string, dynamoBaseURL string, cacheTTL time.Duration) ld.Config {
	datastoreBuilder := lddynamodb.DataStore(dynamoTable)

	if dynamoBaseURL != "" {
		datastoreBuilder.ClientConfig(aws.NewConfig().WithEndpoint(dynamoBaseURL))
	}

	return ld.Config{
		DataSource: ldcomponents.ExternalUpdatesOnly(),
		DataStore: ldcomponents.PersistentDataStore(
			datastoreBuilder,
		).CacheTime(cacheTTL),
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

func (c *Client) Shutdown() error {
	return c.wrappedClient.Close()
}
