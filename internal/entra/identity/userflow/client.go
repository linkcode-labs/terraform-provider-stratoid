package userflow

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/msgraph"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserFlowClient struct {
	Client *msgraph.Client
}

func NewUserFlowClientWithBaseURI(sdkApi sdkEnv.Api) (*UserFlowClient, error) {
	client, err := msgraph.NewClient(sdkApi, "userflow", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating userflow client: %+v", err)
	}

	return &UserFlowClient{
		Client: client,
	}, nil
}