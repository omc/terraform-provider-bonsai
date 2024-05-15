package cluster

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

const (
	dataSourceMarkdownDescription     = "Cluster represents a single cluster on your account."
	listDataSourceMarkdownDescription = "A list of all **active** clusters on " +
		"your account."

	resourceMarkdownDescription = "Provides and manages a Cluster on your account."
)

type planModel struct {
	Slug types.String `tfsdk:"slug"`
	URI  types.String `tfsdk:"uri"`
}

type releaseModel struct {
	ServiceType types.String `tfsdk:"service_type"`
	PackageName types.String `tfsdk:"package_name"`
	Version     types.String `tfsdk:"version"`
	Slug        types.String `tfsdk:"slug"`
	URI         types.String `tfsdk:"uri"`
}

type spaceModel struct {
	Path   types.String `tfsdk:"path"`
	Region types.String `tfsdk:"region"`
	URI    types.String `tfsdk:"uri"`
}

type statsModel struct {
	Docs          types.Int64 `tfsdk:"docs"`
	ShardsUsed    types.Int64 `tfsdk:"shards_used"`
	DataBytesUsed types.Int64 `tfsdk:"data_bytes_used"`
}

type stateModel struct {
	State types.String `tfsdk:"state"`
}

type accessModel struct {
	Host   types.String `tfsdk:"host"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`

	Username types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	URL      types.String `tfsdk:"url"`
}

// resourceModel maps clusters schema data.
type resourceModel struct {
	// ID is a unique identifier, only set for terraform's management.
	// For Cluster, this is set to the Slug.
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
	URI  types.String `tfsdk:"uri"`

	// Message received during Cluster.Create
	Message types.String `tfsdk:"message"`
	// Monitor received during Cluster.Create
	Monitor types.String `tfsdk:"monitor"`

	Plan    planModel    `tfsdk:"plan"`
	Release releaseModel `tfsdk:"release"`
	Space   spaceModel   `tfsdk:"space"`
	Stats   types.Object `tfsdk:"stats"`
	Access  types.Object `tfsdk:"access"`
	State   types.Object `tfsdk:"state"`
}

// dataSourceModel maps clusters schema data.
type dataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
	URI  types.String `tfsdk:"uri"`

	Plan    planModel    `tfsdk:"plan"`
	Release releaseModel `tfsdk:"release"`
	Space   spaceModel   `tfsdk:"space"`
	Stats   statsModel   `tfsdk:"stats"`
	Access  accessModel  `tfsdk:"access"`
	State   stateModel   `tfsdk:"state"`
}

func dataSourceConvert(c bonsai.Cluster) dataSourceModel {
	return dataSourceModel{
		Name: types.StringValue(c.Name),
		Slug: types.StringValue(c.Slug),
		URI:  types.StringValue(c.URI),
		Plan: planModel{
			Slug: types.StringValue(c.Plan.Slug),
			URI:  types.StringValue(c.Plan.URI),
		},
		Release: releaseModel{
			ServiceType: types.StringValue(c.Release.ServiceType),
			PackageName: types.StringValue(c.Release.PackageName),
			Version:     types.StringValue(c.Release.Version),
			Slug:        types.StringValue(c.Release.Slug),
			URI:         types.StringValue(c.Release.URI),
		},
		Space: spaceModel{
			Path:   types.StringValue(c.Space.Path),
			Region: types.StringValue(c.Space.Region),
			URI:    types.StringValue(c.Space.URI),
		},
		Stats: statsModel{
			Docs:          types.Int64Value(c.Stats.Docs),
			ShardsUsed:    types.Int64Value(c.Stats.ShardsUsed),
			DataBytesUsed: types.Int64Value(c.Stats.DataBytesUsed),
		},
		Access: accessModel{
			Host:     types.StringValue(c.Access.Host),
			Port:     types.Int64Value(int64(c.Access.Port)),
			Scheme:   types.StringValue(c.Access.Scheme),
			Username: types.StringValue(c.Access.Username),
			Password: types.StringValue(c.Access.Password),
			URL:      types.StringValue(c.Access.URL),
		},
		State: stateModel{State: types.StringValue(string(c.State))},
	}
}

func dataSourceSchemaAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"name": dschema.StringAttribute{
			MarkdownDescription: "The human-readable name of the cluster.",
			Computed:            true,
		},
		"slug": dschema.StringAttribute{
			MarkdownDescription: "A unique, machine-readable name for the " +
				"cluster. A cluster slug is based its name at creation, to " +
				"which a random integer is concatenated.",
			Computed: true,
			Optional: true,
		},
		"uri": dschema.StringAttribute{
			MarkdownDescription: "A URI to retrieve more information about this cluster.",
			Computed:            true,
		},
		"plan": dschema.SingleNestedAttribute{
			MarkdownDescription: "Plan holds some information about the cluster's current subscription plan.",
			Computed:            true,
			Attributes: map[string]dschema.Attribute{
				"slug": dschema.StringAttribute{
					MarkdownDescription: "A machine-readable name for the plan.",
					Computed:            true,
				},
				"uri": dschema.StringAttribute{
					MarkdownDescription: "A URI to retrieve more information about this Plan.",
					Computed:            true,
				},
			},
		},
		"release": dschema.SingleNestedAttribute{
			MarkdownDescription: "Release holds some information about the cluster's current release.",
			Computed:            true,
			Attributes: map[string]dschema.Attribute{
				"service_type": dschema.StringAttribute{
					MarkdownDescription: "The service type of the deployment - for example, \"elasticsearch\".",
					Computed:            true,
				},
				"package_name": dschema.StringAttribute{
					MarkdownDescription: "PackageName is the package name of the release.",
					Computed:            true,
				},
				"version": dschema.StringAttribute{
					MarkdownDescription: "The version of the release.",
					Computed:            true,
				},
				"slug": dschema.StringAttribute{
					MarkdownDescription: "The machine-readable name for the deployment.",
					Computed:            true,
				},
				"uri": dschema.StringAttribute{
					MarkdownDescription: "A URI to retrieve more information about this Release.",
					Computed:            true,
				},
			},
		},
		"space": dschema.SingleNestedAttribute{
			MarkdownDescription: "Space holds some information about where the cluster is running.",
			Computed:            true,
			Attributes: map[string]dschema.Attribute{
				"path": dschema.StringAttribute{
					MarkdownDescription: "A machine-readable name for the server group.",
					Computed:            true,
				},
				"region": dschema.StringAttribute{
					MarkdownDescription: "The geographic region in which the cluster is running.",
					Computed:            true,
				},
				"uri": dschema.StringAttribute{
					MarkdownDescription: "A URI to retrieve more information about this Space.",
					Computed:            true,
				},
			},
		},
		"stats": dschema.SingleNestedAttribute{
			MarkdownDescription: "Stats holds *some* statistics about the cluster. \n\n" +
				"This attribute should not be used for real-time monitoring! " +
				"Stats are updated every 10-15 minutes. To monitor real-time " +
				"metrics, monitor your cluster directly, via the Index Stats " +
				"API.",
			Computed: true,
			Attributes: map[string]dschema.Attribute{
				"docs": dschema.Int64Attribute{
					MarkdownDescription: "Number of documents in the index.",
					Computed:            true,
				},
				"shards_used": dschema.Int64Attribute{
					MarkdownDescription: "Number of shards the cluster is using.",
					Computed:            true,
				},
				"data_bytes_used": dschema.Int64Attribute{
					MarkdownDescription: "Number of bytes the cluster is using on-disk.",
					Computed:            true,
				},
			},
		},
		"access": dschema.SingleNestedAttribute{
			MarkdownDescription: "Access holds information about connecting to " +
				"the cluster.",
			Computed: true,
			Attributes: map[string]dschema.Attribute{
				"host": dschema.StringAttribute{
					MarkdownDescription: "Host name of the cluster.",
					Computed:            true,
				},
				"port": dschema.Int64Attribute{
					MarkdownDescription: "HTTP Port the cluster is running on.",
					Computed:            true,
				},
				"scheme": dschema.StringAttribute{
					MarkdownDescription: "HTTP Scheme needed to access the " +
						"cluster. Default: \"https\".",
					Computed: true,
				},
				"user": dschema.StringAttribute{
					MarkdownDescription: "User holds the username to access the " +
						"cluster with.\n\n " +
						"Only shown once, during cluster creation.",
					Computed:  true,
					Sensitive: true,
					Optional:  true,
				},
				"password": dschema.StringAttribute{
					MarkdownDescription: "Pass holds the password to access the " +
						"cluster with. \n\n" +
						"Only shown once, during cluster creation.",
					Computed:  true,
					Sensitive: true,
					Optional:  true,
				},
				"url": dschema.StringAttribute{
					MarkdownDescription: "URL is the Cluster endpoint for " +
						"access.\n\n" +
						"Only shown once, during cluster creation.",
					Computed: true,
				},
			},
		},
		"state": dschema.SingleNestedAttribute{
			MarkdownDescription: "State represents the current state of the " +
				"cluster. This indicates what the cluster is doing at " +
				"any given moment.",
			Computed: true,
			Attributes: map[string]dschema.Attribute{
				"state": dschema.StringAttribute{
					MarkdownDescription: "The state of the cluster.",
					Computed:            true,
				},
			},
		},
	}
}
