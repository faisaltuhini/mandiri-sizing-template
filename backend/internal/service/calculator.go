package service

import (
	"math"

	"mandiri/sizing-apps-backend/internal/data"
	"mandiri/sizing-apps-backend/internal/model"
)

type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

// CalculateCapacity computes capacity planning metrics based on Bank Mandiri EAWG standards
func (s *CalculatorService) CalculateCapacity(req model.CalculateRequest) model.CalculateResponse {
	// Defaults if 0
	if req.DailyTransactions <= 0 {
		req.DailyTransactions = 10000000 // 10 million transactions
	}
	if req.PeakHours <= 0 {
		req.PeakHours = 4.0
	}
	if req.PeakRatio <= 0 {
		req.PeakRatio = 2.5
	}
	if req.PodTPSCapacity <= 0 {
		req.PodTPSCapacity = 150.0
	}
	if req.CoresPerPod <= 0 {
		req.CoresPerPod = 2.0
	}
	if req.WorkerCores <= 0 {
		req.WorkerCores = 64.0
	}
	if req.AvgTxKB <= 0 {
		req.AvgTxKB = 4.5
	}
	if req.DaysRetention <= 0 {
		req.DaysRetention = 1825 // 5 years
	}
	if req.AnnualGrowth <= 0 {
		req.AnnualGrowth = 0.20 // 20% annual volume growth
	}

	// 1. Peak TPS
	avgTps := req.DailyTransactions / (req.PeakHours * 3600.0)
	peakTPS := avgTps * req.PeakRatio

	// 2. Pods & Worker Nodes
	totalPods := int(math.Ceil(peakTPS / req.PodTPSCapacity))
	if totalPods < 2 {
		totalPods = 2 // Minimum HA
	}
	requiredCores := float64(totalPods) * req.CoresPerPod

	// Usable cores per worker with 70% threshold & OS headroom
	usableCoresPerWorker := req.WorkerCores * 0.70
	workerNodeQty := int(math.Ceil(requiredCores / usableCoresPerWorker))
	if workerNodeQty < 2 {
		workerNodeQty = 2 // Minimum 2 workers
	}

	// 3. Storage Projection (5-Year Compounding)
	// Day 1 daily GB = (DailyTransactions * AvgTxKB) / (1024 * 1024)
	dailyGB := (req.DailyTransactions * req.AvgTxKB) / (1024.0 * 1024.0)

	var yearlyProjections []model.YearlyProjection
	cumulativeStorage := 0.0
	currentDailyTx := req.DailyTransactions

	for yr := 1; yr <= 5; yr++ {
		annualTx := currentDailyTx * 365.0
		yrDailyGB := (currentDailyTx * req.AvgTxKB) / (1024.0 * 1024.0)
		yrStorageGB := yrDailyGB * 365.0
		cumulativeStorage += yrStorageGB

		yrPeakTPS := (currentDailyTx / (req.PeakHours * 3600.0)) * req.PeakRatio
		yrPods := int(math.Ceil(yrPeakTPS / req.PodTPSCapacity))
		if yrPods < 2 {
			yrPods = 2
		}
		yrCores := float64(yrPods) * req.CoresPerPod
		yrWorkers := int(math.Ceil(yrCores / usableCoresPerWorker))
		if yrWorkers < 2 {
			yrWorkers = 2
		}

		yearlyProjections = append(yearlyProjections, model.YearlyProjection{
			Year:        yr,
			TotalTx:     annualTx,
			PeakTPS:     math.Round(yrPeakTPS*100) / 100,
			StorageGB:   math.Round(cumulativeStorage*10) / 10,
			WorkerNodes: yrWorkers,
		})

		currentDailyTx = currentDailyTx * (1.0 + req.AnnualGrowth)
	}

	dbStorage5YrGB := math.Round(cumulativeStorage*10) / 10
	netBackupGB := math.Round(dbStorage5YrGB*0.50*0.80*10) / 10

	_ = dailyGB

	return model.CalculateResponse{
		PeakTPS:          math.Round(peakTPS*100) / 100,
		TotalPods:        totalPods,
		RequiredCores:    requiredCores,
		WorkerNodeQty:    workerNodeQty,
		DBStorage5YrGB:   dbStorage5YrGB,
		NetBackupGB:      netBackupGB,
		YearlyProjection: yearlyProjections,
	}
}

// CalculateBOQ generates complete Bill of Quantities matching Excel formulas
func (s *CalculatorService) CalculateBOQ(tier string) model.BOQEstimateResponse {
	catalog := data.GetPricingCatalog()
	pricingMap := make(map[string]float64)
	for _, p := range catalog {
		key := p.Item + "_" + p.Environment
		pricingMap[key] = p.PriceAnnual
	}

	priceVcpuProd := 8250000.0
	priceVcpuDev := 7150000.0
	priceStgProd := 33000.0
	priceStgDev := 28600.0
	priceNBU := 27000.0
	priceBM := 794000000.0

	var items []model.BOQEstimateItem

	if tier == "moderate" {
		// Tier Moderate components (IHSP & RMT standard)
		modRows := []struct {
			name   string
			env    string
			qty    int
			vcpu   int
			ram    int
			stg    int
			isBM   bool
		}{
			{"NGINX Plus Gateway DC", "PROD", 2, 8, 16, 100, false},
			{"Red Hat OCP Worker DC", "PROD", 1, 64, 1024, 300, false},
			{"OCP Infra Router & Bastion DC", "PROD", 4, 8, 32, 150, false},
			{"Oracle DB Standalone DG Primary", "PROD", 1, 16, 32, 1000, false},
			{"Redis Caching Cluster DC", "PROD", 3, 8, 16, 200, false},
			{"Apache Kafka Cluster DC", "PROD", 3, 16, 64, 600, false},
			{"Elasticsearch & Logstash Shared DC", "PROD", 3, 8, 32, 500, false},
			{"SFTP Document Server DC", "PROD", 1, 8, 16, 1500, false},
			// DRC
			{"NGINX Plus Gateway DRC", "DRC", 2, 8, 16, 100, false},
			{"Red Hat OCP Worker DRC", "DRC", 1, 64, 1024, 300, false},
			{"Oracle DB Standalone DG Standby", "DRC", 1, 16, 32, 1000, false},
			{"Redis Caching Cluster DRC", "DRC", 3, 8, 16, 200, false},
			{"Apache Kafka Cluster DRC", "DRC", 3, 16, 64, 600, false},
			// Dev / PT
			{"NGINX Plus Dev Gateway", "DEV", 1, 8, 16, 100, false},
			{"OCP Worker Dev Pods", "DEV", 1, 16, 32, 200, false},
			{"Oracle DB Dev Instance", "DEV", 1, 16, 32, 500, false},
		}

		for idx, r := range modRows {
			totCpu := r.qty * r.vcpu
			totStg := r.qty * r.stg
			totRam := r.qty * r.ram

			isDev := r.env == "DEV"
			rateCpu := priceVcpuProd
			rateStg := priceStgProd
			if isDev {
				rateCpu = priceVcpuDev
				rateStg = priceStgDev
			}

			costCpu := float64(totCpu) * rateCpu
			costStg := float64(totStg) * rateStg
			costNbu := float64(totStg) * 0.50 * 0.80 * priceNBU
			costBM := 0.0
			if r.isBM {
				costBM = float64(r.qty) * priceBM
			}

			totCost := costCpu + costStg + costNbu + costBM

			items = append(items, model.BOQEstimateItem{
				No:        idx + 1,
				Category:  "Infrastructure VM",
				ItemName:  r.name,
				Env:       r.env,
				Qty:       r.qty,
				VCPU:      totCpu,
				RAMGB:     totRam,
				StorageGB: totStg,
				CostVCPU:  costCpu,
				CostSTG:   costStg,
				CostNBU:   costNbu,
				CostBM:    costBM,
				TotalCost: totCost,
			})
		}
	} else {
		// Tier Critical (Livin' Merchant standard with Baremetal & Dedicated Segregation)
		critRows := []struct {
			name   string
			env    string
			qty    int
			vcpu   int
			ram    int
			stg    int
			isBM   bool
		}{
			{"NGINX Plus Gateway DC", "PROD", 2, 8, 16, 100, false},
			{"OCP Bare Metal Worker 1 DC", "PROD", 1, 64, 1024, 300, true},
			{"OCP Bare Metal Worker 2 DC", "PROD", 1, 64, 1024, 300, true},
			{"Oracle DB RAC Node 1 (Baremetal)", "PROD", 1, 64, 1024, 500, true},
			{"Oracle DB RAC Node 2 (Baremetal)", "PROD", 1, 64, 1024, 500, true},
			{"Oracle ASM Shared SAN Storage 5.5TB", "PROD", 1, 0, 0, 5500, false},
			{"Redis Enterprise Multi-Shard DC", "PROD", 3, 8, 32, 200, false},
			{"Kafka Zookeeper DC", "PROD", 3, 4, 32, 100, false},
			{"Kafka Broker DC", "PROD", 3, 16, 32, 1000, false},
			{"Elasticsearch Master Nodes DC", "PROD", 3, 8, 32, 50, false},
			{"Elasticsearch Hot Data Nodes DC", "PROD", 2, 16, 32, 1900, false},
			{"Logstash Dedicated Shipper DC", "PROD", 2, 8, 32, 100, false},
			// DRC
			{"NGINX Plus Gateway DRC", "DRC", 2, 8, 16, 100, false},
			{"OCP Bare Metal Worker 1 DRC", "DRC", 1, 64, 1024, 300, true},
			{"OCP Bare Metal Worker 2 DRC", "DRC", 1, 64, 1024, 300, true},
			{"Oracle DB RAC Node 1 DRC (Baremetal)", "DRC", 1, 64, 1024, 500, true},
			{"Oracle DB RAC Node 2 DRC (Baremetal)", "DRC", 1, 64, 1024, 500, true},
			{"Oracle ASM Shared SAN DRC 5.5TB", "DRC", 1, 0, 0, 5500, false},
			{"Redis Enterprise Multi-Shard DRC", "DRC", 3, 8, 32, 200, false},
			{"Logstash Dedicated Shipper DRC", "DRC", 2, 8, 32, 100, false},
			// DEV
			{"NGINX Plus PT Gateway", "DEV", 2, 8, 16, 100, false},
			{"OCP Worker PT Pods", "DEV", 1, 16, 32, 200, false},
			{"Oracle DB PT Standalone VM", "DEV", 1, 16, 32, 1000, false},
		}

		for idx, r := range critRows {
			totCpu := r.qty * r.vcpu
			totStg := r.qty * r.stg
			totRam := r.qty * r.ram

			isDev := r.env == "DEV"
			rateCpu := priceVcpuProd
			rateStg := priceStgProd
			if isDev {
				rateCpu = priceVcpuDev
				rateStg = priceStgDev
			}

			costCpu := float64(totCpu) * rateCpu
			costStg := float64(totStg) * rateStg
			costNbu := float64(totStg) * 0.50 * 0.80 * priceNBU
			costBM := 0.0
			if r.isBM {
				costBM = float64(r.qty) * priceBM
			}

			totCost := costCpu + costStg + costNbu + costBM

			items = append(items, model.BOQEstimateItem{
				No:        idx + 1,
				Category:  "Infrastructure BM & VM",
				ItemName:  r.name,
				Env:       r.env,
				Qty:       r.qty,
				VCPU:      totCpu,
				RAMGB:     totRam,
				StorageGB: totStg,
				CostVCPU:  costCpu,
				CostSTG:   costStg,
				CostNBU:   costNbu,
				CostBM:    costBM,
				TotalCost: totCost,
			})
		}
	}

	// Compute totals
	totVm := 0
	totVcpu := 0
	totRam := 0
	totStg := 0
	totCostHw := 0.0
	totCostNbu := 0.0
	totCostBm := 0.0
	grandTot := 0.0

	for _, it := range items {
		totVm += it.Qty
		totVcpu += it.VCPU
		totRam += it.RAMGB
		totStg += it.StorageGB
		totCostHw += (it.CostVCPU + it.CostSTG)
		totCostNbu += it.CostNBU
		totCostBm += it.CostBM
		grandTot += it.TotalCost
	}

	return model.BOQEstimateResponse{
		Tier:               tier,
		TotalVM:            totVm,
		TotalVCPU:          totVcpu,
		TotalRAMGB:         totRam,
		TotalStorageGB:     totStg,
		AnnualHardwareCost: totCostHw,
		AnnualBackupCost:   totCostNbu,
		AnnualBareMetalCost: totCostBm,
		GrandTotalAnnual:   grandTot,
		Items:              items,
	}
}
