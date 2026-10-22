package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"mandiri/sizing-apps-backend/internal/data"
	"mandiri/sizing-apps-backend/internal/model"
)

type SizingService struct{}

func NewSizingService() *SizingService {
	return &SizingService{}
}

// GetTechStackItems returns the 13 layers of enterprise architecture baseline
func (s *SizingService) GetTechStackItems() []model.TechStackItem {
	return data.GetTechStackItems()
}

// GetKPISummaries returns executive KPIs
func (s *SizingService) GetKPISummaries() []model.KPISummary {
	return data.GetKPISummaries()
}

// CalculatePowerLaw performs non-linear power-law sizing extrapolation
func (s *SizingService) CalculatePowerLaw(req model.PowerLawRequest) model.PowerLawResponse {
	// Fallback defaults if zero
	if req.PeakHours <= 0 {
		req.PeakHours = 4.0
	}
	if req.ConcurrencyFactor <= 0 {
		req.ConcurrencyFactor = 1.5
	}
	if req.MaxCpuUtil <= 0 {
		req.MaxCpuUtil = 0.40 // 40% enterprise bank standard
	}
	if req.WorkerNodeCores <= 0 {
		req.WorkerNodeCores = 32.0
	}
	if req.PT1TPS <= 0 {
		req.PT1TPS = 400.0
	}
	if req.PT2TPS <= 0 {
		req.PT2TPS = 1200.0
	}
	if req.PT1NodeCount <= 0 {
		req.PT1NodeCount = 4
	}
	if req.PT1NodeCores <= 0 {
		req.PT1NodeCores = 16.0
	}
	if req.PT1CpuPercent <= 0 {
		req.PT1CpuPercent = 45.0
	}
	if req.PT2NodeCount <= 0 {
		req.PT2NodeCount = 4
	}
	if req.PT2NodeCores <= 0 {
		req.PT2NodeCores = 16.0
	}
	if req.PT2CpuPercent <= 0 {
		req.PT2CpuPercent = 85.0
	}

	// 1. Peak TPS calculation
	peakTPS := (req.TrxDaily / (req.PeakHours * 3600.0)) * req.ConcurrencyFactor

	// 2. Cores consumed during PT 1 & PT 2
	pt1Cores := float64(req.PT1NodeCount) * req.PT1NodeCores * (req.PT1CpuPercent / 100.0)
	pt2Cores := float64(req.PT2NodeCount) * req.PT2NodeCores * (req.PT2CpuPercent / 100.0)

	// 3. Power-law scaling exponents (b) and constant (a)
	// b = ln(C2/C1) / ln(TPS2/TPS1)
	scalingB := 1.0
	if req.PT2TPS > req.PT1TPS && pt2Cores > 0 && pt1Cores > 0 {
		scalingB = math.Log(pt2Cores/pt1Cores) / math.Log(req.PT2TPS/req.PT1TPS)
	}

	// a = C1 / (TPS1 ^ b)
	scalingA := pt1Cores / math.Pow(req.PT1TPS, scalingB)

	// 4. Predicted raw cores = a * (PeakTPS ^ b)
	predictedCores := scalingA * math.Pow(peakTPS, scalingB)
	if math.IsNaN(predictedCores) || math.IsInf(predictedCores, 0) || predictedCores < 1.0 {
		predictedCores = 1.0
	}

	// 5. OpenShift Worker Node sizing
	allocatedCores := predictedCores / req.MaxCpuUtil
	workerNodeCount := int(math.Ceil(allocatedCores / req.WorkerNodeCores))
	if workerNodeCount < 2 {
		workerNodeCount = 2 // Minimum HA 2 Nodes
	}

	// 6. Database Storage 5-Year Projection
	if req.AvgTxKB <= 0 {
		req.AvgTxKB = 12.0 // 12 KB per transaction default
	}
	dbDailyGB := (req.TrxDaily * req.AvgTxKB) / (1024.0 * 1024.0)

	growthRate := req.AnnualGrowth / 100.0
	yearlyProjections := make([]float64, 5)
	cumulative5YrGB := 0.0

	for yr := 1; yr <= 5; yr++ {
		// Annual raw storage + index/WAL multiplier (2.5x)
		annualGB := (dbDailyGB * 365.0 * 2.5) * math.Pow(1.0+growthRate, float64(yr-1))
		yearlyProjections[yr-1] = math.Round(annualGB*100) / 100
		cumulative5YrGB += annualGB
	}

	isEfficient := scalingB <= 1.10
	scalingTag := "✅ Efficient Linear Scaling (b <= 1.10)"
	if !isEfficient {
		scalingTag = "⚠️ Non-Linear Resource Contention (b > 1.10)"
	}

	return model.PowerLawResponse{
		PeakTPS:         math.Round(peakTPS*100) / 100,
		PT1Cores:        math.Round(pt1Cores*100) / 100,
		PT2Cores:        math.Round(pt2Cores*100) / 100,
		ScalingB:        math.Round(scalingB*10000) / 10000,
		ScalingA:        scalingA,
		PredictedCores:  math.Round(predictedCores*100) / 100,
		AllocatedCores:  math.Round(allocatedCores*100) / 100,
		WorkerNodeCount: workerNodeCount,
		DbDailyGB:       math.Round(dbDailyGB*100) / 100,
		Db5YearGB:       math.Round(cumulative5YrGB*100) / 100,
		DbYearlyGB:      yearlyProjections,
		IsEfficient:     isEfficient,
		ScalingTag:      scalingTag,
	}
}

// CalculateEAWG performs standard EAWG 5-year projections
func (s *SizingService) CalculateEAWG(req model.CalculateRequest) model.CalculateResponse {
	if req.PeakHours <= 0 {
		req.PeakHours = 4.0
	}
	if req.PeakRatio <= 0 {
		req.PeakRatio = 1.5
	}
	if req.PodTPSCapacity <= 0 {
		req.PodTPSCapacity = 150.0
	}
	if req.CoresPerPod <= 0 {
		req.CoresPerPod = 2.0
	}
	if req.WorkerCores <= 0 {
		req.WorkerCores = 32.0
	}
	if req.AvgTxKB <= 0 {
		req.AvgTxKB = 12.0
	}

	peakTPS := (req.DailyTransactions * req.PeakRatio) / (req.PeakHours * 3600.0)
	totalPods := int(math.Ceil(peakTPS / req.PodTPSCapacity))
	if totalPods < 2 {
		totalPods = 2
	}

	requiredCores := float64(totalPods) * req.CoresPerPod
	workerNodeQty := int(math.Ceil(requiredCores / (req.WorkerCores * 0.40))) // 40% threshold
	if workerNodeQty < 2 {
		workerNodeQty = 2
	}

	dailyGB := (req.DailyTransactions * req.AvgTxKB) / (1024.0 * 1024.0)
	growthRate := req.AnnualGrowth / 100.0

	projections := make([]model.YearlyProjection, 5)
	cumulativeStorage := 0.0

	for yr := 1; yr <= 5; yr++ {
		multiplier := math.Pow(1.0+growthRate, float64(yr-1))
		yrTx := req.DailyTransactions * 365.0 * multiplier
		yrPeakTPS := peakTPS * multiplier
		yrStorageGB := (dailyGB * 365.0 * 2.5) * multiplier
		cumulativeStorage += yrStorageGB

		yrPods := int(math.Ceil(yrPeakTPS / req.PodTPSCapacity))
		if yrPods < 2 {
			yrPods = 2
		}
		yrCores := float64(yrPods) * req.CoresPerPod
		yrWorkers := int(math.Ceil(yrCores / (req.WorkerCores * 0.40)))
		if yrWorkers < 2 {
			yrWorkers = 2
		}

		projections[yr-1] = model.YearlyProjection{
			Year:        2024 + yr,
			TotalTx:     math.Round(yrTx),
			PeakTPS:     math.Round(yrPeakTPS*100) / 100,
			StorageGB:   math.Round(yrStorageGB*100) / 100,
			WorkerNodes: yrWorkers,
		}
	}

	netBackupGB := cumulativeStorage * 1.5 // 1.5x backup retention

	return model.CalculateResponse{
		PeakTPS:          math.Round(peakTPS*100) / 100,
		TotalPods:        totalPods,
		RequiredCores:    math.Round(requiredCores*100) / 100,
		WorkerNodeQty:    workerNodeQty,
		DBStorage5YrGB:   math.Round(cumulativeStorage*100) / 100,
		NetBackupGB:      math.Round(netBackupGB*100) / 100,
		YearlyProjection: projections,
	}
}

// AnalyzeSolarwindsCSV reads raw telemetry CSV and generates analytics
func (s *SizingService) AnalyzeSolarwindsCSV(metricName, unit string, r io.Reader) (model.SolarwindsMetricSummary, error) {
	csvReader := csv.NewReader(r)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return model.SolarwindsMetricSummary{}, fmt.Errorf("failed to read csv: %w", err)
	}

	var values []float64
	for i, row := range records {
		if i == 0 || len(row) == 0 {
			continue // skip header
		}
		// Try last column or 2nd column as numeric value
		valStr := strings.TrimSpace(row[len(row)-1])
		val, err := strconv.ParseFloat(valStr, 64)
		if err == nil {
			values = append(values, val)
		}
	}

	if len(values) == 0 {
		return model.SolarwindsMetricSummary{
			MetricName: metricName,
			Unit:       unit,
			Status:     "NO_DATA",
		}, nil
	}

	sort.Float64s(values)

	total := 0.0
	maxVal := values[len(values)-1]
	for _, v := range values {
		total += v
	}
	avgVal := total / float64(len(values))

	p95Idx := int(float64(len(values)) * 0.95)
	if p95Idx >= len(values) {
		p95Idx = len(values) - 1
	}
	p95Val := values[p95Idx]

	p99Idx := int(float64(len(values)) * 0.99)
	if p99Idx >= len(values) {
		p99Idx = len(values) - 1
	}
	p99Val := values[p99Idx]

	status := "NORMAL"
	if p95Val > 85.0 {
		status = "CRITICAL"
	} else if p95Val > 70.0 {
		status = "WARNING"
	}

	return model.SolarwindsMetricSummary{
		MetricName: metricName,
		Unit:       unit,
		Avg:        math.Round(avgVal*100) / 100,
		Max:        math.Round(maxVal*100) / 100,
		P95:        math.Round(p95Val*100) / 100,
		P99:        math.Round(p99Val*100) / 100,
		Status:     status,
	}, nil
}
