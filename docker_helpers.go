package dkrrmtag

import (
	"os"
	"os/exec"
)

func dockerPull(imgTag string) error {
	cmd := exec.Command("docker", "pull", imgTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerPush(imgTag string) error {
	cmd := exec.Command("docker", "push", imgTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
