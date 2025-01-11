package zulip

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
)

type (
	Zuliprc     map[string]SectionData
	SectionData struct {
		Key   string
		Email string
		Site  string
	}
)

func ParseZuliprc(file string) (Zuliprc, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(f)
	return ParseZuliprcContent(r)
}

func ParseZuliprcContent(b io.Reader) (Zuliprc, error) {
	s := bufio.NewScanner(b)

	rxSection := regexp.MustCompile(`\[(.*)\]`)
	rxKeyVal := regexp.MustCompile(`(.*)=(.*)`)

	currentSection := "unknown"

	z := Zuliprc{}
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		if rxSection.MatchString(line) {
			section := rxSection.FindStringSubmatch(line)[1]
			z[section] = SectionData{}
			currentSection = section
			continue
		}

		if rxKeyVal.MatchString(line) {
			kv := rxKeyVal.FindStringSubmatch(line)
			key := kv[1]
			val := kv[2]

			sectionData := z[currentSection]
			switch strings.TrimSpace(key) {
			case "email":
				sectionData.Email = strings.TrimSpace(val)
			case "key":
				sectionData.Key = strings.TrimSpace(val)
			case "site":
				sectionData.Site = strings.TrimSpace(val)
			}
			z[currentSection] = sectionData
		}
	}

	return z, nil
}
