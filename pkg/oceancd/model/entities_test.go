package model

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestConvertToRolloutSpecDetails(t *testing.T) {
	cases := map[string]struct {
		payload  map[string]interface{}
		expected RolloutSpecDetails
	}{
		"empty payload": {
			payload:  map[string]interface{}{},
			expected: RolloutSpecDetails{},
		},
		"empty values": {
			payload: map[string]interface{}{
				"createdAt":      "",
				"failurePolicy":  map[string]interface{}{},
				"name":           "",
				"spotDeployment": map[string]interface{}{},
				"strategy":       map[string]interface{}{},
				"traffic":        map[string]interface{}{},
				"updatedAt":      "",
			},
			expected: RolloutSpecDetails{},
		},
		"fulfilled payload": {
			payload: map[string]interface{}{
				"createdAt": "2022-01-01T00:00:00.000Z",
				"failurePolicy": map[string]interface{}{
					"action": "abort",
				},
				"name": "test",
				"spotDeployment": map[string]interface{}{
					"clusterId": "cluster",
					"name":      "test-name",
					"namespace": "test-namespace",
				},
				"strategy": map[string]interface{}{
					"name": "strategy-name",
				},
				"traffic": map[string]interface{}{
					"canaryService": "canary",
					"stableService": "stable",
				},
				"updatedAt": "2022-01-02T00:00:00.000Z",
			},
			expected: RolloutSpecDetails{
				Name:             "test",
				Strategy:         "strategy-name",
				StableService:    "stable",
				HasFailurePolicy: true,
				UpdatedAt:        "2022-01-02T00:00:00.000Z",
			},
		},
	}

	for name, tc := range cases {
		got := ConvertToRolloutSpecDetails(tc.payload)
		if diff := cmp.Diff(tc.expected, got); diff != "" {
			t.Fatalf(name+"\n%s", diff)
		}
	}
}

func TestConvertToStrategyDetails(t *testing.T) {
	cases := map[string]struct {
		payload  map[string]interface{}
		expected StrategyDetails
	}{
		"empty payload": {
			payload:  map[string]interface{}{},
			expected: StrategyDetails{},
		},
		"empty values": {
			payload: map[string]interface{}{
				"canary":    map[string]interface{}{},
				"createdAt": "",
				"name":      "",
				"updatedAt": "",
			},
			expected: StrategyDetails{
				Type: "Canary",
			},
		},
		"has no strategy": {
			payload: map[string]interface{}{
				"undefined": map[string]interface{}{},
				"createdAt": "2022-01-01T00:00:00.000Z",
				"name":      "test",
				"updatedAt": "2022-01-02T00:00:00.000Z",
			},
			expected: StrategyDetails{
				UpdatedAt: "2022-01-02T00:00:00.000Z",
				Name:      "test",
			},
		},
		"empty background verification": {
			payload: map[string]interface{}{
				"canary": map[string]interface{}{
					"backgroundVerification": map[string]interface{}{
						"templatesNames": []interface{}{},
					},
					"steps": []interface{}{},
				},
			},
			expected: StrategyDetails{
				Type: "Canary",
			},
		},
		"steps are present": {
			payload: map[string]interface{}{
				"canary": map[string]interface{}{
					"steps": []interface{}{
						map[string]interface{}{},
						map[string]interface{}{},
						map[string]interface{}{},
						map[string]interface{}{},
					},
				},
				"name":      "test",
				"updatedAt": "2022-01-02T00:00:00.000Z",
			},
			expected: StrategyDetails{
				Type:       "Canary",
				StepsCount: 4,
				Name:       "test",
				UpdatedAt:  "2022-01-02T00:00:00.000Z",
			},
		},
		"steps are present && rolling update": {
			payload: map[string]interface{}{
				"rolling": map[string]interface{}{
					"steps": []interface{}{
						map[string]interface{}{},
						map[string]interface{}{},
						map[string]interface{}{},
						map[string]interface{}{},
					},
				},
				"name":      "test",
				"updatedAt": "2022-01-02T00:00:00.000Z",
			},
			expected: StrategyDetails{
				Type:       "Rolling Update",
				StepsCount: 4,
				Name:       "test",
				UpdatedAt:  "2022-01-02T00:00:00.000Z",
			},
		},
	}

	for name, tc := range cases {
		got := ConvertToStrategyDetails(tc.payload)
		if diff := cmp.Diff(tc.expected, got); diff != "" {
			t.Fatalf(name+"\n%s", diff)
		}
	}
}

func TestConvertToClusterDetails(t *testing.T) {
	cases := map[string]struct {
		payload  map[string]interface{}
		expected ClusterDetails
	}{
		"empty payload": {
			payload:  map[string]interface{}{},
			expected: ClusterDetails{},
		},
		"empty values": {
			payload: map[string]interface{}{
				"clusterInfo":       map[string]interface{}{},
				"controllerInfo":    map[string]interface{}{},
				"createdAt":         "",
				"id":                "",
				"lastHeartbeatTime": "",
				"updatedAt":         "",
			},
			expected: ClusterDetails{},
		},
		"fulfilled payload": {
			payload: map[string]interface{}{
				"clusterInfo": map[string]interface{}{
					"cloudProvider":          "test-provide",
					"kubernetesOrchestrator": "test-orchestrator",
					"kubernetesVersion":      "test-version",
				},
				"controllerInfo": map[string]interface{}{
					"nodeName":          "test-node",
					"podName":           "test-pod",
					"controllerVersion": "test-controller-version",
				},
				"createdAt":         "2022-01-01T00:00:00.000Z",
				"id":                "test-id",
				"lastHeartbeatTime": "2022-01-03T00:00:00.000Z",
				"updatedAt":         "2022-01-02T00:00:00.000Z",
			},
			expected: ClusterDetails{
				Name:              "test-id",
				K8sVersion:        "test-version",
				ControllerVersion: "test-controller-version",
				LastHeartbeat:     "2022-01-03T00:00:00.000Z",
				UpdatedAt:         "2022-01-02T00:00:00.000Z",
			},
		},
	}

	for name, tc := range cases {
		got := ConvertToClusterDetails(tc.payload)
		if diff := cmp.Diff(tc.expected, got); diff != "" {
			t.Fatalf(name+"\n%s", diff)
		}
	}
}

func TestConvertToVerificationProviderDetails(t *testing.T) {
	cases := map[string]struct {
		payload  map[string]interface{}
		expected VerificationProviderDetails
	}{
		"empty payload": {
			payload:  map[string]interface{}{},
			expected: VerificationProviderDetails{},
		},
		"empty values": {
			payload: map[string]interface{}{
				"clusterIds": []interface{}{},
				"name":       "",
				"newRelic":   map[string]interface{}{},
				"updatedAt":  "",
			},
			expected: VerificationProviderDetails{
				Types: "newRelic",
			},
		},
		"unexpected provider type": {
			payload: map[string]interface{}{
				"unexpectedName": map[string]interface{}{},
			},
			expected: VerificationProviderDetails{},
		},
		"fulfilled payload": {
			payload: map[string]interface{}{
				"clusterIds": []interface{}{"test, test1"},
				"name":       "test-name",
				"newRelic": map[string]interface{}{
					"accountId":      "1111",
					"personalApiKey": "1111",
					"region":         "test",
				},
				"updatedAt": "2022-01-02T00:00:00.000Z",
			},
			expected: VerificationProviderDetails{
				ClusterIDs: "test, test1",
				Types:      "newRelic",
				UpdatedAt:  "2022-01-02T00:00:00.000Z",
				Name:       "test-name",
			},
		},
	}

	for name, tc := range cases {
		got := ConvertToVerificationProviderDetails(tc.payload)
		if diff := cmp.Diff(tc.expected, got); diff != "" {
			t.Fatalf(name+"\n%s", diff)
		}
	}
}

func TestConvertToVerificationTemplateDetails(t *testing.T) {
	cases := map[string]struct {
		payload  map[string]interface{}
		expected VerificationTemplateDetails
	}{
		"empty payload": {
			payload:  map[string]interface{}{},
			expected: VerificationTemplateDetails{},
		},
		"empty values": {
			payload: map[string]interface{}{
				"args":      []interface{}{},
				"createdAt": "",
				"metrics":   []interface{}{},
				"name":      "",
				"updatedAt": "",
			},
			expected: VerificationTemplateDetails{},
		},
		"wrong type": {
			payload: map[string]interface{}{
				"args": []interface{}{"arg"},
			},
			expected: VerificationTemplateDetails{},
		},
		"no provider field": {
			payload: map[string]interface{}{
				"metrics": []interface{}{
					map[string]interface{}{
						"name": "avg-cpu-total-usage",
					},
				},
			},
			expected: VerificationTemplateDetails{
				Metrics: "avg-cpu-total-usage",
			},
		},
		"no metric name": {
			payload: map[string]interface{}{
				"metrics": []interface{}{
					map[string]interface{}{
						"provider": map[string]interface{}{
							"datadog": map[string]interface{}{
								"duration": "90s",
								"query":    "avg:kubernetes.cpu.usage.total{*} by {pod_name}",
							},
						},
					},
				},
			},
			expected: VerificationTemplateDetails{},
		},
		"unexpected provider": {
			payload: map[string]interface{}{
				"metrics": []interface{}{
					map[string]interface{}{
						"provider": map[string]interface{}{
							"unexpected": map[string]interface{}{
								"duration": "90s",
								"query":    "avg:kubernetes.cpu.usage.total{*} by {pod_name}",
							},
						},
					},
				},
			},
			expected: VerificationTemplateDetails{},
		},
		"fulfilled payload": {
			payload: map[string]interface{}{
				"args": []interface{}{
					map[string]interface{}{"name": "arg1"},
					map[string]interface{}{"name": "arg2"},
				},
				"createdAt": "2022-01-01T00:00:00.000Z",
				"metrics": []interface{}{
					map[string]interface{}{
						"name": "avg-cpu-total-usage",
						"provider": map[string]interface{}{
							"datadog": map[string]interface{}{
								"duration": "90s",
								"query":    "avg:kubernetes.cpu.usage.total{*} by {pod_name}",
							},
						},
					},
					map[string]interface{}{
						"name": "cpu-total-usage",
						"provider": map[string]interface{}{
							"prometheus": map[string]interface{}{
								"query": "sum(container_cpu_usage_seconds_total{namespace=\"prod-core\", endpoint=\"{{args.metric-name}}\"})",
							},
						},
					},
				},
				"name":      "test-name",
				"updatedAt": "2022-01-02T00:00:00.000Z",
			},
			expected: VerificationTemplateDetails{
				Name:      "test-name",
				Args:      "arg1, arg2",
				Metrics:   "avg-cpu-total-usage(datadog), cpu-total-usage(prometheus)",
				UpdatedAt: "2022-01-02T00:00:00.000Z",
			},
		},
	}

	for name, tc := range cases {
		got := ConvertToVerificationTemplateDetails(tc.payload)
		if diff := cmp.Diff(tc.expected, got); diff != "" {
			t.Fatalf(name+"\n%s", diff)
		}
	}
}
