package main

import (
	"io"
	"strings"
	"time"
)

// Field and record separators match the format string documented in
// README.md. Commit bodies routinely contain blank lines and stray
// punctuation, so we can't split records on "\n\n" the way a naive parser
// would - control characters that never show up in commit text are the only
// safe delimiter.
const (
	recordSep = "\x1e"
	fieldSep  = "\x1f"
)

type Commit struct {
	Hash        string
	AuthorName  string
	AuthorEmail string
	Date        time.Time
	Subject     string
	Body        string
}

func parseLog(r io.Reader) ([]Commit, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	for _, record := range strings.Split(string(raw), recordSep) {
		record = strings.Trim(record, "\n")
		if record == "" {
			continue
		}

		fields := strings.Split(record, fieldSep)
		if len(fields) != 6 {
			// Doesn't match the expected format; skip rather than abort
			// the whole run over one bad line.
			continue
		}

		date, _ := time.Parse(time.RFC3339, fields[3])

		commits = append(commits, Commit{
			Hash:        fields[0],
			AuthorName:  fields[1],
			AuthorEmail: fields[2],
			Date:        date,
			Subject:     fields[4],
			Body:        strings.TrimSpace(fields[5]),
		})
	}

	return commits, nil
}
