package utils

import (
	"spot-oceancd-cli/pkg/oceancd/model"
)

func GetVerificationProviderEntitiesDetails(entities []interface{}) []model.VerificationProviderDetails {
	retVal := make([]model.VerificationProviderDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = model.ConvertToVerificationProviderDetails(entity.(map[string]interface{}))
	}

	return retVal
}

func GetVerificationTemplateEntitiesDetails(entities []interface{}) []model.VerificationTemplateDetails {
	retVal := make([]model.VerificationTemplateDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = model.ConvertToVerificationTemplateDetails(entity.(map[string]interface{}))
	}

	return retVal
}

func GetStrategyEntitiesDetails(entities []interface{}) []model.StrategyDetails {
	retVal := make([]model.StrategyDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = model.ConvertToStrategyDetails(entity.(map[string]interface{}))
	}

	return retVal
}

func GetRolloutSpecEntitiesDetails(entities []interface{}) []model.RolloutSpecDetails {
	retVal := make([]model.RolloutSpecDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = model.ConvertToRolloutSpecDetails(entity.(map[string]interface{}))
	}

	return retVal
}

func GetClusterEntitiesDetails(entities []interface{}) []model.ClusterDetails {
	retVal := make([]model.ClusterDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = model.ConvertToClusterDetails(entity.(map[string]interface{}))
	}

	return retVal
}
