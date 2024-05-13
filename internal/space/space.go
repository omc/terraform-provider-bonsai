package space

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

const (
	dataSourceMarkdownDescription = "A Space represents the server groups and " +
		"geographic regions available to a [Bonsai.io](https://bonsai.io) " +
		"account, where clusters may be provisioned."
	listDataSourceMarkdownDescription = "A list of all **available** spaces on " +
		"your account."
)

// cloudProviderModel maps space cloud provider details.
type cloudProviderModel struct {
	Provider types.String `tfsdk:"provider"`
	Region   types.String `tfsdk:"region"`
}

// model maps spaces schema data.
type model struct {
	Path           types.String `tfsdk:"path"`
	PrivateNetwork types.Bool   `tfsdk:"private_network"`

	Cloud cloudProviderModel `tfsdk:"cloud"`
}

func convert(s bonsai.Space) model {
	return model{
		Path:           types.StringValue(s.Path),
		PrivateNetwork: types.BoolValue(s.PrivateNetwork),

		Cloud: cloudProviderModel{
			Provider: types.StringValue(s.Cloud.Provider),
			Region:   types.StringValue(s.Cloud.Region),
		},
	}
}

func schemaAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"path": dschema.StringAttribute{
			MarkdownDescription: "A machine-readable name for the server group.",
			Computed:            true,
			Optional:            true,
		},
		"private_network": dschema.BoolAttribute{
			Computed: true,
			MarkdownDescription: "Indicates whether the space is isolated and " +
				"inaccessible from the public Internet. A VPC connection will " +
				"be needed to communicate with a private cluster.",
		},
		"cloud": dschema.SingleNestedAttribute{
			MarkdownDescription: "Details about the cloud provider and region attributes.",
			Computed:            true,
			Attributes: map[string]dschema.Attribute{
				"provider": dschema.StringAttribute{
					MarkdownDescription: "A machine-readable name for the cloud provider in which this space is deployed.",
					Computed:            true,
				},
				"region": dschema.StringAttribute{
					MarkdownDescription: "A machine-readable name for the geographic region of the server group.",
					Computed:            true,
				},
			},
		},
	}
}
