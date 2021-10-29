package conventional

import (
	"github.com/Masterminds/semver"
	"github.com/hooliganlin/versioning/conventional-wisdom/git"
)
const initialTag = "v0.0.0"
// DetermineNextVersion leverages the conventional commit style logs to determine
// the next semantic version based on the commits since the latest tag.
// Any commits with a fix will always warrant a patch.
func DetermineNextVersion(workDir string) (semver.Version, error) {
	g := git.New(workDir)
	latestTag, err := g.GetLatestTag()
	if err != nil {
		return semver.Version{}, err
	}
	if latestTag == "" {
		v, err := semver.NewVersion(initialTag)
		if err != nil {
			return semver.Version{}, err
		}
		return v.IncPatch(), nil
	}

	v, err := semver.NewVersion(latestTag)
	if err != nil {
		return semver.Version{}, err
	}
	commits, err := g.GetCommitsSinceLatestTag()
	if err != nil {
		return semver.Version{}, err
	}

	breaking, nonBreaking := partitionCommits(mapCommits(commits, NewCommit), isBreakingCommit)
	if len(breaking) > 0 {
		return v.IncMajor(), nil
	}
	fixes, _ :=  partitionCommits(nonBreaking, hasFixCommit)
	if len(fixes) > 0 {
		return v.IncPatch(), nil
	}

	return v.IncMinor(), nil
}

func isBreakingCommit(c Commit) bool {
	return c.IsBreaking
}
func hasFixCommit(c Commit) bool {
	return c.Type == Fix
}
// partitionCommits takes in a slice c and applies a conditional function, f to return the results as two slices. The left
// slice is the result of a truthy result from f. A falsy result yields the right slice.
func partitionCommits(c []Commit, f func(c Commit) bool) ([]Commit, []Commit) {
	var left []Commit
	var right []Commit
	for _, e := range c {
		if f(e) {
			left = append(left, e)
		} else {
			right = append(right, e)
		}
	}
	return left, right
}

func mapCommits(c []git.Commit, f func(c git.Commit) Commit) []Commit {
	r := make([]Commit, len(c), cap(c))
	for i, e := range c {
		r[i] = f(e)
	}
	return r
}
