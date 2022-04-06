package cdnfrontdoorsecurityparams

import (
	"fmt"

	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecurityParameters struct {
	TypeName   track1.Type
	ConfigName string
}

type CdnFrontdoorSecurityMappings struct {
	Firewall CdnFrontdoorSecurityParameters
}

func InitializeCdnFrontdoorSecurityMappings() *CdnFrontdoorSecurityMappings {
	m := new(CdnFrontdoorSecurityMappings)

	m.Firewall = CdnFrontdoorSecurityParameters{
		TypeName:   track1.TypeWebApplicationFirewall,
		ConfigName: "firewall",
	}

	return m
}

func ExpandCdnFrontdoorFirewallPolicyParameters(input []interface{}, isStandardSku bool) (track1.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := track1.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return results, nil
	}

	m := InitializeCdnFrontdoorSecurityMappings()
	associations := make([]track1.SecurityPolicyWebApplicationFirewallAssociation, 0)

	// pull off only the firewall policy from the security_policies list
	policyType := input[0].(map[string]interface{})
	firewallPolicy := policyType[m.Firewall.ConfigName].([]interface{})
	v := firewallPolicy[0].(map[string]interface{})

	if id := v["cdn_frontdoor_firewall_policy_id"].(string); id != "" {
		results.WafPolicy = &track1.ResourceReference{
			ID: utils.String(id),
		}
	}

	configAssociations := v["association"].([]interface{})

	for _, item := range configAssociations {
		v := item.(map[string]interface{})
		domains := expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{}))

		if isStandardSku {
			if len(*domains) > 100 {
				return results, fmt.Errorf("the %q sku is only allowed to have 100 or less domains associated with the firewall policy, got %d", track1.SkuNameStandardAzureFrontDoor, len(*domains))
			}
		} else {
			if len(*domains) > 500 {
				return results, fmt.Errorf("the %q sku is only allowed to have 500 or less domains associated with the firewall policy, got %d", track1.SkuNamePremiumAzureFrontDoor, len(*domains))
			}
		}

		association := track1.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         domains,
			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations
	results.Type = m.Firewall.TypeName

	return results, nil
}

func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]track1.ActivatedResourceReference {
	results := make([]track1.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		activatedResourceReference := track1.ActivatedResourceReference{}

		if id := v["cdn_frontdoor_custom_domain_id"].(string); id != "" {
			activatedResourceReference.ID = utils.String(id)

			enabled := v["active"].(bool)

			if !enabled {
				activatedResourceReference.IsActive = utils.Bool(enabled)
			}

			results = append(results, activatedResourceReference)
		}
	}

	return &results
}

func FlattenCdnFrontdoorFirewallPolicyParameters(input track1.BasicSecurityPolicyPropertiesParameters) ([]interface{}, error) {
	// securityPolicy, ok := input.AsSecurityPolicyWebApplicationFirewallParameters()
	// if !ok {
	// 	return nil, fmt.Errorf("expected a security policy web application firewall parameters")
	// }
	// associations := make([]interface{}, 0)
	// domains := make(map[string]interface{})
	// ///////

	// for _, item := range *securityPolicy.Associations {
	// 	foo := make(map[string]interface{})

	// 	item.Domains

	// v := item.(map[string]interface{})
	// 	domains := expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{}))

	// 	}

	// 	association := track1.SecurityPolicyWebApplicationFirewallAssociation{
	// 		Domains:         domains,
	// 		PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
	// 	}

	// 	associations = append(associations, association)
	// }

	// ////

	// // put this inside of association
	// // patterns_to_match
	// // ptm := utils.FlattenStringSlice(securityPolicy.Associations)
	// for _, v := range *securityPolicy.Associations {
	// 	association := make([]interface{}, 0)
	// 	associations = append(associations, domains)
	// 	domains["domain"] = append()

	// 	association = append(association, flattenSecurityPoliciesActivatedResourceReference(v.Domains))
	// 	association = append(association, flattenSecurityPolicyPatternsToMatch(v.PatternsToMatch))
	// }

	// //
	// if wafId := securityPolicy.WafPolicy.ID; wafId != nil {
	// 	// TODO
	// }

	return nil, nil
}

func flattenSecurityPoliciesActivatedResourceReference(input *[]track1.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// missing the domain outter level map...

	for _, item := range *input {
		domain := make(map[string]interface{})
		domain["id"] = *item.ID
		domain["active"] = *item.IsActive

		results = append(results, domain)
	}

	return results
}

func flattenSecurityPolicyPatternsToMatch(input *[]string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	patternsToMatch := make(map[string]interface{})
	patternsToMatch["patterns_to_match"] = utils.FlattenStringSlice(input)
	results = append(results, patternsToMatch)

	return results
}