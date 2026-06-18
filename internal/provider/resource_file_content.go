package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                   = &FileContentResource{}
	_ resource.ResourceWithImportState    = &FileContentResource{}
	_ resource.ResourceWithValidateConfig = &FileContentResource{}
)

func NewFileContentResource() resource.Resource {
	return &FileContentResource{}
}

// FileContentResource manages the content stored against an existing file.
type FileContentResource struct {
	client *Client
}

// FileContentResourceModel describes the resource data model.
type FileContentResourceModel struct {
	ID            types.String `tfsdk:"id"`
	FileID        types.String `tfsdk:"file_id"`
	Content       types.String `tfsdk:"content"`
	Source        types.String `tfsdk:"source"`
	SourceURL     types.String `tfsdk:"source_url"`
	SourceHash    types.String `tfsdk:"source_hash"`
	ContentType   types.String `tfsdk:"content_type"`
	ContentSHA256 types.String `tfsdk:"content_sha256"`
}

// Metadata returns the resource type name.
func (r *FileContentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_content"
}

// Schema defines the schema for the resource.
func (r *FileContentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Uploads content to an existing file. Provide exactly one of `content`, `source`, or `source_url`. Large files are uploaded directly to storage via a presigned request.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the resource (equal to the file ID)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the file to upload content to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Inline content to upload as the file body",
			},
			"source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to a local file whose contents are uploaded",
			},
			"source_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An HTTP(S) or data URL the platform fetches and stores server-side",
			},
			"source_hash": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional hash of the source content (e.g. `filesha256(...)`) used to trigger re-uploads when a local `source` file changes",
			},
			"content_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The MIME type of the content. Detected automatically when omitted.",
			},
			"content_sha256": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SHA-256 digest of the uploaded content",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *FileContentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// ValidateConfig ensures exactly one content source is provided.
func (r *FileContentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data FileContentResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	count := 0
	for _, set := range []bool{
		!data.Content.IsNull(),
		!data.Source.IsNull(),
		!data.SourceURL.IsNull(),
	} {
		if set {
			count++
		}
	}

	if count != 1 {
		resp.Diagnostics.AddError(
			"Invalid content source",
			"Exactly one of `content`, `source`, or `source_url` must be set.",
		)
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *FileContentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FileContentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apply(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *FileContentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Content is write-mostly: the upload endpoints do not expose a cheap way to
	// read the stored bytes back, so we keep the prior state as-is and rely on
	// configuration changes (and source_hash) to drive updates.
	var data FileContentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *FileContentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FileContentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apply(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes the resource from Terraform state. The file's content is left
// in place; it is replaced on the next upload and removed when the file itself
// is deleted (via chatbotkit_file).
func (r *FileContentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState imports the resource using the file ID.
func (r *FileContentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("file_id"), req.ID)...)
}

// apply uploads the configured content to the file and updates the computed
// fields on the model.
func (r *FileContentResource) apply(ctx context.Context, data *FileContentResourceModel, diags *diag.Diagnostics) {
	fileID := data.FileID.ValueString()
	endpoint := r.client.fileUploadEndpoint(fileID)

	switch {
	case !data.SourceURL.IsNull():
		if err := r.client.uploadContentFromURL(ctx, endpoint, data.SourceURL.ValueString()); err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to upload file content: %s", err))
			return
		}
		data.ContentSHA256 = types.StringValue(sha256Hex([]byte(data.SourceURL.ValueString())))

	case !data.Source.IsNull():
		body, err := os.ReadFile(data.Source.ValueString())
		if err != nil {
			diags.AddError("Source Error", fmt.Sprintf("Unable to read source file: %s", err))
			return
		}
		name := filepath.Base(data.Source.ValueString())
		if err := r.client.uploadContent(ctx, endpoint, body, data.ContentType.ValueString(), name); err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to upload file content: %s", err))
			return
		}
		data.ContentSHA256 = types.StringValue(sha256Hex(body))

	default:
		body := []byte(data.Content.ValueString())
		if err := r.client.uploadContent(ctx, endpoint, body, data.ContentType.ValueString(), ""); err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to upload file content: %s", err))
			return
		}
		data.ContentSHA256 = types.StringValue(sha256Hex(body))
	}

	data.ID = types.StringValue(fileID)
}
