package model

import (
	"fmt"

	"github.com/fatih/color"
)

const ClusterEntity = "cluster"
const ServiceEntity = "microservice"
const EnvEntity = "environment"
const RolloutSpecEntity = "rolloutSpec"

type EntityPrinter interface {
	Format(string, interface{}) []string
}
type EntitySpec interface {
	GetEntitySpec() interface{}
}
type EntityMeta interface {
	GetEntityKind() string
	GetEntityName() string
}

type ServicesList struct {
	Services []ServiceRequest `json:"microservices"`
}
type Entities struct {
	Microservices []ServiceRequest `json:"microservices"`
}

type ServiceRequest struct {
	EntitySpec   `json:"-"`
	Microservice *Service `json:"microservice"`
	//ServiceRolloutSpecsList
}
type Service struct {
	EntityMeta //`json:"-"`
	ServiceMetadata
	EntityPrinter `json:"-"`
	//ServiceRolloutSpecsList
}

type ServiceWorkload struct {
	Type            string         `json:"type"`
	Labels          []ServiceLabel `json:"labels"`
	VersionLabelKey string         `json:"versionLabelKey"`
}
type ServiceK8sResources struct {
	ServiceWorkload `json:"workload"`
}
type ServiceMetadata struct {
	Name string `json:"name"`
	//	Workload     string              `json:"workloadType"`
	//	Labels       []ServiceLabel      `json:"labels"`
	K8sResources ServiceK8sResources `json:"k8sResources"`
}

type ServiceLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ServiceRolloutSpecsList struct {
	Rollouts []ServiceRolloutSpec `json:"rollouts"`
}
type ServiceRolloutSpec struct {
	RolloutSpec string `json:"rolloutSpec"`
	Environment string `json:"environment"`
	Namespace   string `json:"namespace"`
	Strategy    string `json:"strategy"`
}

type RolloutSpecRequest struct {
	EntitySpec `json:"-"`
	Spec       *RolloutSpec `json:"rolloutSpec"`
}

type RolloutSpec struct {
	EntityMeta    `json:"-"`
	EntityPrinter `json:"-"`
	Name          string `json:"name"`
	Microservice  string `json:"microservice"`
	Environment   string `json:"environment"`
	Strategy      struct {
		Rolling struct {
			Verification struct {
				Phases []struct {
					IinitialDelay        string `json:"initialDelay"`
					Name                 string `json:"name"`
					ExternalVerification struct {
						Fallback string `json:"fallback"`
						Timeout  string `json:"timeout"`
					} `json:"externalVerification"`
				} `json:"phases"`
			} `json:"verification"`
		} `json:"rolling"`
	} `json:"strategy"`

	Notification struct {
		Providers []string `json:"providers"`
	} `json:"notification"`
	FailurePolicy struct {
		Rollback struct {
			Mode string `json:"mode"`
		} `json:"rollback"`
	} `json:"failurePolicy"`
}

type EntityList struct {
	Environments []EnvironmentSpec
	Specs        []RolloutSpec
	Services     []ServiceRequest
}

type EnvironmentSpec struct {
	EntityMeta    `json:"-"`
	EntityPrinter `json:"-"`
	Name          string `json:"name"`
	ClusterId     string `json:"clusterId"`
	Namespace     string `json:"namespace"`
}
type EnvironmentRequest struct {
	EntitySpec  `json:"-"`
	Envrionment *EnvironmentSpec `json:"environment"`
}
type StartRolloutDetails struct {
	Microservice  string
	InitiatorUser string
	DetectionId   string
	Operation     string
	Namespace     string
}

type K8sServiceIdentifier struct {
	ServiceName   string            `json:"serviceName"`
	ServiceLabels map[string]string `json:"serviceLabels"`
}

type EntitiesResponse struct {
	Items []Entities `json:"items"`
}

type ClusterSpec struct {
	EntityMeta    `json:"-"`
	EntityPrinter `json:"-"`
	Name          string `json:"id"`
	ClusterInfo   struct {
		KubeVersion   string `json:"kubernetesVersion"`
		CloudProvider string `json:"cloudProvider"`
		KubeEngine    string `json:"kubernetesOrchestrator"`
	} `json:"clusterInfo"`
	ControllerInfo struct {
		NodeName          string `json:"nodeName"`
		ControllerVersion string `json:"controllerVersion"`
		PodName           string `json:"podName"`
	} `json:"controllerInfo"`
}

/*
{
	"lastHeartbeatTime": "2018-11-05T12:55:50.000+0000",
	"controllerInfo": {
	"nodeName": "string",
	"controllerVersion": "string",
	"podName": "string"
	},
	"clusterInfo": {
	"kubernetesVersion": "string",
	"cloudProvider": "string",
	"kubernetesOrchestrator": "string"
	},
	"notification": {
	"minutesWithoutHeartbeat": 4,
	"providers": [
	"provider1"
	]
	},
	"id": "oceanCluster",
	"createdAt": "2018-11-05T12:55:50.000+0000",
	"updatedAt": "2018-11-05T12:58:15.000+0000"
	}

*/

func (s *ServiceRequest) GetEntitySpec() interface{} {
	return s.Microservice
}
func (s *Service) GetEntityKind() string {
	return ServiceEntity
}
func (s *Service) GetEntityName() string {
	return s.Name
}
func (e *EnvironmentRequest) GetEntitySpec() interface{} {
	return e.Envrionment
}
func (e *EnvironmentSpec) GetEntityKind() string {
	return EnvEntity
}
func (e *EnvironmentSpec) GetEntityName() string {
	return e.Name
}

func (r *RolloutSpecRequest) GetEntitySpec() interface{} {
	return r.Spec
}

func (r *RolloutSpec) GetEntityKind() string {
	return RolloutSpecEntity
}
func (r RolloutSpec) GetEntityName() string {
	return r.Name
}
func (s *Service) Format(formatType string, more interface{}) []string {
	labels := ""
	for _, l := range s.K8sResources.Labels {
		label := fmt.Sprintf("%v=%v,", l.Key, l.Value)
		labels = labels + label
	}
	row := []string{s.Name, labels, s.K8sResources.Type}
	return row
}

func (e *EnvironmentSpec) Format(formatType string, more interface{}) []string {
	row := []string{e.Name, e.ClusterId, e.Namespace}
	return row
}

func (c *ClusterSpec) GetEntityKind() string {
	return ClusterEntity
}
func (c *ClusterSpec) GetEntityName() string {
	return c.Name
}
func (r *RolloutSpec) Format(formatType string, more interface{}) []string {
	services, ok := more.([]*Service)
	if !ok {
		services = []*Service{}
	}

	selector := ""
	for _, s := range services {
		if s.Name == r.Microservice {

			for _, l := range s.K8sResources.Labels {
				label := fmt.Sprintf("%v=%v,", l.Key, l.Value)
				selector = selector + label
			}
		}
	}
	row := []string{}
	if selector != "" {

		row = []string{r.Name, r.Environment, r.Microservice, color.New(color.FgGreen).Sprint(selector)}
	} else {
		row = []string{r.Name, r.Environment, color.New(color.FgGreen).Sprint(r.Microservice)}

	}
	return row
}

func (c *ClusterSpec) Format(formatType string, more interface{}) []string {
	row := []string{c.Name, c.ClusterInfo.KubeVersion,
		c.ControllerInfo.ControllerVersion,
		c.ControllerInfo.NodeName, c.ControllerInfo.PodName}
	return row
}
