package flags_test

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags"
	"github.com/stretchr/testify/require"
)

func TestInitialisationSingleton(t *testing.T) {
	t.Run("errors if an SDK key is not supplied", func(t *testing.T) {
		err := flags.Configure()
		require.Error(t, err)

		_, err = flags.GetDefaultClient()
		require.Error(t, err)
	})

	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		os.Setenv("LAUNCHDARKLY_SDK_KEY", "foobar")
		defer os.Unsetenv("LAUNCHDARKLY_SDK_KEY")
		err := flags.Configure()
		require.NoError(t, err)
	})

	t.Run("does not error if SDK key supplied as config option", func(t *testing.T) {
		err := flags.Configure(flags.WithSDKKey("foobar"))
		require.NoError(t, err)
	})

	t.Run("allows an initialisation wait time to be specified", func(t *testing.T) {
		err := flags.Configure(
			flags.WithSDKKey("foobar"),
			flags.WithInitWait(2*time.Second))
		require.NoError(t, err)
	})

	t.Run("allows a Relay Proxy URL to be specified", func(t *testing.T) {
		proxyURL, err := url.Parse("http://localhost:8030")
		require.NoError(t, err)

		err = flags.Configure(
			flags.WithSDKKey("foobar"),
			flags.WithProxyMode(proxyURL))
		require.NoError(t, err)
	})

	t.Run("allows daemon mode to be configured", func(t *testing.T) {
		err := flags.Configure(
			flags.WithSDKKey("foobar"),
			flags.WithDaemonMode("dynamo-table-name", 10*time.Second),
		)
		require.NoError(t, err)
	})

	t.Run("allows daemon mode to be configured with an alternate Dynamo base URL", func(t *testing.T) {
		baseURL, err := url.Parse("http://localhost:6789")
		require.NoError(t, err)

		err = flags.Configure(
			flags.WithSDKKey("foobar"),
			flags.WithDaemonMode("dynamo-table-name", 10*time.Second),
			flags.WithDynamoBaseURL(baseURL),
		)
		require.NoError(t, err)

		err = flags.Configure(
			flags.WithSDKKey("foobar"),
			flags.WithDynamoBaseURL(baseURL),
			flags.WithDaemonMode("dynamo-table-name", 10*time.Second),
		)
		require.NoError(t, err, "should not error when base URL provided before Dynamo table")
	})

	t.Run("does not allow both proxy and daemon modes", func(t *testing.T) {
		proxyURL, err := url.Parse("http://localhost:8030")
		require.NoError(t, err)

		err = flags.Configure(
			flags.WithSDKKey("foobar"),
			flags.WithDaemonMode("dynamo-table-name", 10*time.Second),
			flags.WithProxyMode(proxyURL),
		)
		require.Error(t, err)
	})
}
