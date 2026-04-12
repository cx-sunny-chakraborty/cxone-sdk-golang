package models

// FeatureFlag is one entry in the feature-flags response.
type FeatureFlag struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

// Common Checkmarx One feature flag identifiers. Add new ones here as the
// platform exposes them so callers can reference well-known constants
// instead of stringly-typed flag names.
const (
	FeatureFlagCustomStates             = "CUSTOM_STATES_ENABLED"
	FeatureFlagPackageEnforcement       = "PACKAGE_ENFORCEMENT_ENABLED"
	FeatureFlagCVSSV3                   = "CVSS_V3_ENABLED"
	FeatureFlagMinIO                    = "MINIO_ENABLED"
	FeatureFlagSastCustomStates         = "SAST_CUSTOM_STATES_ENABLED"
	FeatureFlagRiskManagement           = "RISK_MANAGEMENT_IDES_PROJECT_RESULTS_SCORES_API_ENABLED"
	FeatureFlagOSSRealtime              = "OSS_REALTIME_ENABLED"
	FeatureFlagSCSLicensingV2           = "SSCS_NEW_LICENSING_ENABLED"
	FeatureFlagSSCSCommitHistory        = "SSCS_COMMIT_HISTORY_ENABLED"
	FeatureFlagDirectAppAssociation     = "DIRECT_APP_ASSOCIATION_ENABLED"
	FeatureFlagDAMigration              = "DA_MIGRATION_ENABLED"
	FeatureFlagIncreaseFileUploadLimit  = "INCREASE_FILE_UPLOAD_LIMIT"
	FeatureFlagSCADeltaScan             = "SCA_DELTASCAN_ENABLED"
)
