package userflows

import "terraform-provider-stratoid/internal/helpers/tf/pluginsdk"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "User Flows"
}

// AssociatedGitHubLabel is the issue/PR label which can be applied to PRs that include changes to this service package
func (r Registration) AssociatedGitHubLabel() string {
	return "feature/user-flows"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"User Flows",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"stratoid_user_flow": userFlowDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"stratoid_user_flow": userFlowResource(),
	}
}
