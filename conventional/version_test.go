package conventional

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hooliganlin/versioning/semversioner/git"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
)
func(s GitTestSuite) TestDetermineNextVersion() {
	s.Run("patch", func() {
		s.runDetermineNextVersionTest("no previous tag",
			[]GenCommitFunction{genFixCommit},
			"",
			"v0.0.1",
		)
		s.runDetermineNextVersionTest("previous patch tag",
			[]GenCommitFunction{genFixCommit},
			"v0.0.1",
			"v0.0.2",
		)
		s.runDetermineNextVersionTest("minor tag",
			[]GenCommitFunction{genFixCommit},
			"v0.1.0",
			"v0.1.1",
		)

		s.runDetermineNextVersionTest("major tag",
			[]GenCommitFunction{genFixCommit},
			"v1.0.0",
			"v1.0.1",
		)

		s.runDetermineNextVersionTest("multiple fixes and one feat",
			[]GenCommitFunction{genFixCommit, genFixCommit, genFeatCommit, genFixCommit},
			"v1.0.0",
			"v1.0.1",
		)
		s.runDetermineNextVersionTest("multiple feats and one fix",
			[]GenCommitFunction{genFeatCommit, genFixCommit, genFeatCommit, genFeatCommit},
			"v1.0.0",
			"v1.0.1",
		)
	})

	s.Run("minor", func() {
		s.runDetermineNextVersionTest("previous patch",
			[]GenCommitFunction{genFeatCommit},
			"v0.0.1",
			"v0.1.0",
		)

		s.runDetermineNextVersionTest("previous minor tag",
			[]GenCommitFunction{genFeatCommit},
			"v0.11.0",
			"v0.12.0",
		)

		s.runDetermineNextVersionTest("previous major tag",
			[]GenCommitFunction{genFeatCommit},
			"v1.11.5",
			"v1.12.0",
		)
	})

	s.Run("major", func() {
		s.runDetermineNextVersionTest("previous patch",
			[]GenCommitFunction{genBreakingCommit},
			"v0.0.1",
			"v1.0.0",
		)

		s.runDetermineNextVersionTest("previous minor tag",
			[]GenCommitFunction{genBreakingCommit},
			"v0.11.0",
			"v1.0.0",
		)

		s.runDetermineNextVersionTest("previous major tag",
			[]GenCommitFunction{genBreakingCommit},
			"v1.2.5",
			"v2.0.0",
		)

		s.runDetermineNextVersionTest("breaking and fix commit",
			[]GenCommitFunction{genBreakingCommit, genFixCommit},
			"v1.2.5",
			"v2.0.0",
		)
	})
}
type GenCommitFunction func(num int) string
func (s GitTestSuite) runDetermineNextVersionTest(
	name string,
	commitFunc []GenCommitFunction,
	initialTag string,
	expectedVersion string) {
	s.SetupTest()
	s.Run(name, func() {
		if _, err := s.Git.CreateCommit(genFeatCommit(0), genCommitBody(), true); err != nil {
			s.Error(err, "could not create initial commit")
		}

		if initialTag != "" {
			if err := s.Git.CreateTag(initialTag, false); err != nil {
				s.Errorf(err, "could not create %s initialTag", initialTag)
			}
		}
		for i, fn := range commitFunc {
			if _, err := s.Git.CreateCommit(fn(i+1), genCommitBody(), true); err != nil {
				s.Errorf(err, "could not create commit %d", i+1)
			}
		}

		v, err := DetermineNextVersion(s.Git.WorkDirectory)
		if err != nil {
			s.Error(err, "could not determine next version")
		}

		s.Equal(*semver.MustParse(expectedVersion), v)
	})
	s.TearDownTest()
}

func TestRunSemverTest(t *testing.T) {
	suite.Run(t, new(GitTestSuite))
}

type GitTestSuite struct {
	suite.Suite
	Git git.Git
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
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
	s.Git = git.New(tempDir)
}

func (s *GitTestSuite) TearDownTest() {
	if err := os.RemoveAll(s.Git.WorkDirectory); err != nil {
		s.T().Fatalf("could not remove temporary directory %s err=%v", s.Git.WorkDirectory, err)
	}
}

func genFixCommit(num int) string {
	prefix := "fix"
	subject := fmt.Sprintf("this is fix number %d", num)
	return fmt.Sprintf("%s: %s", prefix, subject)
}

func genFeatCommit(num int) string {
	prefix := "feat"
	subject := fmt.Sprintf("this is feat number %d", num)
	return fmt.Sprintf("%s: %s", prefix, subject)
}

func genBreakingCommit(num int) string {
	prefix := "feat!"
	subject := fmt.Sprintf("this is a breaking change %d", num)
	return fmt.Sprintf("%s: %s", prefix, subject)
}

func genCommitBody() string {
	words := []string{"lorem", "ipsum", "happy", "brew", "must", "long", "for", "", "a", "something", "\n"}
	randLen := rand.Intn(len(words))
	var sb strings.Builder
	for i := 0; i < randLen; i++ {
		j := rand.Intn(randLen)
		sb.WriteString(words[j] + " ")
	}
	return strings.TrimSuffix(sb.String(), " ")
}