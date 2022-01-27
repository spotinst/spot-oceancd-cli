package model

import (
	"fmt"
	"time"
)

const (
	ClusterEntity              = "cluster"
	EnvEntity                  = "environment"
	ServiceEntity              = "microservice"
	NotificationProviderEntity = "notificationProvider"
	RolloutSpecEntity          = "rolloutSpec"
)

//region Microservice
type MicroserviceMeta interface {
	GetMicroserviceDetails() MicroserviceDetails
}
type Microservice struct {
	MicroserviceMeta `json:"-"`
	ServiceMetadata
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
	Name         string              `json:"name"`
	K8sResources ServiceK8sResources `json:"k8sResources"`
	CreatedAt    time.Time           `json:"createdAt"`
}
type ServiceLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type MicroserviceDetails struct {
	Name      string `header:"Name"`
	Labels    string `header:"Labels"`
	CreatedAt string `header:"Created At"`
}

func (s *Microservice) GetMicroserviceDetails() MicroserviceDetails {
	createdAt := s.CreatedAt.Format("2006-01-02 15:04:05")
	msLabel := ""

	if len(s.K8sResources.Labels) > 0 {
		firstLabel := s.K8sResources.Labels[0]
		msLabel = fmt.Sprintf("%v=%v", firstLabel.Key, firstLabel.Value)

		if len(s.K8sResources.Labels) > 1 {
			msLabel = fmt.Sprintf("%v...", msLabel)
		}
	}

	retVal := MicroserviceDetails{
		Name:      s.Name,
		Labels:    msLabel,
		CreatedAt: createdAt,
	}

	return retVal
}

//endregion

//region rollout spec
type RolloutSpecMeta interface {
	GetRolloutSpecDetails() RolloutSpecDetails
}
type RolloutSpec struct {
	RolloutSpecMeta `json:"-"`
	Name            string `json:"name"`
	Microservice    string `json:"microservice"`
	Environment     string `json:"environment"`
	Strategy        struct {
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
	CreatedAt time.Time `json:"createdAt"`
}
type RolloutSpecDetails struct {
	Name             string `header:"Name"`
	Microservice     string `header:"Microservice"`
	Environment      string `header:"Environment"`
	HasVerification  bool   `header:"Has Verification"`
	HasFailurePolicy bool   `header:"Has Failure Policy"`
	CreatedAt        string `header:"Created At"`
}

func (r *RolloutSpec) GetRolloutSpecDetails() RolloutSpecDetails {
	createdAt := r.CreatedAt.Format("2006-01-02 15:04:05")
	retVal := RolloutSpecDetails{
		Name:             r.Name,
		Microservice:     r.Microservice,
		Environment:      r.Environment,
		HasVerification:  len(r.Strategy.Rolling.Verification.Phases) > 0,
		HasFailurePolicy: r.FailurePolicy.Rollback.Mode != "",
		CreatedAt:        createdAt,
	}

	return retVal
}

//endregion

//region environment
type EnvironmentMeta interface {
	GetEnvironmentDetails() EnvironmentDetails
}
type EnvironmentSpec struct {
	EnvironmentMeta `json:"-"`
	Name            string    `json:"name"`
	ClusterId       string    `json:"clusterId"`
	Namespace       string    `json:"namespace"`
	CreatedAt       time.Time `json:"createdAt"`
}
type EnvironmentDetails struct {
	Name      string `header:"Name"`
	ClusterId string `header:"Cluster Id"`
	Namespace string `header:"Namespace"`
	CreatedAt string `header:"Created At"`
}

func (e *EnvironmentSpec) GetEnvironmentDetails() EnvironmentDetails {
	createdAt := e.CreatedAt.Format("2006-01-02 15:04:05")
	retVal := EnvironmentDetails{
		Name:      e.Name,
		ClusterId: e.ClusterId,
		Namespace: e.Namespace,
		CreatedAt: createdAt,
	}
	return retVal
}

//endregion

//region cluster
type ClusterMeta interface {
	GetClusterDetails() ClusterDetails
}
type ClusterSpec struct {
	ClusterMeta       `json:"-"`
	Name              string    `json:"id"`
	LastHeartbeatTime time.Time `json:"lastHeartbeatTime"`
	ClusterInfo       struct {
		KubeVersion   string `json:"kubernetesVersion"`
		CloudProvider string `json:"cloudProvider"`
		KubeEngine    string `json:"kubernetesOrchestrator"`
	} `json:"clusterInfo"`
	ControllerInfo struct {
		NodeName          string `json:"nodeName"`
		ControllerVersion string `json:"controllerVersion"`
		PodName           string `json:"podName"`
	} `json:"controllerInfo"`
	CreatedAt time.Time `json:"createdAt"`
}
type ClusterDetails struct {
	Name              string `header:"Name"`
	K8sVersion        string `header:"Kubernetes Version"`
	ControllerVersion string `header:"Controller Version"`
	LastHeartbeat     string `header:"Last Heartbeat"`
	CreatedAt         string `header:"Created At"`
}

func (c *ClusterSpec) GetClusterDetails() ClusterDetails {
	lastHeartbeat := c.LastHeartbeatTime.Format("2006-01-02 15:04:05")
	createdAt := c.CreatedAt.Format("2006-01-02 15:04:05")
	retVal := ClusterDetails{
		Name:              c.Name,
		K8sVersion:        c.ClusterInfo.KubeVersion,
		ControllerVersion: c.ControllerInfo.ControllerVersion,
		LastHeartbeat:     lastHeartbeat,
		CreatedAt:         createdAt,
	}
	return retVal
}

//endregion

//region notification provider
type NotificationProviderMeta interface {
	GetNotificationProviderDetails() NotificationProviderDetails
}
type NotificationProviderSpec struct {
	NotificationProviderMeta `json:"-"`
	Name                     string `json:"name"`
	Description              string `json:"description"`
	Webhook                  struct {
		Url string `json:"url"`
	} `json:"webhook"`
	CreatedAt time.Time `json:"createdAt"`
}
type NotificationProviderDetails struct {
	Name        string `header:"Name"`
	Url         string `header:"Url"`
	CreatedAt   string `header:"Created At"`
}

func (c *NotificationProviderSpec) GetNotificationProviderDetails() NotificationProviderDetails {
	createdAt := c.CreatedAt.Format("2006-01-02 15:04:05")
	retVal := NotificationProviderDetails{
		Name:        c.Name,
		Url:         c.Webhook.Url,
		CreatedAt:   createdAt,
	}
	return retVal
}
//endregion
