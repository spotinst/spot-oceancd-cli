package operator

type ClusterManifestsMetadataResponse struct {
	Operator []ManifestMetadata `json:"operator"`
	Argo     []ManifestMetadata `json:"argo"`
	OM       []ManifestMetadata `json:"om"`
	Secrets  []SecretMetadata   `json:"secrets"`
}

type ManifestMetadata struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type SecretMetadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DeleteClusterResponse struct {
	IsDeleted bool `json:"isDeleted"`
}
