package data

import "enterprise/sizing-apps-backend/internal/model"

// GetPricingCatalog returns the exact pricing items matching Template Sizing.xlsx
func GetPricingCatalog() []model.LegacyPricingItem {
	return []model.LegacyPricingItem{
		// Hardware Standard (110% multiplier for standard pricing)
		{No: 1, Category: "Hardware", Item: "Compute vCPU", Environment: "Production", Unit: "vCPU", PriceAnnual: 8250000},
		{No: 2, Category: "Hardware", Item: "Compute vCPU", Environment: "Dev / Non-Prod", Unit: "vCPU", PriceAnnual: 7150000},
		{No: 3, Category: "Hardware", Item: "Storage SAN Tier 2", Environment: "Production", Unit: "GB", PriceAnnual: 33000},
		{No: 4, Category: "Hardware", Item: "Storage SAN Tier 1", Environment: "Production", Unit: "GB", PriceAnnual: 48000},
		{No: 5, Category: "Hardware", Item: "Storage SAN Tier 2", Environment: "Dev / Non-Prod", Unit: "GB", PriceAnnual: 28600},
		{No: 6, Category: "Hardware", Item: "Backup Netbackup (NBU)", Environment: "All", Unit: "GB", PriceAnnual: 27000},
		{No: 7, Category: "Hardware", Item: "Storage SAN Tier 1", Environment: "Dev / Non-Prod", Unit: "GB", PriceAnnual: 35000},
		{No: 8, Category: "Hardware", Item: "Server Baremetal (Standard)", Environment: "Production", Unit: "Unit", PriceAnnual: 794000000},
		{No: 9, Category: "Hardware", Item: "Server Baremetal (High-Perf)", Environment: "Production", Unit: "Unit", PriceAnnual: 794000000},

		// License Subscriptions
		{No: 1, Category: "License", Item: "NGINX Plus (F5)", Environment: "Production", Unit: "Instance", PriceAnnual: 75000000},
		{No: 2, Category: "License", Item: "Red Hat OCP Platform Plus Baremetal", Environment: "Production", Unit: "Node", PriceAnnual: 240000000},
		{No: 3, Category: "License", Item: "Oracle Database Enterprise Edition", Environment: "Production", Unit: "Proc Core", PriceAnnual: 134000000},
		{No: 4, Category: "License", Item: "Confluent Kafka Enterprise", Environment: "Production", Unit: "Broker", PriceAnnual: 125000000},
		{No: 5, Category: "License", Item: "Elasticsearch Platinum Edition", Environment: "Production", Unit: "Node", PriceAnnual: 75000000},
		{No: 6, Category: "License", Item: "Redis Enterprise Software", Environment: "Production", Unit: "RAM Shard", PriceAnnual: 65000000},
	}
}
