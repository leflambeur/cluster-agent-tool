package cattle

import "io/ioutil"

func createWorkDir(OutputDirectory string)(string,error){
	workingDirectory, err := ioutil.TempDir(OutputDirectory, "deployment-")
	if err != nil {
		panic(err)
	}
	return workingDirectory, err
}

func applyDeployment()  {

}
