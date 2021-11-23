package flags

import (
	"errors"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	lddynamodb "github.com/launchdarkly/go-server-sdk-dynamodb"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

var errClientNotConfigured = errors.New("client not configured")

type proxyModeConfig struct {
	relayProxyURL string
}

type daemonModeConfig struct {
	dynamoTableName string
	dynamoBaseURL   string
	cacheTTL        time.Duration
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
