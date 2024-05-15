package plan

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

const (
	dataSourceMarkdownDescription     = "Plan represents a subscription plan."
	listDataSourceMarkdownDescription = "A list of all **available** plans on " +
		"your account."
)

// availableReleaseModel maps Plan Available Release details.
type availableReleaseModel struct {
	Slug types.String `tfsdk:"slug"`
}

func (m availableReleaseModel) ObjectType() attr.Type {
	return types.ObjectType{AttrTypes: m.AttrType()}
}

func (m availableReleaseModel) AttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"slug": types.StringType,
	}
}

func (m availableReleaseModel) Value() attr.Value {
	return types.ObjectValueMust(m.AttrType(), map[string]attr.Value{
		"slug": types.StringValue(m.Slug.ValueString()),
	})
}

// availableSpaceModel maps Plan Available Space details.
type availableSpaceModel struct {
	Path types.String `tfsdk:"path"`
}

func (m availableSpaceModel) ObjectType() attr.Type {
	return types.ObjectType{AttrTypes: m.AttrType()}
}

func (m availableSpaceModel) AttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"path": types.StringType,
	}
}

func (m availableSpaceModel) Value() attr.Value {
	return types.ObjectValueMust(m.AttrType(), map[string]attr.Value{
		"path": types.StringValue(m.Path.ValueString()),
	})
}

// model maps plans schema data.
type model struct {
	Name                    types.String `tfsdk:"name"`
	Slug                    types.String `tfsdk:"slug"`
	PriceInCents            types.Int64  `tfsdk:"price_in_cents"`
	BillingIntervalInMonths types.Int64  `tfsdk:"billing_interval_months"`
	SingleTenant            types.Bool   `tfsdk:"single_tenant"`
	PrivateNetwork          types.Bool   `tfsdk:"private_network"`
	AvailableReleases       types.List   `tfsdk:"available_releases"`
	AvailableSpaces         types.List   `tfsdk:"available_spaces"`
}

func convert(ctx context.Context, r bonsai.Plan) (model, error) {
	m := model{
		Name:                    types.StringValue(r.Name),
		Slug:                    types.StringValue(r.Slug),
		PriceInCents:            types.Int64Value(r.PriceInCents),
		BillingIntervalInMonths: types.Int64Value(int64(r.BillingIntervalInMonths)),
		AvailableReleases:       types.ListNull(availableReleaseModel{}.ObjectType()),
		AvailableSpaces:         types.ListNull(availableSpaceModel{}.ObjectType()),
	}

	if r.PrivateNetwork != nil {
		m.PrivateNetwork = types.BoolValue(*r.PrivateNetwork)
	}

	if r.SingleTenant != nil {
		m.SingleTenant = types.BoolValue(*r.SingleTenant)
	}

	if r.AvailableReleases != nil && len(r.AvailableReleases) > 0 {
		availableReleases := make([]availableReleaseModel, len(r.AvailableReleases))
		for i, release := range r.AvailableReleases {
			availableReleases[i] = availableReleaseModel{Slug: types.StringValue(release.Slug)}
		}

		tfAvailableReleases, diags := types.ListValueFrom(ctx, availableReleaseModel{}.ObjectType(), availableReleases)
		if diags.HasError() {
			return m, fmt.Errorf("failed reading available releases: %s - %s", diags[0].Summary(), diags[0].Detail())
		}
		m.AvailableReleases = tfAvailableReleases
	}

	if r.AvailableSpaces != nil && len(r.AvailableSpaces) > 0 {
		availableSpaces := make([]availableSpaceModel, len(r.AvailableSpaces))
		for i, space := range r.AvailableSpaces {
			availableSpaces[i] = availableSpaceModel{Path: types.StringValue(space.Path)}
		}

		tfAvailableSpaces, diags := types.ListValueFrom(ctx, availableSpaceModel{}.ObjectType(), availableSpaces)
		if diags.HasError() {
			return m, fmt.Errorf("failed reading available spaces: %s - %s", diags[0].Summary(), diags[0].Detail())
		}
		m.AvailableSpaces = tfAvailableSpaces
	}

	return m, nil
}

func schemaAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"name": dschema.StringAttribute{
			MarkdownDescription: "The human-readable name of the plan.",
			Computed:            true,
		},
		"slug": dschema.StringAttribute{
			MarkdownDescription: "The machine-readable name for the plan.",
			Computed:            true,
			Optional:            true,
		},
		"price_in_cents": dschema.Int64Attribute{
			MarkdownDescription: "Represents the plan price in cents.",
			Computed:            true,
		},
		"billing_interval_months": dschema.Int64Attribute{
			MarkdownDescription: "The plan billing interval in months.",
			Computed:            true,
		},
		"single_tenant": dschema.BoolAttribute{
			MarkdownDescription: "Indicates whether the plan is single-tenant or not. " +
				"A value of false indicates the Cluster will share hardware " +
				"with other Clusters. Single tenant environments can be reached" +
				" via the public Internet.",
			Computed: true,
		},
		"private_network": dschema.BoolAttribute{
			MarkdownDescription: "Indicates whether the plan is on a publicly " +
				"addressable network. Private plans provide environments that " +
				"cannot be reached by the public Internet. A VPC connection " +
				"will be needed to communicate with a private cluster.",
			Computed: true,
		},

		"available_releases": dschema.ListNestedAttribute{
			MarkdownDescription: "A collection of search release slugs " +
				"available for the plan.",
			Computed: true,
			NestedObject: dschema.NestedAttributeObject{
				Attributes: map[string]dschema.Attribute{
					"slug": dschema.StringAttribute{
						MarkdownDescription: "A machine-readable name for the release.",
						Computed:            true,
					},
				},
			},
		},
		"available_spaces": dschema.ListNestedAttribute{
			MarkdownDescription: "A collection of Space paths available for the " +
				"plan.",
			Computed: true,
			NestedObject: dschema.NestedAttributeObject{
				Attributes: map[string]dschema.Attribute{
					"path": dschema.StringAttribute{
						MarkdownDescription: "A machine-readable name for the " +
							"server group.",
						Computed: true,
					},
				},
			},
		},
	}
}
