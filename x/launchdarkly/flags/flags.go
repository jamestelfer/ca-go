package flags

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

var flagsClient *Client

const defaultSDKKeyEnvironmentVariable = "LAUNCHDARKLY_SDK_KEY"

// FlagName establishes a type for flag names.
type FlagName string

// Client is a wrapper around the LaunchDarkly client.
type Client struct {
	sdkKey           string
	initWait         time.Duration
	proxyModeConfig  *proxyModeConfig
	daemonModeConfig *daemonModeConfig
	wrappedConfig    ld.Config
	wrappedClient    *ld.LDClient
}

// Configure configures the client as a managed singleton. An error is returned
// if mandatory ConfigOptions are not supplied, or an invalid combination of
// options is provided.
func Configure(opts ...ConfigOption) error {
	c, err := NewClient(opts...)
	if err != nil {
		return fmt.Errorf("configure client: %w", err)
	}

	flagsClient = c
	return nil
}

// Connect attempts to connect the managed singleton to LaunchDarkly. An error
// is returned if the singleton is not yet configured, a connection has already
// been established, or a connection error occurs.
func Connect() error {
	if flagsClient == nil {
		return errClientNotConfigured
	}

	return flagsClient.Connect()
}

// GetDefaultClient returns the managed singleton client. An error is returned
// if the client is not yet configured.
func GetDefaultClient() (*Client, error) {
	if flagsClient == nil {
		return nil, errClientNotConfigured
	}

	return flagsClient, nil
}

// NewClient configures and returns an instance of the client. An error is
// returned if mandatory ConfigOptions are not supplied, or an invalid
// combination of options is provided.
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

// Connect attempts to establish the initial connection to LaunchDarkly. An
// error is returned if a connection has already been established, or a
// connection error occurs.
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

// QueryBool retrieves the value of a boolean flag. User attributes are
// extracted from the context. The supplied default value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryBool(ctx context.Context, key FlagName, defaultValue bool) (bool, error) {
	user, err := UserFromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.BoolVariation(string(key), user.ldUser, defaultValue)
}

// QueryBoolWithUser retrieves the value of a boolean flag. A User object must
// be supplied manually. The supplied default value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryBoolWithUser(key FlagName, user User, defaultValue bool) (bool, error) {
	return c.wrappedClient.BoolVariation(string(key), user.ldUser, defaultValue)
}

// QueryString retrieves the value of a string flag. User attributes are
// extracted from the context. The supplied default value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryString(ctx context.Context, key FlagName, defaultValue string) (string, error) {
	user, err := UserFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.StringVariation(string(key), user.ldUser, defaultValue)
}

// QueryStringWithUser retrieves the value of a string flag. A User object must
// be supplied manually. The supplied default value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryStringWithUser(key FlagName, user User, defaultValue string) (string, error) {
	return c.wrappedClient.StringVariation(string(key), user.ldUser, defaultValue)
}

// QueryInt retrieves the value of an integer flag. User attributes are
// extracted from the context. The supplied default value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryInt(ctx context.Context, key FlagName, defaultValue int) (int, error) {
	user, err := UserFromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.IntVariation(string(key), user.ldUser, defaultValue)
}

// QueryIntWithUser retrieves the value of an integer flag. A User object must
// be supplied manually. The supplied default value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryIntWithUser(key FlagName, user User, defaultValue int) (int, error) {
	return c.wrappedClient.IntVariation(string(key), user.ldUser, defaultValue)
}

// RawClient returns the wrapped LaunchDarkly client. The return value should be
// casted to an *ld.LDClient instance.
func (c *Client) RawClient() interface{} {
	return c.wrappedClient
}

// Shutdown instructs the wrapped LaunchDarkly client to close any open
// connections and flush any flag evaluation events.
func (c *Client) Shutdown() error {
	return c.wrappedClient.Close()
}
