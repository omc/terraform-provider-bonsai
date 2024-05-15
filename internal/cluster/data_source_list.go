package cluster

import (
	"context"
	"fmt"

	tfds "github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

// listDataSourceModel maps the data source schema data.
type listDataSourceModel struct {
	Clusters []dataSourceModel `tfsdk:"clusters"`
}

// listDataSource is the data source implementation.
type listDataSource struct {
	client *bonsai.ClusterClient
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ tfds.DataSource = &listDataSource{}
)

// NewListDataSource is a helper function to simplify the provider implementation.
func NewListDataSource() tfds.DataSource {
	return &listDataSource{}
}

// Metadata returns the data source type name.
func (d *listDataSource) Metadata(_ context.Context, req tfds.MetadataRequest, resp *tfds.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

// Schema defines the schema for the data source.
func (d *listDataSource) Schema(_ context.Context, _ tfds.SchemaRequest, resp *tfds.SchemaResponse) {
	resp.Schema = dschema.Schema{
		MarkdownDescription: listDataSourceMarkdownDescription,
		Attributes: map[string]dschema.Attribute{
			"clusters": dschema.ListNestedAttribute{
				MarkdownDescription: dataSourceMarkdownDescription,
				Computed:            true,
				NestedObject: dschema.NestedAttributeObject{
					Attributes: dataSourceSchemaAttributes(),
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *listDataSource) Read(ctx context.Context, req tfds.ReadRequest, resp *tfds.ReadResponse) {
	var state listDataSourceModel

	clusters, err := d.client.Cluster.All(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Bonsai Clusters from the Bonsai API",
			err.Error(),
		)
		return
	}

	// Map response body to dataSourceModel
	for _, c := range clusters {
		clusterState := dataSourceConvert(c)
		state.Clusters = append(state.Clusters, clusterState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the data source.
func (d *listDataSource) Configure(_ context.Context, req tfds.ConfigureRequest, resp *tfds.ConfigureResponse) {
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

	d.client = &client.Cluster
}