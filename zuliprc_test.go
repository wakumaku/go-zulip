package zulip

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZuliprcParser(t *testing.T) {
	fileContent := `[api]
// empty line

email=user@localhost
key=apikey
site=https://localhost
`

	f, err := os.CreateTemp("", "zuliprc")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = f.WriteString(fileContent)
	assert.NoError(t, err)

	z, err := ParseZuliprc(f.Name())
	assert.NoError(t, err)

	apiSection := z["api"]
	assert.Equal(t, "user@localhost", apiSection.Email)
	assert.Equal(t, "apikey", apiSection.APIKey)
	assert.Equal(t, "https://localhost", apiSection.Site)

	unknownSection := z["unknown"]
	assert.Equal(t, "", unknownSection.Email)
	assert.Equal(t, "", unknownSection.APIKey)
	assert.Equal(t, "", unknownSection.Site)
}

func TestZuliprcParseFileNotFound(t *testing.T) {
	_, err := ParseZuliprc("non-existing-file")
	assert.Error(t, err)
}
