# Markdown Table Generator

This utility converts tab-delimited data into a markdown table format, omitting the row number column as requested.

## Usage

Run the utility:

```bash
go run main.go
```

## Input Format

The utility processes tab-delimited data with the following format:
- Column 1: Row numbers (omitted in output)
- Column 2: Event ID
- Column 3: Count value

## Output

Generates a markdown table with two columns:
- `eventid`: The event identifier
- `c`: The count value

## Example

Input:
```
#	eventid	c
1	FINDERSALES_SavedSearchNotification	130
2	PCCM_inboxNewMessageSummaryNotification	23
```

Output:
```markdown
| eventid | c |
|---|---|
| FINDERSALES_SavedSearchNotification | 130 |
| PCCM_inboxNewMessageSummaryNotification | 23 |
```