package model

/*
{"request":{"url":"/saas/configuration","method":"GET",
"timestamp":"2021-10-20T09:34:01.083+00:00"},
"response":{"items":
[{"microservices":
[{"name":"Inventory-Service","workloadType":"deployment","labels":
[{"key":"app","value":"inventoryService"}],
"rollouts":[{"rolloutSpec":"Inventory-Dev-Rolling","environment":"Dev","namespace":null,"strategy":"rolling"}]}
,{"name":"spotcd-demo","workloadType":"deployment","labels":[{"key":"app","value":"spotcd-demo"}],
"rollouts":[{"rolloutSpec":"spotcd-demo-dev-rolling","environment":"Dev","namespace":null,"strategy":"rolling"}]}]}]
,"count":1}}
*/
type ServicesList struct {
	Services []Service `json:"microservices"`
}
type Entities struct {
	Microservices []Service `json:"microservices"`
}

type Service struct {
	Microservice ServiceMetadata `json:"microservice"`
	//ServiceRolloutSpecsList
}

/*
 "k8sResources": {
    "workload": {
      "type": "deployment",
      "labels": [
        {
          "key": "app",
          "value": "AmirFirstMicroservice"
        }
      ],
      "versionLabelKey": "ms-version"
    }
  },
*/
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

type RolloutSpec struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Namespace   string `json:"namespace"`
	Strategy    struct {
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
		Providers []string
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
	Services     []Service
}

/*

{
	"environment": {
	  "name": "prod",
	  "clusterId": "cluster-prod",
	  "namespace": "prod-ns"
	}
  }
*/
type EnvironmentSpec struct {
	Name      string `json:"name"`
	ClusterId string `json:"clusterId"`
	Namespace string `json:"namespace"`
}
type Environment struct {
	Envrionment EnvironmentSpec `json:"environment"`
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
