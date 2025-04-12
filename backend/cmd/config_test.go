package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigEnvVariablesImpl(t *testing.T) {
	vars := map[string]any{}
	env := map[string]string{
		"APP_PASSWORDENCRYPTIONKEY": "WxGGJwCCNcJdxK9idqbeU9iwOteaw0kiU+/oa8KQZl4=",
		"APP_AWSREGIONS":            "foo,bar",
	}

	require.NoError(t, loadConfigEnvVariablesImpl(func(key string) string {
		return env[key]
	}, func(key string, value any) {
		vars[key] = value
	}))

	assert.Equal(t, map[string]any{
		"App.PasswordEncryptionKey": []byte{
			0x5b, 0x11, 0x86, 0x27, 0x0, 0x82, 0x35, 0xc2,
			0x5d, 0xc4, 0xaf, 0x62, 0x76, 0xa6, 0xde, 0x53,
			0xd8, 0xb0, 0x3a, 0xd7, 0x9a, 0xc3, 0x49, 0x22,
			0x53, 0xef, 0xe8, 0x6b, 0xc2, 0x90, 0x66, 0x5e,
		},
		"App.AWSRegions": []string{"foo", "bar"},
	}, vars)
}
