package git

import (
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"testing"
)

func (s GitTestSuite) TestIsValidGitDir() {
	s.True(s.Git.IsValidGitDir())

	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		s.FailNow("could not create temporary directory", err)
	}
	g := New(tempDir)
	s.False(g.IsValidGitDir())
}

func TestRunGit(t *testing.T) {
	suite.Run(t, new(GitTestSuite))
}

type GitTestSuite struct {
	suite.Suite
	Git Git
}

// SetupTest creates the initial git repository in a temporary directory for subsequent tests to use in testing.
func (s *GitTestSuite) SetupTest() {
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
	s.Git = New(tempDir)
}

func (s *GitTestSuite) TearDownTest() {
	if err := os.RemoveAll(s.Git.WorkDirectory); err != nil {
		s.T().Fatalf("could not remove temporary directory %s err=%v", s.Git.WorkDirectory, err)
	}
}