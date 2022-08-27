package provider

import (
	"context"
	"fmt"
	"github.com/mattn/go-mastodon"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = accountDataSourceType{}
var _ datasource.DataSource = accountDataSource{}

type accountDataSourceType struct{}

func (t accountDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Get account info",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "identifier",
				Required:            true,
				Type:                types.StringType,
			},

			"username": {
				MarkdownDescription: "Instance Contact Email",
				Type:                types.StringType,
				Computed:            true,
			},
			"account": {
				MarkdownDescription: "Instance Thumbnail",
				Type:                types.StringType,
				Computed:            true,
			},
			"display_name": {
				MarkdownDescription: "Instance Title",
				Type:                types.StringType,
				Computed:            true,
			},
			"created_at": {
				MarkdownDescription: "Instance URI",
				Type:                types.StringType,
				Computed:            true,
			},
			"url": {
				MarkdownDescription: "Instance Version",
				Type:                types.StringType,
				Computed:            true,
			},
			"discoverable": {
				MarkdownDescription: "Instance Version",
				Type:                types.BoolType,
				Computed:            true,
			},
		},
	}, nil
}

func (t accountDataSourceType) NewDataSource(_ context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	prov, diags := convertProviderType(in)

	return accountDataSource{
		provider: prov,
	}, diags
}

type accountDataSourceData struct {
	ID types.String `tfsdk:"id"`

	Username     types.String `tfsdk:"username"`
	Acct         types.String `tfsdk:"account"`
	DisplayName  types.String `tfsdk:"display_name"`
	CreatedAt    types.String `tfsdk:"created_at"`
	URL          types.String `tfsdk:"url"`
	Discoverable types.Bool   `tfsdk:"discoverable"`
}

type accountDataSource struct {
	provider mastodonProvider
}

func (d accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	account, err := d.provider.newUnauthenticatedClient().GetAccount(ctx, mastodon.ID(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account, got error: %s", err))

		return
	}

	data.Username = types.String{Value: account.Username}
	data.Acct = types.String{Value: account.Acct}
	data.DisplayName = types.String{Value: account.DisplayName}
	data.CreatedAt = types.String{Value: account.CreatedAt.String()}
	data.URL = types.String{Value: account.URL}
	data.Discoverable = types.Bool{Value: account.Discoverable}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
