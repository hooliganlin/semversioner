package git

import (
	"fmt"
	"regexp"
	"strings"
)

// CreateTag creates a git lightweight tag by default. Setting annotated to true will create
// an annotated tag.
func (g Git) CreateTag(tag string, annotated bool) error {
	args := []string{tag}
	if annotated {
		args = append(args, "-a")
	}
	err := g.exec("tag", args...).Run()
	if err != nil {
		return fmt.Errorf("could not create git tag=%s err=%v", tag, err)
	}
	return nil
}

// GetLatestPreReleaseTag fetches the latest abbreviated tag from the git repository.
// Note: This function removes the "g" from the human-readable tag (ie. 1.0.2-4-g123aefd)
func (g Git) GetLatestPreReleaseTag() (string, error) {
	// check if there are any tags
	out, err := g.exec("tag", "--list").Output()
	if err != nil {
		return "", err
	}
	if strings.Trim(string(out), "\n") == "" {
		return "", nil
	}

	out, err = g.exec("describe", "--tags").Output()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	_ , err = b.Write(out)
	if err != nil {
		return "", err
	}
	sanitizedOutput := strings.TrimSuffix(b.String(), "\n")
	//The length of the abbreviation scales as the repository grows, using the approximate number of objects in
	//the repository and a bit of math around the birthday paradox, and defaults to a minimum of 7.
	//The "g" prefix stands for "git" and is used to allow describing the version of a software
	//depending on the SCM the software is managed with. This is useful in an environment where
	//people may use different SCMs.
	hashRegex := regexp.MustCompile(`-g\w+`)
	res := hashRegex.FindString(sanitizedOutput)
	if res != "" {
		sanitizedOutput = strings.ReplaceAll(sanitizedOutput, "-g", "-")
	}

	return sanitizedOutput, nil
}

// GetLatestTag fetches the latest git tag in the tree.
func (g Git) GetLatestTag() (string, error) {
	if ok := g.hasTagHistory(); !ok{
		return "", nil
	}

	out, err := g.exec("describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	_ , err = b.Write(out)
	if err != nil {
		return "", err
	}
	sanitizedOutput := strings.TrimSuffix(b.String(), "\n")
	return sanitizedOutput, nil
}

func (g Git) hasTagHistory() bool {
	out, err := g.exec("tag", "--list").Output()
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(out), "\n") != ""
}