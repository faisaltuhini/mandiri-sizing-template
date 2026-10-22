package quote

import (
	"fmt"
	"strings"

	"enterprise/sizing-apps-backend/internal/model"
)

// Engine handles the pricing calculation for Enterprise Architecture infrastructure & licenses
type Engine struct {
	pricingMap   map[string]model.PricingItem
	componentMap map[string]model.CatalogComponent
}

// NewEngine creates an initialized calculation engine with pricing and catalog items
func NewEngine(pricing []model.PricingItem, components []model.CatalogComponent) *Engine {
	pMap := make(map[string]model.PricingItem)
	for _, p := range pricing {
		pMap[p.SKU] = p
	}

	cMap := make(map[string]model.CatalogComponent)
	for _, c := range components {
		cMap[c.ID.String()] = c
		// Also index by layer no or component name for convenience
		cMap[fmt.Sprintf("layer-%d", c.LayerNo)] = c
	}

	return &Engine{
		pricingMap:   pMap,
		componentMap: cMap,
	}
}

// CalculateQuote processes a quote request and returns detailed line-item costs & totals
func (e *Engine) CalculateQuote(req model.QuoteRequest) model.QuoteResponse {
	isDev := strings.Contains(strings.ToLower(req.Environment), "dev")
	var grandCompute, grandStorage, grandBackup, grandLicense int64
	var totalNodes, totalVCPU, totalRAM, totalStorage int

	var itemCosts []model.QuoteItemCost
	var warnings []string
	warningSet := make(map[string]bool)

	// Price lookup helpers
	vcpuPrice := int64(8250000) // Default Prod
	if isDev {
		vcpuPrice = 7150000
	}
	if p, ok := e.pricingMap["VCPU_PROD"]; ok && !isDev && p.UnitPriceIDR != nil {
		vcpuPrice = *p.UnitPriceIDR
	} else if p, ok := e.pricingMap["VCPU_DEV"]; ok && isDev && p.UnitPriceIDR != nil {
		vcpuPrice = *p.UnitPriceIDR
	}

	stgT2Price := int64(33000)
	if isDev {
		stgT2Price = 28600
	}
	if p, ok := e.pricingMap["STG_T2_PROD"]; ok && !isDev && p.UnitPriceIDR != nil {
		stgT2Price = *p.UnitPriceIDR
	} else if p, ok := e.pricingMap["STG_T2_DEV"]; ok && isDev && p.UnitPriceIDR != nil {
		stgT2Price = *p.UnitPriceIDR
	}

	stgT3Price := int64(48000)
	if p, ok := e.pricingMap["STG_T3_PROD"]; ok && p.UnitPriceIDR != nil {
		stgT3Price = *p.UnitPriceIDR
	}

	nbuPrice := int64(27000)
	if p, ok := e.pricingMap["NBU_ALL"]; ok && p.UnitPriceIDR != nil {
		nbuPrice = *p.UnitPriceIDR
	}

	bmPrice := int64(794000000)
	if p, ok := e.pricingMap["BM_STD_PROD"]; ok && p.UnitPriceIDR != nil {
		bmPrice = *p.UnitPriceIDR
	}

	for _, item := range req.Items {
		comp, hasComp := e.componentMap[item.ComponentID]
		compName := "Komponen"
		layerName := "Layer"
		if hasComp {
			compName = comp.ComponentName
			layerName = comp.LayerName
		}

		itemNodes := item.NodeCount
		if itemNodes <= 0 {
			itemNodes = 1
		}
		itemVCPU := itemNodes * item.VCPUPerNode
		itemRAM := itemNodes * item.RAMGBPerNode
		itemStorage := itemNodes * item.StorageGBPerNode

		// 1. Compute Cost
		var computeCost int64
		if strings.ToLower(item.ComputeType) == "baremetal" {
			computeCost = int64(itemNodes) * bmPrice
		} else if strings.ToLower(item.ComputeType) == "shared" {
			computeCost = 0
		} else {
			// VM
			computeCost = int64(itemVCPU) * vcpuPrice
		}

		// 2. Storage Cost
		var storageCost int64
		if strings.ToLower(item.StorageClass) == "tier 3" {
			storageCost = int64(itemStorage) * stgT3Price
		} else {
			storageCost = int64(itemStorage) * stgT2Price
		}

		// 3. Backup Cost
		var backupCost int64
		if item.BackupEnabled {
			backupCost = int64(itemStorage) * nbuPrice
		}

		// 4. License Cost
		var licenseCost int64
		var itemWarnings []string
		licenseProductName := "Community (Free OSS)"
		licenseKind := strings.ToLower(item.LicenseKind)
		if licenseKind == "" {
			licenseKind = "free"
		}

		if licenseKind == "enterprise" {
			// Find enterprise license option from component or lookup SKU
			var entSKU string
			var entProductName string
			var entMetric string

			if hasComp {
				for _, lic := range comp.Licenses {
					if strings.ToLower(lic.Kind) == "enterprise" {
						entSKU = lic.PricingSKU
						entProductName = lic.ProductName
						entMetric = lic.Metric
						break
					}
				}
			}

			if entProductName != "" {
				licenseProductName = entProductName
			} else {
				licenseProductName = "Enterprise License"
			}

			if entSKU != "" {
				pItem, exists := e.pricingMap[entSKU]
				if exists {
					if pItem.UnitPriceIDR != nil && *pItem.UnitPriceIDR > 0 {
						unitRate := *pItem.UnitPriceIDR
						switch entMetric {
						case "per_shard":
							licenseCost = int64(itemNodes) * unitRate
							if itemNodes < 3 {
								warnMsg := fmt.Sprintf("%s (%s): Redis Enterprise Cluster memerlukan minimal 3 node untuk Quorum High Availability", licenseProductName, compName)
								itemWarnings = append(itemWarnings, warnMsg)
								if !warningSet[warnMsg] {
									warnings = append(warnings, warnMsg)
									warningSet[warnMsg] = true
								}
							}
						case "per_node", "per_worker", "per_instance", "per_node_db":
							licenseCost = int64(itemNodes) * unitRate
						case "per_5_pods":
							factor := (itemNodes + 4) / 5
							if factor < 1 {
								factor = 1
							}
							licenseCost = int64(factor) * unitRate
						case "per_core":
							licenseCost = int64(itemVCPU) * unitRate
						case "per_4vcpu":
							factor := (itemVCPU + 3) / 4
							licenseCost = int64(factor) * unitRate
						default:
							licenseCost = int64(itemNodes) * unitRate
						}
					} else {
						// Price is null or 0 (Unconfirmed STA price)
						warnMsg := fmt.Sprintf("%s (%s): Estimasi harga belum dikonfirmasi oleh Tim STA/Procurement", licenseProductName, compName)
						itemWarnings = append(itemWarnings, warnMsg)
						if !warningSet[warnMsg] {
							warnings = append(warnings, warnMsg)
							warningSet[warnMsg] = true
						}
					}
				}
			}
		}

		// Item Total
		itemTotal := computeCost + storageCost + backupCost + licenseCost

		totalNodes += itemNodes
		totalVCPU += itemVCPU
		totalRAM += itemRAM
		totalStorage += itemStorage

		grandCompute += computeCost
		grandStorage += storageCost
		grandBackup += backupCost
		grandLicense += licenseCost

		itemCosts = append(itemCosts, model.QuoteItemCost{
			ComponentID:        item.ComponentID,
			ComponentName:      compName,
			LayerName:          layerName,
			NodeCount:          itemNodes,
			TotalVCPU:          itemVCPU,
			TotalRAMGB:         itemRAM,
			TotalStorageGB:     itemStorage,
			ComputeCostAnnual:  computeCost,
			StorageCostAnnual:  storageCost,
			BackupCostAnnual:   backupCost,
			LicenseCostAnnual:  licenseCost,
			TotalCostAnnual:    itemTotal,
			LicenseProductName: licenseProductName,
			LicenseKind:        licenseKind,
			UnconfirmedPrices:  itemWarnings,
		})
	}

	// Apply DRC multiplier if selected
	drcMultiplier := int64(1)
	if req.IncludeDRC {
		drcMultiplier = 2
		grandCompute *= drcMultiplier
		grandStorage *= drcMultiplier
		grandBackup *= drcMultiplier
		// Note: Software licenses in Enterprise Architecture are negotiated based on active/passive or dual-DC policy.
		// For transparency, we double infrastructure costs while keeping license per core/cluster rule.
	}

	grandTotalAnnual := grandCompute + grandStorage + grandBackup + grandLicense

	return model.QuoteResponse{
		Environment:           req.Environment,
		IncludeDRC:            req.IncludeDRC,
		TotalNodes:            totalNodes * int(drcMultiplier),
		TotalVCPU:             totalVCPU * int(drcMultiplier),
		TotalRAMGB:            totalRAM * int(drcMultiplier),
		TotalStorageGB:        totalStorage * int(drcMultiplier),
		SubtotalComputeAnnual: grandCompute,
		SubtotalStorageAnnual: grandStorage,
		SubtotalBackupAnnual:  grandBackup,
		SubtotalLicenseAnnual: grandLicense,
		GrandTotalAnnual:      grandTotalAnnual,
		TCO3Year:              grandTotalAnnual * 3,
		TCO5Year:              grandTotalAnnual * 5,
		Items:                 itemCosts,
		Warnings:              warnings,
	}
}
