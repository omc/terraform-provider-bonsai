// Copyright (c) 2024 One More Cloud, Inc.
// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/omc/bonsai-api-go/v2/bonsai"
	"github.com/omc/terraform-provider-bonsai/internal/cluster"
	"github.com/omc/terraform-provider-bonsai/internal/plan"
	"github.com/omc/terraform-provider-bonsai/internal/release"
	"github.com/omc/terraform-provider-bonsai/internal/space"
)

// Ensure bonsaiProvider satisfies various provider interfaces.
var _ provider.Provider = &bonsaiProvider{}
var _ provider.ProviderWithFunctions = &bonsaiProvider{}

// bonsaiProvider defines the provider implementation.
type bonsaiProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	// bonsaiAPIClient is an optional override of the API Client used by
	// Terraform to perform requests. Defaults to nil, and will use a
	// default provided API Client.
	bonsaiAPIClient *bonsai.Client
}

// bonsaiProviderModel maps provider schema data to a Go type.
type bonsaiProviderModel struct {
	APIKey   types.String `tfsdk:"api_key"`
	APIToken types.String `tfsdk:"api_token"`
}

func (p *bonsaiProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bonsai"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *bonsaiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Bonsai provider is used to create and manage resources on the Bonsai.io platform." + "\n" +
			"To use the provider, you must provide both an API Access Key and Token, obtainable from within the " +
			"[Bonsai.io](https://bonsai.io) management panel!",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				// First line is at the bullet-point, following must be indented
				MarkdownDescription: "Bonsai.io API Access Key." + "\n\n" +
					"   - If not set, terraform will look for the `BONSAI_API_KEY` " +
					"   environment variable." + "\n\n" +
					"   - Obtainable from within the management panel at " +
					"   [Bonsai.io](https://bonsai.io)",
			},
			"api_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				// First line is at the bullet-point, following must be indented
				MarkdownDescription: "Bonsai.io API Access Token." + "\n\n" +
					"   - If not set, terraform will look for the `BONSAI_API_TOKEN` " +
					"   environment variable." + "\n\n" +
					"   - Obtainable from within the management panel at " +
					"   [Bonsai.io](https://bonsai.io)",
			},
		},
	}
}

func (p *bonsaiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		cluster.NewResource,
	}
}

func (p *bonsaiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cluster.NewDataSource,
		cluster.NewListDataSource,
		plan.NewDataSource,
		plan.NewListDataSource,
		release.NewDataSource,
		release.NewListDataSource,
		space.NewDataSource,
		space.NewListDataSource,
	}
}

func (p *bonsaiProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

// BonsaiProviderOption is a functional option, used to configure Client.
type BonsaiProviderOption func(*bonsaiProvider)

// WithAPIClient configures the Bonsai API Client used to perform
// terraform action HTTP requests.
func WithAPIClient(c *bonsai.Client) BonsaiProviderOption {
	return func(p *bonsaiProvider) {
		p.bonsaiAPIClient = c
	}
}

func WithVersion(version string) BonsaiProviderOption {
	return func(p *bonsaiProvider) {
		p.version = version
	}
}

func New(options ...BonsaiProviderOption) func() provider.Provider {
	return func() provider.Provider {
		p := &bonsaiProvider{}

		// apply options
		for _, option := range options {
			option(p)
		}

		return p
	}
}

func (p *bonsaiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config bonsaiProviderModel

	// Bonsai API Client has already been configured; skip all client configuration
	if p.bonsaiAPIClient != nil {
		// Make the Bonsai client available during DataSource and resource
		// type Configure methods.
		resp.DataSourceData = p.bonsaiAPIClient
		resp.ResourceData = p.bonsaiAPIClient
		return
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Bonsai API Key",
			"The provider cannot create the Bonsai API client as there is an unknown configuration value for the Bonsai API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BONSAI_API_KEY environment variable.",
		)
	}

	if config.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Bonsai API Token",
			"The provider cannot create the Bonsai API client as there is an unknown configuration value for the Bonsai API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BONSAI_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiKey := os.Getenv("BONSAI_API_KEY")
	apiToken := os.Getenv("BONSAI_API_TOKEN")

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Bonsai API Key",
			"The provider cannot create the Bonsai API client as there is a missing or empty value for the Bonsai API key. "+
				"Set the API Key value in the configuration or use the BONSAI_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Bonsai API Token",
			"The provider cannot create the Bonsai API client as there is a missing or empty value for the Bonsai API token. "+
				"Set the API Token value in the configuration or use the BONSAI_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	accessKey, err := bonsai.NewAccessKey(apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Bonsai API Client",
			"An unexpected error occurred when creating the Bonsai API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Bonsai Client Error: "+err.Error(),
		)
		return
	}

	accessToken, err := bonsai.NewAccessToken(apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Bonsai API Client",
			"An unexpected error occurred when creating the Bonsai API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Bonsai Client Error: "+err.Error(),
		)
		return
	}

	// Create a new Bonsai client using the configuration values
	client := bonsai.NewClient(
		bonsai.WithCredentialPair(
			bonsai.CredentialPair{
				AccessKey:   accessKey,
				AccessToken: accessToken,
			},
		),
	)

	// Make the Bonsai client available during DataSource and resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}
