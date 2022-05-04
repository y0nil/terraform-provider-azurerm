package validate

import "testing"

func TestDataStoreID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing DatastoreName
			Input: "/",
			Valid: false,
		},

		{
			// missing value for DatastoreName
			Input: "/datastore/",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/datastore/v1.0/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/datastore/v1.0/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Valid: false,
		},

		{
			// missing WorkspaceName
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/",
			Valid: false,
		},

		{
			// missing value for WorkspaceName
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/",
			Valid: false,
		},

		{
			// missing Name
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/",
			Valid: false,
		},

		{
			// missing value for Name
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/datastores/",
			Valid: false,
		},

		{
			// valid
			Input: "/datastore/v1.0/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/datastores/datastore1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/DATASTORE/V1.0/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/WORKSPACE1/DATASTORES/DATASTORE1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := DataStoreID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
