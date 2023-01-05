package rollout

import (
	oceancd "spot-oceancd-cli/pkg/oceancd/model"
	"spot-oceancd-cli/pkg/oceancd/model/phase"
	"spot-oceancd-cli/pkg/oceancd/model/verification"
)

const (
	Pending              Status = "pending"
	InProgress           Status = "inProgress"
	Paused               Status = "paused"
	Failed               Status = "failed"
	Aborted              Status = "aborted"
	Aborting             Status = "aborting"
	Finished             Status = "finished"
	Canceled             Status = "canceled"
	ManualPaused         Status = "manualPaused"
	ManualPausing        Status = "manualPausing"
	InvalidSpec          Status = "invalidSpec"
	Deallocating         Status = "deallocating"
	Verifying            Status = "verifying"
	FailurePolicyPaused  Status = "failurePolicyPaused"
	FailurePolicyPausing Status = "failurePolicyPausing"
	BackgroundVerifying  Status = "backgroundVerifying"
)

var CompletedStatuses = []Status{Failed, Aborted, Finished, Canceled}

type Status string

func (s Status) IsCompleted() bool {
	for _, completedStatus := range CompletedStatuses {
		if s == completedStatus {
			return true
		}
	}

	return false
}

type ReplicasInfo struct {
	Desired    int `json:"desired"`
	Ready      int `json:"ready"`
	InProgress int `json:"inProgress"`
	Failed     int `json:"failed"`
}

type VersionStatus struct {
	Version           string       `json:"version"`
	K8sService        string       `json:"k8sService"`
	TrafficPercentage int          `json:"trafficPercentage"`
	Replicas          ReplicasInfo `json:"replicas"`
}

type Rollout struct {
	Id                        string        `json:"id"`
	Status                    Status        `json:"status"`
	SpotDeployment            string        `json:"spotDeployment"`
	OriginalRolloutId         string        `json:"originalRolloutId"`
	NewRolloutId              string        `json:"newRolloutId"`
	StartTime                 string        `json:"startTime"`
	EndTime                   string        `json:"endTime"`
	ClusterId                 string        `json:"clusterId"`
	Namespace                 string        `json:"namespace"`
	Strategy                  string        `json:"strategy"`
	HasBackgroundVerification bool          `json:"hasBackgroundVerification"`
	NewVersionStatus          VersionStatus `json:"newVersionStatus"`
	StableVersionStatus       VersionStatus `json:"stableVersionStatus"`
}

type DetailedRollout struct {
	Rollout
	Phases        []phase.Phase               `json:"phases"`
	Verifications []verification.Verification `json:"verifications"`
}

func (d *DetailedRollout) GetBackgroundVerifications() []verification.Verification {
	backgroundVerifications := make([]verification.Verification, 0)

	for _, verificationItem := range d.Verifications {
		if verificationItem.Step == oceancd.BackgroundVerificationLabel {
			backgroundVerifications = append(backgroundVerifications, verificationItem)
		}
	}
	return backgroundVerifications
}
