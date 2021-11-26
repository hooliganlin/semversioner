package conventional

import (
	"github.com/hooliganlin/versioning/semversioner/git"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCommit(t *testing.T) {
	commit := git.Commit{
		Subject: "fix: something happened and this should be fixed",
		Body: "I found this bug and fixed it\n\nSigned-of-by: megatron",
		Hash:    "deadbeef",
		Author:  git.Author{ Name: "megatron", Email: "megatron@email.com"},
		Date:    time.Time{},
	}
	c := NewCommit(commit)

	assert.Equal(t, Commit{
		Type:       Fix,
		Scope:      "",
		Title:      "something happened and this should be fixed",
		Body:       "I found this bug and fixed it\n\nSigned-of-by: megatron",
		IsBreaking: false,
		Commit:     commit,
	}, c)

	commit.Subject = "This is not conventionalCommit"
	commit.Body = "I found this bug and fixed it\n\nwhoa whoa!"
	c = NewCommit(commit)
	assert.Equal(t, Commit{
		Type:       "",
		Scope:      "",
		Title:      "This is not conventionalCommit",
		Body:       "I found this bug and fixed it\n\nwhoa whoa!",
		IsBreaking: false,
		Commit:     commit,
	}, c)

	emptyCommit := git.Commit{
		Subject: "",
		Body: "",
		Hash:    "",
		Author:  git.Author{},
		Date:    time.Time{},
	}
	c = NewCommit(emptyCommit)
	assert.Equal(t, Commit{}, c)
}

func TestParseCommitTitle(t *testing.T) {
	commitWithScope := "chore(logging): First commitTitle"
	commitWithNoScope := "feat: First feature"
	commitWithBreaking := "feat!: this is going to break"
	commitWithFix := "fix: remove apples from oranges"
	nonConventionalCommit := "this is my first commit"
	unmappedType := "build: updating the ci"
	multiLines := "build: multiple lines\nfix: this is not correct"

	commitTitle := ParseCommitSubject(commitWithScope)
	assert.Equal(t, Commit{Type: "chore", Scope: "logging", Title: "First commitTitle"}, commitTitle)

	commitTitle = ParseCommitSubject(commitWithNoScope)
	assert.Equal(t, Commit{ Type: Feature, Title: "First feature"}, commitTitle)

	commitTitle = ParseCommitSubject(commitWithBreaking)
	assert.Equal(t, Commit{ Type: Feature, Title: "this is going to break", IsBreaking: true}, commitTitle)

	commitTitle = ParseCommitSubject(commitWithFix)
	assert.Equal(t, Commit{ Type: Fix, Title: "remove apples from oranges"}, commitTitle)

	commitTitle = ParseCommitSubject(nonConventionalCommit)
	assert.Equal(t, Commit{ Type: "", Title: "this is my first commit"}, commitTitle)

	commitTitle = ParseCommitSubject(unmappedType)
	assert.Equal(t, Commit{ Type: "build", Title: "updating the ci"}, commitTitle)

	commitTitle = ParseCommitSubject(multiLines)
	assert.Equal(t, Commit{ Type: "build", Title: "multiple lines"}, commitTitle)
}

func TestParseCommitType(t *testing.T) {
	validCommitTypeWithScope := "feat(build)"
	c := parseCommitType(validCommitTypeWithScope)
	assert.Equal(t, Commit {Scope: "build", Type: Feature}, c)

	noScopeCommit := "fix"
	c = parseCommitType(noScopeCommit)
	assert.Equal(t, Commit {Type: Fix}, c)

	otherType := "chore"
	c = parseCommitType(otherType)
	assert.Equal(t, Commit {Type: "chore"}, c)

	breakingCommit := "chore!"
	c = parseCommitType(breakingCommit)
	assert.Equal(t, Commit {Type: "chore", IsBreaking: true}, c)

	manyExclamations := "feat!!!!!!"
	c = parseCommitType(manyExclamations)
	assert.Equal(t, Commit {Type: "feat!!!!!", IsBreaking: true}, c)

	breaking := "feat(something)!"
	c = parseCommitType(breaking)
	assert.Equal(t, Commit {Type: Feature, IsBreaking: true, Scope: "something"}, c)

	misplacedBreaking := "feat!(something)"
	c = parseCommitType(misplacedBreaking)
	assert.Equal(t, Commit {Type: "feat!", IsBreaking: false, Scope: "something"}, c)
}

func TestParseBodyBreakingChange(t *testing.T) {
	commit := git.Commit{
		Subject: "feat: square peg is now a circle",
		Body: "I found this bug and fixed it\n\nBREAKING CHANGE: This requires a circle\nSigned-of-by: megatron",
		Hash:    "deadbeef",
		Author:  git.Author{ Name: "megatron", Email: "megatron@email.com"},
		Date:    time.Time{},
	}
	c := NewCommit(commit)
	assert.True(t, c.IsBreaking)
}