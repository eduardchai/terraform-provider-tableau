package provider

import (
	"context"
	"fmt"
	"terraform-provider-tableau/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

type userDataSource struct {
	client *client.TableauClient
}

type userDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	SiteRole    types.String `tfsdk:"site_role"`
	AuthSetting types.String `tfsdk:"auth_setting"`
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve user details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the user",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "User email",
			},
			"site_role": schema.StringAttribute{
				Computed:    true,
				Description: "Site role for the user",
			},
			"auth_setting": schema.StringAttribute{
				Computed:    true,
				Description: "Auth setting for the user",
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	user, err := d.client.GetUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau User",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(user.ID)
	state.Email = types.StringValue(user.Email)
	state.SiteRole = types.StringValue(user.SiteRole)
	state.AuthSetting = types.StringValue(user.AuthSetting)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}
