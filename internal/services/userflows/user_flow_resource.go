package userflows

import (
	"context"
	"errors"
	"log"
	"time"

	"terraform-provider-stratoid/internal/clients"
	"terraform-provider-stratoid/internal/entra/identity/userflow"
	"terraform-provider-stratoid/internal/entra/identity/userflow/model"

	"terraform-provider-stratoid/internal/helpers/tf"
	"terraform-provider-stratoid/internal/helpers/tf/pluginsdk"

	"github.com/hashicorp/go-azure-helpers/lang/response"
)

func userFlowResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		// CreateContext: userFlowResourceCreate,
		ReadContext: userFlowResourceRead,
		// UpdateContext: userFlowResourceUpdate,
		DeleteContext: userFlowResourceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		// Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
		// 	if _, errs := stable.ValidateIdentityUserFlowAttributeID(id, "id"); len(errs) > 0 {
		// 		out := ""
		// 		for _, err := range errs {
		// 			out += err.Error()
		// 		}
		// 		return fmt.Errorf(out)
		// 	}
		// 	return nil
		// }),

		SchemaVersion: 1,
		// StateUpgraders: []pluginsdk.StateUpgrader{
		// 	{
		// 		Type:    migrations.ResourceUserFlowAttributeInstanceResourceV0().CoreConfigSchema().ImpliedType(),
		// 		Upgrade: migrations.ResourceUserFlowAttributeInstanceStateUpgradeV0,
		// 		Version: 0,
		// 	},
		// },

		Schema: map[string]*pluginsdk.Schema{
			"display_name": {
				Description: "The display name of the user flow attribute.",
				Type:        pluginsdk.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

// func userFlowResourceCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) pluginsdk.Diagnostics {
// 	client := meta.(*clients.Client).UserFlows.UserFlowClient

// 	displayName := d.Get("display_name").(string)

// 	options := userflow.ListUserFlowOperationOptions{
// 		Filter: pointer.To(fmt.Sprintf("displayName eq '%s'", displayName)),
// 	}
// 	if resp, err := client.ListUserFlows(ctx, options); err != nil {
// 		return tf.ErrorDiagF(err, "Checking for existing user flow attribute")
// 	} else if resp.Model != nil {
// 		for _, r := range *resp.Model {
// 			model := r.IdentityUserFlow()
// 			if model.Id != nil && strings.EqualFold(model.DisplayName.GetOrZero(), displayName) {
// 				return tf.ImportAsExistsDiag("entra_user_flow", *model.Id)
// 			}
// 		}
// 	}

// 	oDataType := "#microsoft.graph.externalUsersSelfServiceSignUpEventsFlow"

// 	flow := model.BaseIdentityUserFlowImpl{
// 		ODataType:   &oDataType,
// 		DisplayName: nullable.NoZero(displayName),
// 		OnInteractiveAuthFlowStart: &model.OnInteractiveAuthFlowStart{
// 			ODataType:       "#microsoft.graph.onInteractiveAuthFlowStartExternalUsersSelfServiceSignUp",
// 			IsSignUpAllowed: true,
// 		},
// 		OnAuthenticationMethodLoadStart: &model.OnAuthenticationMethodLoadStart{
// 			IdentityProviders: &[]model.IdentityProvider{{Id: nullable.NoZero("EmailPassword-OAUTH")}},
// 			ODataType:         "#microsoft.graph.onAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp",
// 		},
// 	}

// 	resp, err := client.CreateUserFlow(ctx, flow, userflow.DefaultCreateUserFlowOperationOptions())
// 	if err != nil {
// 		return tf.ErrorDiagF(err, "Creating user flow")
// 	}

// 	if resp.Model == nil {
// 		return tf.ErrorDiagF(errors.New("model was nil"), "Creating user flow")
// 	}

// 	userFlow := resp.Model.IdentityUserFlow()

// 	if userFlow.Id == nil || *userFlow.Id == "" {
// 		return tf.ErrorDiagF(errors.New("API returned user flow with nil ID"), "Bad API Response")
// 	}

// 	id := model.NewIdentityUserFlowID(*userFlow.Id)
// 	d.SetId(id.ID())

// 	// Now ensure we can retrieve the attribute consistently
// 	if err = consistency.WaitForUpdate(ctx, func(ctx context.Context) (*bool, error) {
// 		resp, err := client.GetUserFlow(ctx, id, userflow.DefaultGetUserFlowOperationOptions())
// 		if err != nil {
// 			if response.WasNotFound(resp.HttpResponse) {
// 				return pointer.To(false), nil
// 			}
// 			return pointer.To(false), err
// 		}
// 		return pointer.To(resp.Model != nil), nil
// 	}); err != nil {
// 		return tf.ErrorDiagF(err, "Waiting for creation of %s", id)
// 	}

// 	return userFlowResourceRead(ctx, d, meta)
// }

// func userFlowResourceUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) pluginsdk.Diagnostics {
// 	client := meta.(*clients.Client).UserFlows.UserFlowAttributeClient

// 	id, err := stable.ParseIdentityUserFlowAttributeID(d.Id())
// 	if err != nil {
// 		return tf.ErrorDiagPathF(err, "id", "Parsing ID")
// 	}

// 	attr := stable.BaseIdentityUserFlowAttributeImpl{
// 		Description: nullable.NoZero(d.Get("description").(string)),
// 	}

// 	if _, err := client.UpdateUserFlowAttribute(ctx, *id, attr, userflowattribute.DefaultUpdateUserFlowAttributeOperationOptions()); err != nil {
// 		return tf.ErrorDiagF(err, "Could not update user flow attribute with ID: %q", id)
// 	}

// 	return userFlowResourceRead(ctx, d, meta)
// }

func userFlowResourceRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) pluginsdk.Diagnostics {
	client := meta.(*clients.Client).UserFlows.UserFlowClient

	id, err := model.ParseIdentityUserFlowID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing ID")
	}

	resp, err := client.GetUserFlow(ctx, *id, userflow.DefaultGetUserFlowOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagF(err, "Retrieving %s", id)
	}

	if resp.Model == nil {
		return tf.ErrorDiagF(errors.New("model was nil"), "Creating user flow attribute")
	}

	userFlow := resp.Model.IdentityUserFlow()

	tf.Set(d, "display_name", userFlow.DisplayName.GetOrZero())

	return nil
}

func userFlowResourceDelete(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) pluginsdk.Diagnostics {
	// client := meta.(*clients.Client).UserFlows.UserFlowAttributeClient

	// id, err := stable.ParseIdentityUserFlowAttributeID(d.Id())
	// if err != nil {
	// 	return tf.ErrorDiagPathF(err, "id", "Parsing ID")
	// }

	// if _, err := client.DeleteUserFlowAttribute(ctx, *id, userflowattribute.DefaultDeleteUserFlowAttributeOperationOptions()); err != nil {
	// 	return tf.ErrorDiagF(err, "Deleting %s", id)
	// }

	// if err := consistency.WaitForDeletion(ctx, func(ctx context.Context) (*bool, error) {
	// 	if resp, err := client.GetUserFlowAttribute(ctx, *id, userflowattribute.DefaultGetUserFlowAttributeOperationOptions()); err != nil {
	// 		if response.WasNotFound(resp.HttpResponse) {
	// 			return pointer.To(false), nil
	// 		}
	// 		return nil, err
	// 	}
	// 	return pointer.To(true), nil
	// }); err != nil {
	// 	return tf.ErrorDiagF(err, "Waiting for deletion of %s", id)
	// }

	return nil
}
