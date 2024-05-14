package release

import (
	"context"
	"fmt"

	tfds "github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

// dataSource is the data source implementation.
type dataSource struct {
	client *bonsai.ReleaseClient
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ tfds.DataSource = &dataSource{}
)

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() tfds.DataSource {
	return &dataSource{}
}

// Metadata returns the data source type name.
func (d *dataSource) Metadata(_ context.Context, req tfds.MetadataRequest, resp *tfds.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_release"
}

// Schema defines the schema for the data source.
func (d *dataSource) Schema(_ context.Context, _ tfds.SchemaRequest, resp *tfds.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Attributes:          schemaAttributes(),
		MarkdownDescription: dataSourceMarkdownDescription,
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dataSource) Read(ctx context.Context, req tfds.ReadRequest, resp *tfds.ReadResponse) {
	var state model

	// Fetch requested Slug from context
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("slug"), &state.Slug)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Slug.IsNull() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read Bonsai Release (%s) from the Bonsai API", state.Slug.ValueString()),
			"expected 'slug' option to be set",
		)
		return
	}

	s, err := d.client.Release.GetBySlug(ctx, state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read Bonsai Release (%s) from the Bonsai API", state.Slug.ValueString()),
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = convert(s)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the data source.
func (d *dataSource) Configure(_ context.Context, req tfds.ConfigureRequest, resp *tfds.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*bonsai.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *bonsai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = &client.Release
}
