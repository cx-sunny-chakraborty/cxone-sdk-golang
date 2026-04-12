package models

// ScanResultsCollection is the response of GET /api/results — a paginated
// envelope of cross-engine results for a single scan.
type ScanResultsCollection struct {
	Results    []*ScanResult `json:"results"`
	TotalCount uint          `json:"totalCount"`
	ScanID     string        `json:"scanID"`
}

// ScanResult is one finding from any engine. The Type field discriminates
// between SAST, SCA, KICS, container, secrets, etc.; engine-specific fields
// live nested under [ScanResultData] or [VulnerabilityDetails].
type ScanResult struct {
	Type                 string               `json:"type,omitempty"`
	ScaType              string               `json:"scaType,omitempty"`
	Label                string               `json:"label,omitempty"`
	ID                   string               `json:"id,omitempty"`
	SimilarityID         string               `json:"similarityId,omitempty"`
	AlternateID          string               `json:"alternateId,omitempty"`
	Status               string               `json:"status,omitempty"`
	State                string               `json:"state,omitempty"`
	Severity             string               `json:"severity,omitempty"`
	Created              string               `json:"created,omitempty"`
	FirstFoundAt         string               `json:"firstFoundAt,omitempty"`
	FoundAt              string               `json:"foundAt,omitempty"`
	FirstScan            string               `json:"firstScan,omitempty"`
	FirstScanID          string               `json:"firstScanId,omitempty"`
	PublishedAt          string               `json:"publishedAt,omitempty"`
	Recommendations      string               `json:"recommendations,omitempty"`
	Description          string               `json:"description,omitempty"`
	DescriptionHTML      string               `json:"descriptionHTML,omitempty"`
	ScanResultData       ScanResultData       `json:"data,omitempty"`
	Comments             ResultComments       `json:"comments,omitempty"`
	VulnerabilityDetails VulnerabilityDetails `json:"vulnerabilityDetails,omitempty"`
}

// ScanResultData is the engine-specific payload of a [ScanResult]. Most
// fields are populated only for one engine type; consult Type on the parent
// to know which fields are meaningful.
//
// QueryID is intentionally typed as `any` because the wire encoding mixes
// integers and strings depending on the engine; let the caller assert.
type ScanResultData struct {
	QueryID               any                      `json:"queryId,omitempty"`
	QueryName             string                   `json:"queryName,omitempty"`
	Group                 string                   `json:"group,omitempty"`
	ResultHash            string                   `json:"resultHash,omitempty"`
	LanguageName          string                   `json:"languageName,omitempty"`
	Redundancy            string                   `json:"redundancy,omitempty"`
	Description           string                   `json:"description,omitempty"`
	Nodes                 []*ScanResultNode        `json:"nodes,omitempty"`
	PackageData           []*ScanResultPackageData `json:"packageData,omitempty"`
	PackageID             []*ScanResultPackageData `json:"packageId,omitempty"`
	PackageIdentifier     string                   `json:"packageIdentifier,omitempty"`
	PublishedAt           string                   `json:"publishedAt,omitempty"`
	ScaPackageCollection  *ScaPackageCollection    `json:"scaPackageData,omitempty"`
	RecommendedVersion    any                      `json:"recommendedVersion,omitempty"`
	Line                  uint                     `json:"line,omitempty"`
	Platform              string                   `json:"platform,omitempty"`
	IssueType             string                   `json:"issueType,omitempty"`
	ExpectedValue         string                   `json:"expectedValue,omitempty"`
	Value                 string                   `json:"value,omitempty"`
	Filename              string                   `json:"filename,omitempty"`
	PackageName           string                   `json:"packageName,omitempty"`
	PackageVersion        string                   `json:"packageVersion,omitempty"`
	ImageName             string                   `json:"imageName,omitempty"`
	ImageTag              string                   `json:"imageTag,omitempty"`
	ImageFilePath         string                   `json:"imageFilePath,omitempty"`
	ImageOrigin           string                   `json:"imageOrigin,omitempty"`
	RuleID                *string                  `json:"ruleId,omitempty"`
	RuleName              string                   `json:"ruleName,omitempty"`
	Snippet               string                   `json:"snippet,omitempty"`
	SlsaStep              string                   `json:"slsaStep,omitempty"`
	RuleDescription       string                   `json:"ruleDescription,omitempty"`
	Remediation           string                   `json:"remediation,omitempty"`
	RemediationLink       string                   `json:"remediationLink,omitempty"`
	RemediationAdditional string                   `json:"remediationAdditional,omitempty"`
	Validity              string                   `json:"validity,omitempty"`
	IsInSource            bool                     `json:"isInSource,omitempty"`
	CommitURL             string                   `json:"commitUrl,omitempty"`
}

// ScanResultNode is one node in a SAST data flow.
type ScanResultNode struct {
	ID          string `json:"id,omitempty"`
	Line        uint   `json:"line"`
	Name        string `json:"name,omitempty"`
	Column      uint   `json:"column"`
	Length      uint   `json:"length,omitempty"`
	Method      string `json:"method,omitempty"`
	NodeID      int    `json:"nodeID,omitempty"`
	DomType     string `json:"domType,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	FullName    string `json:"fullName,omitempty"`
	TypeName    string `json:"typeName,omitempty"`
	MethodLine  uint   `json:"methodLine,omitempty"`
	Definitions string `json:"definitions,omitempty"`
}

// ScanResultPackageData is one entry in the per-result package metadata
// (typically a remediation link or comment).
type ScanResultPackageData struct {
	Comment string `json:"comment,omitempty"`
	Type    string `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
}

// ResultComments holds free-form comments attached to a result.
type ResultComments struct {
	Comments string `json:"comments,omitempty"`
}

// VulnerabilityDetails carries SCA/KICS vulnerability metadata.
type VulnerabilityDetails struct {
	CweID       any               `json:"cweId,omitempty"`
	CvssScore   float64           `json:"cvssScore,omitempty"`
	CveName     string            `json:"cveName,omitempty"`
	CVSS        VulnerabilityCVSS `json:"cvss,omitempty"`
	Compliances []*string         `json:"compliances,omitempty"`
}

// VulnerabilityCVSS is a CVSS vector breakdown.
type VulnerabilityCVSS struct {
	Version            int    `json:"version,omitempty"`
	AttackVector       string `json:"attackVector,omitempty"`
	Availability       string `json:"availability,omitempty"`
	Confidentiality    string `json:"confidentiality,omitempty"`
	AttackComplexity   string `json:"attackComplexity,omitempty"`
	IntegrityImpact    string `json:"integrityImpact,omitempty"`
	Scope              string `json:"scope,omitempty"`
	PrivilegesRequired string `json:"privilegesRequired,omitempty"`
	UserInteraction    string `json:"userInteraction,omitempty"`
}

// ScaPackageCollection is the SCA-specific package detail nested in
// [ScanResultData] for SCA findings.
type ScaPackageCollection struct {
	ID                      string             `json:"id,omitempty"`
	FixLink                 string             `json:"fixLink,omitempty"`
	BestPackageLink         string             `json:"bestPackageLink,omitempty"`
	Locations               []*string          `json:"locations,omitempty"`
	DependencyPathArray     [][]DependencyPath `json:"dependencyPaths,omitempty"`
	Outdated                bool               `json:"outdated,omitempty"`
	SupportsQuickFix        bool               `json:"supportsQuickFix"`
	IsDirectDependency      bool               `json:"isDirectDependency"`
	TypeOfDependency        string             `json:"typeOfDependency"`
	IsDevelopmentDependency bool               `json:"isDevelopmentDependency"`
	IsTestDependency        bool               `json:"isTestDependency"`
}

// DependencyPath is one node in an SCA dependency path.
type DependencyPath struct {
	ID               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Version          string    `json:"version,omitempty"`
	IsResolved       bool      `json:"isResolved,omitempty"`
	IsDevelopment    bool      `json:"isDevelopment,omitempty"`
	Locations        []*string `json:"locations,omitempty"`
	SupportsQuickFix bool      `json:"supportsQuickFix,omitempty"`
}

// ScanSummariesModel is the response of the per-scan summary endpoint.
type ScanSummariesModel struct {
	ScansSummaries []ScanSummaries `json:"scansSummaries,omitempty"`
	TotalCount     int             `json:"totalCount,omitempty"`
}

// ScanSummaries holds the per-engine counters for one scan.
type ScanSummaries struct {
	SastCounters          SastCounters          `json:"sastCounters,omitempty"`
	KicsCounters          KicsCounters          `json:"kicsCounters,omitempty"`
	ScaCounters           ScaCounters           `json:"scaCounters,omitempty"`
	ScaContainersCounters ScaContainersCounters `json:"scaContainersCounters,omitempty"`
}

// SastCounters / KicsCounters / ScaCounters share the same shape — kept as
// distinct types because the API may diverge in future versions and the
// type names communicate intent at call sites.
type SastCounters struct {
	SeverityCounters    []SeverityCounters `json:"SeverityCounters,omitempty"`
	TotalCounter        int                `json:"totalCounter,omitempty"`
	FilesScannedCounter int                `json:"filesScannedCounter,omitempty"`
}

type KicsCounters struct {
	SeverityCounters    []SeverityCounters `json:"SeverityCounters,omitempty"`
	TotalCounter        int                `json:"totalCounter,omitempty"`
	FilesScannedCounter int                `json:"filesScannedCounter,omitempty"`
}

type ScaCounters struct {
	SeverityCounters    []SeverityCounters `json:"SeverityCounters,omitempty"`
	TotalCounter        int                `json:"totalCounter,omitempty"`
	FilesScannedCounter int                `json:"filesScannedCounter,omitempty"`
}

type ScaContainersCounters struct {
	SeverityCounters            []SeverityCounters `json:"severityVulnerabilitiesCounters,omitempty"`
	TotalPackagesCounter        int                `json:"totalPackagesCounter,omitempty"`
	TotalVulnerabilitiesCounter int                `json:"totalVulnerabilitiesCounter,omitempty"`
}

// SeverityCounters is one (severity, counter) pair inside an engine counters
// envelope.
type SeverityCounters struct {
	Severity string `json:"severity,omitempty"`
	Counter  int    `json:"counter,omitempty"`
}
