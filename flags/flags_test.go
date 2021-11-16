package flags_test

import (
	"os"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/flags"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("errors if an SDK key is not supplied", func(t *testing.T) {
		_, err := flags.NewConfig()
		require.Error(t, err)
	})

	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		os.Setenv("LAUNCHDARKLY_SDK_KEY", "foobar")
		_, err := flags.NewConfig()
		require.NoError(t, err)
	})

	t.Run("does not error if SDK key supplied as config option", func(t *testing.T) {
		_, err := flags.NewConfig(flags.WithSDKKey("foobar"))
		require.NoError(t, err)
	})

	t.Run("allows an initialisation wait time to be specified", func(t *testing.T) {
		_, err := flags.NewConfig(
			flags.WithSDKKey("foobar"),
			flags.WithInitWait(2*time.Second))
		require.NoError(t, err)
	})
}
