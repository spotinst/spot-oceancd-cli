package model

import (
	"fmt"
	"strings"
)

//region rollout spec
type RolloutSpecDetails struct {
	Name             string `header:"Name"`
	Strategy         string `header:"Strategy"`
	StableService    string `header:"Stable Service"`
	HasFailurePolicy bool   `header:"Has Failure Policy"`
	UpdatedAt        string `header:"Updated At"`
}

func ConvertToRolloutSpecDetails(rolloutSpec map[string]interface{}) RolloutSpecDetails {
	var hasFailurePolicy bool
	strategyName := ""
	stableService := ""

	if strategyRef, ok := rolloutSpec["strategy"].(map[string]interface{}); ok {
		strategyName, _ = strategyRef["name"].(string)
	}

	if trafficDef, ok := rolloutSpec["traffic"].(map[string]interface{}); ok {
		if val, ok := trafficDef["stableService"]; ok {
			stableService, _ = val.(string)
		}
	}

	name, _ := rolloutSpec["name"].(string)
	updatedAt, _ := rolloutSpec["updatedAt"].(string)
	if failurePolicy, ok := rolloutSpec["failurePolicy"].(map[string]interface{}); ok {
		_, hasFailurePolicy = failurePolicy["action"]
	}

	retVal := RolloutSpecDetails{
		Name:             name,
		Strategy:         strategyName,
		StableService:    stableService,
		HasFailurePolicy: hasFailurePolicy,
		UpdatedAt:        updatedAt,
	}

	return retVal
}

//endregion

//region strategy
type StrategyDetails struct {
	Name                       string `header:"Name"`
	Type                       string `header:"Type"`
	HasBackgroundVerifications bool   `header:"Has Background Verifications"`
	StepsCount                 int    `header:"Steps Count"`
	UpdatedAt                  string `header:"Updated At"`
}

func ConvertToStrategyDetails(strategy map[string]interface{}) StrategyDetails {
	strategyType := ""
	hasBackgroundVerifications := false
	stepsCount := 0

	if canary, ok := strategy["canary"].(map[string]interface{}); ok {
		strategyType = "Canary"
		if backgroundVerification, ok := canary["backgroundVerification"].(map[string]interface{}); ok {
			if templateNames, ok := backgroundVerification["templateNames"].([]interface{}); ok {
				hasBackgroundVerifications = len(templateNames) > 0
			}
		}

		if steps, ok := canary["steps"].([]interface{}); ok {
			stepsCount = len(steps)
		}
	}
	strategyName, _ := strategy["name"].(string)
	updatedAt, _ := strategy["updatedAt"].(string)

	retVal := StrategyDetails{
		Name:                       strategyName,
		Type:                       strategyType,
		HasBackgroundVerifications: hasBackgroundVerifications,
		StepsCount:                 stepsCount,
		UpdatedAt:                  updatedAt,
	}

	return retVal
}

//endregion

//region cluster
type ClusterDetails struct {
	Name              string `header:"Name"`
	K8sVersion        string `header:"Kubernetes Version"`
	ControllerVersion string `header:"Controller Version"`
	LastHeartbeat     string `header:"Last Heartbeat"`
	UpdatedAt         string `header:"Updated At"`
}

func ConvertToClusterDetails(cluster map[string]interface{}) ClusterDetails {
	k8sVersion := ""
	controllerVersion := ""

	lastHeartbeat, _ := cluster["lastHeartbeatTime"].(string)
	updatedAt, _ := cluster["updatedAt"].(string)

	if clusterInfo, ok := cluster["clusterInfo"].(map[string]interface{}); ok {
		k8sVersion, _ = clusterInfo["kubernetesVersion"].(string)
	}

	if controllerInfo, ok := cluster["controllerInfo"].(map[string]interface{}); ok {
		controllerVersion, _ = controllerInfo["controllerVersion"].(string)
	}

	clusterId, _ := cluster["id"].(string)

	retVal := ClusterDetails{
		Name:              clusterId,
		K8sVersion:        k8sVersion,
		ControllerVersion: controllerVersion,
		LastHeartbeat:     lastHeartbeat,
		UpdatedAt:         updatedAt,
	}
	return retVal
}

//endregion

type VerificationProviderDetails struct {
	Name       string `header:"Name"`
	ClusterIDs string `header:"Cluster ID"`
	Types      string `header:"Type"`
	UpdatedAt  string `header:"Updated At"`
}

func ConvertToVerificationProviderDetails(verificationProvider map[string]interface{}) VerificationProviderDetails {
	clusterIDs := make([]string, 0)
	types := make([]string, 0)

	verificationProviderName, _ := verificationProvider["name"].(string)
	updatedAt, _ := verificationProvider["updatedAt"].(string)

	ids, _ := verificationProvider["clusterIds"].([]interface{})
	for _, id := range ids {
		if clusterID, ok := id.(string); ok {
			clusterIDs = append(clusterIDs, clusterID)
		}
	}

	for verificationProviderType := range verificationProvider {
		switch verificationProviderType {
		case Prometheus, Datadog, NewRelic:
			types = append(types, verificationProviderType)
		}
	}

	retVal := VerificationProviderDetails{
		Name:       verificationProviderName,
		ClusterIDs: strings.Join(clusterIDs, ", "),
		Types:      strings.Join(types, ", "),
		UpdatedAt:  updatedAt,
	}

	return retVal
}

type VerificationTemplateDetails struct {
	Name      string `header:"Name"`
	Args      string `header:"Arg (Name)"`
	Metrics   string `header:"Metric (Provider)"`
	UpdatedAt string `header:"Updated At"`
}

func ConvertToVerificationTemplateDetails(verificationTemplate map[string]interface{}) VerificationTemplateDetails {
	argsRes := make([]string, 0)
	metricsRes := make([]string, 0)

	verificationTemplateName, _ := verificationTemplate["name"].(string)
	updatedAt, _ := verificationTemplate["updatedAt"].(string)

	if args, ok := verificationTemplate["args"].([]interface{}); ok {
		for _, arg := range args {
			arg, _ := arg.(map[string]interface{})
			if name, ok := arg["name"].(string); ok {
				argsRes = append(argsRes, name)
			}
		}
	}

	if metrics, ok := verificationTemplate["metrics"].([]interface{}); ok {
		for _, metric := range metrics {
			metric, _ := metric.(map[string]interface{})
			if name, ok := metric["name"].(string); ok {
				providerNames := make([]string, 0)
				row := name

				providers, _ := metric["provider"].(map[string]interface{})
				for providerName := range providers {
					switch providerName {
					case Prometheus, Datadog, NewRelic, CloudWatch, Web:
						providerNames = append(providerNames, providerName)
					}
				}

				if len(providerNames) > 0 {
					row += fmt.Sprintf("(%s)", strings.Join(providerNames, ", "))
				}

				metricsRes = append(metricsRes, row)
			}
		}
	}

	retVal := VerificationTemplateDetails{
		Name:      verificationTemplateName,
		Args:      strings.Join(argsRes, ", "),
		Metrics:   strings.Join(metricsRes, ", "),
		UpdatedAt: updatedAt,
	}

	return retVal
}
