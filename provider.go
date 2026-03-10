package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &namegenProvider{}

type namegenProvider struct {
	version    string
	apiBaseURL string
}

type namegenProviderModel struct {
	APIBaseURL            types.String `tfsdk:"api_base_url"`
	DefaultCloudregion    types.String `tfsdk:"default_cloudregion"`
	DefaultPlatformCode   types.String `tfsdk:"default_platform_code"`
	DefaultEnvironment    types.String `tfsdk:"default_environment"`
}

func (p *namegenProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "namegen"
	resp.Version = p.version
}

func (p *namegenProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL for the name generation API",
				Optional:            true,
			},
			"default_cloudregion": schema.StringAttribute{
				MarkdownDescription: "Default cloud region for all resources",
				Optional:            true,
			},
			"default_platform_code": schema.StringAttribute{
				MarkdownDescription: "Default platform code for all resources",
				Optional:            true,
			},
			"default_environment": schema.StringAttribute{
				MarkdownDescription: "Default environment for all resources",
				Optional:            true,
			},
		},
	}
}

func (p *namegenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
       var data namegenProviderModel

       resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

       if resp.Diagnostics.HasError() {
	       return
       }

       if !data.APIBaseURL.IsNull() && !data.APIBaseURL.IsUnknown() {
	       p.apiBaseURL = data.APIBaseURL.ValueString()
       } else {
	       // Default value if not set
	       p.apiBaseURL = "https://bie-cih-d-csc-apim.azure-api.net/bie-cih-d-fa-namegen/namegenerator"
       }

       // Pass provider data to resources
       resp.DataSourceData = p.apiBaseURL
       resp.ResourceData = p.apiBaseURL
}

func (p *namegenProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNameResource,
	}
}

func (p *namegenProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &namegenProvider{
			version: version,
		}
	}
}