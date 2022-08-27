package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"sync"
)

var _ provider.Provider = &mastodonProvider{}

type mastodonProvider struct {
	accessToken     string
	accessTokenLock *sync.RWMutex
	domain          string
	schema          string

	configured bool
	version    string
}

type providerData struct {
	Domain   types.String `tfsdk:"domain"`
	UseHTTPS types.Bool   `tfsdk:"use_https"`
}

func (p *mastodonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	p.accessTokenLock = &sync.RWMutex{}
	p.domain = data.Domain.Value
	p.schema = "https"
	if !data.UseHTTPS.IsNull() && !data.UseHTTPS.Value {
		p.schema = "http"
	}
	p.configured = true
}

func (p *mastodonProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"mastodon_register_app": registerAppResourceType{},
	}, nil
}

func (p *mastodonProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"mastodon_instance_self": instanceSelfDataSourceType{},
	}, nil
}

func (p *mastodonProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domain": {
				MarkdownDescription: "Domain",
				Required:            true,
				Type:                types.StringType,
			},
			"use_https": {
				MarkdownDescription: "Should we use https to connect to the instance",
				Optional:            true,
				Type:                types.BoolType,
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &mastodonProvider{
			version: version,
		}
	}
}

func convertProviderType(in provider.Provider) (mastodonProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*mastodonProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return mastodonProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return mastodonProvider{}, diags
	}

	return *p, diags
}
