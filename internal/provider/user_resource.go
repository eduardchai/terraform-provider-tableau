package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-tableau/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

type userResource struct {
	client *client.TableauClient
}

type userResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	SiteRole    types.String `tfsdk:"site_role"`
	AuthSetting types.String `tfsdk:"auth_setting"`
}

func NewUserResource() resource.Resource {
	return &userResource{}
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "User email",
			},
			"site_role": schema.StringAttribute{
				Required:    true,
				Description: "Site role for the user",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Creator",
						"Explorer",
						"Interactor",
						"Publisher",
						"ExplorerCanPublish",
						"ServerAdministrator",
						"SiteAdministratorExplorer",
						"SiteAdministratorCreator",
						"Unlicensed",
						"Viewer",
					}...),
				},
			},
			"auth_setting": schema.StringAttribute{
				Required:    true,
				Description: "Auth setting for the user",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"SAML",
						"ServerDefault",
						"OpenID",
						"TableauIDWithMFA",
					}...),
				},
			},
		},
	}
}

// Create a new resource.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create user
	user, err := r.client.CreateUser(
		plan.Email.ValueString(),
		plan.SiteRole.ValueString(),
		plan.AuthSetting.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tableau User",
			err.Error(),
		)
		return
	}

	// Set ID
	plan.ID = types.StringValue(user.ID)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var user *client.User
	var err error

	userID := state.ID.ValueString()
	if strings.HasPrefix(userID, "email/") {
		email := strings.Split(userID, "/")[1]
		user, err = r.client.GetUserByEmail(email)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Tableau User",
				"Could not read Tableau user email "+email+": "+err.Error(),
			)
			return
		}
	} else {
		// Get refreshed values
		user, err = r.client.GetUser(userID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Tableau User",
				"Could not read Tableau user ID "+userID+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(user.ID)
	state.Email = types.StringValue(user.Email)
	state.SiteRole = types.StringValue(user.SiteRole)
	state.AuthSetting = types.StringValue(user.AuthSetting)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update user
	_, err := r.client.UpdateUser(
		plan.ID.ValueString(),
		plan.Email.ValueString(),
		plan.SiteRole.ValueString(),
		plan.AuthSetting.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Tableau User",
			err.Error(),
		)
		return
	}

	// Fetch updated user from server
	updatedUser, err := r.client.GetUser(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau User",
			err.Error(),
		)
		return
	}

	// Update resource state with updated values
	plan.ID = types.StringValue(updatedUser.ID)
	plan.Email = types.StringValue(updatedUser.Email)
	plan.SiteRole = types.StringValue(updatedUser.SiteRole)
	plan.AuthSetting = types.StringValue(updatedUser.AuthSetting)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete user
	err := r.client.DeleteUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Tableau User",
			err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.TableauClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.TableauClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
