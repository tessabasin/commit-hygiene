# commit-hygiene

`git log` gives you commit history but not a clean way to ask questions like
"how many of the last 500 commits have a subject line over 50 characters" or
"who actually wrote most of this repo". You end up piping through `awk` and
`sort` and reinventing the same brittle one-liner every time. commit-hygiene
takes a plain-text dump of a log and prints the numbers.

It doesn't shell out to `git` itself - it reads a log already produced with a
specific `--format`, either from a file or from stdin. That keeps the tool
usable against a log someone else generated, a log from a repo you don't have
checked out locally, or a log filtered through `git log -- some/path` first.

## Usage

Pipe a log straight from git:

```
git log --format='%H%x1f%an%x1f%ae%x1f%aI%x1f%s%x1f%b%x1e' | commit-hygiene
```

Or save the dump and pass it as a file argument:

```
git log --format='%H%x1f%an%x1f%ae%x1f%aI%x1f%s%x1f%b%x1e' > history.log
commit-hygiene history.log
```

The subject-length threshold defaults to 50 characters and can be overridden:

```
commit-hygiene -max-subject-len 72 history.log
```

Sample output:

```
commits analyzed:    842
empty subject:       0
subject > 50 chars: 61
subject ends in '.': 12
not imperative:      37

authors (7):
   401  Jamie Ostrander
   289  Priya Nair
    98  Dev Coleman
   ...
```

## Why the odd format string

Commit bodies can contain blank lines, so splitting records on `\n\n` breaks
on real repos. The format string uses `%x1e` (record separator) between
commits and `%x1f` (field separator) between fields - both are ASCII control
characters that never show up in commit text, so parsing never has to guess
where one commit ends and the next begins.

## Building

```
go build -o commit-hygiene .
```

No third-party dependencies - standard library only.

## Status

Early. Currently reports subject-line length, trailing-period style, empty
messages, non-imperative subjects, and a per-author commit count. The
imperative check is a heuristic on the first word: it flags past tense
("Fixed"), gerunds ("Fixing") and third-person forms of common verbs
("Fixes"), and would rather miss a case than flag a good subject. See the issue tracker for what's
planned next.
