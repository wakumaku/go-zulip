package zulip

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZuliprcParser(t *testing.T) {
	fileContent := `[api]
// empty line

email=user@localhost
key=apikey
site=https://localhost

[bot]
email=bot@localhost
key=botapikey
site=https://localhost
`

	f, err := os.CreateTemp("", "zuliprc")
	require.NoError(t, err)

	defer func() { _ = os.Remove(f.Name()) }()

	_, err = f.WriteString(fileContent)
	require.NoError(t, err)

	z, err := ParseZuliprc(f.Name())
	require.NoError(t, err)

	apiSection := z["api"]
	assert.Equal(t, "user@localhost", apiSection.Email)
	assert.Equal(t, "apikey", apiSection.APIKey)
	assert.Equal(t, "https://localhost", apiSection.Site)

	botSection := z["bot"]
	assert.Equal(t, "bot@localhost", botSection.Email)
	assert.Equal(t, "botapikey", botSection.APIKey)
	assert.Equal(t, "https://localhost", botSection.Site)
}

func TestZuliprcParseFileNotFound(t *testing.T) {
	_, err := ParseZuliprc("non-existing-file")
	assert.Error(t, err)
}
