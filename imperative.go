package main

import (
	"strings"
	"unicode"
)

// Git's own convention is "Fix bug", not "Fixed bug" or "Fixes bug". Doing
// this properly needs a part-of-speech tagger, so this is a heuristic that
// leans toward missing things over flagging good subjects: past tense and
// gerunds are caught by suffix, third-person forms only when the stem is a
// verb commonly seen in commit subjects (plain "-s" matching would flag
// "access", "process" and "focus").

var irregularPast = map[string]bool{
	"made": true, "built": true, "ran": true, "wrote": true, "took": true,
	"gave": true, "got": true, "went": true, "began": true, "broke": true,
	"chose": true, "found": true, "kept": true, "left": true, "lost": true,
	"sent": true, "saw": true,
}

// Words that end in -ed or -ing but are already base-form verbs or nouns.
var edIngExceptions = map[string]bool{
	"need": true, "feed": true, "speed": true, "embed": true, "seed": true,
	"bleed": true, "shed": true, "proceed": true, "succeed": true,
	"exceed": true, "breed": true, "weed": true, "heed": true,
	"bring": true, "string": true, "thing": true, "spring": true,
	"swing": true, "sing": true, "ring": true,
}

var commonVerbs = map[string]bool{
	"add": true, "fix": true, "update": true, "remove": true, "change": true,
	"use": true, "make": true, "handle": true, "move": true, "rename": true,
	"create": true, "delete": true, "implement": true, "refactor": true,
	"improve": true, "support": true, "allow": true, "avoid": true,
	"clean": true, "correct": true, "drop": true, "enable": true,
	"disable": true, "ensure": true, "extract": true, "bump": true,
	"merge": true, "replace": true, "simplify": true, "apply": true,
	"set": true, "skip": true, "tighten": true, "document": true,
	"test": true, "reject": true, "return": true, "show": true,
	"print": true, "read": true, "write": true, "parse": true,
	"check": true, "catch": true, "expose": true, "revert": true,
	"wrap": true, "split": true, "load": true, "save": true, "send": true,
	"run": true, "build": true, "start": true, "stop": true,
	"prevent": true, "reduce": true, "convert": true, "rewrite": true,
}

// subjectVerb returns the lowercased first word of a subject, ignoring a
// conventional-commit style prefix such as "fix(parser): ".
func subjectVerb(subject string) string {
	subject = strings.TrimSpace(subject)
	if i := strings.Index(subject, ": "); i > 0 && !strings.ContainsAny(subject[:i], " \t") {
		subject = strings.TrimSpace(subject[i+2:])
	}
	fields := strings.Fields(subject)
	if len(fields) == 0 {
		return ""
	}
	word := strings.TrimFunc(fields[0], func(r rune) bool { return !unicode.IsLetter(r) })
	return strings.ToLower(word)
}

func isNonImperative(subject string) bool {
	w := subjectVerb(subject)
	if w == "" {
		return false
	}
	if irregularPast[w] {
		return true
	}
	if edIngExceptions[w] {
		return false
	}
	if len(w) > 3 && strings.HasSuffix(w, "ed") {
		return true
	}
	if len(w) > 4 && strings.HasSuffix(w, "ing") {
		return true
	}

	// Third person: "adds", "fixes", "applies".
	var stems []string
	if strings.HasSuffix(w, "ies") {
		stems = append(stems, strings.TrimSuffix(w, "ies")+"y")
	}
	if strings.HasSuffix(w, "es") {
		stems = append(stems, strings.TrimSuffix(w, "es"))
	}
	if strings.HasSuffix(w, "s") {
		stems = append(stems, strings.TrimSuffix(w, "s"))
	}
	for _, s := range stems {
		if commonVerbs[s] {
			return true
		}
	}
	return false
}
