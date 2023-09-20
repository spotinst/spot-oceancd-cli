package operator

import (
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"spot-oceancd-cli/commons/configs"
)

type InstallationConfig struct {
	OceanCDConfig      OceanCDConfig              `json:"oceancd"`
	ArgoRolloutsConfig configs.ArgoRolloutsConfig `json:"argo"`
}

func (c *InstallationConfig) GetOperatorManagerConfig() *configs.OperatorManagerConfig {
	omConfig := &configs.OperatorManagerConfig{
		ArgoRolloutsConfig: configs.ArgoRolloutsConfig{
			General:    c.ArgoRolloutsConfig.General,
			Controller: configs.ArgoRolloutsControllerConfig{},
		},
		OceanCDConfig: configs.OceanCDConfig{
			Operator: c.OceanCDConfig.Operator,
		},
	}

	omConfig.OceanCDConfig.Operator.Namespace = c.OceanCDConfig.Namespace

	return omConfig
}

type OceanCDConfig struct {
	Namespace     string                        `json:"namespace"`
	ManagerConfig ManagerConfig                 `json:"manager"`
	Operator      configs.OceanCDOperatorConfig `json:"operator"`
}

type ManagerConfig struct {
	PodLabels        map[string]string             `json:"podLabels"`
	PodAnnotations   map[string]string             `json:"podAnnotations"`
	Labels           map[string]string             `json:"labels"`
	NodeSelector     map[string]string             `json:"nodeSelector"`
	Affinity         corev1.Affinity               `json:"affinity"`
	Tolerations      []corev1.Toleration           `json:"tolerations"`
	ExtraEnv         []corev1.EnvVar               `json:"extraEnv"`
	Resources        corev1.ResourceRequirements   `json:"resources"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets"`
}

func NewInstallationConfig(config map[string]interface{}) (*InstallationConfig, error) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal installation config: %w", err)
	}

	installationConfig := &InstallationConfig{}

	if err := json.Unmarshal(configBytes, installationConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installation config: %w", err)
	}

	return installationConfig, nil
}

func DefaultInstallationConfig() InstallationConfig {
	return InstallationConfig{
		OceanCDConfig: OceanCDConfig{
			Namespace: "oceancd",
		},
		ArgoRolloutsConfig: configs.ArgoRolloutsConfig{
			General: configs.ArgoRolloutsGeneralConfig{
				Namespaced: configs.Namespaced{Namespace: "argo-rollouts"},
			},
			Controller: configs.ArgoRolloutsControllerConfig{
				DeploymentSpecConfig: configs.DeploymentSpecConfig{
					Replicas: pointer.Int64(1),
				},
				Probes: configs.Probes{
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/healthz",
								Port: intstr.FromString("health"),
							},
						},
						InitialDelaySeconds: 30,
						PeriodSeconds:       20,
						FailureThreshold:    3,
						SuccessThreshold:    1,
						TimeoutSeconds:      10,
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/metrics",
								Port: intstr.FromString("metrics"),
							},
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       5,
						FailureThreshold:    3,
						SuccessThreshold:    1,
						TimeoutSeconds:      4,
					},
				},
			},
			Dashboard: configs.ArgoRolloutsDashboardConfig{
				Enabled:              false,
				DeploymentSpecConfig: configs.DeploymentSpecConfig{Replicas: pointer.Int64(1)},
			},
		},
	}
}

type InstallationPayload struct {
	Namespace string        `json:"namespace"`
	Manager   ManagerConfig `json:"manager"`
}

func NewInstallationPayload(config *InstallationConfig) InstallationPayload {
	return InstallationPayload{Namespace: config.OceanCDConfig.Namespace, Manager: config.OceanCDConfig.ManagerConfig}
}

type InstallationOutput struct {
	Argo    ManifestSet `json:"argo"`
	OceanCD ManifestSet `json:"oceancd"`
}

type ManifestSet struct {
	Appliable []string `json:"apply"`
	Patchable []string `json:"patch"`
}
