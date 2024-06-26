// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cog

import (
	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
)

// CognitiveScanner - Scanner for Cognitive Services Accounts
type CognitiveScanner struct {
	config *scanners.ScannerConfig
	client *armcognitiveservices.AccountsClient
}

// Init - Initializes the CognitiveScanner
func (a *CognitiveScanner) Init(config *scanners.ScannerConfig) error {
	a.config = config
	var err error
	a.client, err = armcognitiveservices.NewAccountsClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	return err
}

// Scan - Scans all Cognitive Services Accounts in a Resource Group
func (c *CognitiveScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	scanners.LogResourceGroupScan(c.config.SubscriptionID, resourceGroupName, "Cognitive Services")

	eventHubs, err := c.listEventHubs(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := c.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, eventHub := range eventHubs {
		rr := engine.EvaluateRules(rules, eventHub, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID:   c.config.SubscriptionID,
			SubscriptionName: c.config.SubscriptionName,
			ResourceGroup:    resourceGroupName,
			ServiceName:      *eventHub.Name,
			Type:             *eventHub.Type,
			Location:         *eventHub.Location,
			Rules:            rr,
		})
	}
	return results, nil
}

func (c *CognitiveScanner) listEventHubs(resourceGroupName string) ([]*armcognitiveservices.Account, error) {
	pager := c.client.NewListByResourceGroupPager(resourceGroupName, nil)

	namespaces := make([]*armcognitiveservices.Account, 0)
	for pager.More() {
		resp, err := pager.NextPage(c.config.Ctx)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, resp.Value...)
	}
	return namespaces, nil
}
