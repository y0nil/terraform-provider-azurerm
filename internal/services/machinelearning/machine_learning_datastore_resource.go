package machinelearning

import (
	"context"
	"fmt"
	"time"

	frsUUID "github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	dataPlane "github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/sdk/dataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataStoreResource struct{}

type DataStoreResourceModel struct {
	Name                       string `tfschema:"name"`
	MachineLearningWorkspaceId string `tfschema:"machine_learning_workspace_id"`
	BlobStorageId              string `tfschema:"blob_storage_id"`
	AccountKey                 string `tfschema:"account_key"`
	//ClientId                   string `tfschema:"client_id"`
	//TenantId                   string `tfschema:"tenant_id"`
	//ClientSecret               string `tfschema:"client_secret"`
}

func (r DataStoreResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"machine_learning_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"blob_storage_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"account_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		//"client_secret": {
		//	Type:         pluginsdk.TypeString,
		//	Required:     true,
		//	ValidateFunc: validation.StringIsNotEmpty,
		//},
		//
		//"tenant_id": {
		//	Type:         pluginsdk.TypeString,
		//	Required:     true,
		//	ValidateFunc: validation.StringIsNotEmpty,
		//},
	}
}

func (r DataStoreResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataStoreResource) ModelObject() interface{} {
	return &DataStoreResource{}
}

func (r DataStoreResource) ResourceType() string {
	return "azurerm_machine_learning_datastore"
}

func (r DataStoreResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataStoreID
}

func (r DataStoreResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DataStoreResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.MachineLearning.DataStoreClient

			workspaceId, err := parse.WorkspaceID(model.MachineLearningWorkspaceId)
			if err != nil {
				return err
			}

			id := parse.NewDataStoreID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, model.Name)

			subscriptionId, err := frsUUID.FromString(id.SubscriptionId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, subscriptionId, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id.ID(), err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			containerId, err := storageParse.StorageContainerDataPlaneID(model.BlobStorageId)
			if err != nil {
				return fmt.Errorf("parsing %q: %s", model.BlobStorageId, err)
			}

			props := &dataPlane.DataStore{
				Name:          utils.String(id.Name),
				DataStoreType: dataPlane.DataStoreTypeAzureBlob,
				AzureStorageSection: &dataPlane.AzureStorage{
					AccountKey:     utils.String(containerId.AccountName),
					ContainerName:  utils.String(containerId.Name),
					CredentialType: dataPlane.AzureStorageCredentialTypesAccountKey,
					Credential:     utils.String(model.AccountKey),
					SubscriptionID: &subscriptionId,
					ResourceGroup:  utils.String(workspaceId.ResourceGroup),
				},
			}

			createIfNotExists := true
			skipValidation := true

			_, err = client.Create(ctx, subscriptionId, id.ResourceGroup, id.Name, props, &createIfNotExists, &skipValidation)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", "some_id", err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DataStoreResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.DataStoreClient
			storageClient := metadata.Client.Storage
			id, err := parse.DataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionId, err := frsUUID.FromString(id.SubscriptionId)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, subscriptionId, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

			state := DataStoreResourceModel{
				Name:                       id.Name,
				MachineLearningWorkspaceId: workspaceId.ID(),
				AccountKey:                 metadata.ResourceData.Get("account_key").(string),
			}

			if props := resp.AzureStorageSection; props != nil {
				//subscriptionId := ""
				//if v := props.SubscriptionID; v != nil {
				//	subscriptionId = v.String()
				//}
				//
				//resourceGroup := ""
				//if v := props.ResourceGroup; v != nil {
				//	resourceGroup = *v
				//}

				accountName := ""
				if v := props.AccountName; v != nil {
					accountName = *v
				}

				containerName := ""
				if v := props.ContainerName; v != nil {
					containerName = *v
				}

				account, err := storageClient.FindAccount(ctx, accountName)
				if err != nil {
					return fmt.Errorf("retrieving Account %q for Container %q: %s", accountName, containerName, err)
				}
				if account == nil {
					return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
				}

				containerId := storageParse.NewStorageContainerDataPlaneId(accountName, storageClient.Environment.StorageEndpointSuffix, containerName)
				state.BlobStorageId = containerId.ID()
			}

			return nil
		},
	}
}

func (r DataStoreResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			return nil
		},
	}
}

func (r DataStoreResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.DataStoreClient
			id, err := parse.DataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionId, err := frsUUID.FromString(id.SubscriptionId)
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, subscriptionId, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
				if !response.WasNotFound(resp.Response) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}
