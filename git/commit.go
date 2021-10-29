package git

import (
	"fmt"
	"github.com/Masterminds/semver"
	"strings"
	"time"
)

type Author struct {
	Name string
	Email string
}

type Commit struct {
	Subject string
	Body string
	Hash string
	Author Author
	Date time.Time
}

// commitLogFormat is a tab delimited format of a git commit message.
const commitSeparator = "~~"
const commitLogFormat = "%+cI%+H%+an%+ae%+s%+b" + commitSeparator

// GetCommitsSinceLatestTag fetches all the commits since the latest tag.
func (g Git) GetCommitsSinceLatestTag() ([]Commit, error) {
	describeArgs := []string{"--tags", "--abbrev=0"}
	output, err := g.exec("describe", describeArgs...).Output()
	if err != nil {
		return []Commit{}, err
	}
	tag := strings.TrimSuffix(string(output), "\n")
	_, err = semver.NewVersion(tag)
	if err != nil {
		return []Commit{}, err
	}

	commits, err := g.parseRawCommits([]string{fmt.Sprintf("%s..HEAD", tag)})
	if err != nil {
		return []Commit{}, err
	}
	return commits, nil
}

// Add stages a file to be tracked by git.
func (g Git) Add(file string) error{
	err := g.exec("add", file).Run()
	if err != nil {
		return fmt.Errorf("could not stage file err=%v", err)
	}
	return nil
}

// CreateCommit creates a commit with a passed in message and whether or not to allow it be to empty
func (g Git) CreateCommit(subject string, body string, allowEmpty bool) (Commit, error) {
	args := []string{"--quiet", "--cleanup", "strip", "-m", subject, "-m", body}
	if allowEmpty {
		args = append(args, "--allow-empty")
	}
	err := g.exec("commit", args...).Run()
	if err != nil {
		return Commit{}, fmt.Errorf("could not create commit err=%v", err)
	}
	commits, err := g.parseRawCommits([]string{"-n1"})
	if err != nil {
		return Commit{}, err
	}
	if len(commits) != 1 {
		return Commit{}, nil
	}
	return commits[0], nil
}

// parseRawCommits takes a list of git log arguments and parses each commit from the git log output
// and converts them to a list of Commit.
func (g Git) parseRawCommits(args []string) ([]Commit, error) {
	args = append([]string{fmt.Sprintf(`--format=%s`, commitLogFormat)}, args...)
	out, err := g.exec("log", args...).Output()
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	_ , err = b.Write(out)
	if err != nil {
		return nil, err
	}

	rawCommits := splitAndFilter(b.String(), commitSeparator)
	commits := make([]Commit, 0)
	for _, line := range rawCommits {
		tokens := strings.Split(strings.Trim(line, "\n"), "\n")
		date, err := time.Parse(time.RFC3339, tokens[0])
		if err != nil {
			return nil, err
		}
		body := ""
		if len(tokens) > 5 {
			for _, s := range tokens[5:] {
				body += fmt.Sprintln(s)
			}
		}
		c := Commit{
			Hash: tokens[1],
			Author: Author {
				Name:  tokens[2],
				Email: tokens[3],
			},
			Subject: tokens[4],
			Body:    strings.TrimSuffix(body, "\n"),
			Date:    date,
		}
		commits = append(commits, c)
	}
	return commits, nil
}

// splitAndFilter takes in a string and filters out any empty string or new line
func splitAndFilter(s string, separator string)[]string {
	lines := strings.Split(s, separator)
	results := make([]string, 0)
	for _, l := range lines {
		if l == "\n" || l == "" || l == "\t" {
			continue
		}
		results = append(results, l)
	}
	return results
}