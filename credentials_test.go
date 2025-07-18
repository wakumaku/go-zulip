package zulip_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip"
)

func TestCredentials(t *testing.T) {
	expectedSite := "the url"
	expectedEmail := "the email"
	expectedAPIKey := "the key"

	credentials := zulip.Credentials(expectedSite, expectedEmail, expectedAPIKey)

	c, err := credentials()
	require.NoError(t, err)
	assert.Equal(t, expectedSite, c.Site)
	assert.Equal(t, expectedEmail, c.Email)
	assert.Equal(t, expectedAPIKey, c.APIKey)
}

func TestCredentialsFromZuliprc(t *testing.T) {
	fileContent := `[api]
email=user@localhost
key=apikey
site=https://localhost
`

	f, err := os.CreateTemp("", "zuliprc")
	require.NoError(t, err)

	_, err = f.WriteString(fileContent)
	require.NoError(t, err)

	credentials := zulip.CredentialsFromZuliprc(f.Name(), "api")

	c, err := credentials()
	require.NoError(t, err)
	assert.Equal(t, "user@localhost", c.Email)
	assert.Equal(t, "apikey", c.APIKey)
	assert.Equal(t, "https://localhost", c.Site)
}
