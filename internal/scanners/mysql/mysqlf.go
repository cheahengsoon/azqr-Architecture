// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package mysql

import (
	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers"
)

// MySQLFlexibleScanner - Scanner for PostgreSQL
type MySQLFlexibleScanner struct {
	config         *scanners.ScannerConfig
	flexibleClient *armmysqlflexibleservers.ServersClient
}

// Init - Initializes the MySQLFlexibleScanner
func (c *MySQLFlexibleScanner) Init(config *scanners.ScannerConfig) error {
	c.config = config
	var err error
	c.flexibleClient, err = armmysqlflexibleservers.NewServersClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	return err
}

// Scan - Scans all MySQL in a Resource Group
func (c *MySQLFlexibleScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	scanners.LogResourceGroupScan(c.config.SubscriptionID, resourceGroupName, "MySQL Flexible")

	flexibles, err := c.listFlexiblePostgre(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := c.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, postgre := range flexibles {
		rr := engine.EvaluateRules(rules, postgre, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID:   c.config.SubscriptionID,
			ResourceGroup:    resourceGroupName,
			SubscriptionName: c.config.SubscriptionName,
			ServiceName:      *postgre.Name,
			Type:             *postgre.Type,
			Location:         *postgre.Location,
			Rules:            rr,
		})
	}

	return results, nil
}
func (c *MySQLFlexibleScanner) listFlexiblePostgre(resourceGroupName string) ([]*armmysqlflexibleservers.Server, error) {
	pager := c.flexibleClient.NewListByResourceGroupPager(resourceGroupName, nil)

	servers := make([]*armmysqlflexibleservers.Server, 0)
	for pager.More() {
		resp, err := pager.NextPage(c.config.Ctx)
		if err != nil {
			return nil, err
		}
		servers = append(servers, resp.Value...)
	}
	return servers, nil
}
