package models

// SCMIntegration is one tenant SCM integration entry returned by the
// repos-manager v2 listing.
type SCMIntegration struct {
	ID          uint64 `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	RepoBaseURL string `json:"repoBaseUrl,omitempty"`
	AuthBaseURL string `json:"authBaseUrl,omitempty"`
	APIBaseURL  string `json:"apiBaseUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	ClientID    string `json:"clientId,omitempty"`
	Scope       string `json:"scope,omitempty"`
	OnPrem      bool   `json:"onPrem,omitempty"`
	RepoCount   int    `json:"repoCount,omitempty"`
}

// SCMRepository is the per-repository configuration entry returned by the
// repos-manager repo lookup. The Branches and Scm sub-shapes are kept loose
// (any-typed) for forward-compatibility — the wire format has been evolving.
type SCMRepository struct {
	ID                             string `json:"id,omitempty"`
	URL                            string `json:"url,omitempty"`
	WebhookID                      string `json:"webhookId,omitempty"`
	WebhookEnabled                 bool   `json:"webhookEnabled,omitempty"`
	PrDecorationEnabled            bool   `json:"prDecorationEnabled,omitempty"`
	IsRepoAdmin                    bool   `json:"isRepoAdmin,omitempty"`
	SastIncrementalScan            bool   `json:"sastIncrementalScan,omitempty"`
	SastScannerEnabled             bool   `json:"sastScannerEnabled,omitempty"`
	ScaScannerEnabled              bool   `json:"scaScannerEnabled,omitempty"`
	KicsScannerEnabled             bool   `json:"kicsScannerEnabled,omitempty"`
	APISecScannerEnabled           bool   `json:"apiSecScannerEnabled,omitempty"`
	ContainerScannerEnabled        bool   `json:"containerScannerEnabled,omitempty"`
	OssfScoreCardScannerEnabled    bool   `json:"ossfScoreCardScannerEnabled,omitempty"`
	SecretsDetectionScannerEnabled bool   `json:"secretsDetectionScannerEnabled,omitempty"`
	PrivatePackage                 bool   `json:"privatePackage,omitempty"`
	Criticality                    int    `json:"criticality,omitempty"`
	Branches                       []SCMRepositoryBranch `json:"branches,omitempty"`
	RemediationSeverities          []any                 `json:"remediationSeverities,omitempty"`
	SSHRepoURL                     string                `json:"sshRepoUrl,omitempty"`
	ScmID                          int                   `json:"scmId,omitempty"`
}

// SCMRepositoryBranch is one branch entry on an SCM repository.
type SCMRepositoryBranch struct {
	Pattern         string         `json:"pattern,omitempty"`
	IsDefaultBranch bool           `json:"isDefaultBranch,omitempty"`
	Tags            map[string]any `json:"tags,omitempty"`
}
