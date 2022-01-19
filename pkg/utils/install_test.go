package utils

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func TestInstall(result *testing.T) {
	jobExec, err := NewJobExecuter()
	if err != nil {
		result.Fatal(err)
	}
	jobExec.ExecuteJob("default")
}

func HandleError(result *testing.T, err error) {
	if err != nil {
		result.Fatal(err)
	}
}
func TestBringInstallScript(result *testing.T) {
	//curl -X POST 'https://api.spotinst.io/ocean/cd/clusterInstaller?clusterId=<MY_CLUSTER_ID>' --header 'Authorization: Bearer <API_TOKEN>'
	token := "79b8b542e613a96ae282c2e10cc328ef98afd89bd5a778078605e7808b8892ec"
	url := "https://api.spotinst.io/ocean/cd/clusterInstaller"
	clusterId := "temp_oleg_1"
	manifest, err := BringInstallScript(url, clusterId, token)

	if err != nil {
		result.Fatal(err)
	}

	fmt.Println(manifest)
}
func TestJobLogs(result *testing.T) {
	fmt.Printf("break cache1")
	jobExec, err := NewJobExecuter()
	HandleError(result, err)

	jobBytes, err := ioutil.ReadFile("/Users/olegv/dev/tools/applier/samples/job.yaml")
	HandleError(result, err)
	job, err := jobExec.RunJobFromManifest(string(jobBytes))
	HandleError(result, err)

	defer jobExec.DeleteJob(job.Namespace, job.Name)
	HandleError(result, err)

	pods, err := jobExec.GetJobPods()
	HandleError(result, err)
	if len(pods.Items) == 0 {
		fmt.Printf("no pods assosiated with job %v/%v", job.Namespace, job.Name)
	}
	stopChannel := make(<-chan struct{})
	var condFunc wait.ConditionFunc

	condFunc = func() (bool, error) {
		done, status, err := jobExec.IsJobCompleted()
		fmt.Printf("job status is %v", status)
		if len(pods.Items) > 0 {
			jobExec.ReadLogs(job.Namespace, pods.Items[0].Name)
		}
		return done, err
	}
	err = wait.PollImmediateUntil(5*time.Second, condFunc, stopChannel)

}
