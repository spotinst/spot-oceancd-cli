package verification

const (
	Running    Status = "running"
	Successful Status = "successful"
	Failed     Status = "failed"
	Error      Status = "error"
	Canceled   Status = "cancel"
)

var StatusOrder = map[Status]int{
	Failed:     1,
	Error:      2,
	Running:    3,
	Successful: 4,
}

type Status string

type Verification struct {
	MetricName       string      `json:"metricName"`
	StartTime        string      `json:"startTime"`
	Status           Status      `json:"status"`
	FailureCondition string      `json:"failureCondition"`
	Query            string      `json:"query"`
	FailureLimit     int         `json:"failureLimit"`
	Interval         string      `json:"interval"`
	Count            int         `json:"count"`
	DataPoints       []DataPoint `json:"dataPoints"`
	Provider         string      `json:"provider"`
	Step             string      `json:"step"`
}

type DataPoint struct {
	Timestamp string `json:"timestamp"`
	Value     string `json:"value"`
	Status    string `json:"status"`
}
