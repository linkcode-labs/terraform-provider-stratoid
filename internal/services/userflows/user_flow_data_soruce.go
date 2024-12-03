package userflows

import (
	"context"
	"errors"
	"fmt"
	"log"
	"terraform-provider-stratoid/internal/clients"
	"terraform-provider-stratoid/internal/entra/identity/userflow"
	"terraform-provider-stratoid/internal/entra/identity/userflow/model"
	"terraform-provider-stratoid/internal/helpers/tf"
	"terraform-provider-stratoid/internal/helpers/tf/pluginsdk"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func userFlowDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		ReadContext: userFlowResourceReadByDisplayName,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"id": {
				Description:  "The employee identifier assigned to the user by the organisation",
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "display_name"},

				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Description:  "The display name of the user",
				Type:         pluginsdk.TypeString,
				Computed:     true,
				ExactlyOneOf: []string{"id", "display_name"},
				Optional:     true,
			},
		},
	}
}

func userFlowResourceReadByDisplayName(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) pluginsdk.Diagnostics {
	client := meta.(*clients.Client).UserFlows.UserFlowClient

	var foundObjectId *string

	if id, ok := d.Get("id").(string); ok && id != "" {
		options := userflow.ListUserFlowOperationOptions{
			Filter: pointer.To(fmt.Sprintf("id eq '%s'", odata.EscapeSingleQuote(id))),
		}

		resp, err := client.ListUserFlows(ctx, options)

		if err != nil {
			return tf.ErrorDiagF(err, "Finding user flow with id: %q  %s", id, resp.HttpResponse.Request.URL.String())
		}

		if resp.Model == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}

		count := len(*resp.Model)
		if count > 1 {
			return tf.ErrorDiagPathF(nil, "id", "More than one user found with id: %q", id)
		} else if count == 0 {
			return tf.ErrorDiagPathF(err, "id", "User with id %q was not found", id)
		}

		foundObjectId = (*resp.Model)[0].Entity().Id
	}

	userFlowId := model.NewIdentityUserFlowID(*foundObjectId)

	// if err != nil {
	// 	return tf.ErrorDiagPathF(err, "id", "Parsing ID")
	// }

	resp, err := client.GetUserFlow(ctx, userFlowId, userflow.DefaultGetUserFlowOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", userFlowId.ID())
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagF(err, "Retrieving %s", userFlowId.ID())
	}

	if resp.Model == nil {
		return tf.ErrorDiagF(errors.New("model was nil"), "Get user flow")
	}

	userFlow := resp.Model.IdentityUserFlow()
	tf.Set(d, "id", userFlow.Id)
	tf.Set(d, "display_name", userFlow.DisplayName.GetOrZero())

	return nil
}
