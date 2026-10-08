package cluster

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	tfrsc "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/omc/bonsai-api-go/v2/bonsai"
	"github.com/stretchr/testify/require"
)

// TestResource_DeleteSendsClusterName checks that Delete confirms the
// destroy with the cluster's exact name in the JSON body, which the API
// requires.
func TestResource_DeleteSendsClusterName(t *testing.T) {
	const (
		slug = "my-cluster-1234567890"
		name = "My Cluster "
	)

	var deleteBodies []string
	mux := http.NewServeMux()
	mux.HandleFunc("/clusters/"+slug, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodDelete:
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			deleteBodies = append(deleteBodies, string(body))

			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"message": "Your cluster is being deprovisioned.",
				"monitor": "https://api.bonsai.io/clusters/" + slug,
				"status":  http.StatusAccepted,
			})
		case http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"errors": []string{"Cluster not found."},
				"status": http.StatusNotFound,
			})
		default:
			t.Errorf("unexpected %s request", r.Method)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := bonsai.NewClient(bonsai.WithEndpoint(server.URL))
	r := &resource{client: &client.Cluster}

	ctx := context.Background()
	schemaResp := &tfrsc.SchemaResponse{}
	r.Schema(ctx, tfrsc.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), "schema: %v", schemaResp.Diagnostics)

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	diags := state.Set(ctx, &resourceModel{
		ID:     types.StringValue(slug),
		Name:   types.StringValue(name),
		Slug:   types.StringValue(slug),
		Stats:  types.ObjectNull(statsModelTypes),
		Access: types.ObjectNull(accessModelTypes),
		State:  types.ObjectNull(stateModelTypes),
	})
	require.False(t, diags.HasError(), "state: %v", diags)

	resp := &tfrsc.DeleteResponse{State: state}
	r.Delete(ctx, tfrsc.DeleteRequest{State: state}, resp)

	require.False(t, resp.Diagnostics.HasError(), "delete: %v", resp.Diagnostics)
	require.Len(t, deleteBodies, 1)
	require.JSONEq(t, `{"name": "My Cluster "}`, deleteBodies[0], "sends the exact cluster name")
}
