package utils

import (
	"fmt"
	"testing"
)

const exampleManifest = `
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: mutating1
---
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfigurationList
items:
- apiVersion: admissionregistration.k8s.io/v1
  kind: MutatingWebhookConfiguration
  metadata:
    name: mutating2
- apiVersion: admissionregistration.k8s.io/v1
  kind: MutatingWebhookConfiguration
  metadata:
    name: mutating3
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: validating1
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfigurationList
items:
- apiVersion: admissionregistration.k8s.io/v1
  kind: ValidatingWebhookConfiguration
  metadata:
    name: validating2
- apiVersion: admissionregistration.k8s.io/v1
  kind: ValidatingWebhookConfiguration
  metadata:
    name: validating3
---
apiVersion: v1
kind: List
items:
- apiVersion: admissionregistration.k8s.io/v1
  kind: MutatingWebhookConfiguration
  metadata:
    name: mutating4
- apiVersion: admissionregistration.k8s.io/v1
  kind: ValidatingWebhookConfiguration
  metadata:
    name: validating4
---
apiVersion: batch/v1
kind: Job
metadata:

  labels:
    app: spot-installer
  name: spot-installer
  #uid: b7c738c8-5b16-4c62-83b8-3a1684845e93
spec:
  backoffLimit: 1
  completions: 1
  parallelism: 1
 #selector:
    #matchLabels:
      #controller-uid: b7c738c8-5b16-4c62-83b8-3a1684845e93
  template:
    metadata:
      creationTimestamp: null
      labels:
     #   controller-uid: b7c738c8-5b16-4c62-83b8-3a1684845e93
        job-name: spot-installer
   
    spec:
      containers:
      - args:
        - --service
        - spot-oceancd-controller-svc
        - --webhook
        - controller.oceancd.spot.io
        - --secret
        - spot-oceancd-controller-tls
        - --namespace
        - oceancd
        command:
        - ./installation.sh
        env:
        - name: TOKEN
          value: YmRiZTdlOGMzNDA0OGI5ZGYzYzQ1NTdhMTg0YjdhZmRhZGM2Y2Q2NDYzMjMzODRlYmRmNGM5NjliNzg2OGU0NzUzMGU3NTkyOWZlYTczZDNhYTE1ZDBlYmRkYmMzOWIwZjhjODViY2QzMjQ0MWVlNjhkZGQxMWUyZGE3N2JmMGM=
        - name: INSTALL_URL
          value: https://raw.githubusercontent.com/spotinst/spot-oceancd-releases/main
        - name: SAAS_URL
          value: aHR0cHM6Ly9jbHVzdGVyLWdhdGV3YXkub2NlYW5jZC5pbw==
        image: spotinst/spot-oceancd-controller-installer
        imagePullPolicy: Always
        name: spot-oceancd-controller-installer
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Never
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: spot-oceancd-controller-installer-service-account
      serviceAccountName: spot-oceancd-controller-installer-service-account
      terminationGracePeriodSeconds: 30
status:
  completionTime: "2021-12-13T03:47:02Z"
  conditions:
  - lastProbeTime: "2021-12-13T03:47:02Z"
    lastTransitionTime: "2021-12-13T03:47:02Z"
    status: "True"
    type: Complete
  startTime: "2021-12-13T03:46:41Z"
  succeeded: 1
---
`

func TestManifestReader(result *testing.T) {
	fmt.Println("1")
	r, err := RunLocalBuilder(exampleManifest)
	if err != nil {
		result.Fatal(err)
	}
	res := NewResources()
	r.Visit(res.NewObjVisitor("Job"))
	for _, info := range res.Infos {
		fmt.Printf("%v\n", info.Name)
	}
}

func TestManifestFromUrl(result *testing.T) {
	token := "79b8b542e613a96ae282c2e10cc328ef98afd89bd5a778078605e7808b8892ec"
	url := "https://api.spotinst.io/ocean/cd/clusterInstaller"
	clusterId := "temp_oleg_1"
	manifest, err := BringInstallScript(url, clusterId, token)

	if err != nil {
		result.Fatal(err)
	}
	res := NewResources()
	r, err := RunLocalBuilder(manifest)
	if err != nil {
		result.Fatal(err)
	}
	r.Visit(res.NewObjVisitor("Job"))

	for _, info := range res.Infos {
		fmt.Printf("%v\n", info.Name)
	}
}
