package git

import (
	"log"
	"os/exec"
)

type Git struct {
	WorkDirectory string
}

func New(workDir string) Git {
	return Git {
		WorkDirectory: workDir,
	}
}

// IsValidGitDir checks if the current working directory contains a git repository.
func (g Git) IsValidGitDir() bool {
	cmd := g.exec("rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] Invalid git repository. error=%v", err)
		return false
	}
	return true
}

// exec runs the underlying git command with the targeted WorkDirectory.
func (g Git) exec(action string, args... string) *exec.Cmd{
	args = append([]string{"-C", g.WorkDirectory, action}, args...)
	return exec.Command("git", args...)
}
