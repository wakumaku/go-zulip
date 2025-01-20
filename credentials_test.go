package zulip_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
)

func TestCredentials(t *testing.T) {
	expectedSite := "the url"
	expectedEmail := "the email"
	expectedKey := "the key"

	credentials := zulip.Credentials(expectedSite, expectedEmail, expectedKey)

	c, err := credentials()
	assert.NoError(t, err)
	assert.Equal(t, expectedSite, c.Site)
	assert.Equal(t, expectedEmail, c.Email)
	assert.Equal(t, expectedKey, c.Key)
}

func TestCredentialsFromZuliprc(t *testing.T) {
	fileContent := `[api]
email=user@localhost
key=apikey
site=https://localhost
`

	f, err := os.CreateTemp("", "zuliprc")
	assert.NoError(t, err)

	_, err = f.WriteString(fileContent)
	assert.NoError(t, err)

	credentials := zulip.CredentialsFromZuliprc(f.Name(), "api")

	c, err := credentials()
	assert.NoError(t, err)
	assert.Equal(t, "user@localhost", c.Email)
	assert.Equal(t, "apikey", c.Key)
	assert.Equal(t, "https://localhost", c.Site)
}
