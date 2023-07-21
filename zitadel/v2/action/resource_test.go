package action_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccAction(t *testing.T) {
	resourceName := "zitadel_action"
	initialProperty := "initialproperty"
	updatedProperty := "updatedproperty"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(configProperty, _ interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id          = "%s"
  name            = "testaction"
  script          = "%s"
  timeout         = "10s"
  allowed_to_fail = true
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(interface{}) resource.TestCheckFunc {
	return func(expect interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetAction(frame, &management.GetActionRequest{Id: frame.StateID(state)})
			if err != nil {
				return err
			}
			actual := remoteResource.GetAction().GetScript()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
