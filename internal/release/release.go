package release

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

const (
	dataSourceMarkdownDescription = "A Release is a version of Elasticsearch " +
		"available to your account."
	listDataSourceMarkdownDescription = "A list of all **available** releases on " +
		"your account."
)

// model maps releases schema data.
type model struct {
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	ServiceType types.String `tfsdk:"service_type"`
	Version     types.String `tfsdk:"version"`
	MultiTenant types.Bool   `tfsdk:"multitenant"`
}

func convert(r bonsai.Release) model {
	return model{
		Name:        types.StringValue(r.Name),
		Slug:        types.StringValue(r.Slug),
		ServiceType: types.StringValue(r.ServiceType),
		Version:     types.StringValue(r.Version),
		MultiTenant: types.BoolValue(r.MultiTenant),
	}
}

func schemaAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"name": dschema.StringAttribute{
			MarkdownDescription: "The name for the release.",
			Computed:            true,
		},
		"slug": dschema.StringAttribute{
			MarkdownDescription: "The machine-readable name for the deployment.",
			Computed:            true,
			Optional:            true,
		},
		"service_type": dschema.StringAttribute{
			MarkdownDescription: "The service type of the deployment - for " +
				"example, \"elasticsearch\".",
			Computed: true,
		},
		"version": dschema.StringAttribute{
			MarkdownDescription: "The version of the release.",
			Computed:            true,
		},
		"multitenant": dschema.BoolAttribute{
			Computed: true,
			MarkdownDescription: "Whether the release is available on " +
				"multitenant deployments.",
		},
	}
}
