package main

import (
	"context"
	"encoding/json"
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NameResource{}

type NameResource struct{
	apiBaseURL string
}

type NameResourceModel struct {
	Id           types.String `tfsdk:"id"`
	ProjectId    types.String `tfsdk:"project_id"`
	Assettag     types.String `tfsdk:"assettag"`
	NameContext  types.String `tfsdk:"name_context"`
	ResourceType types.String `tfsdk:"resource_type"`
	Environment  types.String `tfsdk:"environment"`
	Cloudregion  types.String `tfsdk:"cloudregion"`
	PlatformCode types.String `tfsdk:"platform_code"`
	Name         types.String `tfsdk:"name"`
	CacheKey     types.String `tfsdk:"cache_key"`
	Cached       types.Bool   `tfsdk:"cached"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}

type APIRequest struct {
	ResourceProperties struct {
		Type         string `json:"type"`
		ResourceType string `json:"resource_type"`
		Cloudregion  string `json:"cloudregion"`
		PlatformCode string `json:"platform_code"`
		Environment  string `json:"environment"`
		Assettag     string `json:"assettag"`
		NameContext  string `json:"name_context"`
	} `json:"ResourceProperties"`
}

type APIResponse struct {
	Result string `json:"Result"`
}

func NewNameResource() resource.Resource {
	return &NameResource{}
}
// Implement the Configure method to receive provider data
func (r *NameResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
       if req.ProviderData == nil {
	       r.apiBaseURL = "https://bie-cih-d-csc-apim.azure-api.net/bie-cih-d-fa-namegen/namegenerator"
	       return
       }
       apiBaseURL, ok := req.ProviderData.(string)
       if ok && apiBaseURL != "" {
	       r.apiBaseURL = apiBaseURL
       } else {
	       r.apiBaseURL = "https://bie-cih-d-csc-apim.azure-api.net/bie-cih-d-fa-namegen/namegenerator"
       }
}

func (r *NameResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

func (r *NameResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Name generation resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "GCP project ID",
				Required:            true,
			},
			"assettag": schema.StringAttribute{
				MarkdownDescription: "6-digit asset tag",
				Required:            true,
			},
			"name_context": schema.StringAttribute{
				MarkdownDescription: "Resource context/purpose",
				Required:            true,
			},
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Resource type code",
				Optional:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Environment (d/t/p)",
				Optional:            true,
			},
			"cloudregion": schema.StringAttribute{
				MarkdownDescription: "Cloud region code",
				Optional:            true,
			},
			"platform_code": schema.StringAttribute{
				MarkdownDescription: "Platform code",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Generated resource name",
			},
			"cache_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique cache key",
			},
			"cached": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether result was cached",
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp",
			},
		},
	}
}

func (r *NameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
       var data NameResourceModel

       // Read Terraform plan data into the model
       resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
       if resp.Diagnostics.HasError() {
	       return
       }

       // Generate cache key
       cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s-%s",
	       data.Cloudregion.ValueString(),
	       data.ProjectId.ValueString(),
	       data.Assettag.ValueString(),
	       data.ResourceType.ValueString(),
	       data.NameContext.ValueString(),
	       data.Environment.ValueString())

       // Call name generation API
       apiReq := APIRequest{}
       apiReq.ResourceProperties.Type = "gcpname"
       apiReq.ResourceProperties.ResourceType = data.ResourceType.ValueString()
       apiReq.ResourceProperties.Cloudregion = data.Cloudregion.ValueString()
       apiReq.ResourceProperties.PlatformCode = data.PlatformCode.ValueString()
       apiReq.ResourceProperties.Environment = data.Environment.ValueString()
       apiReq.ResourceProperties.Assettag = data.Assettag.ValueString()
       apiReq.ResourceProperties.NameContext = data.NameContext.ValueString()

       jsonData, err := json.Marshal(apiReq)
       if err != nil {
	       resp.Diagnostics.AddError("JSON Marshal Error", fmt.Sprintf("Unable to marshal API request: %s", err))
	       return
       }

	       // Use API URL from resource struct
	       apiURL := r.apiBaseURL

	       // Make HTTP request to name generation API
	       httpResp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
       if err != nil {
	       resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to call name generation API: %s", err))
	       return
       }
       defer httpResp.Body.Close()

       var apiResp APIResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&apiResp); err != nil {
		resp.Diagnostics.AddError("API Response Error", fmt.Sprintf("Unable to decode API response: %s", err))
		return
	}

	// Set computed values
	data.Id = types.StringValue(cacheKey)
	data.Name = types.StringValue(apiResp.Result)
	data.CacheKey = types.StringValue(cacheKey)
	data.Cached = types.BoolValue(false) // New generation
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NameResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No additional API calls needed for read - data is in state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NameResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update timestamp
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NameResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No API cleanup needed - just remove from state
}