package cattle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type KubeCtl struct {
	Path string
	Opts string
}

func (k *KubeCtl) GetPath() (string, error) {
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	whichCmd := exec.Command("which", "kubectl")
	whichCmd.Stdout = stdoutBuffer
	whichCmd.Stderr = stderrBuffer
	err := whichCmd.Run()
	if err != nil{
		return stderrBuffer.String(), err
	}
	fmt.Println("Kubectl Found:" + stdoutBuffer.String())

	return stdoutBuffer.String(), err
}

func (k *KubeCtl) Apply(arg string)(string, error){
	if k.Path == "" {
		path, err := k.GetPath()
		if err != nil {
			return "", err
		}
		k.Path = path
	}

	fmt.Println("Path Confirmed" + k.Path)

	k.Opts = "apply"
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	//allArgs := fmt.Sprintf("%s %s", k.Opts, arg)

	//cmd := exec.Command(k.Path, "apply", "-f", arg)
    fmt.Println(arg)
	cmd := exec.Command(k.Path, "apply", "-f", arg)
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err != nil {
		return stderrBuffer.String(), err
	}

	return stdoutBuffer.String(), err

}

func createWorkDir(OutputDirectory string)(string, error){
	workingDirectory, err := ioutil.TempDir(OutputDirectory, "deployment-")
	if err != nil {
		panic(err)
	}
	return workingDirectory, err
}

func ApplyDeployment(fullPath string) (string, error) {
	kubeCtl := &KubeCtl{}
	kubectlApply, err := kubeCtl.Apply(fullPath)
	pApply := fmt.Sprintf("%v", kubectlApply)
	if err != nil{
		panic(err)
	}
	return pApply, err
}
