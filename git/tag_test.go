package git

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func (s TagTestSuite) GetLatestPreReleaseTag() {
	if _, err := s.Git.CreateCommit("first commit", "", true); err != nil {
		s.Error(err, "could not create commit")
	}
	if _, err := s.Git.CreateCommit("second commit", "", true); err != nil {
		s.Error(err, "could not create commit")
	}

	if err := s.Git.CreateTag("v0.1.0", false); err != nil {
		s.Error(err, "could not create tag")
	}
	if _, err := s.Git.CreateCommit("third commit", "", true); err != nil {
		s.Error(err, "could not create commit")
	}

	if err := s.Git.CreateTag("v0.2.0", false); err != nil {
		s.Error(err, "could not create tag")
	}

	tag, err := s.Git.GetLatestPreReleaseTag()
	if err != nil {
		s.Error(err, "could not get latest tag")
	}
	s.Equal("v0.2.0", tag)

	_, err = s.Git.CreateCommit("fourth commit", "", true)
	if err != nil {
		s.Error(err, "could not create commit")
	}

	c, err := s.Git.CreateCommit("fifth commit", "", true)
	if err != nil {
		s.Error(err, "could not create commit")
	}

	tag, err = s.Git.GetLatestPreReleaseTag()
	if err != nil {
		s.Error(err, "could not get latest tag")
	}
	s.Equal(fmt.Sprintf("v0.2.0-2-%s", c.Hash[0:7]), tag)
}

func TestTagTestSuite(t *testing.T) {
	suite.Run(t, new(TagTestSuite))
}

type TagTestSuite struct {
	GitTestSuite
}

