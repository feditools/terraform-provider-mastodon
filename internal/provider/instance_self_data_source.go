package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = instanceSelfDataSourceType{}
var _ datasource.DataSource = instanceSelfDataSource{}

type instanceSelfDataSourceType struct{}

func (t instanceSelfDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Instance self",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"email": {
				MarkdownDescription: "Instance Contact Email",
				Type:                types.StringType,
				Computed:            true,
			},
			"thumbnail": {
				MarkdownDescription: "Instance Thumbnail",
				Type:                types.StringType,
				Computed:            true,
			},
			"title": {
				MarkdownDescription: "Instance Title",
				Type:                types.StringType,
				Computed:            true,
			},
			"uri": {
				MarkdownDescription: "Instance URI",
				Type:                types.StringType,
				Computed:            true,
			},
			"version": {
				MarkdownDescription: "Instance Version",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t instanceSelfDataSourceType) NewDataSource(_ context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	prov, diags := convertProviderType(in)

	return instanceSelfDataSource{
		provider: prov,
	}, diags
}

type instanceSelfDataSourceData struct {
	ID types.String `tfsdk:"id"`

	Email     types.String `tfsdk:"email"`
	Thumbnail types.String `tfsdk:"thumbnail"`
	Title     types.String `tfsdk:"title"`
	URI       types.String `tfsdk:"uri"`
	Version   types.String `tfsdk:"version"`
}

type instanceSelfDataSource struct {
	provider scaffoldingProvider
}

func (d instanceSelfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data instanceSelfDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := d.provider.client.GetInstance(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))

		return
	}

	data.ID = types.String{Value: instance.URI}

	data.Email = types.String{Value: instance.EMail}
	data.Thumbnail = types.String{Value: instance.Thumbnail}
	data.Title = types.String{Value: instance.Title}
	data.URI = types.String{Value: instance.URI}
	data.Version = types.String{Value: instance.Version}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
