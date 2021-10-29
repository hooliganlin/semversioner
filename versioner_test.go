package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hooliganlin/versioning/conventional-wisdom/git"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"testing"
)

func(s *VersionerTestSuite) TestGetVersion()  {
	v := newVersioner(s.Git)

	version := v.getVersion(Major, "v2.3.4")
	s.Equal(semver.MustParse("v3.0.0"), &version)

	version = v.getVersion(Minor, "v2.3.4")
	s.Equal(semver.MustParse("v2.4.0"), &version)

	version = v.getVersion(Patch, "v2.3.4")
	s.Equal(semver.MustParse("v2.3.5"), &version)

	_, _ = v.git.CreateCommit("feat: feature 1", "this is body", true)
	err := v.git.CreateTag("v0.0.1", false)
	if err != nil {
		s.FailNow("np no no nonono")
	}

	_, _ = v.git.CreateCommit("doc: docs 1", "", true)
	c, _ := v.git.CreateCommit("chore: chore 1", "", true)
	tag, err := v.git.GetLatestTag()
	if err != nil {
		s.FailNow("np no no nonono")
	}

	version = v.getVersion(Conventional, tag)
	s.Equal(semver.MustParse("v0.1.0"), &version)

	version = v.getVersion("snapshot", tag)
	s.Equal(semver.MustParse(fmt.Sprintf("v0.0.1-%d-%s-SNAPSHOT", 2, c.Hash[:7])), &version)
}

func TestRunVersioner(t *testing.T) {
	suite.Run(t, new(VersionerTestSuite))
}

type VersionerTestSuite struct {
	suite.Suite
	Git git.Git
}

func (s *VersionerTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		s.FailNow("could not create temporary directory", err)
	}
	if err = os.Chdir(tempDir); err != nil {
		s.FailNow("could not change to temporary directory", err)
	}

	if err != nil {
		s.FailNow("cannot setup working directory", err)
	}
	if err = exec.Command("git", "init").Run(); err != nil {
		s.FailNow("could not initialize git repository err", err)
	}
	if err != nil {
		if err = os.RemoveAll(tempDir); err != nil {
			s.FailNowf("could not remove working directory", "dir=%s err=%v", tempDir, err)
		}
		s.FailNow("cannot create new git repo", err)
	}
	s.Git = git.New(tempDir)
}

func (s *VersionerTestSuite) TearDownTest() {
	if err := os.RemoveAll(s.Git.WorkDirectory); err != nil {
		s.T().Fatalf("could not remove temporary directory %s err=%v", s.Git.WorkDirectory, err)
	}
}
