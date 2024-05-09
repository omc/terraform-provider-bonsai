// Copyright (c) 2024 One More Cloud, Inc.
// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"os"

	"github.com/omc/bonsai-api-go/v2/bonsai"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
}

// hashicupsProviderModel maps provider schema data to a Go type.
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
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"api_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *bonsaiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *bonsaiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *bonsaiProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bonsaiProvider{
			version: version,
		}
	}
}

func (p *bonsaiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config bonsaiProviderModel

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

	// Create a new HashiCups client using the configuration values
	client := bonsai.NewClient(
		bonsai.WithCredentialPair(
			bonsai.CredentialPair{
				AccessKey:   accessKey,
				AccessToken: accessToken,
			},
		),
	)

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}
