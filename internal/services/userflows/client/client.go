package client

import (
	"terraform-provider-stratoid/internal/common"
	"terraform-provider-stratoid/internal/entra/identity/userflow"
	// "github.com/hashicorp/go-azure-sdk/microsoft-graph/identity/stable/userflowattribute"
)

type Client struct {
	UserFlowClient *userflow.UserFlowClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	userFlowClient, err := userflow.NewUserFlowClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(userFlowClient.Client)

	return &Client{
		UserFlowClient: userFlowClient,
	}, nil
}
