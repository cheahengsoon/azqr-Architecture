// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cosmos

import (
	"strings"

	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
)

// GetRules - Returns the rules for the CosmosDBScanner
func (a *CosmosDBScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"cosmos-001": {
			Id:             "cosmos-001",
			Category:       scanners.RulesCategoryMonitoringAndAlerting,
			Recommendation: "CosmosDB should have diagnostic settings enabled",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armcosmos.DatabaseAccountGetResults)
				_, ok := scanContext.DiagnosticsSettings[strings.ToLower(*service.ID)]
				return !ok, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cosmos-db/monitor-resource-logs",
		},
		"cosmos-002": {
			Id:             "cosmos-002",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "CosmosDB should have availability zones enabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcosmos.DatabaseAccountGetResults)
				availabilityZones := false
				availabilityZonesNotEnabledInALocation := false
				numberOfLocations := 0
				for _, location := range i.Properties.Locations {
					numberOfLocations++
					if *location.IsZoneRedundant {
						availabilityZones = true
					} else {
						availabilityZonesNotEnabledInALocation = true
					}
				}

				zones := availabilityZones && numberOfLocations >= 2 && !availabilityZonesNotEnabledInALocation

				return !zones, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cosmos-db/high-availability",
		},
		"cosmos-003": {
			Id:             "cosmos-003",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "CosmosDB should have a SLA",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcosmos.DatabaseAccountGetResults)
				sla := "99.99%"
				availabilityZones := false
				availabilityZonesNotEnabledInALocation := false
				numberOfLocations := 0
				for _, location := range i.Properties.Locations {
					numberOfLocations++
					if *location.IsZoneRedundant {
						availabilityZones = true
						sla = "99.995%"
					} else {
						availabilityZonesNotEnabledInALocation = true
					}
				}

				if availabilityZones && numberOfLocations >= 2 && !availabilityZonesNotEnabledInALocation {
					sla = "99.999%"
				}
				return false, sla
			},
			Url: "https://learn.microsoft.com/en-us/azure/cosmos-db/high-availability#slas",
		},
		"cosmos-004": {
			Id:             "cosmos-004",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "CosmosDB should have private endpoints enabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcosmos.DatabaseAccountGetResults)
				pe := len(i.Properties.PrivateEndpointConnections) > 0
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-configure-private-endpoints",
		},
		"cosmos-005": {
			Id:             "cosmos-005",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "CosmosDB SKU",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcosmos.DatabaseAccountGetResults)
				return false, string(*i.Properties.DatabaseAccountOfferType)
			},
			Url: "https://azure.microsoft.com/en-us/pricing/details/cosmos-db/autoscale-provisioned/",
		},
		"cosmos-006": {
			Id:             "cosmos-006",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "CosmosDB Name should comply with naming conventions",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcosmos.DatabaseAccountGetResults)
				caf := strings.HasPrefix(*c.Name, "cosmos")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"cosmos-007": {
			Id:             "cosmos-007",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "CosmosDB should have tags",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcosmos.DatabaseAccountGetResults)
				return len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
		"cosmos-008": {
			Id:             "cosmos-008",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "CosmosDB should have local authentication disabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcosmos.DatabaseAccountGetResults)
				localAuth := c.Properties.DisableLocalAuth != nil && *c.Properties.DisableLocalAuth
				return !localAuth, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-rbac#disable-local-auth",
		},
		"cosmos-009": {
			Id:             "cosmos-009",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "CosmosDB: disable write operations on metadata resources (databases, containers, throughput) via account keys",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcosmos.DatabaseAccountGetResults)
				disabled := c.Properties.DisableKeyBasedMetadataWriteAccess != nil && *c.Properties.DisableKeyBasedMetadataWriteAccess
				return !disabled, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cosmos-db/role-based-access-control#set-via-arm-template",
		},
	}
}
