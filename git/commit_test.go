package git

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)


func (s *CommitTestSuite) TestNoPreviousCommits() {
	_, err := s.Git.GetCommitsSinceLatestTag()
	assert.Error(s.T(), err)
}

func (s *CommitTestSuite) TestCreateCommit() {
	//create a commits and a tag
	tempFile := "foobar"
	file, err := ioutil.TempFile(s.Git.WorkDirectory, tempFile)
	if err != nil {
		s.Failf("could not create temp file", "error=%v", err)
	}

	if err = s.Git.Add(file.Name()); err != nil {
		s.Fail("could not stage file", err)
	}

	c, err := s.Git.CreateCommit("this is my first commit", "", false)
	if err != nil {
		s.Error(err)
	}

	s.Equal(Commit{"this is my first commit", "", c.Hash, c.Author, c.Date}, c)
}

func (s *CommitTestSuite) TestGetCommitsSinceLatestTag() {
	_, err := s.Git.CreateCommit("this is my second commit", "", true)
	if err != nil {
		s.Error(err, "could not create second commit")
	}
	err = s.Git.CreateTag("v0.0.1", false)
	if err != nil {
		s.Error(err, "could not create tag")
	}
	c1, _ := s.Git.CreateCommit("this is my third commit", "\nthis is the body\n\nlalala", true)
	c2, _ := s.Git.CreateCommit("this is my fourth commit", "", true)
	c3, _ := s.Git.CreateCommit("this is my fifth commit", "", true)

	commits, err := s.Git.GetCommitsSinceLatestTag()
	if err != nil {
		s.Error(err, "could not get latest ")
	}

	s.Len(commits, 3)
	s.ElementsMatch(commits, []Commit{c1, c2, c3})
}

func TestCommitTestSuite(t *testing.T) {
	suite.Run(t, new(CommitTestSuite))
}

type CommitTestSuite struct {
	GitTestSuite
}
