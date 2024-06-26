// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package sigr

import (
	"strings"

	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
)

// GetRules - Returns the rules for the SignalRScanner
func (a *SignalRScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"sigr-001": {
			Id:             "sigr-001",
			Category:       scanners.RulesCategoryMonitoringAndAlerting,
			Recommendation: "SignalR should have diagnostic settings enabled",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armsignalr.ResourceInfo)
				_, ok := scanContext.DiagnosticsSettings[strings.ToLower(*service.ID)]
				return !ok, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-signalr/signalr-howto-diagnostic-logs",
		},
		"sigr-002": {
			Id:             "sigr-002",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "SignalR should have availability zones enabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armsignalr.ResourceInfo)
				sku := string(*i.SKU.Name)
				zones := false
				if strings.Contains(sku, "Premium") {
					zones = true
				}
				return !zones, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-signalr/availability-zones",
		},
		"sigr-003": {
			Id:             "sigr-003",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "SignalR should have a SLA",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				return false, "99.9%"
			},
			Url: "https://www.azure.cn/en-us/support/sla/signalr-service/",
		},
		"sigr-004": {
			Id:             "sigr-004",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "SignalR should have private endpoints enabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armsignalr.ResourceInfo)
				pe := len(i.Properties.PrivateEndpointConnections) > 0
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-signalr/howto-private-endpoints",
		},
		"sigr-005": {
			Id:             "sigr-005",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "SignalR SKU",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armsignalr.ResourceInfo)
				return false, string(*i.SKU.Name)
			},
			Url: "https://azure.microsoft.com/en-us/pricing/details/signalr-service/",
		},
		"sigr-006": {
			Id:             "sigr-006",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "SignalR Name should comply with naming conventions",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armsignalr.ResourceInfo)
				caf := strings.HasPrefix(*c.Name, "sigr")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"sigr-007": {
			Id:             "sigr-007",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "SignalR should have tags",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armsignalr.ResourceInfo)
				return len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
	}
}
