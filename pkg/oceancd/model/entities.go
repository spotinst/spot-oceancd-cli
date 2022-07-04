package model

const (
	ClusterEntity     = "cluster"
	StrategyEntity    = "strategy"
	RolloutSpecEntity = "rolloutSpec"
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
