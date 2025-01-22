package zulip

import "fmt"

type credentials struct {
	Email  string
	APIKey string
	Site   string
}

type CredentialsProvider func() (*credentials, error)

func Credentials(site, email, apiKey string) CredentialsProvider {
	return func() (*credentials, error) {
		return &credentials{
			Email:  email,
			APIKey: apiKey,
			Site:   site,
		}, nil
	}
}

func CredentialsFromZuliprc(filePath string, section string) CredentialsProvider {
	return func() (*credentials, error) {
		zuliprc, err := ParseZuliprc(filePath)
		if err != nil {
			return nil, err
		}

		apiSection, ok := zuliprc[section]
		if !ok {
			return nil, fmt.Errorf("no '%s' section found in zuliprc file '%s'", section, filePath)
		}

		return &credentials{
			Email:  apiSection.Email,
			APIKey: apiSection.APIKey,
			Site:   apiSection.Site,
		}, nil
	}
}
