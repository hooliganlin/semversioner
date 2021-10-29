package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hooliganlin/versioning/conventional-wisdom/conventional"
	"github.com/hooliganlin/versioning/conventional-wisdom/git"
	"log"
)

type versioner struct {
	git git.Git
}
func newVersioner(g git.Git) versioner {
	return versioner{
		git: g,
	}
}

func(v versioner) getVersion(releaseType string, tag string) semver.Version {
	switch releaseType {
	case Patch :
		return semver.MustParse(tag).IncPatch()
	case Minor :
		return semver.MustParse(tag).IncMinor()
	case Major:
		return semver.MustParse(tag).IncMajor()
	case Conventional:
		version, err := conventional.DetermineNextVersion(v.git.WorkDirectory)
		if err != nil {
			log.Fatalf("could not determine next versioner by convetional commits err=%v", err)
		}
		return version
	default:
		latestTag, err := v.git.GetLatestPreReleaseTag()
		if err != nil {
			log.Fatalf("could not get the latest pre release latestTag (snapshot) err=%v", err)
		}
		return *semver.MustParse(fmt.Sprintf("%s-%s", latestTag, "SNAPSHOT"))
	}
}
