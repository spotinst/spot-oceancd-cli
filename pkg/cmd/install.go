package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/verchol/applier/pkg/utils"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	//logger = log.New(os.Stderr)
}
func InstallFromDir(dir string) error {

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Printf("installing resource from %s\n", color.GreenString(f.Name()))
	}

	return nil
}

type Installer struct {
	Client    *kubernetes.Clientset
	Namespace string
}
type InstallerFunc func(context.Context, *resource.Info) error

var InstalllerMap map[string]InstallerFunc

func init() {
	InstalllerMap = map[string]InstallerFunc{}
}
func NewInstaller() (*Installer, error) {

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	i := &Installer{Client: client}
	if len(InstalllerMap) == 0 {
		InstalllerMap["Namespace"] = i.InstallNamespace
		InstalllerMap["ServiceAccount"] = i.InstallSa
		InstalllerMap["Role"] = i.InstallRole
		InstalllerMap["RoleBinding"] = i.InstallRoleBinding
		InstalllerMap["Job"] = i.InstallJob
		InstalllerMap["ClusterRole"] = i.InstallClusterRole
		InstalllerMap["ClusterRoleBinding"] = i.InstallClusterRoleBinding
	}
	return i, nil
}
func (i Installer) InstallVisitor(kind string) InstallerFunc {

	return func(ctx context.Context, info *resource.Info) error {
		fmt.Printf("visiting %s (%T)\n", info.String(), info.Object)
		//return installer.Install(i)
		gvk := info.Object.GetObjectKind().GroupVersionKind()
		objKind := gvk.Kind
		ns := info.Namespace
		if i.Namespace != "" {
			ns = i.Namespace
		}
		innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(info.Object)
		if err != nil {
			return err
		}
		u := unstructured.Unstructured{Object: innerObj}
		// Now the `u` variable has all the meta info available with simple getters.
		// Sample:
		u.SetNamespace(ns)

		if kind == "*" || kind == objKind {
			installFunc, isExists := InstalllerMap[objKind]
			if !isExists {
				return errors.New(fmt.Sprintf("no installer for kind %s ", objKind))
			}
			err := installFunc(ctx, info)
			return err

		}
		return nil
		// errors.New(fmt.Sprintf("object was not selecteds as does not match expected kind %s", kind))
	}

}
func RunInstall(ctx context.Context, infos []*resource.Info, installer InstallerFunc) error {
	var retError error
	for _, info := range infos {
		err := installer(ctx, info)
		if err != nil {
			retError = err
		}

	}

	return retError
}
func HandleError(err error, msg string) {

	if err != nil {
		fmt.Printf("error occured for %s\n", msg)
	}
}
func (i *Installer) Install(namespace string, resources *resource.Result) error {
	infos, _ := resources.Infos()
	ctx := context.WithValue(context.Background(), "Namespace", namespace)

	err := RunInstall(ctx, infos, i.InstallVisitor("Namespace"))
	HandleError(err, "can't install Namespace")
	err = RunInstall(ctx, infos, i.InstallVisitor("ClusterRole"))
	HandleError(err, "can't  install cluster role ")
	err = RunInstall(ctx, infos, i.InstallVisitor("ClusterRoleBinding"))
	HandleError(err, "can't install ClusterRoleBinding")
	err = RunInstall(ctx, infos, i.InstallVisitor("ServiceAccount"))
	HandleError(err, "can't install ServiceAccount")
	err = RunInstall(ctx, infos, i.InstallVisitor("Role"))
	HandleError(err, "can't install role")
	err = RunInstall(ctx, infos, i.InstallVisitor("RoleBinding"))
	HandleError(err, "can't install RoleBinding")
	err = RunInstall(ctx, infos, i.InstallVisitor("Job"))
	HandleError(err, "can't install RoleBinding")

	return nil
}
func (i Installer) InstallNamespace(ctx context.Context, info *resource.Info) error {

	fmt.Println("installing namespace")
	targetNS := ctx.Value("Namespace").(string)

	options := metav1.CreateOptions{}
	nsObj := info.Object.(*v1.Namespace)
	if targetNS != "" {
		nsObj.Name = targetNS
	}
	nsObj.Finalizers = []string{}
	_, err := i.Client.CoreV1().Namespaces().Create(context.Background(), nsObj, options)

	return err
}
func (i Installer) InstallJob(ctx context.Context, info *resource.Info) error {
	fmt.Println("installing job")

	targetNS := ctx.Value("Namespace").(string)

	options := metav1.CreateOptions{}
	job := info.Object.(*batchv1.Job)
	ns := job.Namespace
	if targetNS != "" {
		ns = targetNS
	}
	createdJob, err := i.Client.BatchV1().Jobs(ns).Create(context.Background(), job, options)
	if err != nil {
		return err
	}
	executer, err := utils.NewJobExecuter()

	if err != nil {
		return err
	}

	executer.SetJob(createdJob)
	stopChannel := make(<-chan struct{})
	var condFunc wait.ConditionFunc

	pods, err := executer.GetJobPods()
	if err != nil {
		return err
	}
	condFunc = func() (bool, error) {
		done, status, err := executer.IsJobCompleted()
		fmt.Printf("job status is %v", status)
		if len(pods.Items) > 0 {
			executer.ReadLogs(job.Namespace, pods.Items[0].Name)
		}
		return done, err
	}
	err = wait.PollImmediateUntil(5*time.Second, condFunc, stopChannel)

	return err

}
func (i Installer) InstallSa(ctx context.Context, info *resource.Info) error {

	fmt.Println("installing sa")

	options := metav1.CreateOptions{}
	saObj := info.Object.(*v1.ServiceAccount)
	targetNS := ctx.Value("Namespace").(string)

	ns := saObj.Namespace
	if targetNS != "" {
		ns = targetNS
	}
	_, err := i.Client.CoreV1().ServiceAccounts(ns).Create(context.Background(), saObj, options)

	return err
}
func (i Installer) InstallRole(ctx context.Context, info *resource.Info) error {
	fmt.Println("installing role")

	return nil
}
func (i Installer) InstallRoleBinding(ctx context.Context, info *resource.Info) error {
	fmt.Println("installing roleBinding")

	return nil
}
func (i Installer) InstallClusterRole(ctx context.Context, info *resource.Info) error {
	fmt.Println("installing cluster role")

	options := metav1.CreateOptions{}
	clusterRole := info.Object.(*rbacv1.ClusterRole)

	_, err := i.Client.RbacV1().ClusterRoles().Create(context.Background(), clusterRole, options)

	return err

}
func (i Installer) InstallClusterRoleBinding(ctx context.Context, info *resource.Info) error {
	fmt.Println("installing cluster role")

	options := metav1.CreateOptions{}
	clusterRole := info.Object.(*rbacv1.ClusterRoleBinding)
	_, err := i.Client.RbacV1().ClusterRoleBindings().Create(context.Background(), clusterRole, options)

	return err

}
func CheckInstall(job string) bool {
	return true
}

func Logs(job string) {

}
func TestInstall() {

}
