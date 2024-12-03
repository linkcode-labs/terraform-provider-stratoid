package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/hashicorp/go-azure-sdk/microsoft-graph/me/stable/me"
	// "github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/serviceprincipal"
	"terraform-provider-stratoid/internal/common"
	userflows "terraform-provider-stratoid/internal/services/userflows/client"

	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Environment environments.Environment
	TenantID    string
	ClientID    string
	ObjectID    string
	Claims      *claims.Claims

	TerraformVersion string

	StopContext context.Context

	UserFlows *userflows.Client
}

func (client *Client) build(ctx context.Context, o *common.ClientOptions) error {
	client.StopContext = ctx

	var err error

	if client.UserFlows, err = userflows.NewClient(o); err != nil {
		return fmt.Errorf("building clients for UserFlows: %v", err)
	}
	// Acquire an access token upfront, so we can decode the JWT and populate the claims
	token, err := o.Authorizer.Token(ctx, &http.Request{})
	if err != nil {
		return fmt.Errorf("unable to obtain access token: %v", err)
	}

	client.Claims, err = claims.ParseClaims(token)
	if err != nil {
		return fmt.Errorf("unable to parse claims in access token: %v", err)
	}

	// Log the claims for debugging
	claimsJson, err := json.Marshal(client.Claims)
	switch {
	case err != nil:
		log.Printf("[DEBUG] AzureAD Provider could not marshal access token claims for log output")
	case claimsJson == nil:
		log.Printf("[DEBUG] AzureAD Provider access token claims was nil")
	default:
		log.Printf("[DEBUG] AzureAD Provider access token claims: %s", claimsJson)
	}

	// Missing object ID of token holder will break many things
	client.ObjectID = client.Claims.ObjectId
	// if client.ObjectID == "" {
	// 	if strings.Contains(strings.ToLower(client.Claims.Scopes), "openid") {
	// 		log.Printf("[DEBUG] Querying Microsoft Graph to discover authenticated user principal object ID because the `oid` claim was missing from the access token")
	// 		resp, err := client.Users.MeClient.GetMe(ctx, me.DefaultGetMeOperationOptions())
	// 		if err != nil {
	// 			return fmt.Errorf("attempting to discover object ID for authenticated user principal: %+v", err)
	// 		}

	// 		if resp.Model != nil {
	// 			return fmt.Errorf("attempting to discover object ID for authenticated user principal: response was nil")
	// 		}

	// 		id := resp.Model.Id
	// 		if id == nil {
	// 			return fmt.Errorf("attempting to discover object ID for authenticated user principal: returned object ID was nil")
	// 		}

	// 		client.ObjectID = *id
	// 	} else {
	// 		log.Printf("[DEBUG] Querying Microsoft Graph to discover authenticated service principal object ID because the `oid` claim was missing from the access token")
	// 		options := serviceprincipal.ListServicePrincipalsOperationOptions{
	// 			Filter: pointer.To(fmt.Sprintf("appId eq '%s'", client.ClientID)),
	// 		}
	// 		resp, err := client.ServicePrincipals.ServicePrincipalClient.ListServicePrincipals(ctx, options)
	// 		if err != nil {
	// 			return fmt.Errorf("attempting to discover object ID for authenticated service principal: %+v", err)
	// 		}

	// 		if resp.Model != nil && len(*resp.Model) != 1 {
	// 			respLen := "nil"
	// 			if resp.Model != nil {
	// 				respLen = strconv.Itoa(len(*resp.Model))
	// 			}
	// 			return fmt.Errorf("attempting to discover object ID for authenticated service principal: unexpected number of results, expected 1, received %s", respLen)
	// 		}

	// 		id := (*resp.Model)[0].Id
	// 		if id == nil {
	// 			return fmt.Errorf("attempting to discover object ID for authenticated service principal: returned object ID was nil")
	// 		}

	// 		client.ObjectID = *id
	// 	}
	// }

	if client.ObjectID == "" {
		return fmt.Errorf("parsing claims in access token: oid claim is empty")
	}

	return nil
}
