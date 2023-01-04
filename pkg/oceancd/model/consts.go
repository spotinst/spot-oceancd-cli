package model

const (
	ClusterEntity       = "cluster"
	ClusterEntityPlural = "clusters"
	ClusterKind         = "Cluster"

	StrategyEntity       = "strategy"
	StrategyEntityPlural = "strategies"
	StrategyKind         = "Strategy"

	RolloutSpecEntity       = "rolloutSpec"
	RolloutSpecEntityPlural = "rolloutSpecs"
	RolloutSpecKind         = "RolloutSpec"
	RolloutSpecShort        = "rs"

	VerificationProviderEntity       = "verificationProvider"
	VerificationProviderEntityPlural = "verificationProviders"
	VerificationProviderKind         = "VerificationProvider"

	VerificationTemplateEntity       = "verificationTemplate"
	VerificationTemplateEntityPlural = "verificationTemplates"
	VerificationTemplateKind         = "VerificationTemplate"

	Prometheus = "prometheus"
	Datadog    = "datadog"
	NewRelic   = "newRelic"
	CloudWatch = "cloudWatch"
	Web        = "web"

	BackgroundVerificationLabel = "Background"
)

var (
	StrategyEntityShorts       = []string{"stg", "stgs"}
	VerificationProviderShorts = []string{"vp", "vps"}
	VerificationTemplateShorts = []string{"vt", "vts"}
)
