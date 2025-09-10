package main

import (
	"testing"
)

func TestGenerateMarkdownTable(t *testing.T) {
	input := `#	eventid	c
1	FINDERSALES_SavedSearchNotification	130
2	PCCM_inboxNewMessageSummaryNotification	23`

	expected := `| eventid | c |
|---|---|
| FINDERSALES_SavedSearchNotification | 130 |
| PCCM_inboxNewMessageSummaryNotification | 23 |
`

	result := generateMarkdownTable(input)
	
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestGenerateMarkdownTableEmptyInput(t *testing.T) {
	result := generateMarkdownTable("")
	if result != "" {
		t.Errorf("Expected empty string for empty input, got: %s", result)
	}
}

func TestGenerateMarkdownTableSingleRow(t *testing.T) {
	input := `#	eventid	c
1	TEST_event	42`

	expected := `| eventid | c |
|---|---|
| TEST_event | 42 |
`

	result := generateMarkdownTable(input)
	
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}