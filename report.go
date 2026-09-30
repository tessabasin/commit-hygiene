package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Report struct {
	Total          int
	ByAuthor       map[string]int
	MaxSubjectLen  int
	OverLength     int
	TrailingPeriod int
	NonImperative  int
	Empty          int
}

func buildReport(commits []Commit, maxSubjectLen int) Report {
	r := Report{ByAuthor: make(map[string]int), MaxSubjectLen: maxSubjectLen}
	r.Total = len(commits)

	for _, c := range commits {
		r.ByAuthor[c.AuthorName]++

		subject := strings.TrimSpace(c.Subject)
		if subject == "" {
			r.Empty++
			continue
		}
		if len(subject) > maxSubjectLen {
			r.OverLength++
		}
		if strings.HasSuffix(subject, ".") {
			r.TrailingPeriod++
		}
		if isNonImperative(subject) {
			r.NonImperative++
		}
	}

	return r
}

func (r Report) print(w io.Writer) {
	fmt.Fprintf(w, "commits analyzed:    %d\n", r.Total)
	fmt.Fprintf(w, "empty subject:       %d\n", r.Empty)
	fmt.Fprintf(w, "subject > %d chars: %d\n", r.MaxSubjectLen, r.OverLength)
	fmt.Fprintf(w, "subject ends in '.': %d\n", r.TrailingPeriod)
	fmt.Fprintf(w, "not imperative:      %d\n", r.NonImperative)

	type authorCount struct {
		name  string
		count int
	}
	authors := make([]authorCount, 0, len(r.ByAuthor))
	for name, count := range r.ByAuthor {
		authors = append(authors, authorCount{name, count})
	}
	sort.Slice(authors, func(i, j int) bool {
		if authors[i].count != authors[j].count {
			return authors[i].count > authors[j].count
		}
		return authors[i].name < authors[j].name
	})

	fmt.Fprintf(w, "\nauthors (%d):\n", len(authors))
	limit := len(authors)
	if limit > 10 {
		limit = 10
	}
	for _, a := range authors[:limit] {
		fmt.Fprintf(w, "  %4d  %s\n", a.count, a.name)
	}
	if len(authors) > limit {
		fmt.Fprintf(w, "  ... and %d more\n", len(authors)-limit)
	}
}
