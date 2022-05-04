package client

import (
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	dataplane "github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/sdk/dataplane"
)

type Client struct {
	DataStoreClient              *dataplane.DataStoresClient
	WorkspacesClient             *machinelearningservices.WorkspacesClient
	MachineLearningComputeClient *machinelearningservices.ComputeClient
}

func NewClient(o *common.ClientOptions) *Client {
	DataStoreClient := dataplane.NewDataStoresClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataStoreClient.Client, o.ResourceManagerAuthorizer)

	WorkspacesClient := machinelearningservices.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	MachineLearningComputeClient := machinelearningservices.NewComputeClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MachineLearningComputeClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DataStoreClient:              &DataStoreClient,
		WorkspacesClient:             &WorkspacesClient,
		MachineLearningComputeClient: &MachineLearningComputeClient,
	}
}
