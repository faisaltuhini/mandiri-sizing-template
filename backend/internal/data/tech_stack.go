package data

import "enterprise/sizing-apps-backend/internal/model"

// GetKPISummaries returns executive KPIs
func GetKPISummaries() []model.KPISummary {
	return []model.KPISummary{
		{
			ID:     "mod_compute",
			Tier:   "moderate",
			Title:  "Tier Moderate Compute",
			Value:  "14 VM Nodes",
			Detail: "112 vCPU | 256 GB RAM (IHSP & RMT)",
			Bg:     "#002060",
			Accent: "#38bdf8",
		},
		{
			ID:     "mod_storage",
			Tier:   "moderate",
			Title:  "Tier Moderate Storage",
			Value:  "4.500 GB (4.5 TB)",
			Detail: "SAN Tier 2 LUN Virtual Disk",
			Bg:     "#1B365D",
			Accent: "#60a5fa",
		},
		{
			ID:     "crit_compute",
			Tier:   "critical",
			Title:  "Tier Critical Compute",
			Value:  "28 Dedicated Nodes",
			Detail: "384 vCPU | 4.8 TB RAM (incl 4 Bare Metal)",
			Bg:     "#78350f",
			Accent: "#fbbf24",
		},
		{
			ID:     "crit_storage",
			Tier:   "critical",
			Title:  "Tier Critical Storage",
			Value:  "24.500 GB (24.5 TB)",
			Detail: "High-Perf ASM SAN Tier 2 Storage",
			Bg:     "#b45309",
			Accent: "#fde047",
		},
	}
}

// GetTechStackItems returns the official 13 layers
func GetTechStackItems() []model.TechStackItem {
	return []model.TechStackItem{
		{
			ID:        "web_ingress",
			No:        1,
			Layer:     "Web Ingress",
			Component: "NGINX Web Server",
			Icon:      "🌐",
			Moderate: model.TechStackSpec{
				Nodes: "2N (2 Nodes)",
				Type:  "Virtual Machine (VM)",
				Spec:  "8 vCPU / 16GB RAM / 100-200GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "2N (2 Nodes)",
				Type:  "Virtual Machine (VM)",
				Spec:  "8-16 vCPU / 16-32GB RAM / 100-200GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "NGINX Community (com)",
				Enterprise: "NGINX Plus (F5 Enterprise)",
			},
			Notes:       "Moderate: NGINX VM 8/16/100 (IHSP) / 8/16/200 (RMT). Critical: NGINX Plus VM 8-16 vCPU / 16-32 GB (Kontrak F5: TIF.DDL/GH.9076/2025).",
			ContractRef: "TIF.DDL/GH.9076/2025",
		},
		{
			ID:        "app_container",
			No:        2,
			Layer:     "App Container",
			Component: "Red Hat OCP Worker",
			Icon:      "📦",
			Moderate: model.TechStackSpec{
				Nodes: "1W (1 Worker)",
				Type:  "Shared / VM Worker",
				Spec:  "64 vCPU / 1024GB RAM / 300GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "2N (2 Nodes)",
				Type:  "Dedicated Baremetal Worker",
				Spec:  "2x 64 vCPU / 1024GB RAM / 300GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "OKD (Non-Prod)",
				Enterprise: "RH OCP Platform Plus BM Premium",
			},
			Notes:       "Moderate: 1 OCP Worker VM (64/1024/300) per DC. Critical: 2 Dedicated Bare Metal Workers per DC (PPWRKR01/02 & RPWRKR01/02, Kontrak Red Hat: TIF.DDL/GH.1993/2025).",
			ContractRef: "TIF.DDL/GH.1993/2025",
		},
		{
			ID:        "container_router",
			No:        3,
			Layer:     "Container Router",
			Component: "OCP Infra Router & Bastion",
			Icon:      "🔀",
			Moderate: model.TechStackSpec{
				Nodes: "3 Infra + 1 Bastion",
				Type:  "VM Router / HA Proxy",
				Spec:  "8 vCPU / 32GB RAM / 150GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Router Pods",
				Type:  "High-Throughput Router Pods",
				Spec:  "8 vCPU / 32GB RAM / 150GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "Included in OCP",
				Enterprise: "Included in OCP Platform Plus",
			},
			Notes:       "Moderate: 3 Master (8/32/150) + 3 Infra Router (8/32/150) + 1 Bastion (4/8/100). Critical: Dedicated high-throughput router pods cross-DC.",
			ContractRef: "Included in OCP",
		},
		{
			ID:        "core_db",
			No:        4,
			Layer:     "Core Database",
			Component: "Database Engine (Oracle)",
			Icon:      "🗄️",
			Moderate: model.TechStackSpec{
				Nodes: "Single Node, VM",
				Type:  "Virtual Machine (VM)",
				Spec:  "16 vCPU / 32GB RAM / 1000GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "RAC, BM (Baremetal)",
				Type:  "Physical Baremetal + ASM",
				Spec:  "2x 64 vCPU BM Server / 1024GB RAM",
			},
			License: model.LicenseSpec{
				Free:       "PostgreSQL / MySQL Com",
				Enterprise: "Oracle DB Enterprise (RAC + ADG)",
			},
			Notes:       "Moderate: Single Node VM + Data Guard (16/32/1000) tanpa RAC Option untuk efisiensi lisensi core. Critical: Oracle RAC Bare Metal 2x64-Core (High TPS & Zero Data Loss, Kontrak: TIF.DDL/SCD.3585/2025).",
			ContractRef: "TIF.DDL/SCD.3585/2025",
		},
		{
			ID:        "db_storage",
			No:        5,
			Layer:     "Database Storage",
			Component: "Oracle ASM Shared SAN",
			Icon:      "💾",
			Moderate: model.TechStackSpec{
				Nodes: "SAN Virtual Disk LUN",
				Type:  "Storage SAN Tier 2",
				Spec:  "1.000 GB - 1.500 GB SAN LUN",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated ASM SAN Cluster",
				Type:  "High-Perf SAN Tier 2",
				Spec:  "5.500 GB (5.5 TB: Reco 500GB, Data 5TB)",
			},
			License: model.LicenseSpec{
				Free:       "SAN LUN Tier 2",
				Enterprise: "SAN LUN + Oracle ASM",
			},
			Notes:       "Moderate: SAN Tier 2 Virtual LUN 1.0 - 1.5 TB. Critical: SAN Tier 2 ASM 5.5 TB (Reco 500GB, Data 5000GB) pada node bare metal.",
			ContractRef: "SAN Tier 2 ASM",
		},
		{
			ID:        "dist_cache",
			No:        6,
			Layer:     "Distributed Cache",
			Component: "Redis Caching Cluster",
			Icon:      "⚡",
			Moderate: model.TechStackSpec{
				Nodes: "1-3 Nodes",
				Type:  "Virtual Machine (VM)",
				Spec:  "8 vCPU / 16GB RAM / 200GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "3 Nodes Multi-Shard",
				Type:  "Virtual Machine (VM)",
				Spec:  "3x 8 vCPU / 32GB RAM / 200GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "Redis Community (com)",
				Enterprise: "Redis Enterprise Software",
			},
			Notes:       "Moderate: 1-3 VM Redis Community (8/16/200). Critical: 3 VM Redis Enterprise Multi-Shard (PVLMRD01-03 8/32/200, Kontrak: TIF.DDL/GH.9080/2024).",
			ContractRef: "TIF.DDL/GH.9080/2024",
		},
		{
			ID:        "cache_sentinel",
			No:        7,
			Layer:     "Cache Sentinel",
			Component: "Redis Sentinel Quorum",
			Icon:      "🛡️",
			Moderate: model.TechStackSpec{
				Nodes: "Co-located / Sentinel",
				Type:  "Shared VM",
				Spec:  "4 vCPU / 8GB RAM / 100GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Sentinel DC2",
				Type:  "Virtual Machine (VM)",
				Spec:  "4-8 vCPU / 16GB RAM / 100GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "Redis Community (com)",
				Enterprise: "Redis Enterprise",
			},
			Notes:       "Moderate: Co-located sentinel. Critical: Dedicated sentinel cross-DC (Wisdan & Rempoa) untuk quorum deteksi failover otomatis.",
			ContractRef: "TIF.DDL/GH.9080/2024",
		},
		{
			ID:        "message_bus",
			No:        8,
			Layer:     "Message Bus",
			Component: "Confluent / Apache Kafka",
			Icon:      "📨",
			Moderate: model.TechStackSpec{
				Nodes: "3 Nodes (Shared/VM)",
				Type:  "VM Broker Cluster",
				Spec:  "16 vCPU / 64GB / 600GB (IHSP)",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Cluster (7 Nodes)",
				Type:  "ZK + Broker + C3 Cluster",
				Spec:  "3 ZK (4/32) + 3 Broker (16/32/1000) + 1 C3",
			},
			License: model.LicenseSpec{
				Free:       "Apache / Confluent Com",
				Enterprise: "Confluent Kafka Enterprise",
			},
			Notes:       "Moderate: 3 VM Kafka Community (16/64/600 atau 8/16/1000). Critical: 7 Dedicated VM Confluent Kafka Enterprise (3 ZK + 3 Broker + 1 C3, Kontrak: TIF.DDL/GH.9075/2024).",
			ContractRef: "TIF.DDL/GH.9075/2024",
		},
		{
			ID:        "log_analytics",
			No:        9,
			Layer:     "Log Analytics",
			Component: "Elasticsearch Analytics Node",
			Icon:      "🔍",
			Moderate: model.TechStackSpec{
				Nodes: "Combined dgn Logstash",
				Type:  "Shared VM / Combined",
				Spec:  "3 Nodes (8 vCPU / 32GB RAM / 500GB)",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Cluster (8 Nodes)",
				Type:  "Master + Hot + Warm + Coord",
				Spec:  "3 Master (8/32/50) + 2 Hot (16/32/1900) + Warm + Coord",
			},
			License: model.LicenseSpec{
				Free:       "Elasticsearch / OpenSearch",
				Enterprise: "Elasticsearch Platinum",
			},
			Notes:       "Moderate: ELK combined (1 Logstash + 3 Master/Data 8/32/500). Critical: Dedicated 8 VM (3 Master 8/32/50, 2 Hot Data 16/32/1900, 1 Warm, 2 Coord, Kontrak: TIF.DDL/GH.9074/2024).",
			ContractRef: "TIF.DDL/GH.9074/2024",
		},
		{
			ID:        "log_shipper",
			No:        10,
			Layer:     "Log Shipper",
			Component: "Logstash Pipeline Shipper",
			Icon:      "🪵",
			Moderate: model.TechStackSpec{
				Nodes: "Combined / Co-located",
				Type:  "Shared VM",
				Spec:  "1 Node (8 vCPU / 32GB RAM / 100GB)",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Shipper Nodes",
				Type:  "Virtual Machine (VM)",
				Spec:  "2 Nodes (8 vCPU / 32GB RAM / 100GB)",
			},
			License: model.LicenseSpec{
				Free:       "Logstash Open Source",
				Enterprise: "Elastic Platinum",
			},
			Notes:       "Moderate: Combined pada node ELK. Critical: Dedicated 2 VM Logstash Shipper (PVLMEL01/02 8/32/100 di DC & RVLMEL01/02 di DRC) terpisah dari data cluster demi kestabilan heap JVM.",
			ContractRef: "TIF.DDL/GH.9074/2024",
		},
		{
			ID:        "log_dashboard",
			No:        11,
			Layer:     "Log Dashboard",
			Component: "Kibana Observability Portal",
			Icon:      "📊",
			Moderate: model.TechStackSpec{
				Nodes: "Shared Portal",
				Type:  "Shared VM Portal",
				Spec:  "Resource sharing on ELK",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Portal Nodes",
				Type:  "Virtual Machine (VM)",
				Spec:  "2 Nodes (4 vCPU / 32GB RAM / 50GB)",
			},
			License: model.LicenseSpec{
				Free:       "Kibana Community (com)",
				Enterprise: "Elastic Platinum",
			},
			Notes:       "Moderate: Shared portal. Critical: Terintegrasi pada coordinator node PVLMEC01/02 (4/32/50 di DC & RVLMEC01/02 di DRC).",
			ContractRef: "TIF.DDL/GH.9074/2024",
		},
		{
			ID:        "doc_storage",
			No:        12,
			Layer:     "Document Storage",
			Component: "SFTP / Object Storage / RDM",
			Icon:      "🗂️",
			Moderate: model.TechStackSpec{
				Nodes: "SFTP Server VM",
				Type:  "Virtual Machine (VM)",
				Spec:  "8 vCPU / 16GB RAM / 1.500GB Disk",
			},
			Critical: model.TechStackSpec{
				Nodes: "SAN RDM / Dell ECS S3",
				Type:  "Dedicated SAN LUN",
				Spec:  "8 vCPU / 16GB RAM / 365GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "MinIO / OpenSFTP (com)",
				Enterprise: "Dell ECS / SAN Enterprise",
			},
			Notes:       "Moderate: VM SFTP Server (8/16/1500) untuk pertukaran file benefit MFT. Critical: Storage Gateway SAN RDM / Dell ECS S3 untuk arsip transaksi & bukti settlement.",
			ContractRef: "Dell ECS / SAN RDM",
		},
		{
			ID:        "batch_engine",
			No:        13,
			Layer:     "Batch Engine",
			Component: "BMC Control-M / OCP Cron",
			Icon:      "⚙️",
			Moderate: model.TechStackSpec{
				Nodes: "Scheduled Cron on OCP",
				Type:  "Shared OCP Worker",
				Spec:  "Resource sharing on OCP Pods",
			},
			Critical: model.TechStackSpec{
				Nodes: "Dedicated Batch Control-M",
				Type:  "Virtual Machine (VM)",
				Spec:  "4 vCPU / 8GB RAM / 100GB Disk",
			},
			License: model.LicenseSpec{
				Free:       "OCP Native Batch Cron",
				Enterprise: "BMC Control-M Enterprise",
			},
			Notes:       "Moderate: Scheduled CronJob pada worker container OCP. Critical: Dedicated VM Agent BMC Control-M (PVLMBC01 / RVLMBC01 4/8/100) untuk batch settlement end-of-day.",
			ContractRef: "BMC Control-M",
		},
	}
}

// GetEAForumReferences returns official EA forum document anchors
func GetEAForumReferences() []model.ForumReference {
	return []model.ForumReference{
		{
			Title: "TIER MODERATE: IHSP (STM2601010) & RMT (RM TOOLS SATU R1)",
			Tier:  "moderate",
			Badge: "Forum C / SLA 99.5% - 99.9%",
			Points: []model.ForumPoint{
				{
					Label:       "IHSP (Integrated HC Suite Platform)",
					Description: "Acuan Dokumen: EA Forum C STM2601010 R1.1. Sizing per site: NGINX 8/16/100, OCP Master 8/32/150, Infra Router 8/32/150, 1 OCP Worker 64/1024/300, DB VM 16/32/1000, Kafka 16/64/600, SFTP 8/16/1500, ELK 8/32/500.",
				},
				{
					Label:       "RMT (RM Tools SATU R1) & Pola Migrasi RAC",
					Description: "Acuan Dokumen: Capacity Plan RM Tools SATU R1 & evaluasi migrasi database (20260824). NGINX 8/16/200, Redis 8/16/200, Kafka 8/16/1000. Database dialihkan dari RAC ke Single Node Data Guard (Wisdan Primary, Rempoa Standby) untuk optimasi lisensi core.",
				},
				{
					Label:       "Ketentuan Arsitektur Tier Moderate",
					Description: "Kategori EA Forum: Forum C. Target SLA Ketersediaan: 99.5% - 99.9% (RTO < 2 jam, RPO < 15 menit). Seluruh server berbasis Virtual Machine (VM). Database Single Node Data Guard. Server logging (ELK & Logstash) diperbolehkan shared/co-located.",
				},
			},
		},
		{
			Title: "TIER CRITICAL: LIVIN' MERCHANT (DRB2502010 R6 SESI 2)",
			Tier:  "critical",
			Badge: "Forum B (Very Critical) / SLA 99.99%",
			Points: []model.ForumPoint{
				{
					Label:       "Livin' Merchant (Release 6 Sesi 2)",
					Description: "Acuan Dokumen: EA Forum B DRB2502010 R6 Sesi 2. Klasifikasi: Very Critical (Forum B) melayani transaksi merchant nasional 24/7 (SLA 99.99%). Wajib 2 Dedicated OCP Bare Metal Workers (64/1024/300) dan Oracle RAC Bare Metal 64-Core (ASM SAN 5.5 TB).",
				},
				{
					Label:       "Pemisahan Node Logging (ELK & Logstash)",
					Description: "Acuan Server List LM (Tabel 108): Elasticsearch dipisah menjadi Master (PVLMEM01-03 8/32/50), Hot Data (PVLMDH01-02 16/32/1900), Warm Data (PVLMDW01), Coordinator (PVLMEC01-02), dan Logstash Shipper (PVLMEL01-02 8/32/100) di server tersendiri demi stabilitas indexing.",
				},
				{
					Label:       "Lisensi Resmi Enterprise di Livin' Merchant",
					Description: "NGINX Plus (TIF.DDL/GH.9076/2025), Red Hat OCP Platform Plus Bare Metal (TIF.DDL/GH.1993/2025), Oracle RAC Enterprise (TIF.DDL/SCD.3585/2025), Confluent Kafka (TIF.DDL/GH.9075/2024), Elasticsearch Platinum (TIF.DDL/GH.9074/2024), Redis Enterprise (TIF.DDL/GH.9080/2024).",
				},
			},
		},
	}
}
