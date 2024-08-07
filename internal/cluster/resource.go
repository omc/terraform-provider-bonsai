package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	tfrsc "github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/omc/bonsai-api-go/v2/bonsai"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ tfrsc.Resource              = &resource{}
	_ tfrsc.ResourceWithConfigure = &resource{}

	// Unavailable Regexp matches fields which are returned as not available
	// during cluster provisioning.
	unavailableRegexp = regexp.MustCompile(`.*?not available.*`)
	// updateRequestProcessingRegexp matches a message response which
	// indicates that the Cluster update request has been successfully
	// received, and is currently processing.
	updateRequestProcessingRegexp = regexp.MustCompile(`Your cluster is being updated`)
)

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

// dataSource is the data source implementation.
type resource struct {
	client *bonsai.ClusterClient
}

// NewResource is a helper function to simplify the provider implementation.
func NewResource() tfrsc.Resource {
	return &resource{}
}

// Metadata returns the resource type name.
func (r *resource) Metadata(_ context.Context, req tfrsc.MetadataRequest, resp *tfrsc.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *resource) Configure(_ context.Context, req tfrsc.ConfigureRequest, resp *tfrsc.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*bonsai.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource Configure Type",
			fmt.Sprintf("Expected *bonsai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &client.Cluster
}

func resourceSchemaAttributes() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": rschema.StringAttribute{
			MarkdownDescription: "The human-readable name of the cluster.",
			Required:            true,
		},
		"slug": rschema.StringAttribute{
			MarkdownDescription: "A unique, machine-readable name for the " +
				"cluster. A cluster slug is based its name at creation, to " +
				"which a random integer is concatenated.",
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"uri": rschema.StringAttribute{
			MarkdownDescription: "A URI to retrieve more information about this cluster.",
			Computed:            true,
		},
		"message": rschema.StringAttribute{
			MarkdownDescription: "Message received during Cluster creation",
			Computed:            true,
		},
		"monitor": rschema.StringAttribute{
			MarkdownDescription: "Monitor received during Cluster creation",
			Computed:            true,
		},
		"plan": rschema.SingleNestedAttribute{
			MarkdownDescription: "Plan holds some information about the cluster's current subscription plan.",
			Optional:            true,
			Attributes: map[string]rschema.Attribute{
				"slug": rschema.StringAttribute{
					MarkdownDescription: "A machine-readable name for the plan.",
					Required:            true,
				},
				"uri": rschema.StringAttribute{
					MarkdownDescription: "A URI to retrieve more information about this Plan.",
					Computed:            true,
				},
			},
		},
		"release": rschema.SingleNestedAttribute{
			MarkdownDescription: "Release holds some information about the cluster's current release.",
			Optional:            true,
			Attributes: map[string]rschema.Attribute{
				"service_type": rschema.StringAttribute{
					MarkdownDescription: "The service type of the deployment - for example, \"elasticsearch\".",
					Computed:            true,
				},
				"package_name": rschema.StringAttribute{
					MarkdownDescription: "PackageName is the package name of the release.",
					Computed:            true,
				},
				"version": rschema.StringAttribute{
					MarkdownDescription: "The version of the release.",
					Computed:            true,
				},
				"slug": rschema.StringAttribute{
					MarkdownDescription: "The machine-readable name for the deployment.",
					Optional:            true,
				},
				"uri": rschema.StringAttribute{
					MarkdownDescription: "A URI to retrieve more information about this Release.",
					Computed:            true,
				},
			},
		},
		"space": rschema.SingleNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Space holds some information about where the cluster is running.",
			Attributes: map[string]rschema.Attribute{
				"path": rschema.StringAttribute{
					MarkdownDescription: "A machine-readable name for the server group.",
					Optional:            true,
				},
				"region": rschema.StringAttribute{
					MarkdownDescription: "The geographic region in which the cluster is running.",
					Computed:            true,
				},
				"uri": rschema.StringAttribute{
					MarkdownDescription: "A URI to retrieve more information about this Space.",
					Computed:            true,
				},
			},
		},
		"stats": rschema.SingleNestedAttribute{
			MarkdownDescription: "Stats holds *some* statistics about the cluster. \n\n" +
				"This attribute should not be used for real-time monitoring! " +
				"Stats are updated every 10-15 minutes. To monitor real-time " +
				"metrics, monitor your cluster directly, via the Index Stats " +
				"API.",
			Computed: true,
			Attributes: map[string]rschema.Attribute{
				"docs": rschema.Int64Attribute{
					MarkdownDescription: "Number of documents in the index.",
					Computed:            true,
				},
				"shards_used": rschema.Int64Attribute{
					MarkdownDescription: "Number of shards the cluster is using.",
					Computed:            true,
				},
				"data_bytes_used": rschema.Int64Attribute{
					MarkdownDescription: "Number of bytes the cluster is using on-disk.",
					Computed:            true,
				},
			},
		},
		"access": rschema.SingleNestedAttribute{
			MarkdownDescription: "Access holds information about connecting to " +
				"the cluster.",
			Computed: true,
			Attributes: map[string]rschema.Attribute{
				"host": rschema.StringAttribute{
					MarkdownDescription: "Host name of the cluster.",
					Computed:            true,
				},
				"port": rschema.Int64Attribute{
					MarkdownDescription: "HTTP Port the cluster is running on.",
					Computed:            true,
				},
				"scheme": rschema.StringAttribute{
					MarkdownDescription: "HTTP Scheme needed to access the " +
						"cluster. Default: \"https\".",
					Computed: true,
				},
				"user": rschema.StringAttribute{
					MarkdownDescription: "User holds the username to access the " +
						"cluster with.\n\n " +
						"Only shown once, during cluster creation.",
					Computed:  true,
					Sensitive: true,
					Optional:  true,
				},
				"password": rschema.StringAttribute{
					MarkdownDescription: "Pass holds the password to access the " +
						"cluster with. \n\n" +
						"Only shown once, during cluster creation.",
					Computed:  true,
					Sensitive: true,
					Optional:  true,
				},
				"url": rschema.StringAttribute{
					MarkdownDescription: "URL is the Cluster endpoint for " +
						"access.\n\n" +
						"Only shown once, during cluster creation.",
					Sensitive: true,
					Computed:  true,
					Optional:  true,
				},
			},
		},
		"state": rschema.SingleNestedAttribute{
			MarkdownDescription: "State represents the current state of the " +
				"cluster. This indicates what the cluster is doing at " +
				"any given moment.",
			Computed: true,
			Attributes: map[string]rschema.Attribute{
				"state": rschema.StringAttribute{
					MarkdownDescription: "The state of the cluster.",
					Computed:            true,
				},
			},
		},
	}
}

func convertResourceClusterToCreateRequest(r resourceModel) bonsai.ClusterCreateOpts {
	return bonsai.ClusterCreateOpts{
		Name:    r.Name.ValueString(),
		Plan:    r.Plan.Slug.ValueString(),
		Space:   r.Space.Path.ValueString(),
		Release: r.Release.Slug.ValueString(),
	}
}

var statsModelTypes = map[string]attr.Type{
	"docs":            types.Int64Type,
	"shards_used":     types.Int64Type,
	"data_bytes_used": types.Int64Type,
}

var accessModelTypes = map[string]attr.Type{
	"host":     types.StringType,
	"port":     types.Int64Type,
	"scheme":   types.StringType,
	"user":     types.StringType,
	"password": types.StringType,
	"url":      types.StringType,
}

var stateModelTypes = map[string]attr.Type{
	"state": types.StringType,
}

func convertCreateResponseModelToResourceCluster(c bonsai.ClustersResultCreate) (resourceModel, error) {
	access, diags := types.ObjectValueFrom(context.TODO(), accessModelTypes, &accessModel{
		Host:     types.StringValue(c.Access.Host),
		Port:     types.Int64Value(int64(c.Access.Port)),
		Scheme:   types.StringValue(c.Access.Scheme),
		Username: types.StringValue(c.Access.Username),
		Password: types.StringValue(c.Access.Password),
		URL:      types.StringValue(c.Access.URL),
	})
	if diags.HasError() {
		return resourceModel{},
			fmt.Errorf("error reading cluster access: %s - %s", diags[0].Summary(), diags[0].Detail())
	}

	m := resourceModel{
		// The Host is the same as the Slug.
		// TODO: Confirm that this is acceptable and never truncated.
		Slug:    types.StringValue(c.Access.Host),
		Message: types.StringValue(c.Message),
		Monitor: types.StringValue(c.Monitor),

		Access: access,
		Stats:  types.ObjectNull(statsModelTypes),
		State:  types.ObjectNull(stateModelTypes),
	}
	return m, nil
}

func resourceConvert(c bonsai.Cluster) (resourceModel, error) {
	access, diags := types.ObjectValueFrom(context.TODO(), accessModelTypes, &accessModel{
		Host:     types.StringValue(c.Access.Host),
		Port:     types.Int64Value(int64(c.Access.Port)),
		Scheme:   types.StringValue(c.Access.Scheme),
		Username: types.StringValue(c.Access.Username),
		Password: types.StringValue(c.Access.Password),
		URL:      types.StringValue(c.Access.URL),
	})
	if diags.HasError() {
		return resourceModel{},
			fmt.Errorf("error reading cluster access: %s - %s", diags[0].Summary(), diags[0].Detail())
	}

	stats, diags := types.ObjectValueFrom(context.TODO(), statsModelTypes, &statsModel{
		Docs:          types.Int64Value(c.Stats.Docs),
		ShardsUsed:    types.Int64Value(c.Stats.ShardsUsed),
		DataBytesUsed: types.Int64Value(c.Stats.DataBytesUsed),
	})
	if diags.HasError() {
		return resourceModel{},
			fmt.Errorf("error reading cluster access: %s - %s", diags[0].Summary(), diags[0].Detail())
	}

	state, diags := types.ObjectValueFrom(context.TODO(), stateModelTypes, &stateModel{
		State: types.StringValue(string(c.State)),
	})
	if diags.HasError() {
		return resourceModel{},
			fmt.Errorf("error reading cluster access: %s - %s", diags[0].Summary(), diags[0].Detail())
	}

	m := resourceModel{
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
		Stats:  stats,
		Access: access,
		State:  state,
	}
	return m, nil
}

// Schema returns the schema information for a cluster create request resource.
func (r *resource) Schema(_ context.Context, _ tfrsc.SchemaRequest, resp *tfrsc.SchemaResponse) {
	resp.Schema = rschema.Schema{
		MarkdownDescription: resourceMarkdownDescription,
		Attributes:          resourceSchemaAttributes(),
	}
}

// Create requests a new Cluster to be created.
func (r *resource) Create(ctx context.Context, req tfrsc.CreateRequest, resp *tfrsc.CreateResponse) {
	var (
		state, createResultState, refreshState resourceModel
		refreshResult                          bonsai.Cluster
	)

	refreshDeadline := 5 * time.Minute
	refreshDelay := 10 * time.Second

	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := convertResourceClusterToCreateRequest(state)
	tflog.Debug(ctx, fmt.Sprintf("create request: %+v", createRequest))

	createResult, err := r.client.Cluster.Create(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Create Bonsai Cluster (%v) from the Bonsai API", state),
			err.Error(),
		)
		tflog.Debug(ctx, fmt.Sprintf("returning error %s", err.Error()))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("received cluster %+v", createResult))

	// Handles null nested objects as well
	createResultState, err = convertCreateResponseModelToResourceCluster(createResult)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to convert response to resource Cluster (%v) from the Bonsai API", createResult),
			err.Error(),
		)
		return
	}

	// Now, we need to ensure the planned state remains.
	// For example, we don't receive back (from the API) the Plan.Slug which was selected
	// by the user, but we need to save it in the state.
	//
	// See convertResourceClusterToCreateRequest for details on
	// what we expect.
	//
	// Note: createResultState.Slug has been set via Access.Host
	createResultState.Name = state.Name
	createResultState.Plan.Slug = state.Plan.Slug
	createResultState.Space.Path = state.Space.Path
	createResultState.Release.Slug = state.Release.Slug
	// And, set the unique identifier
	createResultState.ID = createResultState.Slug

	refreshCtx, refreshCancel := context.WithDeadline(ctx, time.Now().Add(refreshDeadline))
	defer refreshCancel()

DiscoveryLoop:
	for {
		select {
		case <-refreshCtx.Done():
			// On the event of time-out, set the state we *do* know.
			diags = resp.State.Set(ctx, createResultState)
			tflog.Debug(ctx, fmt.Sprintf("context done - set cluster state %+v", createResultState))
			resp.Diagnostics.Append(diags...)
			tflog.Debug(ctx, fmt.Sprintf("context done - appended diags %+v", createResultState))

			resp.Diagnostics.AddError(
				fmt.Sprintf(
					"Timed out while awaiting Bonsai Cluster (%s) provision, in time (%d)",
					state.Slug.ValueString(),
					refreshDeadline,
				),
				ctx.Err().Error(),
			)
			return
		default:
			refreshResult, err = r.client.Cluster.GetBySlug(ctx, createResultState.Slug.ValueString())
			// If we encounter an error, and it's not just that the cluster
			// hasn't been created yet...
			if err != nil {
				if errors.Is(err, bonsai.ErrHTTPStatusNotFound) {
					// Sleep for a little bit; all of these should be refactored at some point
					time.Sleep(refreshDelay)
					continue DiscoveryLoop
				}

				resp.Diagnostics.AddError(
					"Error while fetching new Bonsai Cluster state",
					fmt.Sprintf(
						"Failed while refreshing Cluster (%s) state after create, unexpected error: %s",
						createResultState.Slug.ValueString(),
						err,
					),
				)
				return
			}

			if unavailableRegexp.MatchString(refreshResult.Space.Path) || unavailableRegexp.MatchString(refreshResult.Space.URI) {
				continue DiscoveryLoop
			}

			// We received a good result
			break DiscoveryLoop
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("received refreshed cluster: %+v", refreshResult))
	tflog.Debug(ctx, fmt.Sprintf("created cluster %+v", createResultState))

	refreshState, err = resourceConvert(refreshResult)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"failed to convert refresh response (%v) to resource Cluster from the Bonsai API",
				refreshResult,
			),
			err.Error(),
		)
		return
	}
	// Add credentials back to the refreshed state
	refreshAccessModel := &accessModel{}
	createAccessModel := &accessModel{}

	tflog.Debug(ctx, "converting access as accessModel")
	refreshState.Access.As(context.Background(), refreshAccessModel, basetypes.ObjectAsOptions{})
	createResultState.Access.As(context.Background(), createAccessModel, basetypes.ObjectAsOptions{})

	refreshAccessModel.Username = createAccessModel.Username
	refreshAccessModel.Password = createAccessModel.Password

	// Now, update the Refresh State's Access URL, with the full URI
	accessUrl := url.URL{}
	accessUrl.Host = refreshAccessModel.Host.ValueString()
	accessUrl.Scheme = createAccessModel.Scheme.ValueString()
	accessUrl.User = url.UserPassword(createAccessModel.Username.ValueString(), createAccessModel.Password.ValueString())
	refreshAccessModel.URL = types.StringValue(accessUrl.String())

	// Finally, we're ready to convert access back into an ObjectValue and update the refreshState object
	tflog.Debug(ctx, "converting access back to ObjectValue")
	access, diags := types.ObjectValueFrom(context.TODO(), accessModelTypes, &refreshAccessModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, "Setting Access to access")
	refreshState.Access = access

	// And, set the unique identifier
	refreshState.ID = createResultState.ID

	diags = resp.State.Set(ctx, refreshState)

	tflog.Debug(ctx, fmt.Sprintf("set cluster state %+v", createResultState))
	resp.Diagnostics.Append(diags...)
	tflog.Debug(ctx, "appended diags after cluster state set")
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read is called when the provider must read resource values in order
// to update state.
func (r *resource) Read(ctx context.Context, req tfrsc.ReadRequest, resp *tfrsc.ReadResponse) {
	var (
		state, apiState resourceModel
	)

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Slug.IsNull() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read Bonsai Cluster (%s) from the Bonsai API", state.Slug.ValueString()),
			"expected 'slug' option to be set",
		)
		return
	}

	apiResp, err := r.client.Cluster.GetBySlug(ctx, state.ID.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("received cluster %v", apiResp))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read Bonsai Cluster (%s) from the Bonsai API", state.Slug.ValueString()),
			err.Error(),
		)
		return
	}

	apiState, err = resourceConvert(apiResp)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to convert response to resource Cluster (%+v)", apiResp),
			err.Error(),
		)
		return
	}

	// Set state details
	apiState.ID = state.ID
	apiState.Access = state.Access

	tflog.Debug(ctx, fmt.Sprintf("read state %v", apiState))

	diags = resp.State.Set(ctx, apiState)
	tflog.Debug(ctx, fmt.Sprintf("setting Read api state %v", apiState))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "returning from read")
}

// Update updates the Alias state.
func (r *resource) Update(ctx context.Context, req tfrsc.UpdateRequest, resp *tfrsc.UpdateResponse) {
	var (
		refreshResult                bonsai.Cluster
		desired, state, refreshState resourceModel
		err                          error
	)

	refreshDeadline := 5 * time.Minute
	refreshDelay := 10 * time.Second

	diags := req.Plan.Get(ctx, &desired)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := bonsai.ClusterUpdateOpts{
		Name: desired.Name.ValueString(),
		Plan: desired.Plan.Slug.ValueString(),
	}

	updateResp, err := r.client.Cluster.Update(ctx, state.Slug.ValueString(), updateOpts)
	tflog.Debug(ctx, fmt.Sprintf("received cluster update response: %v", updateResp))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"update: failed to update Cluster (%s) with desired state (%v)",
				state.Slug.ValueString(),
				updateOpts,
			),
			err.Error(),
		)
		return
	}

	if !updateRequestProcessingRegexp.MatchString(updateResp.Message) {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"update: failed to update Cluster (%s) with desired state"+
					" (%v) - update request message didn't indicate success.",
				state.Slug.ValueString(),
				updateOpts,
			),
			updateResp.Message,
		)
		return
	}

	refreshCtx, refreshCancel := context.WithDeadline(ctx, time.Now().Add(refreshDeadline))
	defer refreshCancel()

UpdateLoop:
	for {
		select {
		case <-refreshCtx.Done():
			resp.Diagnostics.AddError(
				fmt.Sprintf(
					"Update timed out while updating Bonsai Cluster (%s) with desired state (%v), in time (%d)",
					state.Slug.ValueString(),
					updateOpts,
					refreshDeadline,
				),
				ctx.Err().Error(),
			)
			return
		default:
			// Sleep for a little bit; all of these should be refactored at some point
			time.Sleep(refreshDelay)
			refreshResult, err = r.client.Cluster.GetBySlug(ctx, state.Slug.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Update: Error refreshing Bonsai Cluster state",
					fmt.Sprintf(
						"Failed while refreshing Cluster (%s) state after update, unexpected error: %s",
						state.Slug.ValueString(),
						err,
					),
				)
				return
			}

			// If we're not updating the plan, and we've got the desired name,
			// we're done!
			if refreshResult.State != bonsai.ClusterStateUpdatingPlan &&
				refreshResult.Name == desired.Name.ValueString() {
				break UpdateLoop
			}

			continue UpdateLoop
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("update: received refreshed cluster: %+v", refreshResult))

	refreshState, err = resourceConvert(refreshResult)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("update: failed to convert response to resource Cluster (%+v)", refreshResult),
			err.Error(),
		)
		return
	}

	// Set state details
	refreshState.ID = state.ID
	refreshState.Access = state.Access

	diags = resp.State.Set(ctx, refreshState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes an Alias.
func (r *resource) Delete(ctx context.Context, req tfrsc.DeleteRequest, resp *tfrsc.DeleteResponse) {
	var state resourceModel

	refreshDeadline := 5 * time.Minute
	refreshDelay := 10 * time.Second

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Cluster.Destroy(ctx, state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"failed to destroy cluster (%s)",
				state.Slug.ValueString(),
			),
			err.Error(),
		)
		return
	}

	// wait until it's been deleted
	refreshCtx, refreshCancel := context.WithDeadline(ctx, time.Now().Add(refreshDeadline))
	defer refreshCancel()

RefreshLoop:
	for {
		select {
		case <-refreshCtx.Done():
			resp.Diagnostics.AddError(
				fmt.Sprintf(
					"Timed out while deleting Bonsai Cluster (%s) in time (%d)",
					state.Slug.ValueString(),
					refreshDeadline,
				),
				ctx.Err().Error(),
			)
			return
		default:
			result, err := r.client.Cluster.GetBySlug(ctx, state.Slug.ValueString())
			tflog.Debug(ctx, fmt.Sprintf("found cluster: %+v", result))
			if err != nil {
				if errors.Is(err, bonsai.ErrHTTPStatusNotFound) {
					break RefreshLoop
				}
				resp.Diagnostics.AddError(
					"Error refreshing Bonsai Cluster state",
					fmt.Sprintf(
						"Failed while refreshing Cluster (%s) state after destroy request, unexpected error: %s",
						state.Slug.ValueString(),
						err,
					),
				)
				return
			}

			// Deprovisioned, but still exists
			if result.State == bonsai.ClusterStateDeprovisioned {
				return
			}
			// Sleep for a little bit; all of these should be refactored at some point
			time.Sleep(refreshDelay)
			continue RefreshLoop
		}
	}
}
