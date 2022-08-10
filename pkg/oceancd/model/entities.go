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
	strategyName := ""
	strategyRef := rolloutSpec["strategy"].(map[string]interface{})

	if strategyRef != nil {
		strategyName = strategyRef["name"].(string)
	}

	stableService := ""
	trafficDef := rolloutSpec["traffic"].(map[string]interface{})

	if trafficDef != nil {
		if val, ok := trafficDef["stableService"]; ok {
			stableService = val.(string)
		}
	}

	retVal := RolloutSpecDetails{
		Name:             rolloutSpec["name"].(string),
		Strategy:         strategyName,
		StableService:    stableService,
		HasFailurePolicy: rolloutSpec["failurePolicy"] != nil,
		UpdatedAt:        rolloutSpec["updatedAt"].(string),
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

	if strategy["canary"] != nil {
		strategyType = "Canary"
		hasBackgroundVerifications = strategy["canary"].(map[string]interface{})["backgroundVerification"] != nil
		stepsCount = len(strategy["canary"].(map[string]interface{})["steps"].([]interface{}))
	}

	retVal := StrategyDetails{
		Name:                       strategy["name"].(string),
		Type:                       strategyType,
		HasBackgroundVerifications: hasBackgroundVerifications,
		StepsCount:                 stepsCount,
		UpdatedAt:                  strategy["updatedAt"].(string),
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
	lastHeartbeat := ""

	if cluster["lastHeartbeatTime"] != nil {
		lastHeartbeat = cluster["lastHeartbeatTime"].(string)
	}

	updatedAt := cluster["updatedAt"].(string)

	k8sVersion := ""

	if cluster["clusterInfo"] != nil {
		k8sVersion = cluster["clusterInfo"].(map[string]interface{})["kubernetesVersion"].(string)
	}

	controllerVersion := ""

	if cluster["controllerInfo"] != nil && cluster["controllerInfo"].(map[string]interface{})["controllerVersion"] != nil {
		controllerVersion = cluster["controllerInfo"].(map[string]interface{})["controllerVersion"].(string)
	}

	retVal := ClusterDetails{
		Name:              cluster["id"].(string),
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
		if stringID, ok := id.(string); ok {
			clusterIDs = append(clusterIDs, stringID)
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

	args, _ := verificationTemplate["args"].([]interface{})
	for _, arg := range args {
		arg, _ := arg.(map[string]interface{})
		if name, ok := arg["name"].(string); ok {
			argsRes = append(argsRes, name)
		}
	}

	metrics, _ := verificationTemplate["metrics"].([]interface{})
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

	retVal := VerificationTemplateDetails{
		Name:      verificationTemplateName,
		Args:      strings.Join(argsRes, ", "),
		Metrics:   strings.Join(metricsRes, ", "),
		UpdatedAt: updatedAt,
	}

	return retVal
}
