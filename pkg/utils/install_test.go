package utils

import (
	"errors"
	"fmt"
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
func TestJobLogs(result *testing.T) {

	jobExec, err := NewJobExecuter()
	HandleError(result, err)

	job, err := jobExec.RunTestJob("")
	defer jobExec.DeleteJob(job.Namespace, job.Name)
	HandleError(result, err)

	pods, err := jobExec.GetJobPods()
	HandleError(result, err)
	if len(pods.Items) == 0 {
		HandleError(result, errors.New("no pods assosiated with job"))
	}
	stopChannel := make(<-chan struct{})
	var condFunc wait.ConditionFunc

	condFunc = func() (bool, error) {
		done, status, err := jobExec.IsJobCompleted()
		fmt.Printf("job status is %v", status)
		jobExec.ReadLogs(job.Namespace, pods.Items[0].Name)
		return done, err
	}
	err = wait.PollImmediateUntil(5*time.Second, condFunc, stopChannel)

}
