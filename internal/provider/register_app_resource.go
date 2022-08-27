package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mattn/go-mastodon"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = registerAppResourceType{}
var _ resource.Resource = registerAppResource{}

//var _ resource.ResourceWithImportState = registerAppResource{}

type registerAppResourceType struct{}

func (t registerAppResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Register Application",

		Attributes: map[string]tfsdk.Attribute{
			// inputs
			"client_name": {
				MarkdownDescription: "Name to register application with",
				Optional:            true,
				Type:                types.StringType,
			},
			"redirect_uris": {
				MarkdownDescription: "Redirect URI to register application with",
				Optional:            true,
				Type:                types.StringType,
			},
			"scopes": {
				MarkdownDescription: "OAuth scopes",
				Optional:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"website": {
				MarkdownDescription: "Website for registered application",
				Optional:            true,
				Type:                types.StringType,
			},

			// outputs
			"id": {
				MarkdownDescription: "identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"app_config": {
				MarkdownDescription: "Application auth config",
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"client_id":     types.StringType,
						"client_secret": types.StringType,
						"redirect_uri":  types.StringType,
					},
				},
				Computed: true,
			},
		},
	}, nil
}

func (t registerAppResourceType) NewResource(_ context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	prov, diags := convertProviderType(in)

	return registerAppResource{
		provider: prov,
	}, diags
}

type registerAppResourceData struct {
	ClientName   types.String `tfsdk:"client_name"`
	RedirectURIs types.String `tfsdk:"redirect_uris"`
	Scopes       types.List   `tfsdk:"scopes"`
	Website     types.String `tfsdk:"website"`

	ID        types.String `tfsdk:"id"`
	AppConfig types.Object `tfsdk:"app_config"`
}

type registerAppResource struct {
	provider mastodonProvider
}

func (r registerAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data registerAppResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// collect inputs
	clientName := "terraform-provider-mastodon"
	if !data.ClientName.IsNull() {
		clientName = data.ClientName.Value
	}

	redirectURIs := "urn:ietf:wg:oauth:2.0:oob"
	if !data.RedirectURIs.IsNull() {
		redirectURIs = data.RedirectURIs.Value
	}

	scopes := "read write follow admin:read admin:write"
	if !data.Scopes.IsNull() {
		var sb strings.Builder
		lastScope := len(data.Scopes.Elems) - 1
		for i, scope := range data.Scopes.Elems {
			scopeString := scope.(types.String)
			sb.WriteString(scopeString.Value)
			if i != lastScope {
				sb.WriteString(" ")
			}
		}

		scopes = sb.String()
	}

	website := "https://github.com/feditools/terraform-provider-mastodon"
	if !data.Website.IsNull() {
		website = data.Website.Value
	}

	// do registration
	app, err := mastodon.RegisterApp(ctx, &mastodon.AppConfig{
		Server:       r.provider.server(),
		ClientName:   clientName,
		Scopes:       scopes,
		Website:      website,
		RedirectURIs: redirectURIs,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to register application, got error: %s", err.Error()))

		return
	}

	data.ID = types.String{Value: string(app.ID)}
	data.AppConfig = types.Object{
		AttrTypes: map[string]attr.Type{
			"client_id":     types.StringType,
			"client_secret": types.StringType,
			"redirect_uri":  types.StringType,
		},
		Attrs: map[string]attr.Value{
			"client_id":     types.String{Value: app.ClientID},
			"client_secret": types.String{Value: app.ClientSecret},
			"redirect_uri":  types.String{Value: app.RedirectURI},
		},
	}

	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r registerAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data registerAppResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientID := data.AppConfig.Attrs["client_id"].(types.String).Value
	clientSecret := data.AppConfig.Attrs["client_secret"].(types.String).Value

	client, err := r.provider.newAuthenticatedClient(ctx, clientID, clientSecret, "")
	if err != nil {
		if strings.HasPrefix("bad authorization: 401 Unauthorized:", err.Error()) {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create new client, got error: %s", err))

		return
	}


	_, err = client.VerifyAppCredentials(ctx)
	if err != nil {
		if strings.HasPrefix("bad request: 401 Unauthorized:", err.Error()) {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))

		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r registerAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data registerAppResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// collect inputs
	clientName := "terraform-provider-mastodon"
	if !data.ClientName.IsNull() {
		clientName = data.ClientName.Value
	}

	redirectURIs := "urn:ietf:wg:oauth:2.0:oob"
	if !data.RedirectURIs.IsNull() {
		redirectURIs = data.RedirectURIs.Value
	}

	scopes := "read write follow admin:read admin:write"
	if !data.Scopes.IsNull() {
		var sb strings.Builder
		lastScope := len(data.Scopes.Elems) - 1
		for i, scope := range data.Scopes.Elems {
			scopeString := scope.(types.String)
			sb.WriteString(scopeString.Value)
			if i != lastScope {
				sb.WriteString(" ")
			}
		}

		scopes = sb.String()
	}

	server := "https://" + r.provider.domain

	website := "https://github.com/feditools/terraform-provider-mastodon"
	if !data.Website.IsNull() {
		website = data.Website.Value
	}

	// do registration
	app, err := mastodon.RegisterApp(ctx, &mastodon.AppConfig{
		Server:       server,
		ClientName:   clientName,
		Scopes:       scopes,
		Website:      website,
		RedirectURIs: redirectURIs,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to register application, got error: %s", err.Error()))

		return
	}

	data.ID = types.String{Value: string(app.ID)}
	data.AppConfig = types.Object{
		AttrTypes: map[string]attr.Type{
			"client_id":     types.StringType,
			"client_secret": types.StringType,
			"redirect_uri":  types.StringType,
		},
		Attrs: map[string]attr.Value{
			"client_id":     types.String{Value: app.ClientID},
			"client_secret": types.String{Value: app.ClientSecret},
			"redirect_uri":  types.String{Value: app.RedirectURI},
		},
	}

	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r registerAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data registerAppResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r registerAppResource) doRegisterApp(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

/*func (r registerAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}*/
