package utils

import "spot-oceancd-cli/pkg/oceancd/model"

func GetEnvironmentEntityDetails(entity interface{}) model.EnvironmentDetails {
	meta, ok := entity.(model.EnvironmentMeta)
	if ok == false {
		return model.EnvironmentDetails{}
	}

	return meta.GetEnvironmentDetails()
}

func GetEnvironmentEntitiesDetails(entities []interface{}) []model.EnvironmentDetails {
	retVal := make([]model.EnvironmentDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = GetEnvironmentEntityDetails(entity)
	}

	return retVal
}

func GetMicroserviceEntityDetails(entity interface{}) model.MicroserviceDetails {
	meta, ok := entity.(model.MicroserviceMeta)
	if ok == false {
		return model.MicroserviceDetails{}
	}

	return meta.GetMicroserviceDetails()
}

func GetMicroserviceEntitiesDetails(entities []interface{}) []model.MicroserviceDetails {
	retVal := make([]model.MicroserviceDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = GetMicroserviceEntityDetails(entity)
	}

	return retVal
}

func GetRolloutSpecEntityDetails(entity interface{}) model.RolloutSpecDetails {
	meta, ok := entity.(model.RolloutSpecMeta)
	if ok == false {
		return model.RolloutSpecDetails{}
	}

	return meta.GetRolloutSpecDetails()
}

func GetRolloutSpecEntitiesDetails(entities []interface{}) []model.RolloutSpecDetails {
	retVal := make([]model.RolloutSpecDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = GetRolloutSpecEntityDetails(entity)
	}

	return retVal
}

func GetClusterEntityDetails(entity interface{}) model.ClusterDetails {
	meta, ok := entity.(model.ClusterMeta)
	if ok == false {
		return model.ClusterDetails{}
	}

	return meta.GetClusterDetails()
}

func GetClusterEntitiesDetails(entities []interface{}) []model.ClusterDetails {
	retVal := make([]model.ClusterDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = GetClusterEntityDetails(entity)
	}

	return retVal
}

func GetNotificationProviderEntityDetails(entity interface{}) model.NotificationProviderDetails {
	meta, ok := entity.(model.NotificationProviderMeta)
	if ok == false {
		return model.NotificationProviderDetails{}
	}

	return meta.GetNotificationProviderDetails()
}

func GetNotificationProviderEntitiesDetails(entities []interface{}) []model.NotificationProviderDetails {
	retVal := make([]model.NotificationProviderDetails, len(entities))

	for i, entity := range entities {
		retVal[i] = GetNotificationProviderEntityDetails(entity)
	}

	return retVal
}
