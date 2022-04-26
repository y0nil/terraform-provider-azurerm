---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via OIDC"
description: |-
  This guide will cover how to use OIDC for Azure resources as authentication for the Azure Provider.
---

# Azure Provider: Authenticating using OIDC for Azure resources

Terraform supports a number of different methods for authenticating to Azure:

- [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
- Authenticating to Azure using OIDC (covered in this guide)
- [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
- [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

Authentication using OIDC is useful when running Terraform in an environment able to issue ID tokens that are trusted by Azure, such as GitHub Actions.

For other use cases, we recommend using a service principal or a managed identity when running Terraform non-interactively (such as when running Terraform in a CI/CD pipeline), and authenticating using the Azure CLI when running Terraform locally.

## What is OIDC?

Azure uses OIDC as the mechanism for [Federated Identity Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation-create-trust) which can be used to authenticate to services that support Azure Active Directory (Azure AD) authentication. When making use of Identity Federation, Terraform will obtain an ID token from a third-party issuer, such as GitHub Actions, and exchange that ID token for an Access Token with Azure. This workflow happens without Terraform knowing any secret credentials and relies on a trust relationship that you establish between the identity provider (e.g. GitHub) and your Azure AD tenant.

Once authenticated, Terraform assumes the identity of a Service Principal within your Azure AD tenant, and is therefore subject to the same role-based access control (RBAC) and access control (IAM) for Azure resources.

Before you can use OIDC authentication, it has to be configured. This broadly follows these steps:

1. Create an Application and Service Principal within your Azure tenant, or select an existing Application and Service Principal whose identity will be assumed for RBAC and IAM.
2. Configure one or more Federated Credentials for the Azure AD Application.
3. Assign a role for the Service Principal, associating it with the subscription that will be used to run Terraform. This step gives the identity permission to access Azure Resource Manager (ARM) resources.
4. Configure Terraform to use OIDC Authentication.

Before you can configure the Application and Service Principal and then assign an IAM role, your account needs sufficient permissions. See the following section for more details.

## Configuring an Application and Service Principal

The (simplified) Terraform configuration below creates an Application and Service Principal, and configures a Federated Credential to enable Terraform to assume the identity of this service principal having obtained an ID token from a third-party. It then grants the Contributor role to the service principal to grant access to the Subscription.

```hcl
data "azuread_client_config" "current" {}

resource "azuread_application" "example" {
  display_name = "example"
}

resource "azuread_application_federated_identity_credential" "example" {
  application_object_id = azuread_application.example.object_id
  display_name          = "my-repo-deploy"
  description           = "Deployments for my-repo"
  audiences             = ["api://AzureADTokenExchange"]
  issuer                = "https://token.actions.githubusercontent.com"
  subject               = "repo:my-organization/my-repo:environment:prod"
}

resource "azuread_service_principal" "example" {
  application_id = azuread_application.example.application_id
}

data "azurerm_subscription" "current" {}

resource "azurerm_role_assignment" "example" {
  scope              = data.azurerm_subscription.current.id
  role_definition_id = "Contributor"
  principal_id       = azuread_service_principal.example.object_id
}

output "tenant_id" {
  value = data.azuread_client_config.current.tenant_id
}

output "client_id" {
  value = azuread_application.example.application_id
}
```

This configuration assumes the principal being used to execute Terraform has the Global Administrator role assigned. For more information on the required Azure AD Roles or API permissions, please see the respective documentation for the [azuread_application](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application) and [azuread_service_principal](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/service_principal) resources.

In order to assign the Contributor role to your Subscription, the principal being used to apply this configuration must be assigned the **Owner** role, or have **Contributor** plus **User Access Administrator** roles, in the scope of the subscription.

## Configuring Terraform to use OIDC Authentication

At this point we assume that a Federated Identity Credentials is configured for an Application in your Azure AD tenant, and that permissions have been assigned to the associated Service Principal via Azure's Identity and Access Management system.

Terraform can be configured to use OIDC for authentication in one of two ways: using environment variables, or by defining the fields within the provider block.

### Configuring with environment variables

Setting the`ARM_USE_OIDC` environment variable (equivalent to provider block argument [`use_oidc`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#use_oidc)) to `true` tells Terraform to use OIDC for authentication.

In order to know which Application it should request an ID token for, you must specify the `ARM_TENANT_ID` environment variable (equivalent to provider block argument [`tenant_id`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#tenant_id)), and the`ARM_CLIENT_ID` environment variable (equivalent to provider block argument [`client_id`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#client_id)).

In addition to these, Terraform also needs to know the subscription ID to identify the full context for the AzureRM provider.

Terraform will automatically consume the `ACTIONS_ID_TOKEN_REQUEST_URL` and `ACTIONS_ID_TOKEN_REQUEST_TOKEN` environment variables which are set automatically by GitHub for a workflow.

The following example GitHub workflow definition shows how to set these environment variables using values from GitHub secrets.

```yaml
name: Plan
on:
  pull_request:
    types: ['opened', 'synchronize']
    paths:
      - '.github/workflows/plan.yml'
      - '**.tf'

concurrency:
  group: 'plan-${{ github.head_ref }}'
  cancel-in-progress: true

permissions:
  contents: read
  id-token: write

jobs:
  plan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: hashicorp/setup-terraform@v2
      - run: terraform init

      - id: plan
        run: terraform plan -no-color
        env:
          ARM_USE_OIDC: true
          ARM_CLIENT_ID: ${{ secrets.ARM_CLIENT_ID }}
          ARM_TENANT_ID: ${{ secrets.ARM_TENANT_ID }}
          ARM_SUBSCRIPTION_ID: ${{ secrets.ARM_SUBSCRIPTION_ID }}
          TF_IN_AUTOMATION: true
```

Note specifically the `permissions` which are required in order for Terraform to both obtain an ID token, and to read your checked out configuration.

We also recommend defining provider blocks so that you can pin or constrain the version of the provider being used, and configure other optional settings.

```hcl
# We strongly recommend using the required_providers block to set the
# AzureRM Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.4.0"
    }
  }
}

# Configure the AzureRM Provider
provider "azurerm" {
  features {}
}
```

### Configuring with the provider block

It's also possible to configure OIDC authentication within the provider block.

```hcl
# We strongly recommend using the required_providers block to set the
# AzureRM Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.4.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}

  use_oidc  = true
  client_id = "00000000-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  tenant_id = "11111111-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

At this time, the Azure Backend does not yet support OIDC Authentication.

More information on [the fields supported in the provider block can be found here](../index.html#argument-reference).

<!-- it's not clear to me that we even need this info; it seems like this is the sort of thing you'd know about if you needed it.

### Custom MSI endpoints

Developers who are using a custom MSI endpoint can specify the endpoint in one of two ways:

- In the provider block using the `msi_endpoint` field
- Using the `ARM_MSI_ENDPOINT` environment variable.

You don't normally need to set the endpoint, because Terraform and the Azure Provider will automatically locate the appropriate endpoint.

-->
