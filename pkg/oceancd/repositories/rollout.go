package repositories

import (
	"encoding/json"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model/rollout"
	strategymodel "spot-oceancd-cli/pkg/oceancd/model/strategy"
)

func NewRolloutRepository() *RolloutRepository {
	return &RolloutRepository{}
}

type RolloutRepository struct {
}

func (r *RolloutRepository) GetStrategy(rolloutId string) (rollout.Strategy, error) {
	var retVal rollout.Strategy
	strategyDefinition := map[string]interface{}{}
	rolloutDefinition, err := oceancd.GetRolloutDefinition(rolloutId)

	if strategyInfo, ok := rolloutDefinition["strategy"]; ok {
		if strategy, ok := strategyInfo.(map[string]interface{}); ok {
			if canaryInfo, ok := strategy["canary"]; ok {
				if canary, ok := canaryInfo.(map[string]interface{}); ok {
					retVal = &strategymodel.CanaryStrategy{}
					strategyDefinition = canary
				}
			}
			if rollingInfo, ok := strategy["rolling"]; ok {
				if rolling, ok := rollingInfo.(map[string]interface{}); ok {
					retVal = &strategymodel.RollingUpdateStrategy{}
					strategyDefinition = rolling
				}
			}
		}
	}

	bytes, err := json.Marshal(strategyDefinition)
	if err != nil {
		return retVal, err
	}

	err = json.Unmarshal(bytes, &retVal)
	if err != nil {
		return retVal, err
	}

	return retVal, nil
}
