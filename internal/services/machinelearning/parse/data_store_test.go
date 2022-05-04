package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = DataStoreId{}

func TestDataStoreIDFormatter(t *testing.T) {
	actual := NewDataStoreID("00000000-0000-0000-0000-000000000000", "resGroup1", "workspace1", "datastore1").ID()
	expected := "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/datastores/datastore1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDataStoreID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DataStoreId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing DatastoreName
			Input: "/",
			Error: true,
		},

		{
			// missing value for DatastoreName
			Input: "/datastore/",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/datastore/v1.0/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/datastore/v1.0/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},

		{
			// missing WorkspaceName
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/",
			Error: true,
		},

		{
			// missing Name
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/datastores/",
			Error: true,
		},

		{
			// valid
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/datastores/datastore1",
			Expected: &DataStoreId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				ResourceGroup:  "resGroup1",
				WorkspaceName:  "workspace1",
				Name:           "datastore1",
			},
		},

		{
			// upper-cased
			Input: "/DATASTORE/V1.0/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/WORKSPACE1/DATASTORES/DATASTORE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DataStoreID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
