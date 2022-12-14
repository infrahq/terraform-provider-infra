package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
	"github.com/infrahq/infra/uid"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Infra user resource creates a user with a specified name. The name must be an email address.",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The user's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:      "The user's email address, e.g. `alice@example.com`.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsEmail(),
			},
			"password": {
				Description:      "The user's password. This password is one-time use and must be changed before the account can be used. If omitted, a password will be randomly generated. Note: this field will be empty for an imported user.",
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: stringMinLength(8),
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	user, err := client.CreateUser(ctx, &api.CreateUserRequest{Name: d.Get("name").(string)})
	if err != nil {
		return diag.FromErr(err)
	}

	if password := d.Get("password").(string); password != "" {
		request := api.UpdateUserRequest{
			ID:       user.ID,
			Password: password,
		}

		if _, err := client.UpdateUser(ctx, &request); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("password", user.OneTimePassword); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(user.ID.String())
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.GetUser(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("password") {
		oldPassword, newPassword := d.GetChange("password")
		request := api.UpdateUserRequest{
			ID:          id,
			OldPassword: oldPassword.(string),
			Password:    newPassword.(string),
		}

		if _, err := client.UpdateUser(ctx, &request); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteUser(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}

func userFromIDOrEmail(ctx context.Context, client *api.Client, d *schema.ResourceData, id, email string) (*api.User, error) {
	if s := d.Get(id).(string); s != "" {
		userID, err := uid.Parse([]byte(s))
		if err != nil {
			return nil, err
		}

		return client.GetUser(ctx, userID)
	}

	if s := d.Get(email).(string); s != "" {
		return userFromEmail(ctx, client, s)
	}

	return nil, fmt.Errorf("one of `%s,%s` must be specified", id, email)
}

func userFromEmail(ctx context.Context, client *api.Client, email string) (*api.User, error) {
	request := api.ListUsersRequest{
		Name:       email,
		ShowSystem: true,
		PaginationRequest: api.PaginationRequest{
			Limit: 1,
		},
	}

	response, err := client.ListUsers(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Count < 1 {
		return nil, fmt.Errorf("user not found: %s", email)
	}

	return &response.Items[0], nil
}
