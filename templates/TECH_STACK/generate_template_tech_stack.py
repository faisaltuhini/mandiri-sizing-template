import openpyxl
import os
from copy import copy
from openpyxl.styles import Font, PatternFill, Alignment, Border, Side
from openpyxl.utils import get_column_letter

def build_template_tech_stack_workbook(output_file):
    wb = openpyxl.Workbook()
    
    # -------------------------------------------------------------------------
    # WORKBOOK SHEETS
    # -------------------------------------------------------------------------
    ws_cover = wb.active
    ws_cover.title = "Cover & HLD Architecture"
    
    ws_matrix = wb.create_sheet(title="Template Tech Stack")
    ws_moderate = wb.create_sheet(title="Sizing Tier Moderate (IHSP RMT)")
    ws_critical = wb.create_sheet(title="Sizing Tier Critical (LM)")
    ws_pricing = wb.create_sheet(title="Pricing")
    ws_formula = wb.create_sheet(title="Formula & Guidelines")
    
    # -------------------------------------------------------------------------
    # STYLES DEFINITION (Enterprise Architecture Corporate Design System)
    # -------------------------------------------------------------------------
    font_main_title = Font(name='Segoe UI', size=15, bold=True, color='FFFFFF')
    font_sub_title = Font(name='Segoe UI', size=10, bold=False, color='E0E8F5')
    font_card_header = Font(name='Calibri', size=11, bold=True, color='FFFFFF')
    font_section_title = Font(name='Calibri', size=13, bold=True, color='1B365D')
    font_table_title = Font(name='Calibri', size=12, bold=True, color='FFFFFF')
    font_title_bold = Font(name='Calibri', size=12, bold=True, color='FFFFFF')
    font_header_white = Font(name='Calibri', size=10, bold=True, color='FFFFFF')
    
    font_bold_12 = Font(name='Calibri', size=12, bold=True)
    font_bold_11 = Font(name='Calibri', size=11, bold=True)
    font_bold_10 = Font(name='Calibri', size=10, bold=True)
    font_regular_11 = Font(name='Calibri', size=11, bold=False)
    font_regular_10 = Font(name='Calibri', size=10, bold=False)
    font_italic_10 = Font(name='Calibri', size=10, italic=True, color='555555')
    font_italic_9 = Font(name='Calibri', size=9, italic=True, color='555555')
    
    font_badge_green = Font(name='Calibri', size=10, bold=True, color='276A3C')
    font_badge_blue = Font(name='Calibri', size=10, bold=True, color='1F4E79')
    
    # Color Fills
    fill_navy = PatternFill(start_color='1B365D', end_color='1B365D', fill_type='solid')       # Enterprise Deep Navy
    fill_blue_hdr = PatternFill(start_color='305496', end_color='305496', fill_type='solid')    # Corporate Blue
    fill_gold_hdr = PatternFill(start_color='C65911', end_color='C65911', fill_type='solid')    # Rust / Enterprise Gold
    fill_accent_blue = PatternFill(start_color='418AB3', end_color='418AB3', fill_type='solid') # Accent Blue
    fill_teal_hdr = PatternFill(start_color='205D67', end_color='205D67', fill_type='solid')
    fill_dc1_dc2 = PatternFill(start_color='DDEBF7', end_color='DDEBF7', fill_type='solid')     # Light Blue
    fill_pt_dev = PatternFill(start_color='FCE4D6', end_color='FCE4D6', fill_type='solid')      # Light Peach
    fill_zebra = PatternFill(start_color='F9FAFB', end_color='F9FAFB', fill_type='solid')       # Subtle Zebra
    fill_subtotal = PatternFill(start_color='E8EEF5', end_color='E8EEF5', fill_type='solid')
    fill_card_bg = PatternFill(start_color='F0F4F8', end_color='F0F4F8', fill_type='solid')
    
    fill_badge_green = PatternFill(start_color='E2EFDA', end_color='E2EFDA', fill_type='solid') # Light Green (Free/Com)
    fill_badge_blue = PatternFill(start_color='D9E1F2', end_color='D9E1F2', fill_type='solid')  # Light Blue (Enterprise)
    
    # Borders
    border_thin = Border(
        left=Side(style='thin', color='D0D5DD'),
        right=Side(style='thin', color='D0D5DD'),
        top=Side(style='thin', color='D0D5DD'),
        bottom=Side(style='thin', color='D0D5DD')
    )
    border_header = Border(
        left=Side(style='thin', color='FFFFFF'),
        right=Side(style='thin', color='FFFFFF'),
        top=Side(style='thin', color='305496'),
        bottom=Side(style='medium', color='305496')
    )
    border_subtotal = Border(
        left=Side(style='thin', color='D0D5DD'),
        right=Side(style='thin', color='D0D5DD'),
        top=Side(style='thin', color='305496'),
        bottom=Side(style='double', color='305496')
    )
    
    align_center = Alignment(horizontal='center', vertical='center', wrap_text=True)
    align_left = Alignment(horizontal='left', vertical='center')
    align_left_wrap = Alignment(horizontal='left', vertical='center', wrap_text=True)
    align_right = Alignment(horizontal='right', vertical='center')
    align_header = Alignment(horizontal='center', vertical='center', wrap_text=True)
    
    num_fmt_curr = '_-[$Rp-3809]* #,##0.00_-;\\-[$Rp-3809]* #,##0.00_-;_-[$Rp-3809]* "-"??_-;_-@_-'
    num_fmt_int = '#,##0'

    # =========================================================================
    # 1. SHEET: COVER & HLD ARCHITECTURE (Exact styling matching 11_SCM)
    # =========================================================================
    ws = ws_cover
    ws.views.sheetView[0].showGridLines = True
    
    col_widths_cover = {
        'A': 4, 'B': 6, 'C': 26, 'D': 36, 'E': 28,
        'F': 24, 'G': 24, 'H': 24, 'I': 32
    }
    for c_let, w in col_widths_cover.items():
        ws.column_dimensions[c_let].width = w

    # Header Banner (Rows 2..4)
    ws.merge_cells('B2:I2')
    ws.merge_cells('B3:I3')
    ws.merge_cells('B4:I4')
    
    ws.cell(2, 2, "BASELINE ARSITEKTUR & SIZING INFRASTRUKTUR IT").font = font_main_title
    ws.cell(3, 2, "STANDAR KOMPARASI: TIER MODERATE (IHSP/RMT) VS TIER CRITICAL (LIVIN' MERCHANT)").font = font_card_header
    ws.cell(4, 2, "Dokumen Acuan: EA Forum C STM2601010 (IHSP), Capacity Plan RMT, & EA Forum B DRB2502010 (Livin' Merchant)").font = font_sub_title
    
    for r in range(2, 5):
        ws.row_dimensions[r].height = 22
        for c in range(2, 10):
            ws.cell(r, c).fill = fill_navy
            ws.cell(r, c).alignment = align_left
            
    # Metadata Box (Rows 6..9)
    meta_data = [
        ("Nama Dokumen", "Baseline Sizing Infrastruktur & Technology Stack", "Target Environment", "Production (DC Wisdan & DRC Rempoa), PT, DEV & UAT"),
        ("Sistem Acuan Moderate", "Integrated HC Suite Platform (IHSP) & RM Tools (RMT)", "Sistem Acuan Critical", "Livin' Merchant (LM) - Merchant App 2026 Release 6"),
        ("Arsitektur Database", "Moderate: Single Node VM Data Guard (Tanpa RAC)", "Database Critical", "Critical: Oracle RAC Bare Metal + 5.5 TB ASM SAN Storage"),
        ("Klasifikasi EA Forum", "Moderate: Forum C (STM2601010) | Critical: Forum B", "Pola Logging & Monitoring", "Moderate: ELK Co-located / Shared | Critical: ELK Dedicated Cluster (Logstash Terpisah)")
    ]
    for idx, (lbl1, val1, lbl2, val2) in enumerate(meta_data, start=6):
        ws.row_dimensions[idx].height = 22
        ws.cell(idx, 2, "").border = border_thin
        ws.cell(idx, 3, lbl1).font = font_bold_11
        ws.cell(idx, 3).fill = fill_subtotal
        ws.cell(idx, 3).border = border_thin
        
        ws.merge_cells(start_row=idx, start_column=4, end_row=idx, end_column=5)
        ws.cell(idx, 4, val1).font = font_regular_11
        ws.cell(idx, 4).border = border_thin
        ws.cell(idx, 5).border = border_thin
        
        ws.cell(idx, 6, lbl2).font = font_bold_11
        ws.cell(idx, 6).fill = fill_subtotal
        ws.cell(idx, 6).border = border_thin
        
        ws.merge_cells(start_row=idx, start_column=7, end_row=idx, end_column=9)
        ws.cell(idx, 7, val2).font = font_regular_11
        ws.cell(idx, 7).border = border_thin
        ws.cell(idx, 8).border = border_thin
        ws.cell(idx, 9).border = border_thin

    # Executive Sizing KPI Summary Cards (Rows 11..14)
    ws.merge_cells('B11:I11')
    ws.cell(11, 2, "RINGKASAN SIZING & ALOKASI RESOURCE (BASELINE)").font = font_section_title
    ws.row_dimensions[11].height = 26
    
    kpis = [
        (3, "Tier Moderate Compute", "14 VM Nodes", "112 vCPU | 256 GB RAM (IHSP/RMT)", fill_navy),
        (5, "Tier Moderate Storage", "4.500 GB (4.5 TB)", "SAN Tier 2 LUN Virtual Disk", fill_blue_hdr),
        (7, "Tier Critical Compute", "28 Dedicated Nodes", "384 vCPU | 4.8 TB RAM (incl 4 BM)", fill_accent_blue),
        (9, "Tier Critical Storage", "24.500 GB (24.5 TB)", "High-Perf ASM SAN Tier 2", fill_gold_hdr)
    ]
    card_cols = [3, 5, 7, 9]
    ws.row_dimensions[12].height = 24
    ws.row_dimensions[13].height = 28
    ws.row_dimensions[14].height = 24
    for idx_k, (col_i, title, val, sub, bg) in enumerate(kpis):
        real_col = card_cols[idx_k]
        c1 = ws.cell(12, real_col, title)
        c1.font = font_card_header
        c1.fill = bg
        c1.alignment = align_center
        
        c2 = ws.cell(13, real_col, val)
        c2.font = font_bold_12
        c2.alignment = align_center
        c2.fill = fill_card_bg
        c2.border = border_thin
        
        c3 = ws.cell(14, real_col, sub)
        c3.font = font_italic_10
        c3.alignment = align_center
        c3.fill = fill_card_bg
        c3.border = border_thin

    # Architecture Components & Layers Table (Rows 16..30)
    ws.merge_cells('B16:I16')
    ws.cell(16, 2, "MATRIKS KOMPARASI KOMPONEN TEKNOLOGI (MODERATE VS CRITICAL)").font = font_section_title
    ws.row_dimensions[16].height = 26
    
    # Headers Row 17 (Notice G17:I17 is properly merged!)
    arch_headers = [
        (2, 2, "No"),
        (3, 3, "Layer Arsitektur"),
        (4, 4, "Komponen Sistem"),
        (5, 5, "Spesifikasi Tier Moderate (IHSP/RMT)"),
        (6, 6, "Spesifikasi Tier Critical (Livin' Merchant)"),
        (7, 9, "Ketentuan Lisensi & Catatan EA")
    ]
    ws.row_dimensions[17].height = 24
    for start_c, end_c, h_text in arch_headers:
        if start_c != end_c:
            ws.merge_cells(start_row=17, start_column=start_c, end_row=17, end_column=end_c)
        for c in range(start_c, end_c + 1):
            cell = ws.cell(17, c)
            if c == start_c:
                cell.value = h_text
            cell.font = font_header_white
            cell.fill = fill_blue_hdr
            cell.alignment = align_header
            cell.border = border_header
        
    arch_rows = [
        (1, "Web Ingress", "NGINX Web Server", "2 VM per site (8 vCPU / 16GB / 100-200GB)", "2 VM per site (8-16 vCPU / 16-32GB / 100-200GB)", "Moderate: NGINX Community (Open Source). Critical: NGINX Plus Enterprise Subscription (Kontrak F5: TIF.DDL/GH.9076/2025)."),
        (2, "App Container", "Red Hat OCP Worker", "1 Worker VM per site (64 vCPU / 1024GB / 300GB)", "2 Dedicated Bare Metal Workers per site (64 vCPU / 1024GB / 300GB)", "Moderate: OCP Core Subscription. Critical: Red Hat OpenShift Platform Plus Bare Metal Node Premium (Kontrak TIF.DDL/GH.1993/2025)."),
        (3, "Container Router", "OCP Infra Router & Bastion", "3 Infra Router (8/32/150) + 1 Bastion HA (4/8/100)", "Dedicated High-Throughput Ingress Router Pods", "Included dalam langganan Red Hat OpenShift Container Platform."),
        (4, "Core Database", "Database Engine (Oracle)", "Single Node VM Data Guard (16 vCPU / 32GB / 1000GB)", "Oracle RAC Bare Metal (2x 64-Core BM per site + Active Data Guard)", "Moderate: Oracle DB Enterprise (tanpa RAC Option untuk efisiensi lisensi core). Critical: Oracle RAC Bare Metal (High TPS & zero data loss, Kontrak TIF.DDL/SCD.3585/2025)."),
        (5, "Database Storage", "Oracle ASM Shared Storage", "Storage SAN Tier 2 Virtual LUN (1.0 TB - 1.5 TB)", "Storage SAN Tier 2 High-Performance ASM (5.5 TB: Reco 500GB, Data 5TB)", "Diskgroup tablespace data dan Fast Recovery Area (FRA) berbasis ASM striping."),
        (6, "Distributed Cache", "Redis Caching Cluster", "1-3 VM (8 vCPU / 16GB / 200GB)", "3 VM Multi-Shard (8 vCPU / 32GB / 200GB)", "Moderate: Redis Community (Open Source). Critical: Redis Enterprise Software Subscription (Kontrak TIF.DDL/GH.9080/2024)."),
        (7, "Cache Sentinel", "Redis Sentinel Quorum", "Co-located / 1 Standalone VM", "Dedicated Sentinel Quorum Arbiter antar-DC (Wisdan & Rempoa)", "Arbiter quorum untuk deteksi failover otomatis dan promosi master Redis cross-DC."),
        (8, "Message Bus", "Confluent / Apache Kafka", "3 VM (16 vCPU / 64GB / 600GB atau 8/16/1000GB)", "Dedicated Cluster 7 VM: 3 ZK (4/32/100) + 3 Broker (16/32/1000) + 1 C3 (8/16/100)", "Moderate: Kafka Community (Shared broker cluster). Critical: Confluent Kafka Enterprise Subscription (Kontrak TIF.DDL/GH.9075/2024)."),
        (9, "Log Analytics", "Elasticsearch Data Cluster", "Combined dengan Logstash (1 Logstash + 3 Master/Data 8/32/500)", "Dedicated Cluster 8 VM: 3 Master (8/32/50), 2 Hot Data (16/32/1900), 1 Warm, 2 Coordinator", "Moderate: OpenSearch / Elasticsearch Basic. Critical: Elasticsearch Platinum (Kontrak TIF.DDL/GH.9074/2024). Pemisahan node data hot & master untuk kestabilan heap JVM."),
        (10, "Log Shipper", "Logstash Pipeline Shipper", "Combined / Co-located pada server logging", "Dedicated 2 VM Logstash Shipper (PVLMEL01/02 8 vCPU / 32GB / 100GB)", "Moderate: Co-located. Critical: Wajib dipisah ke server tersendiri agar proses parsing Grok tidak membebani query indexing Elasticsearch."),
        (11, "Log Dashboard", "Kibana Observability", "Shared Dashboard Portal", "Terintegrasi pada Coordinator Node (PVLMEC01/02 4 vCPU / 32GB / 50GB)", "Visualisasi audit log transaksi, tracing performa APM, dan dashboard operasional."),
        (12, "Document Storage", "Object Storage / SFTP / RDM", "VM SFTP Server (8/16/1500) / MinIO S3 Gateway", "Storage Gateway SAN RDM (8/16/365) / Dell ECS Object Storage", "Media penyimpanan arsip dokumen, file transaksi QRIS, invoice merchant, dan pertukaran file MFT."),
        (13, "Batch Engine", "BMC Control-M / OCP Cron", "Scheduled CronJob pada worker container OCP", "Dedicated VM Agent BMC Control-M (PVLMBC01 / RVLMBC01 4/8/100)", "Moderate: Native OCP CronJob. Critical: BMC Control-M Enterprise untuk otomasi settlement transaksi end-of-day.")
    ]
    
    for idx, (no, lyr, comp, mod_p, crit_p, lic_note) in enumerate(arch_rows, start=18):
        ws.row_dimensions[idx].height = 28
        ws.cell(idx, 2, no).alignment = align_center
        ws.cell(idx, 3, lyr).alignment = align_left
        ws.cell(idx, 4, comp).alignment = align_left
        ws.cell(idx, 5, mod_p).alignment = align_left_wrap
        ws.cell(idx, 6, crit_p).alignment = align_left_wrap
        
        ws.merge_cells(start_row=idx, start_column=7, end_row=idx, end_column=9)
        ws.cell(idx, 7, lic_note).alignment = align_left_wrap
        
        for c in range(2, 10):
            cell = ws.cell(idx, c)
            cell.font = font_regular_10
            cell.border = border_thin
            if idx % 2 == 1:
                cell.fill = fill_zebra

    # =========================================================================
    # 2. SHEET: TEMPLATE TECH STACK (Polished layout, generous widths & heights)
    # =========================================================================
    ws_m = ws_matrix
    ws_m.views.sheetView[0].showGridLines = True
    
    # Generous column widths so text never truncates or squishes
    col_widths_matrix = {
        'A': 4,
        'B': 6,   # No
        'C': 22,  # Layer Arsitektur
        'D': 26,  # Komponen / Software
        'E': 18,  # Moderate Node Alokasi
        'F': 22,  # Moderate Tipe Compute
        'G': 34,  # Moderate Spec per Node (CPU/RAM/Disk)
        'H': 18,  # Critical Node Alokasi
        'I': 25,  # Critical Tipe Compute
        'J': 36,  # Critical Spec per Node (CPU/RAM/Disk)
        'K': 24,  # Free / Community
        'L': 32,  # Enterprise
        'M': 68   # Catatan Teknis & Justifikasi Arsitektur
    }
    for c_let, w in col_widths_matrix.items():
        ws_m.column_dimensions[c_let].width = w

    # Row 1: Title Banner
    ws_m.merge_cells('B1:M1')
    ws_m.cell(1, 2, "BASELINE TECHNOLOGY STACK & SIZING PATTERN (ENTERPRISE)").font = font_title_bold
    for c in range(2, 14):
        ws_m.cell(1, c).fill = fill_navy
    ws_m.row_dimensions[1].height = 28
    
    # Row 2: Subtitle
    ws_m.merge_cells('B2:M2')
    ws_m.cell(2, 2, "Standar Arsitektur: Tier Moderate (IHSP & RMT) vs Tier Critical (Livin' Merchant)").font = font_sub_title
    for c in range(2, 14):
        ws_m.cell(2, c).fill = fill_navy
    ws_m.row_dimensions[2].height = 22

    # Row 3 & 4: Two-Tier Header with generous heights (44pt for Row 4 so (CPU/RAM/Disk) never clips)
    ws_m.row_dimensions[3].height = 28
    ws_m.row_dimensions[4].height = 44
    
    # Merge vertical headers (Row 3..4)
    for col_i, h_txt in [(2, "No"), (3, "Layer Arsitektur"), (4, "Komponen / Software"), (13, "Catatan Arsitektur & Justifikasi Sizing")]:
        col_let = get_column_letter(col_i)
        ws_m.merge_cells(f'{col_let}3:{col_let}4')
        c3 = ws_m.cell(3, col_i, h_txt)
        c3.font = font_header_white
        c3.fill = fill_blue_hdr
        c3.alignment = align_header
        c3.border = border_header
        ws_m.cell(4, col_i).border = border_header
    
    # Moderate Tier Header (Cols E..G)
    ws_m.merge_cells('E3:G3')
    ws_m.cell(3, 5, "TIER MODERATE (IHSP & RMT)").font = font_header_white
    for c in range(5, 8):
        ws_m.cell(3, c).fill = fill_accent_blue
        ws_m.cell(3, c).alignment = align_header
        ws_m.cell(3, c).border = border_header
    
    ws_m.cell(4, 5, "Node Alokasi").font = font_header_white
    ws_m.cell(4, 5).fill = fill_accent_blue
    ws_m.cell(4, 5).alignment = align_header
    ws_m.cell(4, 5).border = border_header
    
    ws_m.cell(4, 6, "Tipe Compute").font = font_header_white
    ws_m.cell(4, 6).fill = fill_accent_blue
    ws_m.cell(4, 6).alignment = align_header
    ws_m.cell(4, 6).border = border_header
    
    ws_m.cell(4, 7, "Spec per Node\n(CPU/RAM/Disk)").font = font_header_white
    ws_m.cell(4, 7).fill = fill_accent_blue
    ws_m.cell(4, 7).alignment = align_header
    ws_m.cell(4, 7).border = border_header
    
    # Critical Tier Header (Cols H..J)
    ws_m.merge_cells('H3:J3')
    ws_m.cell(3, 8, "TIER CRITICAL (LIVIN' MERCHANT)").font = font_header_white
    for c in range(8, 11):
        ws_m.cell(3, c).fill = fill_gold_hdr
        ws_m.cell(3, c).alignment = align_header
        ws_m.cell(3, c).border = border_header
        
    ws_m.cell(4, 8, "Node Alokasi").font = font_header_white
    ws_m.cell(4, 8).fill = fill_gold_hdr
    ws_m.cell(4, 8).alignment = align_header
    ws_m.cell(4, 8).border = border_header
    
    ws_m.cell(4, 9, "Tipe Compute").font = font_header_white
    ws_m.cell(4, 9).fill = fill_gold_hdr
    ws_m.cell(4, 9).alignment = align_header
    ws_m.cell(4, 9).border = border_header
    
    ws_m.cell(4, 10, "Spec per Node\n(CPU/RAM/Disk)").font = font_header_white
    ws_m.cell(4, 10).fill = fill_gold_hdr
    ws_m.cell(4, 10).alignment = align_header
    ws_m.cell(4, 10).border = border_header

    # License Headers (Cols K..L)
    ws_m.merge_cells('K3:L3')
    ws_m.cell(3, 11, "KATEGORI LISENSI").font = font_header_white
    for c in range(11, 13):
        ws_m.cell(3, c).fill = fill_teal_hdr
        ws_m.cell(3, c).alignment = align_header
        ws_m.cell(3, c).border = border_header
        
    ws_m.cell(4, 11, "Free / Community").font = font_header_white
    ws_m.cell(4, 11).fill = fill_teal_hdr
    ws_m.cell(4, 11).alignment = align_header
    ws_m.cell(4, 11).border = border_header
    
    ws_m.cell(4, 12, "Enterprise").font = font_header_white
    ws_m.cell(4, 12).fill = fill_teal_hdr
    ws_m.cell(4, 12).alignment = align_header
    ws_m.cell(4, 12).border = border_header

    # Comprehensive Data Rows (Directly matching IHSP, RMT, and Livin Merchant documents)
    tech_stack_data = [
        (
            1, "Web Ingress", "NGINX Web Server",
            "2N (2 Nodes)", "Virtual Machine (VM)", "8 vCPU / 16GB RAM / 100-200GB Disk",
            "2N (2 Nodes)", "Virtual Machine (VM)", "8-16 vCPU / 16-32GB RAM / 100-200GB Disk",
            "NGINX Community (com)", "NGINX Plus (F5 Enterprise)",
            "Moderate: NGINX VM 8/16/100 (IHSP) / 8/16/200 (RMT). Critical: NGINX Plus VM 8-16 vCPU / 16-32 GB (Kontrak F5: TIF.DDL/GH.9076/2025)."
        ),
        (
            2, "App Container", "Red Hat OCP Worker",
            "1W (1 Worker)", "Shared / VM Worker", "64 vCPU / 1024GB RAM / 300GB Disk",
            "2N (2 Nodes)", "Dedicated Baremetal Worker", "2x 64 vCPU / 1024GB RAM / 300GB Disk",
            "OKD (Non-Prod)", "RH OCP Platform Plus BM Premium",
            "Moderate: 1 OCP Worker VM (64/1024/300) per DC. Critical: 2 Dedicated Bare Metal Workers per DC (PPWRKR01/02 & RPWRKR01/02, Kontrak Red Hat: TIF.DDL/GH.1993/2025)."
        ),
        (
            3, "Container Router", "OCP Infra Router & Bastion",
            "3 Infra + 1 Bastion", "VM Router / HA Proxy", "8 vCPU / 32GB RAM / 150GB Disk",
            "Dedicated Router Pods", "High-Throughput Router Pods", "8 vCPU / 32GB RAM / 150GB Disk",
            "Included in OCP", "Included in OCP Platform Plus",
            "Moderate: 3 Master (8/32/150) + 3 Infra Router (8/32/150) + 1 Bastion (4/8/100). Critical: Dedicated high-throughput router pods cross-DC."
        ),
        (
            4, "Core Database", "Database Engine (Oracle)",
            "Single Node, VM", "Virtual Machine (VM)", "16 vCPU / 32GB RAM / 1000GB Disk",
            "RAC, BM (Baremetal)", "Physical Baremetal + ASM", "2x 64 vCPU BM Server / 1024GB RAM",
            "PostgreSQL / MySQL Com", "Oracle DB Enterprise (RAC + ADG)",
            "Moderate: Single Node VM + Data Guard (16/32/1000) tanpa RAC Option untuk efisiensi lisensi core. Critical: Oracle RAC Bare Metal 2x64-Core (High TPS & Zero Data Loss, Kontrak: TIF.DDL/SCD.3585/2025)."
        ),
        (
            5, "Database Storage", "Oracle ASM Shared SAN",
            "SAN Virtual Disk LUN", "Storage SAN Tier 2", "1.000 GB - 1.500 GB SAN LUN",
            "Dedicated ASM SAN Cluster", "High-Perf SAN Tier 2", "5.500 GB (5.5 TB: Reco 500GB, Data 5TB)",
            "SAN LUN Tier 2", "SAN LUN + Oracle ASM",
            "Moderate: SAN Tier 2 Virtual LUN 1.0 - 1.5 TB. Critical: SAN Tier 2 ASM 5.5 TB (Reco 500GB, Data 5000GB) pada node bare metal."
        ),
        (
            6, "Distributed Cache", "Redis Caching Cluster",
            "1-3 Nodes", "Virtual Machine (VM)", "8 vCPU / 16GB RAM / 200GB Disk",
            "3 Nodes Multi-Shard", "Virtual Machine (VM)", "3x 8 vCPU / 32GB RAM / 200GB Disk",
            "Redis Community (com)", "Redis Enterprise Software",
            "Moderate: 1-3 VM Redis Community (8/16/200). Critical: 3 VM Redis Enterprise Multi-Shard (PVLMRD01-03 8/32/200, Kontrak: TIF.DDL/GH.9080/2024)."
        ),
        (
            7, "Cache Sentinel", "Redis Sentinel Quorum",
            "Co-located / Sentinel", "Shared VM", "4 vCPU / 8GB RAM / 100GB Disk",
            "Dedicated Sentinel DC2", "Virtual Machine (VM)", "4-8 vCPU / 16GB RAM / 100GB Disk",
            "Redis Community (com)", "Redis Enterprise",
            "Moderate: Co-located sentinel. Critical: Dedicated sentinel cross-DC (Wisdan & Rempoa) untuk quorum deteksi failover otomatis."
        ),
        (
            8, "Message Bus", "Confluent / Apache Kafka",
            "3 Nodes (Shared/VM)", "VM Broker Cluster", "16 vCPU / 64GB / 600GB (IHSP)",
            "Dedicated Cluster (7 Nodes)", "ZK + Broker + C3 Cluster", "3 ZK (4/32) + 3 Broker (16/32/1000) + 1 C3",
            "Apache / Confluent Com", "Confluent Kafka Enterprise",
            "Moderate: 3 VM Kafka Community (16/64/600 atau 8/16/1000). Critical: 7 Dedicated VM Confluent Kafka Enterprise (3 ZK + 3 Broker + 1 C3, Kontrak: TIF.DDL/GH.9075/2024)."
        ),
        (
            9, "Log Analytics", "Elasticsearch Analytics Node",
            "Combined dgn Logstash", "Shared VM / Combined", "3 Nodes (8 vCPU / 32GB RAM / 500GB)",
            "Dedicated Cluster (8 Nodes)", "Master + Hot + Warm + Coord", "3 Master (8/32/50) + 2 Hot (16/32/1900) + Warm + Coord",
            "Elasticsearch / OpenSearch", "Elasticsearch Platinum",
            "Moderate: ELK combined (1 Logstash + 3 Master/Data 8/32/500). Critical: Dedicated 8 VM (3 Master 8/32/50, 2 Hot Data 16/32/1900, 1 Warm, 2 Coord, Kontrak: TIF.DDL/GH.9074/2024)."
        ),
        (
            10, "Log Shipper", "Logstash Pipeline Shipper",
            "Combined / Co-located", "Shared VM", "1 Node (8 vCPU / 32GB RAM / 100GB)",
            "Dedicated Shipper Nodes", "Virtual Machine (VM)", "2 Nodes (8 vCPU / 32GB RAM / 100GB)",
            "Logstash Open Source", "Elastic Platinum",
            "Moderate: Combined pada node ELK. Critical: Dedicated 2 VM Logstash Shipper (PVLMEL01/02 8/32/100 di DC & RVLMEL01/02 di DRC) terpisah dari data cluster demi kestabilan heap JVM."
        ),
        (
            11, "Log Dashboard", "Kibana Observability Portal",
            "Shared Portal", "Shared VM Portal", "Resource sharing on ELK",
            "Dedicated Portal Nodes", "Virtual Machine (VM)", "2 Nodes (4 vCPU / 32GB RAM / 50GB)",
            "Kibana Community (com)", "Elastic Platinum",
            "Moderate: Shared portal. Critical: Terintegrasi pada coordinator node PVLMEC01/02 (4/32/50 di DC & RVLMEC01/02 di DRC)."
        ),
        (
            12, "Document Storage", "SFTP / Object Storage / RDM",
            "SFTP Server VM", "Virtual Machine (VM)", "8 vCPU / 16GB RAM / 1.500GB Disk",
            "SAN RDM / Dell ECS S3", "Dedicated SAN LUN", "8 vCPU / 16GB RAM / 365GB Disk",
            "MinIO / OpenSFTP (com)", "Dell ECS / SAN Enterprise",
            "Moderate: VM SFTP Server (8/16/1500) untuk pertukaran file benefit MFT. Critical: Storage Gateway SAN RDM / Dell ECS S3 untuk arsip transaksi & bukti settlement."
        ),
        (
            13, "Batch Engine", "BMC Control-M / OCP Cron",
            "Scheduled Cron on OCP", "Shared OCP Worker", "Resource sharing on OCP Pods",
            "Dedicated Batch Control-M", "Virtual Machine (VM)", "4 vCPU / 8GB RAM / 100GB Disk",
            "OCP Native Batch Cron", "BMC Control-M Enterprise",
            "Moderate: Scheduled CronJob pada worker container OCP. Critical: Dedicated VM Agent BMC Control-M (PVLMBC01 / RVLMBC01 4/8/100) untuk batch settlement end-of-day."
        )
    ]

    # Generous row height (58pt) so all notes and multiline specs breathe without clipping
    for row_idx, item in enumerate(tech_stack_data, start=5):
        ws_m.row_dimensions[row_idx].height = 58
        no, lyr, comp, mod_n, mod_typ, mod_spec, crit_n, crit_typ, crit_spec, lic_free, lic_ent, notes = item
        
        ws_m.cell(row_idx, 2, no).alignment = align_center
        ws_m.cell(row_idx, 3, lyr).alignment = align_left
        ws_m.cell(row_idx, 4, comp).alignment = align_left
        
        ws_m.cell(row_idx, 5, mod_n).alignment = align_center
        ws_m.cell(row_idx, 6, mod_typ).alignment = align_center
        ws_m.cell(row_idx, 7, mod_spec).alignment = align_left_wrap
        
        ws_m.cell(row_idx, 8, crit_n).alignment = align_center
        ws_m.cell(row_idx, 9, crit_typ).alignment = align_center
        ws_m.cell(row_idx, 10, crit_spec).alignment = align_left_wrap
        
        c_free = ws_m.cell(row_idx, 11, lic_free)
        c_free.alignment = align_center
        c_free.fill = fill_badge_green
        c_free.font = font_badge_green
        
        c_ent = ws_m.cell(row_idx, 12, lic_ent)
        c_ent.alignment = align_center
        c_ent.fill = fill_badge_blue
        c_ent.font = font_badge_blue
        
        ws_m.cell(row_idx, 13, notes).alignment = align_left_wrap
        
        for c in range(2, 14):
            cell = ws_m.cell(row_idx, c)
            cell.border = border_thin
            if c not in (11, 12):
                cell.font = font_regular_10
                if row_idx % 2 == 1:
                    cell.fill = fill_zebra

    # Bottom Application Examples Box (Dynamic positioning after 13 rows -> starts at row 19)
    start_app_box = 5 + len(tech_stack_data) + 1  # 5 + 13 + 1 = 19
    ws_m.merge_cells(start_row=start_app_box, start_column=2, end_row=start_app_box, end_column=13)
    ws_m.cell(start_app_box, 2, "CATATAN ARSITEKTUR & ACUAN RESMI EA FORUM").font = font_section_title
    ws_m.row_dimensions[start_app_box].height = 28
    
    header_app_row = start_app_box + 1  # 20
    ws_m.merge_cells(start_row=header_app_row, start_column=2, end_row=header_app_row, end_column=7)
    ws_m.cell(header_app_row, 2, "TIER MODERATE: IHSP (STM2601010) & RMT (RM TOOLS SATU R1)").font = font_card_header
    for c in range(2, 8):
        ws_m.cell(header_app_row, c).fill = fill_accent_blue
        ws_m.cell(header_app_row, c).alignment = align_left
        ws_m.cell(header_app_row, c).border = border_header
        
    ws_m.merge_cells(start_row=header_app_row, start_column=8, end_row=header_app_row, end_column=13)
    ws_m.cell(header_app_row, 8, "TIER CRITICAL: LIVIN' MERCHANT (DRB2502010 R6 SESI 2)").font = font_card_header
    for c in range(8, 14):
        ws_m.cell(header_app_row, c).fill = fill_gold_hdr
        ws_m.cell(header_app_row, c).alignment = align_left
        ws_m.cell(header_app_row, c).border = border_header
    ws_m.row_dimensions[header_app_row].height = 26
    
    mod_apps_desc = [
        ("IHSP (Integrated HC Suite Platform)", "Acuan Dokumen: EA Forum C STM2601010 R1.1. Sizing per site: NGINX 8/16/100, OCP Master 8/32/150, Infra Router 8/32/150, 1 OCP Worker 64/1024/300, DB VM 16/32/1000, Kafka 16/64/600, SFTP 8/16/1500, ELK 8/32/500."),
        ("RMT (RM Tools SATU R1) & Pola Migrasi RAC", "Acuan Dokumen: Capacity Plan RM Tools SATU R1 & evaluasi migrasi database (20260824). NGINX 8/16/200, Redis 8/16/200, Kafka 8/16/1000. Database dialihkan dari RAC ke Single Node Data Guard (Wisdan Primary, Rempoa Standby) untuk optimasi lisensi core."),
        ("Ketentuan Arsitektur Tier Moderate", "Kategori EA Forum: Forum C. Target SLA Ketersediaan: 99.5% - 99.9% (RTO < 2 jam, RPO < 15 menit). Seluruh server berbasis Virtual Machine (VM). Database Single Node Data Guard. Server logging (ELK & Logstash) diperbolehkan shared/co-located.")
    ]
    crit_apps_desc = [
        ("Livin' Merchant (Release 6 Sesi 2)", "Acuan Dokumen: EA Forum B DRB2502010 R6 Sesi 2. Klasifikasi: Very Critical (Forum B) melayani transaksi merchant nasional 24/7 (SLA 99.99%). Wajib 2 Dedicated OCP Bare Metal Workers (64/1024/300) dan Oracle RAC Bare Metal 64-Core (ASM SAN 5.5 TB)."),
        ("Pemisahan Node Logging (ELK & Logstash)", "Acuan Server List LM (Tabel 108): Elasticsearch dipisah menjadi Master (PVLMEM01-03 8/32/50), Hot Data (PVLMDH01-02 16/32/1900), Warm Data (PVLMDW01), Coordinator (PVLMEC01-02), dan Logstash Shipper (PVLMEL01-02 8/32/100) di server tersendiri demi stabilitas indexing."),
        ("Lisensi Resmi Enterprise di Livin' Merchant", "NGINX Plus (TIF.DDL/GH.9076/2025), Red Hat OCP Platform Plus Bare Metal (TIF.DDL/GH.1993/2025), Oracle RAC Enterprise (TIF.DDL/SCD.3585/2025), Confluent Kafka (TIF.DDL/GH.9075/2024), Elasticsearch Platinum (TIF.DDL/GH.9074/2024), Redis Enterprise (TIF.DDL/GH.9080/2024).")
    ]
    
    for i in range(3):
        r_num = header_app_row + 1 + i
        ws_m.row_dimensions[r_num].height = 72  # 72pt accommodates 4-5 lines of detailed reference notes without clipping
        
        m_title, m_desc = mod_apps_desc[i]
        c_title, c_desc = crit_apps_desc[i]
        
        ws_m.merge_cells(start_row=r_num, start_column=2, end_row=r_num, end_column=3)
        ws_m.cell(r_num, 2, m_title).font = font_bold_10
        ws_m.cell(r_num, 2).fill = fill_card_bg
        ws_m.cell(r_num, 2).alignment = align_left_wrap
        
        ws_m.merge_cells(start_row=r_num, start_column=4, end_row=r_num, end_column=7)
        ws_m.cell(r_num, 4, m_desc).font = font_regular_10
        ws_m.cell(r_num, 4).fill = fill_card_bg
        ws_m.cell(r_num, 4).alignment = align_left_wrap
        
        ws_m.merge_cells(start_row=r_num, start_column=8, end_row=r_num, end_column=9)
        ws_m.cell(r_num, 8, c_title).font = font_bold_10
        ws_m.cell(r_num, 8).fill = fill_card_bg
        ws_m.cell(r_num, 8).alignment = align_left_wrap
        
        ws_m.merge_cells(start_row=r_num, start_column=10, end_row=r_num, end_column=13)
        ws_m.cell(r_num, 10, c_desc).font = font_regular_10
        ws_m.cell(r_num, 10).fill = fill_card_bg
        ws_m.cell(r_num, 10).alignment = align_left_wrap
        
        for col_idx in range(2, 14):
            ws_m.cell(r_num, col_idx).border = border_thin

    # =========================================================================
    # 3. SHEET: SIZING TIER MODERATE & CRITICAL GENERATOR
    # =========================================================================
    def populate_sizing_sheet(ws_s, tier_name, sizing_rows, banner_title, banner_fill):
        ws_s.views.sheetView[0].showGridLines = True
        col_widths_s = {
            'A': 4, 'B': 6, 'C': 24, 'D': 38, 'E': 14, 'F': 14,
            'G': 10, 'H': 10, 'I': 10, 'J': 12,
            'K': 12, 'L': 12, 'M': 14,
            'N': 18, 'O': 18, 'P': 18, 'Q': 20, 'R': 22
        }
        for c_let, w in col_widths_s.items():
            ws_s.column_dimensions[c_let].width = w
            
        ws_s.merge_cells('B1:R1')
        ws_s.cell(1, 2, banner_title).font = font_title_bold
        for c in range(2, 19):
            ws_s.cell(1, c).fill = banner_fill
        ws_s.row_dimensions[1].height = 26
        
        headers = [
            "No", "Group", "Description", "Edition", "Environment",
            "# of VM", "CPU", "Memory", "Storage",
            "Total CPU", "Total Memory", "Total Storage",
            "Pricing CPU", "Pricing Storage", "Pricing Backup (NBU)", "Pricing Baremetal", "Total (IDR / Thn)"
        ]
        ws_s.row_dimensions[2].height = 26
        for idx, h in enumerate(headers, start=2):
            c = ws_s.cell(2, idx, h)
            c.font = font_header_white
            c.fill = fill_blue_hdr
            c.alignment = align_header
            c.border = border_header
            
        current_row = 3
        data_rows_indices = []
        for item in sizing_rows:
            if item == "SEPARATOR_DC2":
                ws_s.row_dimensions[current_row].height = 10
                for c in range(2, 19):
                    ws_s.cell(current_row, c).fill = fill_blue_hdr
                current_row += 1
                continue
            elif item == "SEPARATOR_NONPROD":
                ws_s.row_dimensions[current_row].height = 10
                for c in range(2, 19):
                    ws_s.cell(current_row, c).fill = fill_gold_hdr
                current_row += 1
                continue
                
            no, grp, desc, ed, env, vm, cpu, mem, stg, is_bm, env_type = item
            ws_s.row_dimensions[current_row].height = 20
            bg_fill = fill_dc1_dc2 if env_type in ("DC1", "DC2") else fill_pt_dev
            
            f_tot_cpu = f"=G{current_row}*H{current_row}"
            f_tot_mem = f"=G{current_row}*I{current_row}"
            f_tot_stg = f"=G{current_row}*J{current_row}"
            
            if env_type in ("DC1", "DC2"):
                f_p_cpu = f"=K{current_row}*Pricing!$D$3"
                f_p_stg = f"=M{current_row}*Pricing!$D$5"
            else:
                f_p_cpu = f"=K{current_row}*Pricing!$D$4"
                f_p_stg = f"=M{current_row}*Pricing!$D$7"
                
            f_p_nbu = f"=M{current_row}*0.5*0.8*Pricing!$D$8"
            f_p_bm = f"=G{current_row}*Pricing!$D$10" if is_bm else 0
            f_grand_tot = f"=N{current_row}+O{current_row}+P{current_row}+Q{current_row}"
            
            ws_s.cell(current_row, 2, no).alignment = align_center
            ws_s.cell(current_row, 3, grp).alignment = align_left
            ws_s.cell(current_row, 4, desc).alignment = align_left
            ws_s.cell(current_row, 5, ed).alignment = align_center
            ws_s.cell(current_row, 6, env).alignment = align_center
            
            ws_s.cell(current_row, 7, vm).alignment = align_center
            ws_s.cell(current_row, 8, cpu).alignment = align_center
            ws_s.cell(current_row, 9, mem).alignment = align_center
            ws_s.cell(current_row, 10, stg).alignment = align_center
            
            for col_idx, (formula_val, is_curr) in enumerate([
                (f_tot_cpu, False), (f_tot_mem, False), (f_tot_stg, False),
                (f_p_cpu, True), (f_p_stg, True), (f_p_nbu, True), (f_p_bm, True), (f_grand_tot, True)
            ], start=11):
                cell = ws_s.cell(current_row, col_idx, formula_val)
                cell.alignment = align_right
                cell.number_format = num_fmt_curr if is_curr else num_fmt_int
                
            for col_idx in range(2, 19):
                c_obj = ws_s.cell(current_row, col_idx)
                c_obj.fill = bg_fill
                c_obj.border = border_thin
                c_obj.font = font_regular_10
                
            data_rows_indices.append(current_row)
            current_row += 1
            
        # Total Summary Row
        ws_s.row_dimensions[current_row].height = 24
        ws_s.merge_cells(start_row=current_row, start_column=2, end_row=current_row, end_column=6)
        c_lbl = ws_s.cell(current_row, 2, "TOTAL BIAYA INFRASTRUKTUR / TAHUN")
        c_lbl.font = font_bold_11
        c_lbl.alignment = align_center
        
        sum_vm = f"=SUM(G{data_rows_indices[0]}:G{data_rows_indices[-1]})"
        sum_cpu = f"=SUM(K{data_rows_indices[0]}:K{data_rows_indices[-1]})"
        sum_mem = f"=SUM(L{data_rows_indices[0]}:L{data_rows_indices[-1]})"
        sum_stg = f"=SUM(M{data_rows_indices[0]}:M{data_rows_indices[-1]})"
        sum_p_cpu = f"=SUM(N{data_rows_indices[0]}:N{data_rows_indices[-1]})"
        sum_p_stg = f"=SUM(O{data_rows_indices[0]}:O{data_rows_indices[-1]})"
        sum_p_nbu = f"=SUM(P{data_rows_indices[0]}:P{data_rows_indices[-1]})"
        sum_p_bm = f"=SUM(Q{data_rows_indices[0]}:Q{data_rows_indices[-1]})"
        sum_grand = f"=SUM(R{data_rows_indices[0]}:R{data_rows_indices[-1]})"
        
        ws_s.cell(current_row, 7, sum_vm).number_format = num_fmt_int
        ws_s.cell(current_row, 11, sum_cpu).number_format = num_fmt_int
        ws_s.cell(current_row, 12, sum_mem).number_format = num_fmt_int
        ws_s.cell(current_row, 13, sum_stg).number_format = num_fmt_int
        
        ws_s.cell(current_row, 14, sum_p_cpu).number_format = num_fmt_curr
        ws_s.cell(current_row, 15, sum_p_stg).number_format = num_fmt_curr
        ws_s.cell(current_row, 16, sum_p_nbu).number_format = num_fmt_curr
        ws_s.cell(current_row, 17, sum_p_bm).number_format = num_fmt_curr
        ws_s.cell(current_row, 18, sum_grand).number_format = num_fmt_curr
        
        for col_idx in range(2, 19):
            cell = ws_s.cell(current_row, col_idx)
            cell.font = font_bold_11
            cell.fill = fill_subtotal
            cell.border = border_subtotal

    # Moderate Sizing Data (Exactly matching IHSP STM2601010 & RMT SATU R1 tables)
    moderate_rows = [
        # DC1 (Production Wisdan - Ref Dokumen IHSP & RMT)
        (1, "Web Ingress (NGINX)", "NGINX Ingress Node 1 (8/16/100)", "Community", "DC1", 1, 8, 16, 100, False, "DC1"),
        (2, "Web Ingress (NGINX)", "NGINX Ingress Node 2 (8/16/100)", "Community", "DC1", 1, 8, 16, 100, False, "DC1"),
        (3, "App Platform (OCP)", "Red Hat OCP Worker 1W (64/1024/300)", "Enterprise", "DC1", 1, 64, 1024, 300, False, "DC1"),
        (4, "OCP Infra & Master", "OCP Master & Infra Router Quorum", "Enterprise", "DC1", 1, 8, 32, 150, False, "DC1"),
        (5, "Database (Oracle VM)", "Oracle DB Single Node VM Primary (16/32/1000)", "Enterprise", "DC1", 1, 16, 32, 1000, False, "DC1"),
        (6, "Message Bus (Kafka)", "Confluent Kafka Shared Broker (16/64/600)", "Community", "DC1", 3, 16, 64, 600, False, "DC1"),
        (7, "Caching (Redis)", "Redis Cache Standalone (8/16/200)", "Community", "DC1", 1, 8, 16, 200, False, "DC1"),
        (8, "Storage MFT (SFTP)", "SFTP Gateway Server (8/16/1500)", "Community", "DC1", 1, 8, 16, 1500, False, "DC1"),
        (9, "Log Shipper (Logstash)", "ELK Server - Logstash Shipper (8/32/100)", "Community", "DC1", 1, 8, 32, 100, False, "DC1"),
        (10, "Log Data (Elasticsearch)", "ELK Server - Master & Data (8/32/500)", "Community", "DC1", 3, 8, 32, 500, False, "DC1"),
        (11, "Monitoring (Prometheus)", "Monitoring Server Grafana & Prometheus (16/32/300)", "Community", "DC1", 1, 16, 32, 300, False, "DC1"),
        
        "SEPARATOR_DC2",
        # DRC (Production Rempoa - Standby Replica)
        (12, "Web Ingress (NGINX)", "NGINX Ingress Node 1 DRC (8/16/100)", "Community", "DC2", 1, 8, 16, 100, False, "DC2"),
        (13, "Web Ingress (NGINX)", "NGINX Ingress Node 2 DRC (8/16/100)", "Community", "DC2", 1, 8, 16, 100, False, "DC2"),
        (14, "App Platform (OCP)", "Red Hat OCP Worker 1W DRC (64/1024/300)", "Enterprise", "DC2", 1, 64, 1024, 300, False, "DC2"),
        (15, "Database (Oracle VM)", "Oracle DB Single Node VM Standby (16/32/1000)", "Enterprise", "DC2", 1, 16, 32, 1000, False, "DC2"),
        (16, "Message Bus (Kafka)", "Confluent Kafka Standby Broker (16/64/600)", "Community", "DC2", 3, 16, 64, 600, False, "DC2"),
        (17, "Caching (Redis)", "Redis Cache Standby (8/16/200)", "Community", "DC2", 1, 8, 16, 200, False, "DC2"),
        (18, "Storage MFT (SFTP)", "SFTP Gateway Standby (8/16/1500)", "Community", "DC2", 1, 8, 16, 1500, False, "DC2"),
        (19, "Log Data (Elasticsearch)", "ELK Server Standby Replica (8/32/500)", "Community", "DC2", 3, 8, 32, 500, False, "DC2"),
        
        "SEPARATOR_NONPROD",
        # Non-Prod (PT & DEV - Ref IHSP Tabel Nested 3)
        (20, "Web Ingress (NGINX PT)", "NGINX Ingress PT (8/16/100)", "Community", "DEV", 2, 8, 16, 100, False, "DEV"),
        (21, "App Platform (OCP PT)", "OCP Worker PT Pods (64/1024/300)", "Enterprise", "DEV", 1, 64, 1024, 300, False, "DEV"),
        (22, "Database (Oracle PT)", "Oracle DB PT Instance (16/32/1000)", "Enterprise", "DEV", 1, 16, 32, 1000, False, "DEV"),
        (23, "Message Bus (Kafka PT)", "Kafka Broker PT (16/64/600)", "Community", "DEV", 3, 16, 64, 600, False, "DEV"),
        (24, "Caching (Redis PT)", "Redis Cache PT (8/16/200)", "Community", "DEV", 1, 8, 16, 200, False, "DEV"),
        (25, "App Platform (OCP DEV)", "OCP Worker DEV Pods (64/1024/300)", "Enterprise", "DEV", 1, 64, 1024, 300, False, "DEV")
    ]
    populate_sizing_sheet(
        ws_moderate, "Moderate Tier", moderate_rows,
        "BILL OF QUANTITY (BOQ) & SIZING - TIER MODERATE (IHSP & RMT)",
        fill_navy
    )

    # Critical Sizing Data (Exactly matching Livin' Merchant DRB2502010 R6 Nested Table 108/112)
    # Includes Bare Metal Servers, Oracle RAC 64-Core, and dedicated Elasticsearch Master/Hot/Warm/Coord & Logstash!
    critical_rows = [
        # DC Wisdan (Active Production - Ref LM Nested Table 108)
        (1, "Web Ingress (Nginx Plus)", "PVLMNG01ABD24LP - Nginx Plus Gateway", "Enterprise", "DC1", 1, 8, 16, 100, False, "DC1"),
        (2, "Web Ingress (Nginx Plus)", "PVLMNG02ABD24LP - Nginx Plus Gateway", "Enterprise", "DC1", 1, 8, 16, 100, False, "DC1"),
        (3, "OCP Worker (Baremetal)", "PPWRKR01LMC25LP - OCP BM Worker Node 1", "Enterprise", "DC1", 1, 64, 1024, 300, True, "DC1"),
        (4, "OCP Worker (Baremetal)", "PPWRKR02LMC25LP - OCP BM Worker Node 2", "Enterprise", "DC1", 1, 64, 1024, 300, True, "DC1"),
        (5, "Oracle DB RAC 1 (BM)", "PPLMDB1DRX25LP - Oracle RAC BM Node 1", "Enterprise", "DC1", 1, 64, 1024, 500, True, "DC1"),
        (6, "Oracle DB RAC 2 (BM)", "PPLMDB2DRX25LP - Oracle RAC BM Node 2", "Enterprise", "DC1", 1, 64, 1024, 500, True, "DC1"),
        (7, "Oracle DB ASM SAN Storage", "Oracle ASM SAN Storage Tier 2 (Data 5.0 TB)", "Enterprise", "DC1", 1, 0, 0, 5000, False, "DC1"),
        (8, "Redis Enterprise Shard 1", "PVLMRD01ABD24LP - Redis Enterprise Node 1", "Enterprise", "DC1", 1, 8, 32, 200, False, "DC1"),
        (9, "Redis Enterprise Shard 2", "PVLMRD02ABD24LP - Redis Enterprise Node 2", "Enterprise", "DC1", 1, 8, 32, 200, False, "DC1"),
        (10, "Redis Enterprise Shard 3", "PVLMRD03ABD24LP - Redis Enterprise Node 3", "Enterprise", "DC1", 1, 8, 32, 200, False, "DC1"),
        (11, "Kafka Zookeeper Cluster", "PVLMZK01-03 - Confluent Kafka Zookeeper", "Enterprise", "DC1", 3, 4, 32, 100, False, "DC1"),
        (12, "Kafka Broker Cluster", "PVLMKB01-03 - Confluent Kafka Broker", "Enterprise", "DC1", 3, 16, 32, 1000, False, "DC1"),
        (13, "Kafka Connect C3", "PVLMKC01ABD24LP - Confluent Control Center", "Enterprise", "DC1", 1, 8, 16, 100, False, "DC1"),
        (14, "Elasticsearch Master", "PVLMEM01-03 - Elasticsearch Master Node", "Enterprise", "DC1", 3, 8, 32, 50, False, "DC1"),
        (15, "Elasticsearch Data Hot", "PVLMDH01-02 - Elasticsearch Data Hot Node", "Enterprise", "DC1", 2, 16, 32, 1900, False, "DC1"),
        (16, "Elasticsearch Data Warm", "PVLMDW01 - Elasticsearch Data Warm Node", "Enterprise", "DC1", 1, 16, 32, 1900, False, "DC1"),
        (17, "Elasticsearch Coordinator", "PVLMEC01-02 - Elasticsearch Coordinator Node", "Enterprise", "DC1", 2, 4, 32, 50, False, "DC1"),
        (18, "Logstash Shipper", "PVLMEL01-02 - Logstash Shipper Node", "Enterprise", "DC1", 2, 8, 32, 100, False, "DC1"),
        (19, "Batch Server Control-M", "PVLMBC01ABD25LP - Batch Automation Server", "Enterprise", "DC1", 1, 4, 8, 100, False, "DC1"),
        
        "SEPARATOR_DC2",
        # DRC Rempoa (Active/Standby - Ref LM Nested Table 108)
        (20, "Web Ingress (Nginx Plus)", "RVLMNG01ABD24LP - Nginx Plus Gateway DRC", "Enterprise", "DC2", 1, 8, 16, 100, False, "DC2"),
        (21, "Web Ingress (Nginx Plus)", "RVLMNG02ABD24LP - Nginx Plus Gateway DRC", "Enterprise", "DC2", 1, 8, 16, 100, False, "DC2"),
        (22, "OCP Worker (Baremetal)", "RPWRKR01LMC25LP - OCP BM Worker Node 1 DRC", "Enterprise", "DC2", 1, 64, 1024, 300, True, "DC2"),
        (23, "OCP Worker (Baremetal)", "RPWRKR02LMC25LP - OCP BM Worker Node 2 DRC", "Enterprise", "DC2", 1, 64, 1024, 300, True, "DC2"),
        (24, "Oracle DB RAC 1 (BM)", "RPLMDB1DRX25LP - Oracle RAC BM Node 1 DRC", "Enterprise", "DC2", 1, 64, 1024, 500, True, "DC2"),
        (25, "Oracle DB RAC 2 (BM)", "RPLMDB2DRX25LP - Oracle RAC BM Node 2 DRC", "Enterprise", "DC2", 1, 64, 1024, 500, True, "DC2"),
        (26, "Oracle DB ASM SAN Storage", "Oracle ASM SAN Storage Tier 2 DRC (Data 5.0 TB)", "Enterprise", "DC2", 1, 0, 0, 5000, False, "DC2"),
        (27, "Redis Enterprise Shard", "RVLMRD01-03 - Redis Enterprise Replica DRC", "Enterprise", "DC2", 3, 8, 32, 200, False, "DC2"),
        (28, "Kafka Zookeeper Cluster", "RVLMZK01-03 - Confluent Kafka Zookeeper DRC", "Enterprise", "DC2", 3, 4, 32, 100, False, "DC2"),
        (29, "Kafka Broker Cluster", "RVLMKB01-03 - Confluent Kafka Broker DRC", "Enterprise", "DC2", 3, 16, 32, 1000, False, "DC2"),
        (30, "Elasticsearch Master", "RVLMEM01-03 - Elasticsearch Master Node DRC", "Enterprise", "DC2", 3, 8, 32, 50, False, "DC2"),
        (31, "Elasticsearch Data Hot", "RVLMDH01-02 - Elasticsearch Data Hot Node DRC", "Enterprise", "DC2", 2, 16, 32, 1900, False, "DC2"),
        (32, "Elasticsearch Data Warm", "RVLMDW01 - Elasticsearch Data Warm Node DRC", "Enterprise", "DC2", 1, 16, 32, 1900, False, "DC2"),
        (33, "Logstash Shipper", "RVLMEL01-02 - Logstash Shipper Node DRC", "Enterprise", "DC2", 2, 8, 32, 100, False, "DC2"),
        
        "SEPARATOR_NONPROD",
        # Non-Prod (PT, DEV & UAT)
        (34, "Web Ingress (Nginx PT)", "Nginx Plus PT Gateway", "Enterprise", "DEV", 2, 8, 16, 100, False, "DEV"),
        (35, "OCP Worker (PT Pod)", "Red Hat OCP Worker PT Pods", "Enterprise", "DEV", 1, 16, 32, 200, False, "DEV"),
        (36, "Oracle DB Standalone (PT)", "Oracle DB PT Standalone Instance (VM)", "Enterprise", "DEV", 1, 16, 32, 1000, False, "DEV"),
        (37, "Kafka Broker Non-Prod", "Kafka Broker Non-Prod Cluster", "Enterprise", "DEV", 3, 4, 16, 300, False, "DEV"),
        (38, "Elasticsearch Non-Prod", "Elasticsearch PT Cluster Node", "Enterprise", "DEV", 2, 8, 16, 200, False, "DEV"),
        (39, "Logstash Non-Prod", "Logstash PT Shipper Node", "Enterprise", "DEV", 1, 4, 8, 100, False, "DEV")
    ]
    populate_sizing_sheet(
        ws_critical, "Critical Tier", critical_rows,
        "BILL OF QUANTITY (BOQ) & SIZING - TIER CRITICAL (LIVIN' MERCHANT)",
        fill_gold_hdr
    )

    # =========================================================================
    # 4. SHEET: PRICING (Master pricing reference - 100% clone dari Template Sizing.xlsx)
    # =========================================================================
    ws_p = ws_pricing
    ws_p.views.sheetView[0].showGridLines = True
    
    pricing_ref_path = "/Users/mcomacbook/Documents/03_PRIBADI/diagram-ddl/Template Sizing.xlsx"
    if os.path.exists(pricing_ref_path):
        wb_ref = openpyxl.load_workbook(pricing_ref_path)
        ws_ref_p = wb_ref['Pricing']
        
        # Copy column dimensions
        for col, dim in ws_ref_p.column_dimensions.items():
            if dim.width:
                ws_p.column_dimensions[col].width = dim.width
        # Ensure minimum comfortable widths
        ws_p.column_dimensions['A'].width = max(ws_p.column_dimensions['A'].width or 0, 9.5)
        ws_p.column_dimensions['B'].width = max(ws_p.column_dimensions['B'].width or 0, 22.0)
        ws_p.column_dimensions['C'].width = max(ws_p.column_dimensions['C'].width or 0, 16.0)
        ws_p.column_dimensions['D'].width = max(ws_p.column_dimensions['D'].width or 0, 20.0)
        
        for row, dim in ws_ref_p.row_dimensions.items():
            if dim.height:
                ws_p.row_dimensions[row].height = dim.height
                
        # Copy cells, values, formulas, and formatting
        for row in ws_ref_p.iter_rows():
            for cell in row:
                new_cell = ws_p[cell.coordinate]
                new_cell.value = cell.value
                if cell.has_style:
                    new_cell.font = copy(cell.font)
                    new_cell.border = copy(cell.border)
                    new_cell.fill = copy(cell.fill)
                    new_cell.number_format = copy(cell.number_format)
                    new_cell.protection = copy(cell.protection)
                    new_cell.alignment = copy(cell.alignment)
                    
        # Apply exact colors to match Template Sizing.xlsx in all Excel versions
        # Gold/Yellow Theme 7 Accent 4 (#FFC000) for Section Banners
        fill_gold_hdr_exact = PatternFill(fill_type='solid', start_color='FFFFC000', end_color='FFFFC000')
        # Navy Blue (#002060) for Table Column Headers
        fill_navy_hdr_exact = PatternFill(fill_type='solid', start_color='FF002060', end_color='FF002060')
        font_hdr_gold_exact = Font(name='Calibri', size=12, bold=True, color='FF000000')
        font_hdr_navy_exact = Font(name='Calibri', size=12, bold=True, color='FFFFFFFF')
        
        for c in range(1, 5):
            # Row 1 (Hardware banner)
            c1 = ws_p.cell(1, c)
            c1.fill = fill_gold_hdr_exact
            c1.font = font_hdr_gold_exact
            
            # Row 2 (Hardware headers: No, Item, Env, Price)
            c2 = ws_p.cell(2, c)
            c2.fill = fill_navy_hdr_exact
            c2.font = font_hdr_navy_exact
            
            # Row 14 (License banner)
            c14 = ws_p.cell(14, c)
            c14.fill = fill_gold_hdr_exact
            c14.font = font_hdr_gold_exact
            
            # Row 15 (License headers: No, Item, Env, Price)
            c15 = ws_p.cell(15, c)
            c15.fill = fill_navy_hdr_exact
            c15.font = font_hdr_navy_exact
            
        # Copy merged cell ranges
        for merge in ws_ref_p.merged_cells.ranges:
            ws_p.merge_cells(str(merge))

    # =========================================================================
    # 5. SHEET: FORMULA & GUIDELINES
    # =========================================================================
    ws_f = ws_formula
    ws_f.views.sheetView[0].showGridLines = True
    col_widths_f = {'A': 4, 'B': 6, 'C': 26, 'D': 32, 'E': 40, 'F': 44}
    for c_let, w in col_widths_f.items():
        ws_f.column_dimensions[c_let].width = w
        
    ws_f.merge_cells('B1:F1')
    ws_f.cell(1, 2, "KAMUS PERHITUNGAN SIZING & FORMULA BIAYA INFRASTRUKTUR").font = font_title_bold
    for c in range(2, 7):
        ws_f.cell(1, c).fill = fill_navy
    ws_f.row_dimensions[1].height = 26
    
    headers_f = ["No", "Parameter Sizing", "Formula Perhitungan", "Logika Kalkulasi", "Ketentuan Arsitektur & Best Practice"]
    ws_f.row_dimensions[2].height = 26
    for idx, h in enumerate(headers_f, start=2):
        c = ws_f.cell(2, idx, h)
        c.font = font_header_white
        c.fill = fill_blue_hdr
        c.alignment = align_header
        c.border = border_header
        
    formula_guide = [
        (1, "Total Compute vCPU", "# of VM * vCPU", "Perkalian jumlah VM dengan alokasi vCPU per node.", "Moderate NGINX: 2 VM x 8 vCPU = 16 vCPU. Critical NGINX: 2 VM x 8-16 vCPU."),
        (2, "Total Memory RAM", "# of VM * RAM (GB)", "Perkalian jumlah VM dengan alokasi Memory (GB) per node.", "Moderate NGINX: 2 x 16GB = 32GB RAM. Critical NGINX: 2 x 16-32GB RAM."),
        (3, "Total Storage Disk", "# of VM * Disk (GB)", "Perkalian jumlah VM dengan kapasitas storage disk (GB) per node.", "Standar root OS & file logging: 100GB - 200GB per VM."),
        (4, "Pemisahan Logstash & Elasticsearch", "Node Segregation (Tier Critical)", "Logstash shipper dipisahkan ke server tersendiri di luar cluster Elasticsearch.", "Tier Critical memisahkan Logstash (PVLMEL) dari Elasticsearch Master (PVLMEM) dan Hot Data (PVLMDH) agar resource parsing Grok tidak memicu Out of Memory (OOM) pada heap JVM Elasticsearch."),
        (5, "Biaya vCPU Tahunan", "Total vCPU * Tarif vCPU", "Total vCPU dikalikan tarif standar tahunan (Production atau Dev).", "Mengacu sheet Pricing: Production Rp 8.250.000/vCPU/thn (=7.500.000*110%), Dev Rp 7.150.000/vCPU/thn (=6.500.000*110%)."),
        (6, "Biaya Storage Tahunan", "Total Storage * Tarif Storage", "Total GB dikalikan tarif Storage SAN Tier 2 tahunan.", "Mengacu sheet Pricing: Production Tier 2 Rp 33.000/GB/thn, Dev Rp 28.600/GB/thn."),
        (7, "Biaya Backup Netbackup", "Storage * 50% * 80% * Tarif NBU", "Storage dikalikan 50% (asumsi data aktif), 80% (deduplikasi), dan tarif Rp 27.000/GB.", "Formula standar Enterprise Architecture untuk proteksi dan pencadangan data berkala."),
        (8, "Biaya Pengadaan Bare Metal", "# Unit BM * Tarif Server", "Jumlah unit physical server dikalikan tarif pengadaan Rp 794.000.000.", "Digunakan khusus Tier Critical (Oracle RAC Bare Metal & Dedicated OCP Worker BM)."),
        (9, "Pola Database: Single Node vs RAC", "Pola Arsitektur Database", "Moderate: Single Node VM Data Guard | Critical: Oracle RAC Bare Metal", "Tier Moderate mengoptimalkan biaya lisensi core database melalui Single Node Data Guard; Tier Critical wajib Oracle RAC Bare Metal untuk menangani beban transaksi tinggi dengan zero downtime failover.")
    ]
    
    for r_idx, (no, kmp, nts, lgk, pnj) in enumerate(formula_guide, start=3):
        ws_f.row_dimensions[r_idx].height = 32
        ws_f.cell(r_idx, 2, no).alignment = align_center
        ws_f.cell(r_idx, 3, kmp).alignment = align_left
        ws_f.cell(r_idx, 4, nts).alignment = align_left
        ws_f.cell(r_idx, 5, lgk).alignment = align_left_wrap
        ws_f.cell(r_idx, 6, pnj).alignment = align_left_wrap
        
        for c in range(2, 7):
            cell = ws_f.cell(r_idx, c)
            cell.font = font_regular_10
            cell.border = border_thin
            if r_idx % 2 == 1:
                cell.fill = fill_zebra

    # Save workbook
    wb.save(output_file)
    print(f"Successfully generated template workbook at: {output_file}")

if __name__ == "__main__":
    out_path = "/Users/mcomacbook/Documents/03_PRIBADI/diagram-ddl/00_TEMPLATE/Template_Tech_Stack_Architecture_Sizing.xlsx"
    build_template_tech_stack_workbook(out_path)
