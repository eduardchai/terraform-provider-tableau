package provider

import (
	"context"
	"os"
	"terraform-provider-tableau/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &tableauProvider{}
)

// tableauProvider is the provider implementation.
type tableauProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type tableauCloudProviderModel struct {
	ServerURL                 types.String `tfsdk:"server_url"`
	ApiVersion                types.String `tfsdk:"api_version"`
	PersonalAccessTokenName   types.String `tfsdk:"personal_access_token_name"`
	PersonalAccessTokenSecret types.String `tfsdk:"personal_access_token_secret"`
	Site                      types.String `tfsdk:"site"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &tableauProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *tableauProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tableau"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *tableauProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Tableau.",
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Description: "Server URL for Tableau. May also be provided via `TABLEAU_SERVER_URL` environment variable.",
				Optional:    true,
			},
			"api_version": schema.StringAttribute{
				Description: "API version for Tableau. May also be provided via `TABLEAU_API_VERSION` environment variable.",
				Optional:    true,
			},
			"personal_access_token_name": schema.StringAttribute{
				Description: "Personal Access Token (PAT) name for Tableau. May also be provided via `TABLEAU_PAT_NAME` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"personal_access_token_secret": schema.StringAttribute{
				Description: "Personal Access Token (PAT) secret for Tableau. May also be provided via `TABLEAU_PAT_SECRET` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"site": schema.StringAttribute{
				Description: "Site for Tableau. May also be provided via `TABLEAU_SITE` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *tableauProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Tableau client")

	// Retrieve provider data from configuration
	var config tableauCloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ServerURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Unknown Server URL",
			"The provider cannot create the Tableau API client as there is an unknown configuration value for the Tableau server_url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_SERVER_URL environment variable.",
		)
	}

	if config.ApiVersion.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_version"),
			"Unknown API Version",
			"The provider cannot create the Tableau API client as there is an unknown configuration value for the Tableau api_version. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_API_VERSION environment variable.",
		)
	}

	if config.PersonalAccessTokenName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_name"),
			"Unknown Personal Access Token (PAT) Name",
			"The provider cannot create the Tableau API client as there is an unknown configuration value for the Tableau personal_access_token_name. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_PAT_NAME environment variable.",
		)
	}

	if config.PersonalAccessTokenSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_secret"),
			"Unknown Personal Access Token (PAT) Secret",
			"The provider cannot create the Tableau API client as there is an unknown configuration value for the Tableau personal_access_token_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_PAT_SECRET environment variable.",
		)
	}

	if config.Site.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("site"),
			"Unknown Site",
			"The provider cannot create the Tableau API client as there is an unknown configuration value for the Tableau site. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_SITE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	serverURL := os.Getenv("TABLEAU_SERVER_URL")
	apiVersion := os.Getenv("TABLEAU_API_VERSION")
	personalAccessTokenName := os.Getenv("TABLEAU_PAT_NAME")
	personalAccessTokenSecret := os.Getenv("TABLEAU_PAT_SECRET")
	site := os.Getenv("TABLEAU_SITE")

	if !config.ServerURL.IsNull() {
		serverURL = config.ServerURL.ValueString()
	}
	if !config.ApiVersion.IsNull() {
		apiVersion = config.ApiVersion.ValueString()
	}
	if !config.PersonalAccessTokenName.IsNull() {
		personalAccessTokenName = config.PersonalAccessTokenName.ValueString()
	}
	if !config.PersonalAccessTokenSecret.IsNull() {
		personalAccessTokenSecret = config.PersonalAccessTokenSecret.ValueString()
	}
	if !config.Site.IsNull() {
		site = config.Site.ValueString()
	}

	if serverURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing Server URL",
			"The provider cannot create the Tableau API client as there is a missing configuration value for the Tableau server_url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_SERVER_URL environment variable.",
		)
	}

	if apiVersion == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_version"),
			"Missing API Version",
			"The provider cannot create the Tableau API client as there is a missing configuration value for the Tableau api_version. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_API_VERSION environment variable.",
		)
	}

	if personalAccessTokenName == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_name"),
			"Missing Personal Access Token (PAT) Name",
			"The provider cannot create the Tableau API client as there is a missing configuration value for the Tableau personal_access_token_name. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_PAT_NAME environment variable.",
		)
	}

	if personalAccessTokenSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_secret"),
			"Missing Personal Access Token (PAT) Secret",
			"The provider cannot create the Tableau API client as there is a missing configuration value for the Tableau personal_access_token_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_PAT_SECRET environment variable.",
		)
	}

	if site == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("site"),
			"Missing Site",
			"The provider cannot create the Tableau API client as there is a missing configuration value for the Tableau site. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TABLEAU_SITE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Tableau client")

	// Create a new Tableau client using the configuration values
	client, err := client.NewTableauClient(serverURL, apiVersion, site, personalAccessTokenName, personalAccessTokenSecret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tableau API Client",
			"An unexpected error occurred when creating the Tableau API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Tableau Client Error: "+err.Error(),
		)
		return
	}

	// Make the Tableau client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Tableau client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *tableauProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
		NewGroupDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *tableauProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewGroupResource,
		NewGroupMembershipResource,
	}
}
