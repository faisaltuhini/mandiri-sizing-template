package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"enterprise/sizing-apps-backend/internal/model"
)

var DB *gorm.DB

// InitDB initializes connection to SQL Server (default) or PostgreSQL with automatic schema migration
func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default to SQL Server with user sa and password P@ssw0rd.1
		dsn = "sqlserver://sa:P%40ssw0rd.1@localhost:1433?database=sizing&trustServerCertificate=true"
	}

	var dialector gorm.Dialector

	if strings.HasPrefix(strings.ToLower(dsn), "postgres://") || strings.HasPrefix(strings.ToLower(dsn), "postgresql://") {
		dialector = postgres.Open(dsn)
		log.Println("🔌 Connecting to PostgreSQL...")
	} else {
		// Ensure SQL Server database exists before connecting to specific database
		ensureSQLServerDatabaseExists(dsn)
		dialector = sqlserver.Open(dsn)
		log.Println("🔌 Connecting to Microsoft SQL Server...")
	}

	database, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB handle: %w", err)
	}

	// High-Performance Connection Pool (eliminates TCP + TDS handshake latency on concurrent requests)
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(60)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// Auto-migrate schema (Cross-database compatible with rich indexes)
	err = database.AutoMigrate(
		&model.PriceListVersion{},
		&model.PricingItem{},
		&model.CatalogComponent{},
		&model.ComponentBaseline{},
		&model.LicenseOption{},
		&model.LicenseRule{},
		&model.Scenario{},
		&model.ScenarioItem{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	// Ensure fast SQL Server indexing for non-clustered queries
	ensurePerformanceIndexes(database)

	DB = database

	// Seed default catalog & pricing if empty
	if err := SeedMasterData(database); err != nil {
		log.Printf("⚠️ Warning seeding database: %v", err)
	}

	// Seed license rules & missing catalog SKUs if empty
	if err := SeedLicenseRules(database); err != nil {
		log.Printf("⚠️ Warning seeding license rules: %v", err)
	}

	log.Println("✅ Database connected, connection pool optimized (30 idle / 60 max), migrated, and indexed!")
	return database, nil
}

// ensurePerformanceIndexes creates composite and single-column nonclustered indexes if not already present
func ensurePerformanceIndexes(db *gorm.DB) {
	indexes := []string{
		"CREATE NONCLUSTERED INDEX idx_pricing_items_cat_sku ON pricing_items(category, sku)",
		"CREATE NONCLUSTERED INDEX idx_pricing_items_sku ON pricing_items(sku)",
		"CREATE NONCLUSTERED INDEX idx_license_rules_cat_name ON license_rules(category, rule_name)",
		"CREATE NONCLUSTERED INDEX idx_license_rules_sku ON license_rules(pricing_sku)",
		"CREATE NONCLUSTERED INDEX idx_catalog_components_layer ON catalog_components(layer_no)",
		"CREATE NONCLUSTERED INDEX idx_component_baselines_comp_tier ON component_baselines(component_id, tier)",
		"CREATE NONCLUSTERED INDEX idx_license_options_comp_kind ON license_options(component_id, kind)",
	}

	for _, stmt := range indexes {
		_ = db.Exec(stmt).Error
	}
}

// ensureSQLServerDatabaseExists connects to master to create the target database if missing
func ensureSQLServerDatabaseExists(targetDSN string) {
	// Parse database name from DSN
	dbName := "sizing"
	if idx := strings.Index(targetDSN, "database="); idx != -1 {
		sub := targetDSN[idx+9:]
		if ampIdx := strings.Index(sub, "&"); ampIdx != -1 {
			dbName = sub[:ampIdx]
		} else {
			dbName = sub
		}
	}

	if dbName == "master" || dbName == "" {
		return
	}

	// Build master DSN
	masterDSN := strings.Replace(targetDSN, "database="+dbName, "database=master", 1)
	masterDB, err := gorm.Open(sqlserver.Open(masterDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return // If master connection fails, let standard connection attempt proceed
	}

	var count int
	masterDB.Raw("SELECT count(*) FROM sys.databases WHERE name = ?", dbName).Scan(&count)
	if count == 0 {
		log.Printf("📦 Creating SQL Server database '%s'...", dbName)
		if err := masterDB.Exec(fmt.Sprintf("CREATE DATABASE [%s]", dbName)).Error; err != nil {
			log.Printf("⚠️ Warning creating database %s: %v", dbName, err)
		} else {
			log.Printf("✅ SQL Server database '%s' created successfully.", dbName)
		}
	}

	sqlDB, err := masterDB.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// SeedMasterData populates the 13 Enterprise Architecture layers & price catalog
func SeedMasterData(db *gorm.DB) error {
	var count int64
	db.Model(&model.PriceListVersion{}).Count(&count)
	if count > 0 {
		return nil // already seeded
	}

	log.Println("🌱 Seeding Enterprise Architecture Pricing Catalog & 13 Architecture Layers...")

	// 1. Create Active Master Price List Version
	versionID := uuid.New()
	version := model.PriceListVersion{
		ID:            versionID,
		Name:          "Katalog Resmi Enterprise Architecture (Template Sizing.xlsx)",
		EffectiveDate: time.Now(),
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&version).Error; err != nil {
		return err
	}

	// 2. Seed Pricing Items (Exact from Template Sizing.xlsx!Pricing)
	pricingSeed := []struct {
		SKU         string
		Category    string
		Name        string
		Env         string
		Unit        string
		Price       *int64
		IsConfirmed bool
	}{
		// Hardware
		{"VCPU_PROD", "compute", "Compute vCPU", "Production", "vCPU / Thn", int64Ptr(8250000), true},
		{"VCPU_DEV", "compute", "Compute vCPU", "Dev / Non-Prod", "vCPU / Thn", int64Ptr(7150000), true},
		{"STG_T2_PROD", "storage", "Storage SAN Tier 2", "Production", "GB / Thn", int64Ptr(33000), true},
		{"STG_T3_PROD", "storage", "Storage SAN Tier 3", "Production", "GB / Thn", int64Ptr(48000), true},
		{"STG_T2_DEV", "storage", "Storage SAN Tier 2", "Dev / Non-Prod", "GB / Thn", int64Ptr(28600), true},
		{"NBU_ALL", "backup", "Backup Netbackup (NBU)", "Production", "GB / Thn", int64Ptr(27000), true},
		{"RUBRIK_PROD", "backup", "Rubrik Backup", "Production", "GB / Thn", int64Ptr(35000), true},
		{"BM_STD_PROD", "baremetal", "Server Baremetal (Standard)", "Production", "Unit", int64Ptr(794000000), true},
		{"BM_STD_DEV", "baremetal", "Server Baremetal (Standard)", "Dev / Non-Prod", "Unit", int64Ptr(794000000), true},

		// Licenses
		{"LIC_NGINX_PROD", "license", "NGINX Plus (F5 Enterprise)", "Production", "Instance / Thn", nil, false},
		{"LIC_NGINX_DEV", "license", "NGINX Plus (F5 Enterprise)", "Dev / Non-Prod", "Instance / Thn", nil, false},
		{"LIC_OCP_PROD", "license", "Red Hat Openshift Container Platform", "Production", "Node / Thn", int64Ptr(240000000), true},
		{"LIC_OCP_DEV", "license", "Red Hat Openshift Container Platform", "Dev / Non-Prod", "Node / Thn", int64Ptr(134000000), true},
		{"LIC_REDIS_PROD", "license", "Redis Enterprise Software", "Production", "Shard / Thn", int64Ptr(125000000), true},
		{"LIC_REDIS_DEV", "license", "Redis Enterprise Software", "Dev / Non-Prod", "Shard / Thn", int64Ptr(75000000), true},
		{"LIC_KAFKA_PROD", "license", "Confluent Kafka Enterprise", "Production", "Broker / Thn", nil, false},
		{"LIC_KAFKA_DEV", "license", "Confluent Kafka Enterprise", "Dev / Non-Prod", "Broker / Thn", nil, false},
		{"LIC_ORACLE_PROD", "license", "Oracle Database Enterprise Edition", "Production", "Core / Thn", nil, false},
	}

	for _, p := range pricingSeed {
		item := model.PricingItem{
			ID:                 uuid.New(),
			PriceListVersionID: versionID,
			SKU:                p.SKU,
			Category:           p.Category,
			Name:               p.Name,
			Environment:        p.Env,
			Unit:               p.Unit,
			UnitPriceIDR:       p.Price,
			IsConfirmed:        p.IsConfirmed,
			Source:             "Template Sizing.xlsx!Pricing",
		}
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	// 3. Seed 13 Architecture Layers (CatalogComponents + Baselines + LicenseOptions)
	type layerDef struct {
		No         int
		Layer      string
		Component  string
		Notes      string
		ModNode    int
		ModType    string
		ModVCPU    int
		ModRAM     int
		ModSTG     int
		ModSTGCls  string
		CritNode   int
		CritType   string
		CritVCPU   int
		CritRAM    int
		CritSTG    int
		CritSTGCls string
		FreeLic    string
		EntLic     string
		EntVendor  string
		EntSKU     string
		EntMetric  string
	}

	layers := []layerDef{
		{
			No: 1, Layer: "Web Ingress", Component: "NGINX Web Server",
			Notes:   "Moderate: NGINX VM 8/16/100 (IHSP) / 8/16/200 (RMT). Critical: NGINX Plus VM 8-16 vCPU / 16-32 GB (Kontrak F5: TIF.DDL/GH.9076/2025).",
			ModNode: 2, ModType: "vm", ModVCPU: 8, ModRAM: 16, ModSTG: 200, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "vm", CritVCPU: 16, CritRAM: 32, CritSTG: 200, CritSTGCls: "Tier 2",
			FreeLic: "NGINX Community (oss)", EntLic: "NGINX Plus (F5 Enterprise)", EntVendor: "F5", EntSKU: "LIC_NGINX_PROD", EntMetric: "per_instance",
		},
		{
			No: 2, Layer: "App Container", Component: "Red Hat OCP Worker",
			Notes:   "Moderate: 1 OCP Worker VM (64/1024/1200) per DC. Critical: 2 Dedicated Bare Metal Workers per DC (09PWRKR01/02 & RPWRKR01/02, Kontrak Red Hat: TIF.DDL/GH.1993/2025).",
			ModNode: 1, ModType: "vm", ModVCPU: 64, ModRAM: 1024, ModSTG: 1200, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "baremetal", CritVCPU: 64, CritRAM: 1024, CritSTG: 300, CritSTGCls: "Tier 2",
			FreeLic: "OKD (Community / Non-Prod)", EntLic: "RH OCP Platform Plus BM Premium", EntVendor: "Red Hat", EntSKU: "LIC_OCP_PROD", EntMetric: "per_node",
		},
		{
			No: 3, Layer: "Container Router", Component: "OCP Infra Router & Bastion",
			Notes:   "Moderate: 3 Master (8/32/150) + 3 Infra Router (8/32/150) + 1 Bastion (4/8/100). Critical: Dedicated high-throughput router pods cross-DC.",
			ModNode: 7, ModType: "vm", ModVCPU: 8, ModRAM: 32, ModSTG: 150, ModSTGCls: "Tier 2",
			CritNode: 7, CritType: "vm", CritVCPU: 8, CritRAM: 32, CritSTG: 150, CritSTGCls: "Tier 2",
			FreeLic: "Included in OKD", EntLic: "Included in OCP Platform Plus", EntVendor: "Red Hat", EntSKU: "", EntMetric: "free",
		},
		{
			No: 4, Layer: "Core Database", Component: "Database Engine (Oracle)",
			Notes:   "Moderate: Single Node VM + Data Guard (16/32/1000) tanpa RAC. Option untuk efisiensi lisensi core. Critical: Oracle RAC Bare Metal 2-Node Core (High TPS & Zero Data Loss, Kontrak: TIF.DDL/SCD.3565/2025).",
			ModNode: 1, ModType: "vm", ModVCPU: 16, ModRAM: 32, ModSTG: 1000, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "baremetal", CritVCPU: 64, CritRAM: 1024, CritSTG: 1000, CritSTGCls: "Tier 3",
			FreeLic: "PostgreSQL / MySQL Community", EntLic: "Oracle DB Enterprise (RAC + ADG)", EntVendor: "Oracle", EntSKU: "LIC_ORACLE_PROD", EntMetric: "per_core",
		},
		{
			No: 5, Layer: "Database Storage", Component: "Oracle ASM Shared SAN",
			Notes:   "Moderate: SAN Tier 2 Virtual LUN 1.0 - 1.5 TB. Critical: SAN Tier 2 ASM 5.5 TB (Reco 500GB, Data 5000GB) pada node bare metal.",
			ModNode: 1, ModType: "shared", ModVCPU: 0, ModRAM: 0, ModSTG: 1500, ModSTGCls: "Tier 2",
			CritNode: 1, CritType: "baremetal", CritVCPU: 0, CritRAM: 0, CritSTG: 5500, CritSTGCls: "Tier 3",
			FreeLic: "SAN LUN Tier 2 Standard", EntLic: "SAN LUN Tier 3 + Oracle ASM", EntVendor: "Enterprise Storage", EntSKU: "", EntMetric: "free",
		},
		{
			No: 6, Layer: "Distributed Cache", Component: "Redis Caching Cluster",
			Notes:   "Moderate: 1-3 VM Redis Community (8/16/200). Critical: 3 VM Redis Enterprise Multi-Shard (01VLMBD01-03 8/32/200, Kontrak: TIF.DDL/GH.9060/2024).",
			ModNode: 2, ModType: "vm", ModVCPU: 8, ModRAM: 16, ModSTG: 200, ModSTGCls: "Tier 2",
			CritNode: 3, CritType: "vm", CritVCPU: 8, CritRAM: 32, CritSTG: 200, CritSTGCls: "Tier 2",
			FreeLic: "Redis Community (oss)", EntLic: "Redis Enterprise Software", EntVendor: "Redis Ltd", EntSKU: "LIC_REDIS_PROD", EntMetric: "per_shard",
		},
		{
			No: 7, Layer: "Cache Sentinel", Component: "Redis Sentinel Quorum",
			Notes:   "Moderate: 3 Small VM (2/4/50) quorum sentinel. Critical: 3 Small VM (4/8/100) integrated sentinel auto-failover cluster.",
			ModNode: 3, ModType: "vm", ModVCPU: 2, ModRAM: 4, ModSTG: 50, ModSTGCls: "Tier 2",
			CritNode: 3, CritType: "vm", CritVCPU: 4, CritRAM: 8, CritSTG: 100, CritSTGCls: "Tier 2",
			FreeLic: "Redis Sentinel Community", EntLic: "Redis Sentinel Enterprise", EntVendor: "Redis Ltd", EntSKU: "", EntMetric: "free",
		},
		{
			No: 8, Layer: "Message Bus", Component: "Confluent / Apache Kafka",
			Notes:   "Moderate: 3 VM Kafka (8/32/500). Critical: 3 VM Kafka (16/64/1000) High-Throughput streaming event logs.",
			ModNode: 3, ModType: "vm", ModVCPU: 8, ModRAM: 32, ModSTG: 500, ModSTGCls: "Tier 2",
			CritNode: 3, CritType: "vm", CritVCPU: 16, CritRAM: 64, CritSTG: 1000, CritSTGCls: "Tier 2",
			FreeLic: "Apache Kafka (OSS Community)", EntLic: "Confluent Enterprise Platform", EntVendor: "Confluent", EntSKU: "LIC_KAFKA_PROD", EntMetric: "per_node",
		},
		{
			No: 9, Layer: "Log Analytics", Component: "Elasticsearch Analytics Node",
			Notes:   "Moderate: 2 VM ELK (8/32/500). Critical: 3 VM ELK (16/64/2000) Hot-Warm retention log indexing.",
			ModNode: 2, ModType: "vm", ModVCPU: 8, ModRAM: 32, ModSTG: 500, ModSTGCls: "Tier 2",
			CritNode: 3, CritType: "vm", CritVCPU: 16, CritRAM: 64, CritSTG: 2000, CritSTGCls: "Tier 2",
			FreeLic: "OpenSearch / Elastic Basic", EntLic: "Elasticsearch Platinum Edition", EntVendor: "Elastic", EntSKU: "", EntMetric: "free",
		},
		{
			No: 10, Layer: "Log Shipper", Component: "Logstash Pipeline Shipper",
			Notes:   "Moderate: 1 VM Logstash (4/8/100). Critical: 2 VM Logstash (8/16/200) HA pipeline parser.",
			ModNode: 1, ModType: "vm", ModVCPU: 4, ModRAM: 8, ModSTG: 100, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "vm", CritVCPU: 8, CritRAM: 16, CritSTG: 200, CritSTGCls: "Tier 2",
			FreeLic: "Logstash Community", EntLic: "Included in Elastic Enterprise", EntVendor: "Elastic", EntSKU: "", EntMetric: "free",
		},
		{
			No: 11, Layer: "Log Dashboard", Component: "Kibana Observability Portal",
			Notes:   "Moderate: 1 VM Kibana (4/8/100). Critical: 2 VM Kibana (4/8/100) behind internal load balancer.",
			ModNode: 1, ModType: "vm", ModVCPU: 4, ModRAM: 8, ModSTG: 100, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "vm", CritVCPU: 4, CritRAM: 8, CritSTG: 100, CritSTGCls: "Tier 2",
			FreeLic: "Kibana Community", EntLic: "Included in Elastic Enterprise", EntVendor: "Elastic", EntSKU: "", EntMetric: "free",
		},
		{
			No: 12, Layer: "Document Storage", Component: "SFTP / Object Storage / RDM",
			Notes:   "Moderate: 1 VM SFTP (4/8/500). Critical: 2 VM Active-Passive SFTP + Dell ECS Object Storage (8/16/5000).",
			ModNode: 1, ModType: "vm", ModVCPU: 4, ModRAM: 8, ModSTG: 500, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "vm", CritVCPU: 8, CritRAM: 16, CritSTG: 5000, CritSTGCls: "Tier 2",
			FreeLic: "Linux OpenSSH SFTP", EntLic: "Dell EMC ECS Object Storage", EntVendor: "Dell", EntSKU: "", EntMetric: "free",
		},
		{
			No: 13, Layer: "Batch Engine", Component: "BMC Control-M / OCP Cron",
			Notes:   "Moderate: Control-M Agent on OCP VM. Critical: Dedicated Control-M Agent with High Availability on Baremetal.",
			ModNode: 1, ModType: "vm", ModVCPU: 4, ModRAM: 8, ModSTG: 100, ModSTGCls: "Tier 2",
			CritNode: 2, CritType: "vm", CritVCPU: 8, CritRAM: 16, CritSTG: 200, CritSTGCls: "Tier 2",
			FreeLic: "K8s Native CronJob", EntLic: "BMC Control-M Enterprise Agent", EntVendor: "BMC", EntSKU: "", EntMetric: "free",
		},
	}

	for _, l := range layers {
		compID := uuid.New()
		comp := model.CatalogComponent{
			ID:            compID,
			LayerNo:       l.No,
			LayerName:     l.Layer,
			ComponentName: l.Component,
			Description:   fmt.Sprintf("Standar Layer %d (%s) untuk arsitektur sistem enterprise", l.No, l.Layer),
			Notes:         l.Notes,
			SortOrder:     l.No,
		}
		if err := db.Create(&comp).Error; err != nil {
			return err
		}

		// Baseline Moderate
		modBase := model.ComponentBaseline{
			ID:               uuid.New(),
			ComponentID:      compID,
			Tier:             "moderate",
			ComputeType:      l.ModType,
			NodeCount:        l.ModNode,
			VCPUPerNode:      l.ModVCPU,
			RAMGBPerNode:     l.ModRAM,
			StorageGBPerNode: l.ModSTG,
			StorageClass:     l.ModSTGCls,
			VCPUMin:          1, VCPUMax: 128,
			RAMMin: 1, RAMMax: 2048,
		}
		db.Create(&modBase)

		// Baseline Critical
		critBase := model.ComponentBaseline{
			ID:               uuid.New(),
			ComponentID:      compID,
			Tier:             "critical",
			ComputeType:      l.CritType,
			NodeCount:        l.CritNode,
			VCPUPerNode:      l.CritVCPU,
			RAMGBPerNode:     l.CritRAM,
			StorageGBPerNode: l.CritSTG,
			StorageClass:     l.CritSTGCls,
			VCPUMin:          1, VCPUMax: 256,
			RAMMin: 1, RAMMax: 4096,
		}
		db.Create(&critBase)

		// License Free Option
		freeOpt := model.LicenseOption{
			ID:          uuid.New(),
			ComponentID: compID,
			Kind:        "free",
			ProductName: l.FreeLic,
			Vendor:      "Community / Open-Source",
			Metric:      "free",
		}
		db.Create(&freeOpt)

		// License Enterprise Option
		entOpt := model.LicenseOption{
			ID:          uuid.New(),
			ComponentID: compID,
			Kind:        "enterprise",
			ProductName: l.EntLic,
			Vendor:      l.EntVendor,
			PricingSKU:  l.EntSKU,
			Metric:      l.EntMetric,
		}
		db.Create(&entOpt)
	}

	log.Println("✅ 13 Architecture Layers and Pricing Items seeded successfully!")
	return nil
}

func int64Ptr(v int64) *int64 {
	return &v
}

// SeedLicenseRules populates standard license calculation rules and AppDynamics catalog SKUs
func SeedLicenseRules(db *gorm.DB) error {
	// 1. Ensure AppDynamics pricing items exist in active price list
	var appDynCount int64
	db.Model(&model.PricingItem{}).Where("sku = ?", "APPDYN_POD_PROD").Count(&appDynCount)
	if appDynCount == 0 {
		var activeVersion model.PriceListVersion
		if err := db.Where("status = ?", "active").First(&activeVersion).Error; err == nil {
			log.Println("🌱 Seeding missing AppDynamics SKUs into Pricing Items...")
			appDynItems := []model.PricingItem{
				{
					ID:                 uuid.New(),
					PriceListVersionID: activeVersion.ID,
					SKU:                "APPDYN_POD_PROD",
					Category:           "license",
					Name:               "AppDynamics APM Microservices (5-pack)",
					Environment:        "Production",
					Unit:               "5 Pods / Thn",
					UnitPriceIDR:       int64Ptr(8392364),
					IsConfirmed:        true,
					Source:             "Kontrak Cisco / AppDynamics Enterprise Architecture (TIF.DDL/GH.2410/2024)",
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				},
				{
					ID:                 uuid.New(),
					PriceListVersionID: activeVersion.ID,
					SKU:                "APPDYN_DB_PROD",
					Category:           "license",
					Name:               "AppDynamics Database Visibility Agent",
					Environment:        "Production",
					Unit:               "DB Node / Thn",
					UnitPriceIDR:       int64Ptr(7514100),
					IsConfirmed:        true,
					Source:             "Kontrak Cisco / AppDynamics Enterprise Architecture (TIF.DDL/GH.2410/2024)",
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				},
			}
			for _, item := range appDynItems {
				db.Create(&item)
			}
		}
	}

	// 2. Check if LicenseRules already seeded
	var ruleCount int64
	db.Model(&model.LicenseRule{}).Count(&ruleCount)
	if ruleCount > 0 {
		return nil
	}

	log.Println("🌱 Seeding Enterprise Software License Rules...")
	defaultRules := []model.LicenseRule{
		{
			ID:             uuid.New(),
			RuleName:       "AppDynamics APM Microservices (5-pack)",
			Vendor:         "Cisco / AppDynamics",
			Category:       "apm",
			MetricType:     "per_5_pods",
			PackSize:       5,
			MinClusterNode: 1,
			TargetLayer:    "App Container",
			PricingSKU:     "APPDYN_POD_PROD",
			UnitPriceIDR:   int64Ptr(8392364),
			Description:    "Monitoring microservices container application per kelipatan 5 pods. Formula: ceil(pods / 5) * Harga Satuan. Kontrak TIF.DDL/GH.2410/2024.",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			RuleName:       "AppDynamics Database Visibility Agent",
			Vendor:         "Cisco / AppDynamics",
			Category:       "database",
			MetricType:     "per_node_db",
			PackSize:       1,
			MinClusterNode: 1,
			TargetLayer:    "Core Database",
			PricingSKU:     "APPDYN_DB_PROD",
			UnitPriceIDR:   int64Ptr(7514100),
			Description:    "Monitoring performa database per DB Node / Instance. Formula: 1 agent per node database.",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			RuleName:       "Redis Enterprise Software",
			Vendor:         "Redis Ltd",
			Category:       "cache",
			MetricType:     "per_shard",
			PackSize:       1,
			MinClusterNode: 3,
			TargetLayer:    "Distributed Cache",
			PricingSKU:     "LIC_REDIS_PROD",
			UnitPriceIDR:   int64Ptr(125000000),
			Description:    "Lisensi Redis Enterprise dihitung per Shard data. Syarat arsitektur cluster HA mewajibkan minimal 3 node untuk quorum sentinel/raft failover.",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			RuleName:       "Red Hat OpenShift Platform Plus",
			Vendor:         "Red Hat",
			Category:       "platform",
			MetricType:     "per_worker",
			PackSize:       1,
			MinClusterNode: 1,
			TargetLayer:    "App Container",
			PricingSKU:     "LIC_OCP_PROD",
			UnitPriceIDR:   int64Ptr(240000000),
			Description:    "Lisensi RH OCP Platform Plus dihitung per Worker node. Master node, Infra router & Bastion tidak dikenakan biaya lisensi (Free of Charge).",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			RuleName:       "NGINX Plus Enterprise",
			Vendor:         "F5",
			Category:       "ingress",
			MetricType:     "per_instance",
			PackSize:       1,
			MinClusterNode: 1,
			TargetLayer:    "Web Ingress",
			PricingSKU:     "LIC_NGINX_PROD",
			UnitPriceIDR:   nil,
			Description:    "Lisensi NGINX Plus per instance/VM load balancer. Harga menunggu konfirmasi Tim STA / Procurement (Kontrak F5: TIF.DDL/GH.9076/2025).",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			RuleName:       "Oracle Database Enterprise Edition",
			Vendor:         "Oracle",
			Category:       "database",
			MetricType:     "per_core",
			PackSize:       1,
			MinClusterNode: 1,
			TargetLayer:    "Core Database",
			PricingSKU:     "LIC_ORACLE_PROD",
			UnitPriceIDR:   nil,
			Description:    "Lisensi Oracle DB Enterprise (RAC + Active Data Guard) dihitung per Processor Core fisik/vCPU dengan Core Factor Intel x86 (0.5).",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			RuleName:       "Confluent Kafka Enterprise",
			Vendor:         "Confluent",
			Category:       "streaming",
			MetricType:     "per_node",
			PackSize:       1,
			MinClusterNode: 3,
			TargetLayer:    "Message Bus",
			PricingSKU:     "LIC_KAFKA_PROD",
			UnitPriceIDR:   nil,
			Description:    "Platform Confluent Kafka Enterprise dihitung per Broker Node (Rekomendasi minimal 3 node untuk replication factor 3 & in-sync replicas).",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, rule := range defaultRules {
		if err := db.Create(&rule).Error; err != nil {
			return err
		}
	}

	log.Println("✅ Enterprise License Rules seeded successfully!")
	return nil
}

