package userflow

import (
	"context"
	"encoding/json"
	"net/http"

	"terraform-provider-stratoid/internal/entra/identity/userflow/model"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateUserFlowOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        model.IdentityUserFlow
}

type CreateUserFlowOperationOptions struct {
	Metadata  *odata.Metadata
	RetryFunc client.RequestRetryFunc
}

func DefaultCreateUserFlowOperationOptions() CreateUserFlowOperationOptions {
	return CreateUserFlowOperationOptions{}
}

func (o CreateUserFlowOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CreateUserFlowOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	if o.Metadata != nil {
		out.Metadata = *o.Metadata
	}
	return &out
}

func (o CreateUserFlowOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// CreateUserFlowAttribute - Create identityUserFlowAttribute. Create a new custom identityUserFlowAttribute object.
func (c UserFlowClient) CreateUserFlow(ctx context.Context, input model.IdentityUserFlow, options CreateUserFlowOperationOptions) (result CreateUserFlowOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusCreated,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          "/identity/AuthenticationEventsFlows",
		RetryFunc:     options.RetryFunc,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var respObj json.RawMessage
	if err = resp.Unmarshal(&respObj); err != nil {
		return
	}
	model, err := model.UnmarshalIdentityUserFlowImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = model

	return
}
