package models

// TenantConfigurationEntry is one row of the tenant configuration response.
type TenantConfigurationEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FindTenantConfigurationByKey returns the first tenant configuration entry
// whose Key matches key, or nil when no such entry exists. Free function so
// it works directly on the slice type returned by the tenant endpoint.
func FindTenantConfigurationByKey(config []*TenantConfigurationEntry, key string) *TenantConfigurationEntry {
	for _, entry := range config {
		if entry != nil && entry.Key == key {
			return entry
		}
	}
	return nil
}

// FindProjectConfigurationByKey returns the first project configuration
// entry whose Key or Name matches keyOrName, or nil when no such entry
// exists. Mirrors the Cx1ClientGo GetConfigurationByKey helper.
func FindProjectConfigurationByKey(config []ProjectConfiguration, keyOrName string) *ProjectConfiguration {
	for i := range config {
		if config[i].Key == keyOrName || config[i].Name == keyOrName {
			return &config[i]
		}
	}
	return nil
}
