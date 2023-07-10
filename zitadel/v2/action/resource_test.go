package action_test

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccZITADELAction(t *testing.T) {
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
		func(configProperty, _ string) string {
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
		checkDestroy(frame),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			rs := state.RootModule().Resources[frame.TerraformName]
			remoteResource, err := frame.GetAction(frame, &management.GetActionRequest{Id: rs.Primary.ID})
			if err != nil {
				return err
			}
			actual := remoteResource.GetAction().GetScript()
			if actual != expect {
				return fmt.Errorf("expected %s, actual: %s", expect, actual)
			}
			return nil
		}
	}
}

func checkDestroy(frame *test_utils.OrgTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := checkRemoteProperty(frame)("")(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}
