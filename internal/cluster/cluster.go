package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
