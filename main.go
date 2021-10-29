package main

import (
	"fmt"
	"github.com/hooliganlin/versioning/conventional-wisdom/git"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

const (
	Patch = "patch"
	Minor = "minor"
	Major = "major"
	Conventional = "conventional"
)

type Opts struct {
	WorkDir			string 	`long:"directory" description:"Working directory of a git repository" default:"."`
	Type        	string 	`long:"type" description:"The release type" choice:"major" choice:"minor" choice:"patch" choice:"conventional"`
	Prerelease  	string  `long:"prerelease" description:"The name of the pre-release (ie. alpha, rc)"`
}

func main() {
	var opts Opts
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatalf("could not parse %v", err)
	}

	g := git.New(opts.WorkDir)
	if !g.IsValidGitDir() {
		log.Fatalf("no valid git repo for working directory: %s", opts.WorkDir)
	}

	latestTag, err := g.GetLatestTag()
	if err != nil {
		log.Fatalf("could not fetch latest tag err: %v", err)
	}
	v := newVersioner(g)
	version := v.getVersion(opts.Type, latestTag)

	if opts.Prerelease != "" {
		version, err = version.SetPrerelease(opts.Prerelease)
		if err != nil {
			log.Fatalf("could not set pre release name err=%v", err)
		}
	}
	fmt.Println(version.String())
}
