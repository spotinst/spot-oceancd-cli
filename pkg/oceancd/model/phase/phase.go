package phase

import "spot-oceancd-cli/pkg/oceancd/model/verification"

const (
	Finished     Status = "finished"
	Promoted     Status = "promoted"
	FullPromoted Status = "fullPromoted"
	Promoting    Status = "promoting"
	Paused       Status = "paused"
	InProgress   Status = "inProgress"
	Aborted      Status = "aborted"
	Canceled     Status = "canceled"
	Pending      Status = "pending"
	Dropped      Status = "dropped"
	Verifying    Status = "verifying"
	Aborting     Status = "aborting"
)

var UncompletedStatuses = []Status{Promoting, Paused, InProgress, Pending, Verifying, Aborting}

type Status string

type Phase struct {
	Name              string                      `json:"name"`
	Status            Status                      `json:"status"`
	StartTime         string                      `json:"startTime"`
	PausedAt          string                      `json:"pausedAt"`
	VerifiedAt        string                      `json:"verifiedAt"`
	EndTime           string                      `json:"endTime"`
	TrafficPercentage int                         `json:"trafficPercentage"`
	Verifications     []verification.Verification `json:"verifications"`
}

func (p *Phase) IsUncompleted() bool {
	for _, uncompletedStatus := range UncompletedStatuses {
		if p.Status == uncompletedStatus {
			return true
		}
	}

	return false
}
