package utils

import (
	"bytes"
	"fmt"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/scheme"
)

type InstallVistior interface {
	InstallJob(job runtime.Object) error
	InstallSA(job runtime.Object) error
	InstallOther(job runtime.Object) error
	Install(info *resource.Info) error
}

type TestInstaller struct {
	InstallVistior
}

func (i *TestInstaller) Install(info *resource.Info) error {
	gvk := info.Object.GetObjectKind().GroupVersionKind()

	kind := gvk.Kind

	switch kind {
	case "Job":
		return i.InstallJob(info)
	}
	return nil
}
func (i *TestInstaller) InstallJob(info *resource.Info) error {
	fmt.Printf("installing %s \n ", info.ObjectName())
	jobExec, err := NewJobExecuter()
	if err != nil {
		return err
	}
	//TODOL remove it

	err = jobExec.Run(info.Object.(*v1.Job))
	return err
}
func (i *TestInstaller) InstallSA(info *resource.Info) error {
	return nil
}
func (i *TestInstaller) InstallOther(info *resource.Info) error {
	return nil
}

func RunLocalBuilder(manifest string) {
	// Create a local builder...
	builder := resource.NewLocalBuilder().
		// Configure with a scheme to get typed objects in the versions registered with the scheme.
		// As an alternative, could call Unstructured() to get unstructured objects.
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		// Provide input via a Reader.
		// As an alternative, could call Path(false, "/path/to/file") to read from a file.
		Stream(bytes.NewBufferString(manifest), "input").
		// Flatten items contained in List objects
		Flatten().
		// Accumulate as many items as possible
		ContinueOnError()

	// Run the builder
	result := builder.Do()

	if err := result.Err(); err != nil {
		fmt.Println("builder error:", err)
		return
	}
	installer := TestInstaller{}
	result.Visit(func(i *resource.Info, e error) error {
		fmt.Printf("visiting %s (%T)\n", i.String(), i.Object)
		return installer.Install(i)

	})

	// Output:
	// Name: "mutating1", Namespace: "" (*v1.MutatingWebhookConfiguration)
	// Name: "mutating2", Namespace: "" (*v1.MutatingWebhookConfiguration)
	// Name: "mutating3", Namespace: "" (*v1.MutatingWebhookConfiguration)
	// Name: "validating1", Namespace: "" (*v1.ValidatingWebhookConfiguration)
	// Name: "validating2", Namespace: "" (*v1.ValidatingWebhookConfiguration)
	// Name: "validating3", Namespace: "" (*v1.ValidatingWebhookConfiguration)
	// Name: "mutating4", Namespace: "" (*v1.MutatingWebhookConfiguration)
	// Name: "validating4", Namespace: "" (*v1.ValidatingWebhookConfiguration)
}
