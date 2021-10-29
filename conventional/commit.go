package conventional

import (
	"bufio"
	"github.com/hooliganlin/versioning/conventional-wisdom/git"
	"log"
	"regexp"
	"strings"
)

type CommitType string
const (
	Fix      CommitType = "fix"
	Feature  CommitType = "feat"
)

type Commit struct {
	Type  CommitType
	Scope string
	Title      string
	Body       string
	IsBreaking bool
	git.Commit
}

// NewCommit creates a new conventional.Commit from a git.Commit. Parses out the Title of a commit from the conventional
// commit standard.
func NewCommit(commit git.Commit) Commit {
	if commit.Subject == "" {
		return Commit{}
	}
	c := ParseCommitSubject(commit.Subject)
	c.Commit = commit

	scanner := bufio.NewScanner(strings.NewReader(commit.Body))
	var body []string
	for scanner.Scan() {
		body = append(body, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("[error] reading commit message Body error=%v", err)
		return Commit{}
	}
	c.Body = strings.TrimPrefix(strings.Join(body, "\n"), "\n")

	// from the Body check if it's breaking
	c.IsBreaking = c.IsBreaking || strings.Contains(c.Body, "BREAKING CHANGE")

	return c
}

// ParseCommitSubject takes in a Title formatted in the conventional commit paradigm and parses it out as a Commit.
// Parse titles such as:
//	- chore: Remove elastic APM in staging and prod
//	- fix(apm): Add development to environment chart
func ParseCommitSubject(s string) Commit {
	r := regexp.MustCompile(`^([^:]*):\s*(.*)`)
	res := r.FindStringSubmatch(s)

	if res != nil {
		// remove the first element which is the entire matched string
		commitMsg := res[1:]
		//extract the Scope if exists
		c := parseCommitType(commitMsg[0])
		c.Title = commitMsg[1]
		return c
	}

	return Commit{
		Title: s,
	}
}

// parseCommitType Takes in a commit type from the Title of the commit and determines the Scope and type of commit.
// Types would be chore(Scope), or feat! or feat(Scope)!
func parseCommitType(c string) Commit {
	//extract the Scope if exists based on https://www.conventionalcommits.org/en/v1.0.0/#summary
	re := regexp.MustCompile(`([a-zA-Z].*)\((.*?)\)`)
	results := re.FindStringSubmatch(c)
	isBreaking := strings.Contains(c[len(c)-1:], "!")
	if len(results) > 0 {
		commitType := results[1]
		var scope string
		if len(results) > 1 {
			scope = results[2]
		}
		return Commit{
			Type:       CommitType(commitType),
			Scope:      scope,
			IsBreaking: isBreaking,
		}
	}

	if isBreaking {
		c = c[0:len(c)-1]
	}

	return Commit{
		Type:       CommitType(c),
		IsBreaking: isBreaking,
	}
}
