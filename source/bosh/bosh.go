package bosh

import (
	"os/exec"
	"strings"
)

type BoshProxy struct {
	user, password, target string
}

func NewBoshProxy(target, user, password string) (*BoshProxy, string, error) {
	proxy := &BoshProxy{user, password, target}
	cmd := exec.Command("bosh", "--no-color", "-u", user, "-p", password, "target", target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", boshError{cmd.Path, cmd.Args, err, out}
	}
	cmd = exec.Command("bosh", "--no-color", "-u", user, "-p", password, "status", "--uuid")
	out, err = cmd.Output()
	if err != nil {
		return nil, "", boshError{cmd.Path, cmd.Args, err, out}
	}
	return proxy, string(out), nil
}

func (b *BoshProxy) UploadRelease(release string) error {
	cmd := exec.Command("bosh", "--no-color", "-u", b.user, "-p", b.password, "upload", "release", release)
	out, err := cmd.CombinedOutput()
	if strings.Contains(string(out), "already exists") {
		return nil
	}
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}

func (b *BoshProxy) UploadStemcell(stemcell string) error {
	cmd := exec.Command("bosh", "--no-color", "-u", b.user, "-p", b.password, "upload", "stemcell", stemcell)
	out, err := cmd.CombinedOutput()
	if strings.Contains(string(out), "already exists") {
		return nil
	}
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}

func (b *BoshProxy) Deploy(deploymentPath string) (string, error) {
	cmd := exec.Command("bosh", "--no-color", "-u", b.user, "-p", b.password, "-d", deploymentPath, "-n", "-N", "deploy")
	out, err := cmd.CombinedOutput()
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

func (b *BoshProxy) Status(task string) (string, error) {
	cmd := exec.Command("bosh", "--no-color", "-u", b.user, "-p", b.password, "task", task)
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
	return res[2], nil
}
func (b *BoshProxy) DeleteDeployment(name string) error {
	cmd := exec.Command("bosh", "--no-color", "-u", b.user, "-p", b.password, "-n", "delete", "deployment", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, cmd.Args, err, out}
	}
	return nil
}
