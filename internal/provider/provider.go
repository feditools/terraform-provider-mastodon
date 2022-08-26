package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mattn/go-mastodon"
)

var _ provider.Provider = &scaffoldingProvider{}

type scaffoldingProvider struct {
	client *mastodon.Client

	configured bool
	version    string
}

type providerData struct {
	Domain types.String `tfsdk:"domain"`
}

func (p *scaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	p.client = mastodon.NewClient(&mastodon.Config{
		Server: "https://" + data.Domain.Value,
	})

	p.configured = true
}

func (p *scaffoldingProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		//"scaffolding_example": instanceSelfResourceType{},
	}, nil
}

func (p *scaffoldingProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"mastodon_instance_self": instanceSelfDataSourceType{},
	}, nil
}

func (p *scaffoldingProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domain": {
				MarkdownDescription: "Domain",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &scaffoldingProvider{
			version: version,
		}
	}
}

func convertProviderType(in provider.Provider) (scaffoldingProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*scaffoldingProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return scaffoldingProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return scaffoldingProvider{}, diags
	}

	return *p, diags
}
