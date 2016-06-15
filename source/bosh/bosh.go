package bosh

import (
	"os/exec"
	"strings"
)

type BoshProxy struct {
	user, password, target string
}

func (b *BoshProxy) Init(target, user, password string) error {
	cmd := exec.Command("bosh", "--no-color", "-u", user, "-p", password, "target", target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}

func (b *BoshProxy) UploadRelease(release string) error {
	cmd := exec.Command("bosh", "--no-color", "upload release", release)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}

func (b *BoshProxy) UploadStemcell(stemcell string) error {
	cmd := exec.Command("bosh", "--no-color", "-u", user, "-p", password, "upload", "stemcell", stemcell)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}

func (b *BoshProxy) Deploy(deploymentPath string) (string, error) {
	cmd := exec.Command("bosh", "--no-color", "-d", deploymentPath, "-n", "-N", "deploy")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", boshError{cmd.Path, cmd.Args, err, out}
	}
	lines := strings.Split(string(out), "\n")
	res := strings.Split(lines[len(lines)-1], " ")
	return res[2], nil
}

func (b *BoshProxy) Status(task string) (string, error) {
	cmd := exec.Command("bosh", "--no-color", "task", task)
	out, err := cmd.Output()
	if err != nil {
		return "", boshError{cmd.Path, cmd.Args, err, out}
	}
	lines := strings.Split(string(out), "\n")
	var l string
	for _, l = range lines {
		if strings.Contains(l, "Task") {
			break
		}
	}
	res := strings.Split(l, " ")
	return res[1], nil
}
func (b *BoshProxy) DeleteDeployment(name string) error {
	cmd := exec.Command("bosh", "--no-color", "-n", "delete deployment", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}
