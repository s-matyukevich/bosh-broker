package bosh

import (
	"fmt"
	"os/exec"
	"strings"
)

type BoshProxy struct {
}

func (b *BoshProxy) Init(target, user, password string) error {
	cmd := exec.Command(fmt.Sprintf("bosh --no-color -u %s -p %s target %s", user, password, target))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, err, out}
	}
	return nil
}

func (b *BoshProxy) UploadRelease(release string) error {
	cmd := exec.Command(fmt.Sprintf("bosh --no-color upload release %s", release))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, err, out}
	}
	return nil
}

func (b *BoshProxy) UploadStemcell(stemcell string) error {
	cmd := exec.Command(fmt.Sprintf("bosh --no-color upload stemcell %s", stemcell))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, err, out}
	}
	return nil
}

func (b *BoshProxy) Deploy(deploymentPath string) (string, error) {
	cmd := exec.Command(fmt.Sprintf("bosh --no-color -d '%s' -n -N deploy", deploymentPath))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", boshError{cmd.Path, err, out}
	}
	lines := strings.Split(string(out), "\n")
	res := strings.Split(lines[len(lines)-1], " ")
	return res[2], nil
}

func (b *BoshProxy) Status(task string) (string, error) {
	cmd := exec.Command(fmt.Sprintf("bosh --no-color task %s | grep Task | awk '{print $3}'", task))
	out, err := cmd.Output()
	if err != nil {
		return "", boshError{cmd.Path, err, out}
	}
	return "", nil
}
func (b *BoshProxy) DeleteDeployment(name string) error {
	cmd := exec.Command(fmt.Sprintf("bosh --no-color -n delete deployment %s", name))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return boshError{cmd.Path, err, out}
	}
	return nil
}
