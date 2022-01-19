package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

const testJobManifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4`

type JobExecuter struct {
	Client *kubernetes.Clientset
	Job    *v1.Job
}

func NewJobExecuter() (*JobExecuter, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	executer := &JobExecuter{Client: client}

	return executer, nil
}
func (e *JobExecuter) ExecuteJob(namespace string) error {

	jobBytes, err := ioutil.ReadFile("/Users/olegv/dev/tools/applier/samples/job.yaml")

	//bytes := []byte(jobManifest)
	//jobConfig := &v1config.JobApplyConfiguration{}
	var job *v1.Job

	//jobMap := map[string]interface{}{}
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, groupVersionKind, err := decoder.Decode(
		[]byte(jobBytes),
		nil,
		nil)
	options := metav1.CreateOptions{}

	if groupVersionKind.Kind == "Job" {
		job = obj.(*v1.Job)
		job.Namespace = "default"

	}

	//jobConfig = jobConfig.WithName(job.Name).WithNamespace(job.Namespace).WithSpec(job.Spec)
	client := e.Client
	result, err := client.BatchV1().Jobs("default").Create(context.Background(), job, options)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("job %v has status %v", result.Name, result.Status.Ready)
	listOpts := metav1.ListOptions{}

	labelSelector := metav1.LabelSelector{}
	labelSelector.MatchLabels = result.Labels
	selector, _ := metav1.LabelSelectorAsSelector(&labelSelector)

	listOpts.LabelSelector = selector.String()
	//..logstream := client.CoreV1().Pods("default").GetLogs()
	watch, err := client.BatchV1().Jobs(namespace).Watch(context.Background(), listOpts)

	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case ev := <-watch.ResultChan():

			fmt.Printf("ev [%v]\n", ev.Type)
			if ev.Object != nil {
				job, ok := ev.Object.(*v1.Job)
				if !ok {

					return errors.New("can't convet to JJob")
				}
				conds := job.Status.Conditions
				for _, c := range conds {
					fmt.Printf("cond %v\n ", c)
					if c.Type == v1.JobFailed || c.Type == v1.JobSuspended {
						msg := fmt.Sprintf("job status[%v] with reason [%v] , %v", c.Type, c.Reason, c.Message)
						return errors.New(msg)
					}
					if c.Type == v1.JobFailed || c.Type == v1.JobComplete {
						return nil
					}
				}

			}
		}
	}

}
func BringInstallScript(url string, clusterId string, token string) (string, error) {

	c := http.Client{Timeout: time.Duration(10) * time.Second}
	fullUrl := fmt.Sprintf("%s?clusterId=%s", url, clusterId)

	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return "", err

	}
	//req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}
func (e *JobExecuter) RunTestJob() (*v1.Job, error) {
	return e.RunJobFromManifest(testJobManifest)
}

func (e *JobExecuter) Run(job *v1.Job) error {
	options := metav1.CreateOptions{}
	createJob, err := e.Client.BatchV1().Jobs("default").Create(context.Background(), job, options)

	e.Job = createJob
	return err
}
func (e *JobExecuter) RunJobFromManifest(manifest string) (*v1.Job, error) {

	jobBytes := []byte(manifest)

	var job *v1.Job

	//jobMap := map[string]interface{}{}
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, groupVersionKind, err := decoder.Decode(
		[]byte(jobBytes),
		nil,
		nil)
	options := metav1.CreateOptions{}

	if groupVersionKind.Kind == "Job" {
		job = obj.(*v1.Job)
		job.Namespace = "default"

	}
	result, err := e.Client.BatchV1().Jobs("default").Create(context.Background(), job, options)

	e.Job = job
	return result, err

}

func (e *JobExecuter) waitForCompletion() {

}
func (e *JobExecuter) GetJob(fromCache bool) (*v1.Job, error) {
	if fromCache {
		return e.Job, nil
	}
	opts := metav1.GetOptions{}
	job, err := e.Client.BatchV1().Jobs(e.Job.Namespace).Get(context.Background(), e.Job.Name, opts)

	if err != nil {
		e.Job = job
	}

	return job, err
}
func (e *JobExecuter) SetJob(job *v1.Job) *JobExecuter {

	e.Job = job

	return e
}

func (e *JobExecuter) IsJobCompleted() (done bool, status bool, err error) {

	job, err := e.GetJob(false)
	if err != nil {
		return false, false, err
	}
	conds := job.Status.Conditions
	for _, c := range conds {
		fmt.Printf("cond %v\n ", c)

		if c.Type == v1.JobFailed || c.Type == v1.JobSuspended {
			msg := fmt.Sprintf("job status[%v] with reason [%v] , %v", c.Type, c.Reason, c.Message)
			fmt.Printf(msg)
			return true, false, nil
		}
		if c.Type == v1.JobComplete {
			switch c.Status {
			case corev1.ConditionTrue:
				return true, true, nil
			case corev1.ConditionFalse:
				return true, false, nil
			case corev1.ConditionUnknown:
				return true, false, errors.New("unknown status")
			}
		}
	}

	return false, false, nil
}
func (e *JobExecuter) GetJobPods() (*corev1.PodList, error) {

	job, _ := e.GetJob(true) //bring from cache

	labelSelector := metav1.LabelSelector{}
	labelSelector.MatchLabels = job.Spec.Template.Labels
	selector, _ := metav1.LabelSelectorAsSelector(&labelSelector)
	opts := metav1.ListOptions{LabelSelector: selector.String()}
	return e.Client.CoreV1().Pods(job.Namespace).List(context.Background(), opts)

}
func (e *JobExecuter) DeleteJob(ns string, jobName string) error {
	opts := metav1.DeleteOptions{}
	err := e.Client.BatchV1().Jobs(ns).Delete(context.Background(), jobName, opts)

	return err
}
func (e *JobExecuter) ReadLogs(ns string, podName string) error {
	opts := corev1.PodLogOptions{Follow: true}
	logStream := e.Client.CoreV1().Pods(ns).GetLogs(podName, &opts)
	closer, err := logStream.Stream(context.Background())
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(closer)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf(line)
	}

	return nil

}
