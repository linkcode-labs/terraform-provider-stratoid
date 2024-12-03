package userflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-stratoid/internal/entra/identity/userflow/model"
	"terraform-provider-stratoid/internal/helpers/tf"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListUserFlowOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]model.IdentityUserFlow
}

type ListUserFlowCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []model.IdentityUserFlow
}

type ListUserFlowOperationOptions struct {
	Count     *bool
	Expand    *odata.Expand
	Filter    *string
	Metadata  *odata.Metadata
	OrderBy   *odata.OrderBy
	RetryFunc client.RequestRetryFunc
	Search    *string
	Select    *[]string
	Skip      *int64
	Top       *int64
}

func DefaultListUserFlowOperationOptions() ListUserFlowOperationOptions {
	return ListUserFlowOperationOptions{}
}

func (o ListUserFlowOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListUserFlowOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	if o.Count != nil {
		out.Count = *o.Count
	}
	if o.Expand != nil {
		out.Expand = *o.Expand
	}
	if o.Filter != nil {
		out.Filter = *o.Filter
	}
	if o.Metadata != nil {
		out.Metadata = *o.Metadata
	}
	if o.OrderBy != nil {
		out.OrderBy = *o.OrderBy
	}
	if o.Search != nil {
		out.Search = *o.Search
	}
	if o.Select != nil {
		out.Select = *o.Select
	}
	if o.Skip != nil {
		out.Skip = int(*o.Skip)
	}
	if o.Top != nil {
		out.Top = int(*o.Top)
	}
	return &out
}

func (o ListUserFlowOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

type ListUserFlowCustomPager struct {
	NextLink *odata.Link `json:"@odata.nextLink"`
}

func (p *ListUserFlowCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListUserFlowAttributes - List identityUserFlowAttributes. Retrieve a list of identityUserFlowAttribute objects.
func (c UserFlowClient) ListUserFlows(ctx context.Context, options ListUserFlowOperationOptions) (result ListUserFlowOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListUserFlowCustomPager{},
		Path:          "/identity/AuthenticationEventsFlows",
		RetryFunc:     options.RetryFunc,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	tf.ErrorDiagPathF(nil, "[DEBUG] %s was not found - removing from state!", req.Header.Get("Authorization"))

	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]json.RawMessage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	temp := make([]model.IdentityUserFlow, 0)
	if values.Values != nil {
		for i, v := range *values.Values {
			val, err := model.UnmarshalIdentityUserFlowImplementation(v)
			if err != nil {
				err = fmt.Errorf("unmarshalling item %d for model.IdentityUserFlow (%q): %+v", i, v, err)
				return result, err
			}
			temp = append(temp, val)
		}
	}
	result.Model = &temp

	return
}

// ListUserFlowAttributesComplete retrieves all the results into a single object
func (c UserFlowClient) ListUserFlowComplete(ctx context.Context, options ListUserFlowOperationOptions) (ListUserFlowCompleteResult, error) {
	return c.ListUserFlowCompleteMatchingPredicate(ctx, options, IdentityUserFlowOperationPredicate{})
}

// ListUserFlowAttributesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c UserFlowClient) ListUserFlowCompleteMatchingPredicate(ctx context.Context, options ListUserFlowOperationOptions, predicate IdentityUserFlowOperationPredicate) (result ListUserFlowCompleteResult, err error) {
	items := make([]model.IdentityUserFlow, 0)

	resp, err := c.ListUserFlows(ctx, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListUserFlowCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
