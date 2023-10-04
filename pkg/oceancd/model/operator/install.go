package operator

import (
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"spot-oceancd-operator-commons/component_configs"
)

type InstallationConfig struct {
	OceanCDConfig      OceanCDConfig                        `json:"oceancd"`
	ArgoRolloutsConfig component_configs.ArgoRolloutsConfig `json:"argo"`
}

func (c *InstallationConfig) GetData() map[string]interface{} {
	bytes, _ := json.Marshal(c)

	data := map[string]interface{}{}
	_ = json.Unmarshal(bytes, &data)

	return data
}

func (c *InstallationConfig) GetOperatorManagerConfig() *component_configs.OperatorManagerConfig {
	omConfig := &component_configs.OperatorManagerConfig{
		ArgoRolloutsConfig: c.ArgoRolloutsConfig,
		OceanCDConfig: component_configs.OceanCDConfig{
			Operator: c.OceanCDConfig.Operator,
		},
	}

	omConfig.OceanCDConfig.Operator.Namespace = c.OceanCDConfig.Namespace

	return omConfig
}

type OceanCDConfig struct {
	Namespace     string                                  `json:"namespace"`
	ManagerConfig ManagerConfig                           `json:"manager"`
	Operator      component_configs.OceanCDOperatorConfig `json:"operator"`
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

func NewInstallationConfig(data map[string]interface{}) (*InstallationConfig, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal installation config: %w", err)
	}

	config := DefaultInstallationConfig()

	if err := json.Unmarshal(dataBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installation config: %w", err)
	}

	return &config, nil
}

func DefaultInstallationConfig() InstallationConfig {
	return InstallationConfig{
		OceanCDConfig: OceanCDConfig{
			Namespace: "oceancd",
			ManagerConfig: ManagerConfig{
				PodLabels:        map[string]string{},
				PodAnnotations:   map[string]string{},
				Labels:           map[string]string{},
				NodeSelector:     map[string]string{},
				Tolerations:      []corev1.Toleration{},
				ExtraEnv:         []corev1.EnvVar{},
				ImagePullSecrets: []corev1.LocalObjectReference{},
			},
			Operator: component_configs.OceanCDOperatorConfig{
				MetadataConfig: component_configs.MetadataConfig{
					Labeled:   component_configs.Labeled{Labels: map[string]string{}},
					Annotated: component_configs.Annotated{Annotations: map[string]string{}},
				},
				ExtraEnv:                  []corev1.EnvVar{},
				NodeSelector:              map[string]string{},
				Tolerations:               []corev1.Toleration{},
				ImagePullSecrets:          []corev1.LocalObjectReference{},
				ServiceAccountAnnotations: map[string]string{},
			},
		},
		ArgoRolloutsConfig: component_configs.ArgoRolloutsConfig{
			General: component_configs.ArgoRolloutsGeneralConfig{
				Namespaced:         component_configs.Namespaced{Namespace: "argo-rollouts"},
				Labeled:            component_configs.Labeled{Labels: map[string]string{}},
				PodAnnotations:     map[string]string{},
				PodLabels:          map[string]string{},
				ServiceAnnotations: map[string]string{},
				ServiceLabels:      map[string]string{},
			},
			Controller: component_configs.ArgoRolloutsControllerConfig{
				Replicas:                 pointer.Int64(1),
				NodeSelector:             map[string]string{},
				Tolerations:              []corev1.Toleration{},
				ExtraArgs:                map[string]string{},
				ExtraEnv:                 []corev1.EnvVar{},
				ImagePullSecrets:         []corev1.LocalObjectReference{},
				ContainerSecurityContext: &corev1.SecurityContext{},
				Probes: component_configs.Probes{
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
			Dashboard: component_configs.ArgoRolloutsDashboardConfig{
				Enabled:                  false,
				Replicas:                 pointer.Int64(1),
				NodeSelector:             map[string]string{},
				Tolerations:              []corev1.Toleration{},
				ExtraEnv:                 []corev1.EnvVar{},
				ImagePullSecrets:         []corev1.LocalObjectReference{},
				ContainerSecurityContext: &corev1.SecurityContext{},
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
	OM OM `json:"om"`
}

type OM struct {
	Manifests []string `json:"manifests"`
}
