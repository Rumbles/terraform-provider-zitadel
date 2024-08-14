package privacy_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetPrivacyPolicyToDefault(helper.CtxWithID(ctx, d), &management.ResetPrivacyPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset privacy policy: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomPrivacyPolicy(helper.CtxWithID(ctx, d), &management.UpdateCustomPrivacyPolicyRequest{
		TosLink:      d.Get(tosLinkVar).(string),
		PrivacyLink:  d.Get(privacyLinkVar).(string),
		HelpLink:     d.Get(HelpLinkVar).(string),
		SupportEmail: d.Get(supportEmailVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update privacy policy: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomPrivacyPolicy(helper.CtxWithID(ctx, d), &management.AddCustomPrivacyPolicyRequest{
		TosLink:      d.Get(tosLinkVar).(string),
		PrivacyLink:  d.Get(privacyLinkVar).(string),
		HelpLink:     d.Get(HelpLinkVar).(string),
		SupportEmail: d.Get(supportEmailVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create privacy policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPrivacyPolicy(helper.CtxWithID(ctx, d), &management.GetPrivacyPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get privacy policy")
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		helper.OrgIDVar: policy.GetDetails().GetResourceOwner(),
		tosLinkVar:      policy.GetTosLink(),
		privacyLinkVar:  policy.GetPrivacyLink(),
		HelpLinkVar:     policy.GetHelpLink(),
		supportEmailVar: policy.GetSupportEmail(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of privacy policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
