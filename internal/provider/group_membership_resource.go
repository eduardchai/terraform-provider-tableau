package provider

import (
	"context"
	"fmt"
	"terraform-provider-tableau/internal/client"
	"terraform-provider-tableau/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &groupMembershipResource{}
	_ resource.ResourceWithConfigure   = &groupMembershipResource{}
	_ resource.ResourceWithImportState = &groupMembershipResource{}
)

type groupMembershipResource struct {
	client *client.TableauClient
}

type groupMembershipResourceModel struct {
	GroupID    types.String `tfsdk:"group_id"`
	UserEmails types.List   `tfsdk:"users"`
}

func NewGroupMembershipResource() resource.Resource {
	return &groupMembershipResource{}
}

// Metadata returns the resource type name.
func (r *groupMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_membership"
}

// Schema defines the schema for the resource.
func (r *groupMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Group id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"users": schema.ListAttribute{
				Required:    true,
				Description: "List of user ids",
				ElementType: types.StringType,
			},
		},
	}
}

// Create a new resource.
func (r *groupMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupMembershipResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse plan tf list types to go list/slice types
	var userEmails []string
	diags = plan.UserEmails.ElementsAs(ctx, &userEmails, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Add users to group
	for _, email := range userEmails {
		err := r.client.CreateGroupMembershipByUserEmail(
			plan.GroupID.ValueString(),
			email,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add user to Tableau Group",
				err.Error(),
			)
			return
		}
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *groupMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state groupMembershipResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed values
	groupMembershipEmailList, err := r.client.GetGroupMembership(state.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Group Membership",
			"Could not read Tableau Group Membership"+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.GroupID = types.StringValue(groupMembershipEmailList.GroupID)
	state.UserEmails, diags = types.ListValueFrom(ctx, types.StringType, groupMembershipEmailList.UserEmails)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan groupMembershipResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse plan tf list types to go list/slice types
	var userEmails []string
	diags = plan.UserEmails.ElementsAs(ctx, &userEmails, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get actual values
	groupMembershipEmailList, err := r.client.GetGroupMembership(plan.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Group Membership",
			"Could not read Tableau Group Membership"+": "+err.Error(),
		)
		return
	}

	// Delete user if groupMembershipEmailList.UserEmails is not in plan.UserEmails
	for _, email := range groupMembershipEmailList.UserEmails {
		if !utils.StringInSlice(email, userEmails) {
			err = r.client.DeleteGroupMembershipByUserEmail(plan.GroupID.ValueString(), email)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to delete user from Tableau Group",
					err.Error(),
				)
				return
			}
		}
	}

	// Add user if plan.UserEmails is not in groupMembershipEmailList.UserEmails
	for _, email := range userEmails {
		if !utils.StringInSlice(email, groupMembershipEmailList.UserEmails) {
			err = r.client.CreateGroupMembershipByUserEmail(
				plan.GroupID.ValueString(),
				email,
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to add user to Tableau Group",
					err.Error(),
				)
				return
			}
		}
	}

	// Get updated values
	updatedGroupMembershipEmailList, err := r.client.GetGroupMembership(plan.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Group Membership",
			"Could not read Tableau Group Membership"+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated values
	plan.GroupID = types.StringValue(updatedGroupMembershipEmailList.GroupID)
	plan.UserEmails, diags = types.ListValueFrom(ctx, types.StringType, updatedGroupMembershipEmailList.UserEmails)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Deletes the resource and removes the Terraform state on success.
func (r *groupMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state groupMembershipResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse plan tf list types to go list/slice types
	var userEmails []string
	state.UserEmails.ElementsAs(ctx, &userEmails, false)

	// Delete users from group
	for _, email := range userEmails {
		err := r.client.DeleteGroupMembershipByUserEmail(state.GroupID.ValueString(), email)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to delete user from Tableau Group",
				err.Error(),
			)
			return
		}
	}
}

// Configure adds the provider configured client to the resource.
func (r *groupMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("group_id"), req, resp)
}
