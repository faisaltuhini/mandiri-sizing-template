package model

import (
	"time"

	"github.com/google/uuid"
)

// PriceListVersion represents a versioned pricing catalog
type PriceListVersion struct {
	ID            uuid.UUID     `gorm:"primaryKey" json:"id"`
	Name          string        `gorm:"type:varchar(255);not null" json:"name"`
	EffectiveDate time.Time     `json:"effective_date"`
	Status        string        `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, active, archived
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Items         []PricingItem `gorm:"foreignKey:PriceListVersionID" json:"items,omitempty"`
}

// PricingItem represents an infrastructure or license unit price
type PricingItem struct {
	ID                 uuid.UUID `gorm:"primaryKey" json:"id"`
	PriceListVersionID uuid.UUID `gorm:"not null;index:idx_pricing_version" json:"price_list_version_id"`
	SKU                string    `gorm:"type:varchar(100);not null;index:idx_pricing_sku;index:idx_pricing_cat_sku,priority:2" json:"sku"`
	Category           string    `gorm:"type:varchar(50);not null;index:idx_pricing_category;index:idx_pricing_cat_sku,priority:1" json:"category"` // compute, storage, backup, baremetal, license
	Name               string    `gorm:"type:varchar(255);not null;index:idx_pricing_name" json:"name"`
	Environment        string    `gorm:"type:varchar(50);not null;index:idx_pricing_env" json:"environment"` // Production, Dev / Non-Prod, All
	Unit               string    `gorm:"type:varchar(50);not null" json:"unit"`                               // vCPU / Thn, GB / Thn, Shard / Thn, Node / Thn, Unit
	UnitPriceIDR       *int64    `json:"unit_price_idr"`                                                      // Nullable for unconfirmed STA prices
	IsConfirmed        bool      `gorm:"default:true;index:idx_pricing_confirmed" json:"is_confirmed"`
	Source             string    `gorm:"type:varchar(255)" json:"source"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CatalogComponent represents one of the 13 Enterprise Architecture layers
type CatalogComponent struct {
	ID            uuid.UUID           `gorm:"primaryKey" json:"id"`
	LayerNo       int                 `gorm:"not null;index:idx_catalog_layer_no" json:"layer_no"`
	LayerName     string              `gorm:"type:varchar(100);not null;index:idx_catalog_layer_name" json:"layer_name"`
	ComponentName string              `gorm:"type:varchar(255);not null;index:idx_catalog_comp_name" json:"component_name"`
	Description   string              `gorm:"type:text" json:"description"`
	Notes         string              `gorm:"type:text" json:"notes"`
	ContractRef   string              `gorm:"type:varchar(255)" json:"contract_ref,omitempty"`
	SortOrder     int                 `gorm:"default:0;index:idx_catalog_sort_order" json:"sort_order"`
	Baselines     []ComponentBaseline `gorm:"foreignKey:ComponentID" json:"baselines,omitempty"`
	Licenses      []LicenseOption     `gorm:"foreignKey:ComponentID" json:"licenses,omitempty"`
}

// ComponentBaseline represents default baseline specs for moderate vs critical
type ComponentBaseline struct {
	ID               uuid.UUID `gorm:"primaryKey" json:"id"`
	ComponentID      uuid.UUID `gorm:"not null;index:idx_baseline_comp;index:idx_baseline_comp_tier,priority:1" json:"component_id"`
	Tier             string    `gorm:"type:varchar(50);not null;index:idx_baseline_comp_tier,priority:2" json:"tier"` // moderate, critical
	ComputeType      string    `gorm:"type:varchar(50);not null" json:"compute_type"`                                  // vm, baremetal, shared
	NodeCount        int       `gorm:"not null" json:"node_count"`
	VCPUPerNode      int       `gorm:"not null" json:"vcpu_per_node"`
	RAMGBPerNode     int       `gorm:"not null" json:"ram_gb_per_node"`
	StorageGBPerNode int       `gorm:"not null" json:"storage_gb_per_node"`
	StorageClass     string    `gorm:"type:varchar(50);default:'Tier 2'" json:"storage_class"` // Tier 2, Tier 3, Tier 1
	VCPUMin          int       `gorm:"default:1" json:"vcpu_min"`
	VCPUMax          int       `gorm:"default:256" json:"vcpu_max"`
	RAMMin           int       `gorm:"default:1" json:"ram_min"`
	RAMMax           int       `gorm:"default:2048" json:"ram_max"`
}

// LicenseOption defines free community vs enterprise licensed options
type LicenseOption struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	ComponentID uuid.UUID `gorm:"not null;index:idx_licopt_comp;index:idx_licopt_comp_kind,priority:1" json:"component_id"`
	Kind        string    `gorm:"type:varchar(50);not null;index:idx_licopt_comp_kind,priority:2" json:"kind"` // free, enterprise
	ProductName string    `gorm:"type:varchar(255);not null" json:"product_name"`
	Vendor      string    `gorm:"type:varchar(100);index:idx_licopt_vendor" json:"vendor"`
	PricingSKU  string    `gorm:"type:varchar(100);index:idx_licopt_sku" json:"pricing_sku,omitempty"`
	Metric      string    `gorm:"type:varchar(50);not null" json:"metric"` // free, per_node, per_instance, per_core, per_4vcpu, per_shard
}

// LicenseRule defines business rules for software licenses (e.g. per 5 pods, per node db, per shard with min cluster)
type LicenseRule struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	RuleName       string    `gorm:"type:varchar(255);not null;index:idx_license_rule_name;index:idx_license_cat_name,priority:2" json:"rule_name"` // e.g. "AppDynamics APM Microservices (5-pack)"
	Vendor         string    `gorm:"type:varchar(100);not null;index:idx_license_vendor" json:"vendor"`                                             // Cisco / AppDynamics, Redis Ltd, Red Hat, Oracle, F5
	Category       string    `gorm:"type:varchar(50);default:'apm';index:idx_license_category;index:idx_license_cat_name,priority:1" json:"category"` // apm, database, cache, platform, ingress
	MetricType     string    `gorm:"type:varchar(50);not null;index:idx_license_metric" json:"metric_type"`                                         // per_5_pods, per_node_db, per_shard, per_worker, per_core, per_instance
	PackSize       int       `gorm:"default:1" json:"pack_size"`                                                                                     // e.g. 5 for per_5_pods, 1 for single unit
	MinClusterNode int       `gorm:"default:1" json:"min_cluster_node"`                                                                              // e.g. 3 for Redis
	TargetLayer    string    `gorm:"type:varchar(100);index:idx_license_layer" json:"target_layer"`                                                  // App Container, Core Database, Distributed Cache, Web Ingress
	PricingSKU     string    `gorm:"type:varchar(100);index:idx_license_sku" json:"pricing_sku"`                                                     // e.g. APPDYN_POD_PROD, APPDYN_DB_PROD, LIC_REDIS_PROD, LIC_OCP_PROD
	UnitPriceIDR   *int64    `json:"unit_price_idr"`                                                                                                 // e.g. 8392364 for AppDyn 5-pack, 7514100 for DB Agent
	Description    string    `gorm:"type:text" json:"description"`
	IsActive       bool      `gorm:"default:true;index:idx_license_active" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Scenario represents a saved sizing configuration
type Scenario struct {
	ID                 uuid.UUID      `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"type:varchar(255);not null" json:"name"`
	ProjectName        string         `gorm:"type:varchar(255)" json:"project_name"`
	PriceListVersionID uuid.UUID      `gorm:"not null" json:"price_list_version_id"`
	Environment        string         `gorm:"type:varchar(50);default:'Production'" json:"environment"`
	IncludeDRC         bool           `gorm:"default:true" json:"include_drc"`
	TotalAnnualIDR     int64          `json:"total_annual_idr"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Items              []ScenarioItem `gorm:"foreignKey:ScenarioID" json:"items,omitempty"`
}

// ScenarioItem represents one configured component inside a scenario
type ScenarioItem struct {
	ID               uuid.UUID `gorm:"primaryKey" json:"id"`
	ScenarioID       uuid.UUID `gorm:"not null;index" json:"scenario_id"`
	ComponentID      uuid.UUID `gorm:"not null" json:"component_id"`
	Tier             string    `gorm:"type:varchar(50);not null" json:"tier"`
	ComputeType      string    `gorm:"type:varchar(50);not null" json:"compute_type"`
	NodeCount        int       `gorm:"not null" json:"node_count"`
	VCPUPerNode      int       `gorm:"not null" json:"vcpu_per_node"`
	RAMGBPerNode     int       `gorm:"not null" json:"ram_gb_per_node"`
	StorageGBPerNode int       `gorm:"not null" json:"storage_gb_per_node"`
	StorageClass     string    `gorm:"type:varchar(50)" json:"storage_class"`
	BackupEnabled    bool      `gorm:"default:true" json:"backup_enabled"`
	LicenseOptionID  uuid.UUID `json:"license_option_id"`
}

// --- QUOTE ENGINE DTOs ---

type QuoteItemRequest struct {
	ComponentID      string `json:"component_id"`
	Tier             string `json:"tier"`         // moderate, critical
	ComputeType      string `json:"compute_type"` // vm, baremetal
	NodeCount        int    `json:"node_count"`
	VCPUPerNode      int    `json:"vcpu_per_node"`
	RAMGBPerNode     int    `json:"ram_gb_per_node"`
	StorageGBPerNode int    `json:"storage_gb_per_node"`
	StorageClass     string `json:"storage_class"` // Tier 2, Tier 3, Tier 1
	BackupEnabled    bool   `json:"backup_enabled"`
	LicenseKind      string `json:"license_kind"` // free, enterprise
	LicenseOptionID  string `json:"license_option_id"`
}

type QuoteRequest struct {
	PriceListVersionID string             `json:"price_list_version_id,omitempty"`
	Environment        string             `json:"environment"` // Production, Dev / Non-Prod
	IncludeDRC         bool               `json:"include_drc"`
	Items              []QuoteItemRequest `json:"items"`
}

type QuoteItemCost struct {
	ComponentID        string   `json:"component_id"`
	ComponentName      string   `json:"component_name"`
	LayerName          string   `json:"layer_name"`
	NodeCount          int      `json:"node_count"`
	TotalVCPU          int      `json:"total_vcpu"`
	TotalRAMGB         int      `json:"total_ram_gb"`
	TotalStorageGB     int      `json:"total_storage_gb"`
	ComputeCostAnnual  int64    `json:"compute_cost_annual"`
	StorageCostAnnual  int64    `json:"storage_cost_annual"`
	BackupCostAnnual   int64    `json:"backup_cost_annual"`
	LicenseCostAnnual  int64    `json:"license_cost_annual"`
	TotalCostAnnual    int64    `json:"total_cost_annual"`
	LicenseProductName string   `json:"license_product_name"`
	LicenseKind        string   `json:"license_kind"`
	UnconfirmedPrices  []string `json:"unconfirmed_prices,omitempty"`
}

type QuoteResponse struct {
	Environment           string          `json:"environment"`
	IncludeDRC            bool            `json:"include_drc"`
	TotalNodes            int             `json:"total_nodes"`
	TotalVCPU             int             `json:"total_vcpu"`
	TotalRAMGB            int             `json:"total_ram_gb"`
	TotalStorageGB        int             `json:"total_storage_gb"`
	SubtotalComputeAnnual int64           `json:"subtotal_compute_annual"`
	SubtotalStorageAnnual int64           `json:"subtotal_storage_annual"`
	SubtotalBackupAnnual  int64           `json:"subtotal_backup_annual"`
	SubtotalLicenseAnnual int64           `json:"subtotal_license_annual"`
	GrandTotalAnnual      int64           `json:"grand_total_annual"`
	TCO3Year              int64           `json:"tco_3year"`
	TCO5Year              int64           `json:"tco_5year"`
	Items                 []QuoteItemCost `json:"items"`
	Warnings              []string        `json:"warnings,omitempty"`
}

// Backward-compatibility legacy DTOs
type CalculateRequest struct {
	DailyTransactions float64 `json:"daily_transactions"`
	PeakRatio         float64 `json:"peak_ratio"`
	PeakHours         float64 `json:"peak_hours"`
	PodTPSCapacity    float64 `json:"pod_tps_capacity"`
	CoresPerPod       float64 `json:"cores_per_pod"`
	WorkerCores       float64 `json:"worker_cores"`
	AvgTxKB           float64 `json:"avg_tx_kb"`
	DaysRetention     int     `json:"days_retention"`
	AnnualGrowth      float64 `json:"annual_growth"`
}

type YearlyProjection struct {
	Year        int     `json:"year"`
	TotalTx     float64 `json:"total_tx"`
	PeakTPS     float64 `json:"peak_tps"`
	StorageGB   float64 `json:"storage_gb"`
	WorkerNodes int     `json:"worker_nodes"`
}

type CalculateResponse struct {
	PeakTPS          float64            `json:"peak_tps"`
	TotalPods        int                `json:"total_pods"`
	RequiredCores    float64            `json:"required_cores"`
	WorkerNodeQty    int                `json:"worker_node_qty"`
	DBStorage5YrGB   float64            `json:"db_storage_5yr_gb"`
	NetBackupGB      float64            `json:"netbackup_gb"`
	YearlyProjection []YearlyProjection `json:"yearly_projection"`
}

type BOQEstimateItem struct {
	No        int     `json:"no"`
	Category  string  `json:"category"`
	ItemName  string  `json:"item_name"`
	Env       string  `json:"env"`
	Qty       int     `json:"qty"`
	VCPU      int     `json:"vcpu"`
	RAMGB     int     `json:"ram_gb"`
	StorageGB int     `json:"storage_gb"`
	CostVCPU  float64 `json:"cost_vcpu"`
	CostSTG   float64 `json:"cost_stg"`
	CostNBU   float64 `json:"cost_nbu"`
	CostBM    float64 `json:"cost_bm"`
	TotalCost float64 `json:"total_cost"`
}

type BOQEstimateResponse struct {
	Tier                string            `json:"tier"`
	TotalVM             int               `json:"total_vm"`
	TotalVCPU           int               `json:"total_vcpu"`
	TotalRAMGB          int               `json:"total_ram_gb"`
	TotalStorageGB      int               `json:"total_storage_gb"`
	AnnualHardwareCost  float64           `json:"annual_hardware_cost"`
	AnnualBackupCost    float64           `json:"annual_backup_cost"`
	AnnualBareMetalCost float64           `json:"annual_bare_metal_cost"`
	GrandTotalAnnual    float64           `json:"grand_total_annual"`
	Items               []BOQEstimateItem `json:"items"`
}

type LegacyPricingItem struct {
	No          int     `json:"no"`
	Category    string  `json:"category"`
	Item        string  `json:"item"`
	Environment string  `json:"environment"`
	Unit        string  `json:"unit"`
	PriceAnnual float64 `json:"price_annual"`
}

type KPISummary struct {
	ID     string `json:"id"`
	Tier   string `json:"tier"`
	Title  string `json:"title"`
	Value  string `json:"value"`
	Detail string `json:"detail"`
	Bg     string `json:"bg"`
	Accent string `json:"accent"`
}

type TechStackSpec struct {
	Nodes string `json:"nodes"`
	Type  string `json:"type"`
	Spec  string `json:"spec"`
}

type LicenseSpec struct {
	Free       string `json:"free"`
	Enterprise string `json:"enterprise"`
}

type TechStackItem struct {
	ID          string        `json:"id"`
	No          int           `json:"no"`
	Layer       string        `json:"layer"`
	Component   string        `json:"component"`
	Icon        string        `json:"icon"`
	Moderate    TechStackSpec `json:"moderate"`
	Critical    TechStackSpec `json:"critical"`
	License     LicenseSpec   `json:"license"`
	Notes       string        `json:"notes"`
	ContractRef string        `json:"contractRef,omitempty"`
}

type ForumPoint struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

type ForumReference struct {
	Title  string       `json:"title"`
	Tier   string       `json:"tier"`
	Badge  string       `json:"badge"`
	Points []ForumPoint `json:"points"`
}

type PowerLawRequest struct {
	TrxDaily          float64 `json:"trx_daily"`
	PeakHours         float64 `json:"peak_hours"`
	ConcurrencyFactor float64 `json:"concurrency_factor"`
	MaxCpuUtil        float64 `json:"max_cpu_util"`
	WorkerNodeCores   float64 `json:"worker_node_cores"`
	PT1TPS            float64 `json:"pt1_tps"`
	PT2TPS            float64 `json:"pt2_tps"`
	PT1NodeCount      int     `json:"pt1_node_count"`
	PT1NodeCores      float64 `json:"pt1_node_cores"`
	PT1CpuPercent     float64 `json:"pt1_cpu_percent"`
	PT2NodeCount      int     `json:"pt2_node_count"`
	PT2NodeCores      float64 `json:"pt2_node_cores"`
	PT2CpuPercent     float64 `json:"pt2_cpu_percent"`
	AvgTxKB           float64 `json:"avg_tx_kb"`
	AnnualGrowth      float64 `json:"annual_growth"`
}

type PowerLawResponse struct {
	PeakTPS         float64   `json:"peak_tps"`
	PT1Cores        float64   `json:"pt1_cores"`
	PT2Cores        float64   `json:"pt2_cores"`
	ScalingB        float64   `json:"scaling_b"`
	ScalingA        float64   `json:"scaling_a"`
	PredictedCores  float64   `json:"predicted_cores"`
	AllocatedCores  float64   `json:"allocated_cores"`
	WorkerNodeCount int       `json:"worker_node_count"`
	DbDailyGB       float64   `json:"db_daily_gb"`
	Db5YearGB       float64   `json:"db_5year_gb"`
	DbYearlyGB      []float64 `json:"db_yearly_gb"`
	IsEfficient     bool      `json:"is_efficient"`
	ScalingTag      string    `json:"scaling_tag"`
}

type SolarwindsMetricSummary struct {
	MetricName string  `json:"metric_name"`
	Unit       string  `json:"unit"`
	Avg        float64 `json:"avg"`
	Max        float64 `json:"max"`
	P95        float64 `json:"p95"`
	P99        float64 `json:"p99"`
	Status     string  `json:"status"`
}

type TemplateFile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	FileName    string `json:"file_name"`
	SizeBytes   int64  `json:"size_bytes"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

