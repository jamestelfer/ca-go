package launchdarkly_test

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/x/launchdarkly"
	"github.com/stretchr/testify/require"
)

func TestInitialisationSingleton(t *testing.T) {
	t.Run("errors if an SDK key is not supplied", func(t *testing.T) {
		err := launchdarkly.Configure()
		require.Error(t, err)

		_, err = launchdarkly.GetDefaultClient()
		require.Error(t, err)
	})

	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		os.Setenv("LAUNCHDARKLY_SDK_KEY", "foobar")
		defer os.Unsetenv("LAUNCHDARKLY_SDK_KEY")
		err := launchdarkly.Configure()
		require.NoError(t, err)
	})

	t.Run("does not error if SDK key supplied as config option", func(t *testing.T) {
		err := launchdarkly.Configure(launchdarkly.WithSDKKey("foobar"))
		require.NoError(t, err)
	})

	t.Run("allows an initialisation wait time to be specified", func(t *testing.T) {
		err := launchdarkly.Configure(
			launchdarkly.WithSDKKey("foobar"),
			launchdarkly.WithInitWait(2*time.Second))
		require.NoError(t, err)
	})

	t.Run("allows a Relay Proxy URL to be specified", func(t *testing.T) {
		proxyURL, err := url.Parse("http://localhost:8030")
		require.NoError(t, err)

		err = launchdarkly.Configure(
			launchdarkly.WithSDKKey("foobar"),
			launchdarkly.WithProxyMode(proxyURL))
		require.NoError(t, err)
	})

	t.Run("allows daemon mode to be configured", func(t *testing.T) {
		err := launchdarkly.Configure(
			launchdarkly.WithSDKKey("foobar"),
			launchdarkly.WithDaemonMode("dynamo-table-name", 10*time.Second),
		)
		require.NoError(t, err)
	})

	t.Run("allows daemon mode to be configured with an alternate Dynamo base URL", func(t *testing.T) {
		baseURL, err := url.Parse("http://localhost:6789")
		require.NoError(t, err)

		err = launchdarkly.Configure(
			launchdarkly.WithSDKKey("foobar"),
			launchdarkly.WithDaemonMode("dynamo-table-name", 10*time.Second),
			launchdarkly.WithDynamoBaseURL(baseURL),
		)
		require.NoError(t, err)

		err = launchdarkly.Configure(
			launchdarkly.WithSDKKey("foobar"),
			launchdarkly.WithDynamoBaseURL(baseURL),
			launchdarkly.WithDaemonMode("dynamo-table-name", 10*time.Second),
		)
		require.NoError(t, err, "should not error when base URL provided before Dynamo table")
	})

	t.Run("does not allow both proxy and daemon modes", func(t *testing.T) {
		proxyURL, err := url.Parse("http://localhost:8030")
		require.NoError(t, err)

		err = launchdarkly.Configure(
			launchdarkly.WithSDKKey("foobar"),
			launchdarkly.WithDaemonMode("dynamo-table-name", 10*time.Second),
			launchdarkly.WithProxyMode(proxyURL),
		)
		require.Error(t, err)
	})
}
