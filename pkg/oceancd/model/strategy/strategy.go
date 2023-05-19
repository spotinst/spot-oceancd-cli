package strategy

type RollingUpdateStrategy struct {
	Steps []Step `json:"steps"`
}

type CanaryStrategy struct {
	Steps []Step `json:"steps"`
}

func (c *CanaryStrategy) GetHeaderRouteMatchesBySteps() map[string][]Match {
	matches := map[string][]Match{}

	for _, step := range c.Steps {
		matches[step.Name] = append(matches[step.Name], step.SetHeaderRoute.Match...)
	}

	return matches
}

func (r *RollingUpdateStrategy) GetHeaderRouteMatchesBySteps() map[string][]Match {
	matches := map[string][]Match{}

	for _, step := range r.Steps {
		matches[step.Name] = append(matches[step.Name], step.SetHeaderRoute.Match...)
	}

	return matches
}

type Step struct {
	Name           string         `json:"name"`
	SetHeaderRoute SetHeaderRoute `json:"setHeaderRoute"`
}

type SetHeaderRoute struct {
	Match []Match `json:"match"`
}

type Match struct {
	HeaderName  string      `json:"headerName"`
	HeaderValue HeaderValue `json:"headerValue"`
}

type HeaderValue struct {
	Exact  string `json:"exact"`
	Prefix string `json:"prefix"`
	Regex  string `json:"regex"`
}
