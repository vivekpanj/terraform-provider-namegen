package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Id            types.String `tfsdk:"id"`
	Type          types.String `tfsdk:"type"`
	ApiUrl        types.String `tfsdk:"api_url"`
	HostnameType  types.String `tfsdk:"hostname_type"`
	StackId       types.String `tfsdk:"stack_id"`
	ResourceType  types.String `tfsdk:"resource_type"`
	Cloudregion   types.String `tfsdk:"cloudregion"`
	PlatformCode  types.String `tfsdk:"platform_code"`
	Environment   types.String `tfsdk:"environment"`
	Assettag      types.String `tfsdk:"assettag"`
	NameContext   types.String `tfsdk:"name_context"`
	Name          types.String `tfsdk:"name"`
	CacheKey      types.String `tfsdk:"cache_key"`
	Cached        types.Bool   `tfsdk:"cached"`
	LastUpdated   types.String `tfsdk:"last_updated"`
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
	       return
       }
       apiBaseURL, ok := req.ProviderData.(string)
       if ok && apiBaseURL != "" {
	       r.apiBaseURL = apiBaseURL
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
		       "type": schema.StringAttribute{
			       MarkdownDescription: "Type of name generation (host, DB, gcpname)",
			       Required:            true,
		       },
		       "api_url": schema.StringAttribute{
			       MarkdownDescription: "API endpoint URL for name generation",
			       Required:            true,
		       },
		       "hostname_type": schema.StringAttribute{
			       MarkdownDescription: "Hostname type (required for host/DB)",
			       Optional:            true,
		       },
		       "stack_id": schema.StringAttribute{
			       MarkdownDescription: "Stack ID (required for host/DB)",
			       Optional:            true,
		       },
		       "resource_type": schema.StringAttribute{
			       MarkdownDescription: "Resource type code (required for gcpname)",
			       Optional:            true,
		       },
		       "cloudregion": schema.StringAttribute{
			       MarkdownDescription: "Cloud region code (required for gcpname)",
			       Optional:            true,
		       },
		       "platform_code": schema.StringAttribute{
			       MarkdownDescription: "Platform code (required for gcpname)",
			       Optional:            true,
		       },
		       "environment": schema.StringAttribute{
			       MarkdownDescription: "Environment (required for gcpname)",
			       Optional:            true,
		       },
		       "assettag": schema.StringAttribute{
			       MarkdownDescription: "6-digit asset tag (required for gcpname)",
			       Optional:            true,
		       },
		       "name_context": schema.StringAttribute{
			       MarkdownDescription: "Resource context/purpose (required for gcpname)",
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

       resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
       if resp.Diagnostics.HasError() {
	       return
       }

       // Validate required fields based on type
       t := data.Type.ValueString()
       switch t {
       case "host", "DB":
	       if data.HostnameType.IsNull() || data.StackId.IsNull() {
		       resp.Diagnostics.AddError("Missing Required Fields", "'hostname_type' and 'stack_id' are required for type 'host' or 'DB'.")
		       return
	       }
       case "gcpname":
	       if data.ResourceType.IsNull() || data.Cloudregion.IsNull() || data.PlatformCode.IsNull() || data.Environment.IsNull() || data.Assettag.IsNull() || data.NameContext.IsNull() {
		       resp.Diagnostics.AddError("Missing Required Fields", "'resource_type', 'cloudregion', 'platform_code', 'environment', 'assettag', and 'name_context' are required for type 'gcpname'.")
		       return
	       }
       default:
	       resp.Diagnostics.AddError("Invalid Type", "'type' must be one of: host, DB, gcpname.")
	       return
       }

       // Generate cache key (example, can be customized)
       cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s-%s",
	       t,
	       data.HostnameType.ValueString(),
	       data.StackId.ValueString(),
	       data.ResourceType.ValueString(),
	       data.Cloudregion.ValueString(),
	       data.PlatformCode.ValueString(),
	       data.Environment.ValueString(),
	       data.Assettag.ValueString(),
       )

       // Build API request body
       var apiReq map[string]interface{}
       switch t {
       case "host", "DB":
	       apiReq = map[string]interface{}{
		       "ResourceProperties": map[string]interface{}{
			       "HostnameType": data.HostnameType.ValueString(),
			       "StackId":      data.StackId.ValueString(),
			       "type":         t,
		       },
	       }
       case "gcpname":
	       apiReq = map[string]interface{}{
		       "ResourceProperties": map[string]interface{}{
			       "type":          t,
			       "resource_type": data.ResourceType.ValueString(),
			       "cloudregion":   data.Cloudregion.ValueString(),
			       "platform_code": data.PlatformCode.ValueString(),
			       "environment":   data.Environment.ValueString(),
			       "assettag":      data.Assettag.ValueString(),
			       "name_context":  data.NameContext.ValueString(),
		       },
	       }
       }

       jsonData, err := json.Marshal(apiReq)
       if err != nil {
	       resp.Diagnostics.AddError("JSON Marshal Error", fmt.Sprintf("Unable to marshal API request: %s", err))
	       return
       }

       // Make HTTP request to name generation API with retry logic
       apiURL := data.ApiUrl.ValueString()
       
       var httpResp *http.Response
       maxRetries := 3
       retryDelay := time.Second * 2
       
       for attempt := 1; attempt <= maxRetries; attempt++ {
	       httpResp, err = http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	       
	       // Success - break out of retry loop
	       if err == nil && httpResp != nil && httpResp.StatusCode >= 200 && httpResp.StatusCode < 300 {
		       break
	       }
	       
	       // If this was the last attempt, fail
	       if attempt == maxRetries {
		       if err != nil {
			       resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to call name generation API after %d attempts: %s", maxRetries, err))
		       } else if httpResp != nil {
			       resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned HTTP status %d after %d attempts", httpResp.StatusCode, maxRetries))
		       }
		       return
	       }
	       
	       // Close response body if exists and retry
	       if httpResp != nil && httpResp.Body != nil {
		       httpResp.Body.Close()
	       }
	       
	       // Wait before retrying (exponential backoff)
	       time.Sleep(retryDelay * time.Duration(attempt))
       }
       
       defer httpResp.Body.Close()

       var apiResp APIResponse
       if err := json.NewDecoder(httpResp.Body).Decode(&apiResp); err != nil {
	       resp.Diagnostics.AddError("API Response Error", fmt.Sprintf("Unable to decode API response: %s", err))
	       return
       }

       // Check if result is empty
       if apiResp.Result == "" {
	       resp.Diagnostics.AddError("API Error", "API returned empty result")
	       return
       }

       // Check for error patterns in result (more comprehensive)
       resultLower := strings.ToLower(apiResp.Result)
       errorPatterns := []string{"error", "err", "failed", "fail", "invalid", "denied", "forbidden", "unauthorized"}
       for _, pattern := range errorPatterns {
	       if strings.Contains(resultLower, pattern) {
		       resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned error: %s", apiResp.Result))
		       return
	       }
       }

       // Validate that result looks like a valid name (basic sanity check)
       if len(apiResp.Result) < 3 {
	       resp.Diagnostics.AddError("Invalid Name", fmt.Sprintf("Generated name is too short: %s", apiResp.Result))
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
