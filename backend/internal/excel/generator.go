package excel

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"enterprise/sizing-apps-backend/internal/model"
)

// GenerateQuoteWorkbook creates a 100% database-driven spreadsheet matching Sizing_SCM_Inhouse_Architecture.xlsx:
// 1. Cover & HLD Architecture (Dynamic from Catalog Components & Baselines in DB)
// 2. Opsi Sizing SCM (Dynamic from user selection & live pricing lookup formulas)
// 3. Pricing (100% Dynamic from PricingItem table in DB)
// 4. Formula (Enterprise Architecture standard architecture sizing rules)
func GenerateQuoteWorkbook(
	quote model.QuoteResponse,
	rawReq model.QuoteRequest,
	projectName string,
	components []model.CatalogComponent,
	pricing []model.PricingItem,
) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	if projectName == "" {
		projectName = "SCM INHOUSE (SUPPLY CHAIN MANAGEMENT)"
	}

	// 1. Setup Sheet Names
	coverSheet := "Cover & HLD Architecture"
	sizingSheet := "Opsi Sizing SCM"
	pricingSheet := "Pricing"
	formulaSheet := "Formula"

	f.SetSheetName("Sheet1", coverSheet)
	f.NewSheet(sizingSheet)
	f.NewSheet(pricingSheet)
	f.NewSheet(formulaSheet)

	// Build Component & Pricing Lookups
	compMap := make(map[string]model.CatalogComponent)
	for _, c := range components {
		compMap[c.ID.String()] = c
		compMap[fmt.Sprintf("layer-%d", c.LayerNo)] = c
	}

	// Build Styles
	styles := createWorkbookStyles(f)

	// 1. Build Pricing Sheet first (Returns dynamic SKU-to-Row mapping)
	skuRowMap := buildPricingSheet(f, pricingSheet, pricing, styles)

	// 2. Build Cover Sheet (Dynamic from components in DB)
	buildCoverSheet(f, coverSheet, quote, projectName, compMap, styles)

	// 3. Build Sizing Sheet (Dynamic formulas linked to exact rows in Pricing Sheet)
	buildSizingSheet(f, sizingSheet, quote, rawReq, compMap, skuRowMap, styles)

	// 4. Build Formula Sheet
	buildFormulaSheet(f, formulaSheet, styles)

	// Set Cover as default active sheet
	coverIdx, _ := f.GetSheetIndex(coverSheet)
	f.SetActiveSheet(coverIdx)

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// -----------------------------------------------------------------------------
// STYLES & THEME (Bank Enterprise Deep Navy, Corporate Blue, Gold/Rust)
// -----------------------------------------------------------------------------

type workbookStyles struct {
	fontCardTitle int
	fontSubTitle  int
	fontSecTitle  int

	fillNavy      int
	fillBlueHdr   int
	fillOrangeTtl int
	fillSubtotal  int
	fillZebra     int

	cellLeft      int
	cellCenter    int
	cellRightInt  int
	cellRightCurr int

	cardVal       int
	cardSub       int
	kpiHeaderNavy int
	kpiHeaderBlue int
	kpiHeaderAcc  int
	kpiHeaderGold int

	subtotalLabel int
	subtotalCurr  int
	grandTotalTtl int
	grandTotalVal int
}

func createWorkbookStyles(f *excelize.File) *workbookStyles {
	numFmtCurr := `"Rp "#,##0`
	numFmtInt := `#,##0`

	borderThin := []excelize.Border{
		{Type: "top", Color: "D0D5DD", Style: 1},
		{Type: "bottom", Color: "D0D5DD", Style: 1},
		{Type: "left", Color: "D0D5DD", Style: 1},
		{Type: "right", Color: "D0D5DD", Style: 1},
	}
	borderHeader := []excelize.Border{
		{Type: "top", Color: "305496", Style: 1},
		{Type: "bottom", Color: "305496", Style: 2},
		{Type: "left", Color: "FFFFFF", Style: 1},
		{Type: "right", Color: "FFFFFF", Style: 1},
	}
	borderSubtotal := []excelize.Border{
		{Type: "top", Color: "305496", Style: 1},
		{Type: "bottom", Color: "305496", Style: 6}, // Double bottom
		{Type: "left", Color: "D0D5DD", Style: 1},
		{Type: "right", Color: "D0D5DD", Style: 1},
	}

	st := &workbookStyles{}

	st.fillNavy, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Color: "FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1B365D"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	st.fontCardTitle, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1B365D"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	st.fontSubTitle, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "E0E8F5", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1B365D"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	st.fontSecTitle, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13, Color: "1B365D", Family: "Calibri"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	st.fillBlueHdr, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"305496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    borderHeader,
	})
	st.fillOrangeTtl, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "000000", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"FFC000"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    borderThin,
	})

	st.cellLeft, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Calibri"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    borderThin,
	})
	st.cellCenter, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Calibri"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    borderThin,
	})
	st.cellRightInt, _ = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 10, Family: "Calibri"},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &numFmtInt,
		Border:       borderThin,
	})
	st.cellRightCurr, _ = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 10, Family: "Calibri"},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &numFmtCurr,
		Border:       borderThin,
	})

	st.fillZebra, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F9FAFB"}, Pattern: 1},
		Border:    borderThin,
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	st.fillSubtotal, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"E8EEF5"}, Pattern: 1},
		Border:    borderThin,
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	st.kpiHeaderNavy, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1B365D"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	st.kpiHeaderBlue, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"305496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	st.kpiHeaderAcc, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"418AB3"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	st.kpiHeaderGold, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"C65911"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	st.cardVal, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "1B365D", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F0F4F8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    borderThin,
	})
	st.cardSub, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true, Size: 9, Color: "555555", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F0F4F8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    borderThin,
	})

	st.subtotalLabel, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "1B365D", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"E8EEF5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    borderSubtotal,
	})
	st.subtotalCurr, _ = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true, Size: 10, Color: "1B365D", Family: "Calibri"},
		Fill:         excelize.Fill{Type: "pattern", Color: []string{"E8EEF5"}, Pattern: 1},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &numFmtCurr,
		Border:       borderSubtotal,
	})
	st.grandTotalTtl, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF", Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1B365D"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    borderSubtotal,
	})
	st.grandTotalVal, _ = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF", Family: "Calibri"},
		Fill:         excelize.Fill{Type: "pattern", Color: []string{"1B365D"}, Pattern: 1},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &numFmtCurr,
		Border:       borderSubtotal,
	})

	return st
}

// -----------------------------------------------------------------------------
// 1. SHEET: COVER & HLD ARCHITECTURE (Dynamic from Database)
// -----------------------------------------------------------------------------

func buildCoverSheet(
	f *excelize.File,
	sh string,
	quote model.QuoteResponse,
	projectName string,
	compMap map[string]model.CatalogComponent,
	st *workbookStyles,
) {
	colWidths := map[string]float64{
		"A": 4, "B": 8, "C": 26, "D": 38, "E": 26,
		"F": 22, "G": 24, "H": 24, "I": 24, "J": 24,
	}
	for col, w := range colWidths {
		_ = f.SetColWidth(sh, col, col, w)
	}

	// Banner Rows 2..4
	_ = f.MergeCell(sh, "B2", "I2")
	_ = f.MergeCell(sh, "B3", "I3")
	_ = f.MergeCell(sh, "B4", "I4")

	_ = f.SetCellValue(sh, "B2", "INFRASTRUCTURE SIZING & HIGH LEVEL DESIGN (HLD)")
	_ = f.SetCellValue(sh, "B3", fmt.Sprintf("PROJECT: %s - ENTERPRISE", strings.ToUpper(projectName)))
	_ = f.SetCellValue(sh, "B4", "Modernization: Web NGINX, OCP Microservices, Core DB, Redis Sharding, Kafka & Observability")

	_ = f.SetRowHeight(sh, 2, 24)
	_ = f.SetRowHeight(sh, 3, 22)
	_ = f.SetRowHeight(sh, 4, 22)

	_ = f.SetCellStyle(sh, "B2", "I2", st.fillNavy)
	_ = f.SetCellStyle(sh, "B3", "I3", st.fontCardTitle)
	_ = f.SetCellStyle(sh, "B4", "I4", st.fontSubTitle)

	// Metadata Box (Rows 6..9)
	targetEnvDesc := fmt.Sprintf("%s (Wisdan DC1 & Rempoa DC2)", quote.Environment)
	if quote.IncludeDRC {
		targetEnvDesc += " + DRC Active-Active (2x Multiplier)"
	}

	meta := [][]string{
		{"Nama Sistem", projectName, "Lingkungan Target", targetEnvDesc},
		{"Model Arsitektur", "High Availability Active-Standby Cross-DC & Multi-Node Cluster", "Database Replication", "Oracle Real Application Clusters (RAC) + ASM Shared SAN"},
		{"Container Platform", "Red Hat OpenShift (OCP) Container Platform / Baremetal", "Middleware & Caching", "Redis Enterprise 2-Shard Cluster + Confluent Kafka"},
		{"Dokumen Acuan", "Katalog Resmi Enterprise Architecture (Template Sizing.xlsx)", "Status Dokumen", "Final Proposed Infrastructure & License Sizing"},
	}

	for idx, m := range meta {
		r := 6 + idx
		_ = f.SetRowHeight(sh, r, 22)

		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), m[0])
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), st.fillSubtotal)

		_ = f.MergeCell(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("E%d", r))
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), m[1])
		_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("E%d", r), st.cellLeft)

		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), m[2])
		_ = f.SetCellStyle(sh, fmt.Sprintf("F%d", r), fmt.Sprintf("F%d", r), st.fillSubtotal)

		_ = f.MergeCell(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("I%d", r))
		_ = f.SetCellValue(sh, fmt.Sprintf("G%d", r), m[3])
		_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("I%d", r), st.cellLeft)
	}

	// Executive Sizing KPI Summary Cards (Rows 11..14)
	_ = f.MergeCell(sh, "B11", "I11")
	_ = f.SetCellValue(sh, "B11", "RINGKASAN EKSEKUTIF ALOKASI RESOURCE & SIZING (ALL ENVIRONMENTS)")
	_ = f.SetCellStyle(sh, "B11", "I11", st.fontSecTitle)
	_ = f.SetRowHeight(sh, 11, 26)

	kpiCards := []struct {
		col     string
		title   string
		val     string
		sub     string
		styleID int
	}{
		{"C", "Total Compute Nodes", fmt.Sprintf("%d Nodes", quote.TotalNodes), "Prod & Standby/DRC", st.kpiHeaderNavy},
		{"E", "Total Compute vCPU", fmt.Sprintf("%d vCPU", quote.TotalVCPU), "Alokasi Terkonsolidasi", st.kpiHeaderBlue},
		{"G", "Total Memory RAM", fmt.Sprintf("%d GB RAM", quote.TotalRAMGB), "In-Memory Compute Pool", st.kpiHeaderAcc},
		{"I", "Total Storage Tier 2", fmt.Sprintf("%d GB (%.1f TB)", quote.TotalStorageGB, float64(quote.TotalStorageGB)/1000.0), "High-Perf Enterprise SAN", st.kpiHeaderGold},
	}

	_ = f.SetRowHeight(sh, 12, 24)
	_ = f.SetRowHeight(sh, 13, 26)
	_ = f.SetRowHeight(sh, 14, 20)

	for _, k := range kpiCards {
		cHeader := fmt.Sprintf("%s12", k.col)
		cVal := fmt.Sprintf("%s13", k.col)
		cSub := fmt.Sprintf("%s14", k.col)

		_ = f.SetCellValue(sh, cHeader, k.title)
		_ = f.SetCellStyle(sh, cHeader, cHeader, k.styleID)

		_ = f.SetCellValue(sh, cVal, k.val)
		_ = f.SetCellStyle(sh, cVal, cVal, st.cardVal)

		_ = f.SetCellValue(sh, cSub, k.sub)
		_ = f.SetCellStyle(sh, cSub, cSub, st.cardSub)
	}

	// Architecture Components & Layers Table (Rows 16..)
	_ = f.MergeCell(sh, "B16", "I16")
	_ = f.SetCellValue(sh, "B16", "MATRIKS KOMPONEN ARSITEKTUR TARGET ENTERPRISE MODERNIZATION")
	_ = f.SetCellStyle(sh, "B16", "I16", st.fontSecTitle)
	_ = f.SetRowHeight(sh, 16, 26)

	archHeaders := []string{"No", "Layer Arsitektur", "Komponen Sistem", "Spesifikasi Node Target", "Lingkungan DC", "Fungsi & Peran Teknis Sistem"}
	_ = f.SetRowHeight(sh, 17, 24)
	cols := []string{"B", "C", "D", "E", "F", "G"}
	for i, h := range archHeaders {
		cell := fmt.Sprintf("%s17", cols[i])
		_ = f.SetCellValue(sh, cell, h)
		_ = f.SetCellStyle(sh, cell, cell, st.fillBlueHdr)
	}
	_ = f.MergeCell(sh, "G17", "I17")

	startR := 18
	for i, item := range quote.Items {
		r := startR + i
		_ = f.SetRowHeight(sh, r, 22)

		specStr := fmt.Sprintf("%d Node | %d vCPU | %d GB RAM | %d GB Disk",
			item.NodeCount, item.TotalVCPU, item.TotalRAMGB, item.TotalStorageGB)

		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), i+1)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), item.LayerName)
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), item.ComponentName)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), specStr)
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), "DC1 & DC2")

		// Lookup description from DB catalog
		roleDesc := "Standard enterprise workload"
		if comp, exists := compMap[item.ComponentID]; exists {
			if comp.Description != "" {
				roleDesc = comp.Description
			}
		}

		_ = f.MergeCell(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("I%d", r))
		_ = f.SetCellValue(sh, fmt.Sprintf("G%d", r), roleDesc)

		rowStyle := st.cellLeft
		if i%2 == 1 {
			rowStyle = st.fillZebra
		}

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("E%d", r), rowStyle)
		_ = f.SetCellStyle(sh, fmt.Sprintf("F%d", r), fmt.Sprintf("F%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("I%d", r), rowStyle)
	}
}

// -----------------------------------------------------------------------------
// 2. SHEET: PRICING (100% Dynamic from PricingItem table in Database)
// -----------------------------------------------------------------------------

func buildPricingSheet(
	f *excelize.File,
	sh string,
	pricing []model.PricingItem,
	st *workbookStyles,
) map[string]int {
	colWidths := map[string]float64{"A": 4, "B": 8, "C": 36, "D": 22, "E": 24, "F": 42}
	for col, w := range colWidths {
		_ = f.SetColWidth(sh, col, col, w)
	}

	skuRowMap := make(map[string]int)

	// Filter Hardware vs License from DB
	var hwPricing []model.PricingItem
	var licPricing []model.PricingItem

	for _, p := range pricing {
		cat := strings.ToLower(p.Category)
		if cat == "license" {
			licPricing = append(licPricing, p)
		} else {
			hwPricing = append(hwPricing, p)
		}
	}

	// 1. Hardware Standard Table
	_ = f.MergeCell(sh, "B1", "F1")
	_ = f.SetCellValue(sh, "B1", "STANDAR HARGA SATUAN HARDWARE & INFRASTRUKTUR (TERBARU)")
	_ = f.SetCellStyle(sh, "B1", "F1", st.fontCardTitle)
	_ = f.SetRowHeight(sh, 1, 26)

	headersHW := []string{"No", "Item Hardware / Service", "Environment", "Price (IDR)", "Satuan / Keterangan Formula"}
	_ = f.SetRowHeight(sh, 2, 24)
	cols := []string{"B", "C", "D", "E", "F"}
	for i, h := range headersHW {
		c := fmt.Sprintf("%s2", cols[i])
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	hwRowStart := 3
	for idx, p := range hwPricing {
		r := hwRowStart + idx
		skuRowMap[p.SKU] = r

		var priceVal int64
		if p.UnitPriceIDR != nil {
			priceVal = *p.UnitPriceIDR
		}

		_ = f.SetRowHeight(sh, r, 22)
		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), idx+1)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), p.Name)
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), p.Environment)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), priceVal)

		note := p.Unit
		if p.Source != "" {
			note = fmt.Sprintf("%s (Ref %s)", p.Unit, p.Source)
		}
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), note)

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), st.cellLeft)
		_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("E%d", r), fmt.Sprintf("E%d", r), st.cellRightCurr)
		_ = f.SetCellStyle(sh, fmt.Sprintf("F%d", r), fmt.Sprintf("F%d", r), st.cellLeft)
	}

	// 2. Software License Table
	licHeaderR := hwRowStart + len(hwPricing) + 2
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", licHeaderR), fmt.Sprintf("F%d", licHeaderR))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", licHeaderR), "STANDAR HARGA LISENSI SOFTWARE (ACUAN RESMI ENTERPRISE)")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", licHeaderR), fmt.Sprintf("F%d", licHeaderR), st.fontCardTitle)
	_ = f.SetRowHeight(sh, licHeaderR, 26)

	headersLic := []string{"No", "Item Software / Platform", "Environment", "Price (IDR)", "Keterangan Lisensi"}
	_ = f.SetRowHeight(sh, licHeaderR+1, 24)
	for i, h := range headersLic {
		c := fmt.Sprintf("%s%d", cols[i], licHeaderR+1)
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	licRowStart := licHeaderR + 2
	for idx, p := range licPricing {
		r := licRowStart + idx
		skuRowMap[p.SKU] = r

		var priceVal int64
		if p.UnitPriceIDR != nil {
			priceVal = *p.UnitPriceIDR
		}

		_ = f.SetRowHeight(sh, r, 22)
		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), idx+1)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), p.Name)
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), p.Environment)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), priceVal)

		note := p.Unit
		if !p.IsConfirmed {
			note = "Belum confirm ke STA / DBA Enterprise"
		}
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), note)

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), st.cellLeft)
		_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("E%d", r), fmt.Sprintf("E%d", r), st.cellRightCurr)
		_ = f.SetCellStyle(sh, fmt.Sprintf("F%d", r), fmt.Sprintf("F%d", r), st.cellLeft)
	}

	return skuRowMap
}

// -----------------------------------------------------------------------------
// 3. SHEET: OPSI SIZING SCM (Dynamic Formulas Linked to Pricing SKU Rows)
// -----------------------------------------------------------------------------

func buildSizingSheet(
	f *excelize.File,
	sh string,
	quote model.QuoteResponse,
	rawReq model.QuoteRequest,
	compMap map[string]model.CatalogComponent,
	skuRowMap map[string]int,
	st *workbookStyles,
) {
	colWidths := map[string]float64{
		"A": 4.0, "B": 8.0, "C": 34.0, "D": 38.0, "E": 20.0,
		"F": 18.0, "G": 12.0, "H": 12.0, "I": 14.0, "J": 14.0,
		"K": 12.0, "L": 14.0, "M": 14.0, "N": 22.0, "O": 22.0,
		"P": 22.0, "Q": 20.0, "R": 24.0,
	}
	for col, w := range colWidths {
		_ = f.SetColWidth(sh, col, col, w)
	}

	// Dynamic SKU Row Lookups with Fallbacks
	isDev := strings.Contains(strings.ToLower(quote.Environment), "dev")
	vcpuSKU := "VCPU_PROD"
	if isDev {
		vcpuSKU = "VCPU_DEV"
	}
	vcpuRow := skuRowMap[vcpuSKU]
	if vcpuRow == 0 {
		vcpuRow = 3
	}

	stgT2SKU := "STG_T2_PROD"
	if isDev {
		stgT2SKU = "STG_T2_DEV"
	}
	stgT2Row := skuRowMap[stgT2SKU]
	if stgT2Row == 0 {
		stgT2Row = 6
	}

	stgT3Row := skuRowMap["STG_T3_PROD"]
	if stgT3Row == 0 {
		stgT3Row = 6
	}

	nbuRow := skuRowMap["NBU_ALL"]
	if nbuRow == 0 {
		nbuRow = 8
	}

	bmRow := skuRowMap["BM_STD_PROD"]
	if bmRow == 0 {
		bmRow = 10
	}

	// Row 1: Title Banner
	_ = f.MergeCell(sh, "B1", "R1")
	_ = f.SetCellValue(sh, "B1", "Summary Hardware & Bill of Quantity (BOQ)")
	_ = f.SetCellStyle(sh, "B1", "R1", st.fontCardTitle)
	_ = f.SetRowHeight(sh, 1, 26)

	// Row 2: Table Headers (B2..R2)
	headers := []string{
		"No", "Group", "Description", "Edition", "Environment",
		"# of VM", "CPU", "Memory", "Storage",
		"Total CPU", "Total Memory", "Total Storage",
		"Pricing CPU", "Pricing Storage", "Pricing Rubrik / NBU",
		"Total Pricing Baremetal", "Total",
	}
	_ = f.SetRowHeight(sh, 2, 26)
	cols := []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R"}
	for i, h := range headers {
		c := fmt.Sprintf("%s2", cols[i])
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	// Data Rows (Starting at Row 3)
	startR := 3
	currR := startR

	for idx, item := range quote.Items {
		r := currR
		_ = f.SetRowHeight(sh, r, 20)

		edition := "Community"
		if strings.ToLower(item.LicenseKind) == "enterprise" {
			edition = "Enterprise"
		}
		env := "DC1"

		nodeCount := item.NodeCount
		if nodeCount <= 0 {
			nodeCount = 1
		}
		vcpuPerNode := item.TotalVCPU / nodeCount
		ramPerNode := item.TotalRAMGB / nodeCount
		stgPerNode := item.TotalStorageGB / nodeCount

		// Raw spec check from request
		isBaremetal := false
		isShared := false
		isNone := false
		stgClass := "Tier 2"
		backupEnabled := true

		if idx < len(rawReq.Items) {
			raw := rawReq.Items[idx]
			if strings.ToLower(raw.ComputeType) == "baremetal" {
				isBaremetal = true
			} else if strings.ToLower(raw.ComputeType) == "shared" {
				isShared = true
			} else if strings.ToLower(raw.ComputeType) == "none" {
				isNone = true
			}
			if raw.StorageClass != "" {
				stgClass = raw.StorageClass
			}
			backupEnabled = raw.BackupEnabled
		}

		vmCount := nodeCount
		if isNone {
			vmCount = 0
			vcpuPerNode = 0
			ramPerNode = 0
			stgPerNode = 0
		}

		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), idx+1)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), item.LayerName)
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), item.ComponentName)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), edition)
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), env)
		_ = f.SetCellValue(sh, fmt.Sprintf("G%d", r), vmCount)
		_ = f.SetCellValue(sh, fmt.Sprintf("H%d", r), vcpuPerNode)
		_ = f.SetCellValue(sh, fmt.Sprintf("I%d", r), ramPerNode)
		_ = f.SetCellValue(sh, fmt.Sprintf("J%d", r), stgPerNode)

		// Dynamic Formulas
		_ = f.SetCellFormula(sh, fmt.Sprintf("K%d", r), fmt.Sprintf("H%d*G%d", r, r))
		_ = f.SetCellFormula(sh, fmt.Sprintf("L%d", r), fmt.Sprintf("I%d*G%d", r, r))
		_ = f.SetCellFormula(sh, fmt.Sprintf("M%d", r), fmt.Sprintf("J%d*G%d", r, r))

		if isBaremetal {
			_ = f.SetCellValue(sh, fmt.Sprintf("N%d", r), 0)
			_ = f.SetCellFormula(sh, fmt.Sprintf("Q%d", r), fmt.Sprintf("G%d*Pricing!$E$%d", r, bmRow))
		} else if isShared || isNone {
			_ = f.SetCellValue(sh, fmt.Sprintf("N%d", r), 0)
			_ = f.SetCellValue(sh, fmt.Sprintf("Q%d", r), 0)
		} else {
			_ = f.SetCellFormula(sh, fmt.Sprintf("N%d", r), fmt.Sprintf("K%d*Pricing!$E$%d", r, vcpuRow))
			_ = f.SetCellValue(sh, fmt.Sprintf("Q%d", r), 0)
		}

		if isNone {
			_ = f.SetCellValue(sh, fmt.Sprintf("O%d", r), 0)
			_ = f.SetCellValue(sh, fmt.Sprintf("P%d", r), 0)
		} else {
			targetStgRow := stgT2Row
			if strings.ToLower(stgClass) == "tier 3" {
				targetStgRow = stgT3Row
			}
			_ = f.SetCellFormula(sh, fmt.Sprintf("O%d", r), fmt.Sprintf("M%d*Pricing!$E$%d", r, targetStgRow))

			if backupEnabled && item.BackupCostAnnual > 0 {
				_ = f.SetCellFormula(sh, fmt.Sprintf("P%d", r), fmt.Sprintf("M%d*0.5*0.8*Pricing!$E$%d", r, nbuRow))
			} else {
				_ = f.SetCellValue(sh, fmt.Sprintf("P%d", r), 0)
			}
		}

		_ = f.SetCellFormula(sh, fmt.Sprintf("R%d", r), fmt.Sprintf("SUM(N%d:Q%d)", r, r))

		// Cell Styles
		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("E%d", r), st.cellLeft)
		_ = f.SetCellStyle(sh, fmt.Sprintf("F%d", r), fmt.Sprintf("F%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("M%d", r), st.cellRightInt)
		_ = f.SetCellStyle(sh, fmt.Sprintf("N%d", r), fmt.Sprintf("R%d", r), st.cellRightCurr)

		currR++
	}

	lastDataR := currR - 1

	// Row: Grand Total Hardware
	_ = f.SetRowHeight(sh, currR, 24)
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", currR), fmt.Sprintf("Q%d", currR))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", currR), "Total Sizing Hardware (All Environments)")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", currR), fmt.Sprintf("Q%d", currR), st.subtotalLabel)

	_ = f.SetCellFormula(sh, fmt.Sprintf("R%d", currR), fmt.Sprintf("SUM(R%d:R%d)", startR, lastDataR))
	_ = f.SetCellStyle(sh, fmt.Sprintf("R%d", currR), fmt.Sprintf("R%d", currR), st.subtotalCurr)

	// Summary Sections below main grid
	currR += 2
	summaryStartR := currR

	// 1. Summary Hardware Prod Box (Cols B..D) - 100% Dynamic from quote.Items!
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", currR), fmt.Sprintf("D%d", currR))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", currR), "Summary Hardware Prod")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", currR), fmt.Sprintf("D%d", currR), st.fillOrangeTtl)
	_ = f.SetRowHeight(sh, currR, 22)

	_ = f.SetRowHeight(sh, currR+1, 24)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", currR+1), "No")
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", currR+1), "Item")
	_ = f.SetCellValue(sh, fmt.Sprintf("D%d", currR+1), "Price")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", currR+1), fmt.Sprintf("D%d", currR+1), st.fillBlueHdr)

	subHwStart := currR + 2
	for idx, item := range quote.Items {
		r := subHwStart + idx
		_ = f.SetRowHeight(sh, r, 20)
		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), idx+1)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), item.ComponentName)
		_ = f.SetCellFormula(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("R%d", startR+idx))

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), st.cellLeft)
		_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), st.cellRightCurr)
	}

	subHwEnd := subHwStart + len(quote.Items) - 1
	subHwTotalR := subHwEnd + 1
	_ = f.SetRowHeight(sh, subHwTotalR, 24)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", subHwTotalR), "")
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", subHwTotalR), "Total")
	_ = f.SetCellFormula(sh, fmt.Sprintf("D%d", subHwTotalR), fmt.Sprintf("SUM(D%d:D%d)", subHwStart, subHwEnd))
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", subHwTotalR), fmt.Sprintf("C%d", subHwTotalR), st.subtotalLabel)
	_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", subHwTotalR), fmt.Sprintf("D%d", subHwTotalR), st.subtotalCurr)

	// 2. Data Center Resource Allocation Table (Cols G..J, beside Summary Hardware)
	_ = f.MergeCell(sh, fmt.Sprintf("G%d", summaryStartR), fmt.Sprintf("J%d", summaryStartR))
	_ = f.SetCellValue(sh, fmt.Sprintf("G%d", summaryStartR), "Matriks Alokasi Resource Data Center")
	_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", summaryStartR), fmt.Sprintf("J%d", summaryStartR), st.fillOrangeTtl)

	dcH := []string{"DC / Environment", "CPU (vCPU)", "RAM (GB)", "Storage (GB)"}
	for i, h := range dcH {
		colLet := string(rune('G' + i))
		c := fmt.Sprintf("%s%d", colLet, summaryStartR+1)
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	dcRows := []struct {
		name string
		fCPU string
		fRAM string
		fSTG string
	}{
		{"DC1 (Wisdan - Active)", fmt.Sprintf("SUM(K%d:K%d)", startR, lastDataR), fmt.Sprintf("SUM(L%d:L%d)", startR, lastDataR), fmt.Sprintf("SUM(M%d:M%d)", startR, lastDataR)},
		{"DC2 (Rempoa - Standby)", fmt.Sprintf("SUM(K%d:K%d)", startR, lastDataR), fmt.Sprintf("SUM(L%d:L%d)", startR, lastDataR), fmt.Sprintf("SUM(M%d:M%d)", startR, lastDataR)},
	}

	for idx, d := range dcRows {
		r := summaryStartR + 2 + idx
		_ = f.SetRowHeight(sh, r, 20)
		_ = f.SetCellValue(sh, fmt.Sprintf("G%d", r), d.name)
		_ = f.SetCellFormula(sh, fmt.Sprintf("H%d", r), d.fCPU)
		_ = f.SetCellFormula(sh, fmt.Sprintf("I%d", r), d.fRAM)
		_ = f.SetCellFormula(sh, fmt.Sprintf("J%d", r), d.fSTG)

		_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("G%d", r), st.cellLeft)
		_ = f.SetCellStyle(sh, fmt.Sprintf("H%d", r), fmt.Sprintf("J%d", r), st.cellRightInt)
	}

	// 3. Summary License Prod Table (Cols B..G) - 100% Dynamic from Enterprise Licenses chosen!
	licStartR := subHwTotalR + 3
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", licStartR), fmt.Sprintf("G%d", licStartR))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", licStartR), "Summary License Prod")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", licStartR), fmt.Sprintf("G%d", licStartR), st.fillOrangeTtl)
	_ = f.SetRowHeight(sh, licStartR, 22)

	licH := []string{"No", "Item Software / Platform", "#of License", "Harga Satuan", "Total Harga", "Remarks"}
	_ = f.SetRowHeight(sh, licStartR+1, 24)
	for i, h := range licH {
		colLet := string(rune('B' + i))
		c := fmt.Sprintf("%s%d", colLet, licStartR+1)
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	// Collect enterprise license items from quote.Items
	type entLicRow struct {
		name    string
		count   int
		priceF  string
		remarks string
	}
	var activeLic []entLicRow

	for idx, it := range quote.Items {
		if strings.ToLower(it.LicenseKind) == "enterprise" {
			// Find pricing SKU for this component
			var matchedSKU string
			var remarks string

			if idx < len(rawReq.Items) && rawReq.Items[idx].PricingSKU != "" {
				matchedSKU = rawReq.Items[idx].PricingSKU
				if rawReq.Items[idx].LicenseMetric != "" {
					remarks = fmt.Sprintf("Per %s", rawReq.Items[idx].LicenseMetric)
				}
			}
			if matchedSKU == "" && it.PricingSKU != "" {
				matchedSKU = it.PricingSKU
			}

			if matchedSKU == "" {
				if comp, ok := compMap[it.ComponentID]; ok {
					for _, lic := range comp.Licenses {
						if strings.ToLower(lic.Kind) == "enterprise" {
							matchedSKU = lic.PricingSKU
							remarks = fmt.Sprintf("Per %s", lic.Metric)
							break
						}
					}
				}
			}
			if remarks == "" {
				remarks = "Katalog SKU"
			}

			priceFormula := "0"
			if matchedSKU != "" && skuRowMap[matchedSKU] > 0 {
				priceFormula = fmt.Sprintf("Pricing!$E$%d", skuRowMap[matchedSKU])
			} else if it.LicenseCostAnnual > 0 {
				priceFormula = fmt.Sprintf("%d", it.LicenseCostAnnual/int64(it.NodeCount))
			}

			activeLic = append(activeLic, entLicRow{
				name:    it.LicenseProductName,
				count:   it.NodeCount,
				priceF:  priceFormula,
				remarks: remarks,
			})
		}
	}

	subLicDataStart := licStartR + 2
	if len(activeLic) == 0 {
		// No enterprise licenses selected (All Free Open Source)
		_ = f.SetRowHeight(sh, subLicDataStart, 20)
		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", subLicDataStart), 1)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", subLicDataStart), "Community / Open Source (No Paid Licenses)")
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", subLicDataStart), 0)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", subLicDataStart), 0)
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", subLicDataStart), 0)
		_ = f.SetCellValue(sh, fmt.Sprintf("G%d", subLicDataStart), "Free Software Tier")

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", subLicDataStart), fmt.Sprintf("B%d", subLicDataStart), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", subLicDataStart), fmt.Sprintf("C%d", subLicDataStart), st.cellLeft)
		_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", subLicDataStart), fmt.Sprintf("D%d", subLicDataStart), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("E%d", subLicDataStart), fmt.Sprintf("F%d", subLicDataStart), st.cellRightCurr)
		_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", subLicDataStart), fmt.Sprintf("G%d", subLicDataStart), st.cellLeft)
		activeLic = append(activeLic, entLicRow{}) // 1 dummy entry for bounds
	} else {
		for idx, l := range activeLic {
			r := subLicDataStart + idx
			_ = f.SetRowHeight(sh, r, 20)

			_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), idx+1)
			_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), l.name)
			_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), l.count)

			if strings.HasPrefix(l.priceF, "Pricing!") {
				_ = f.SetCellFormula(sh, fmt.Sprintf("E%d", r), l.priceF)
			} else {
				_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), l.priceF)
			}
			_ = f.SetCellFormula(sh, fmt.Sprintf("F%d", r), fmt.Sprintf("D%d*E%d", r, r))
			_ = f.SetCellValue(sh, fmt.Sprintf("G%d", r), l.remarks)

			_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
			_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), st.cellLeft)
			_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), st.cellCenter)
			_ = f.SetCellStyle(sh, fmt.Sprintf("E%d", r), fmt.Sprintf("F%d", r), st.cellRightCurr)
			_ = f.SetCellStyle(sh, fmt.Sprintf("G%d", r), fmt.Sprintf("G%d", r), st.cellLeft)
		}
	}

	subLicDataEnd := subLicDataStart + len(activeLic) - 1
	subLicTotalR := subLicDataEnd + 1
	_ = f.SetRowHeight(sh, subLicTotalR, 24)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", subLicTotalR), "")
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", subLicTotalR), "Total")
	_ = f.SetCellFormula(sh, fmt.Sprintf("F%d", subLicTotalR), fmt.Sprintf("SUM(F%d:F%d)", subLicDataStart, subLicDataEnd))
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", subLicTotalR), fmt.Sprintf("E%d", subLicTotalR), st.subtotalLabel)
	_ = f.SetCellStyle(sh, fmt.Sprintf("F%d", subLicTotalR), fmt.Sprintf("G%d", subLicTotalR), st.subtotalCurr)

	// 4. Total Price Prod (Hardware + License)
	totProdR := subLicTotalR + 3
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", totProdR), fmt.Sprintf("D%d", totProdR))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", totProdR), "Total Price Prod (Hardware + License)")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", totProdR), fmt.Sprintf("D%d", totProdR), st.fillOrangeTtl)
	_ = f.SetRowHeight(sh, totProdR, 22)

	_ = f.SetRowHeight(sh, totProdR+1, 24)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", totProdR+1), "No")
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", totProdR+1), "Komponen Biaya")
	_ = f.SetCellValue(sh, fmt.Sprintf("D%d", totProdR+1), "Total Harga")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", totProdR+1), fmt.Sprintf("D%d", totProdR+1), st.fillBlueHdr)

	_ = f.SetRowHeight(sh, totProdR+2, 20)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", totProdR+2), 1)
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", totProdR+2), "Hardware Compute & Storage")
	_ = f.SetCellFormula(sh, fmt.Sprintf("D%d", totProdR+2), fmt.Sprintf("D%d", subHwTotalR))
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", totProdR+2), fmt.Sprintf("B%d", totProdR+2), st.cellCenter)
	_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", totProdR+2), fmt.Sprintf("C%d", totProdR+2), st.cellLeft)
	_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", totProdR+2), fmt.Sprintf("D%d", totProdR+2), st.cellRightCurr)

	_ = f.SetRowHeight(sh, totProdR+3, 20)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", totProdR+3), 2)
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", totProdR+3), "Software License & Middleware")
	_ = f.SetCellFormula(sh, fmt.Sprintf("D%d", totProdR+3), fmt.Sprintf("F%d", subLicTotalR))
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", totProdR+3), fmt.Sprintf("B%d", totProdR+3), st.cellCenter)
	_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", totProdR+3), fmt.Sprintf("C%d", totProdR+3), st.cellLeft)
	_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", totProdR+3), fmt.Sprintf("D%d", totProdR+3), st.cellRightCurr)

	totProdEndR := totProdR + 4
	_ = f.SetRowHeight(sh, totProdEndR, 24)
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", totProdEndR), "")
	_ = f.SetCellValue(sh, fmt.Sprintf("C%d", totProdEndR), "Total")
	_ = f.SetCellFormula(sh, fmt.Sprintf("D%d", totProdEndR), fmt.Sprintf("SUM(D%d:D%d)", totProdR+2, totProdR+3))
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", totProdEndR), fmt.Sprintf("C%d", totProdEndR), st.subtotalLabel)
	_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", totProdEndR), fmt.Sprintf("D%d", totProdEndR), st.subtotalCurr)

	// 5. Grand Total Banner (Annual & 3-Year TCO)
	grandR := totProdEndR + 2
	_ = f.SetRowHeight(sh, grandR, 28)
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", grandR), fmt.Sprintf("C%d", grandR))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", grandR), "GRAND TOTAL TCO (ANNUAL OPEX)")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", grandR), fmt.Sprintf("C%d", grandR), st.grandTotalTtl)

	_ = f.SetCellFormula(sh, fmt.Sprintf("D%d", grandR), fmt.Sprintf("D%d", totProdEndR))
	_ = f.SetCellStyle(sh, fmt.Sprintf("D%d", grandR), fmt.Sprintf("D%d", grandR), st.grandTotalVal)
}

// -----------------------------------------------------------------------------
// 4. SHEET: FORMULA (Glossary & Sizing Formulas)
// -----------------------------------------------------------------------------

func buildFormulaSheet(f *excelize.File, sh string, st *workbookStyles) {
	colWidths := map[string]float64{"A": 4, "B": 6, "C": 26, "D": 32, "E": 48, "F": 46}
	for col, w := range colWidths {
		_ = f.SetColWidth(sh, col, col, w)
	}

	// Main Title Banner
	_ = f.MergeCell(sh, "B2", "F2")
	_ = f.MergeCell(sh, "B3", "F3")
	_ = f.SetCellValue(sh, "B2", "PANDUAN RUMUS SIZING & GLOSARIUM ISTILAH ARSITEKTUR SCM")
	_ = f.SetCellValue(sh, "B3", "Penjelasan Logika Perhitungan Biaya & Kamus Definisi Teknis SCM Inhouse Modernization")

	_ = f.SetRowHeight(sh, 2, 24)
	_ = f.SetRowHeight(sh, 3, 22)
	_ = f.SetCellStyle(sh, "B2", "F2", st.fillNavy)
	_ = f.SetCellStyle(sh, "B3", "F3", st.fontSubTitle)

	// Section 1: Daftar Rumus
	_ = f.MergeCell(sh, "B5", "F5")
	_ = f.SetCellValue(sh, "B5", "1. DAFTAR RUMUS & LOGIKA PERHITUNGAN BIAYA SIZING")
	_ = f.SetCellStyle(sh, "B5", "F5", st.fontSecTitle)
	_ = f.SetRowHeight(sh, 5, 26)

	headers := []string{"No", "Komponen Sizing", "Notasi Rumus", "Logika Perhitungan Sederhana", "Contoh Perhitungan Riil"}
	_ = f.SetRowHeight(sh, 6, 26)
	cols := []string{"B", "C", "D", "E", "F"}
	for i, h := range headers {
		c := fmt.Sprintf("%s6", cols[i])
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	formulasData := []struct {
		no     int
		komp   string
		rumus  string
		logika string
		contoh string
	}{
		{1, "Total Compute (vCPU)", "CPU * # of VM", "Jumlah VM dikalikan kapasitas CPU per node", "1 VM x 16 vCPU = 16 vCPU"},
		{2, "Total Memory (RAM)", "Memory * # of VM", "Jumlah VM dikalikan kapasitas RAM per node", "1 VM x 32 GB = 32 GB RAM"},
		{3, "Total Storage (GB)", "Storage * # of VM", "Jumlah VM dikalikan kapasitas Disk per node", "1 VM x 200 GB = 200 GB Disk"},
		{4, "Biaya Compute Prod", "Total CPU * Harga vCPU Prod", "Total vCPU dikalikan tarif vCPU katalog resmi", "16 vCPU x Rp 8.250.000 = Rp 132.000.000"},
		{5, "Biaya Storage Prod Tier 2", "Total Storage * Harga Tier 2", "Total Storage GB dikalikan tarif storage tahunan", "4.500 GB x Rp 33.000 = Rp 148.500.000"},
		{6, "Biaya Backup NBU", "Storage * 50% * 80% * Tarif NBU", "Storage x 50% (Data Aktif) x 80% (Dedupe) x Tarif", "4.500 GB x 0.5 x 0.8 x Rp 27.000 = Rp 48.600.000"},
		{7, "Lisensi OCP Worker Node", "Jumlah Worker * Harga OCP", "Jumlah worker node baremetal/VM dikalikan tarif OCP", "2 Node x Rp 240.000.000 = Rp 480.000.000"},
		{8, "Lisensi Redis Enterprise", "Jumlah Shard * Harga per Shard", "Jumlah shard cache RAM dikalikan harga shard", "2 Shards x Rp 125.000.000 = Rp 250.000.000"},
		{9, "Lisensi Core Oracle DB RAC", "Total vCPU DB * 0.5", "vCPU DB dikalikan Core Factor Multiplier (0.5 Intel)", "(16 + 16) vCPU x 0.5 = 16 Processor License"},
		{10, "Total Biaya per Node VM", "SUM(Compute + Storage + Backup)", "Penjumlahan seluruh komponen biaya node per tahun", "Compute + Storage + Backup + Baremetal"},
		{11, "Grand Total TCO", "Hardware + Software License", "Total seluruh investasi infrastruktur tahunan", "Total Hardware Prod + Total Lisensi Prod"},
	}

	for idx, fData := range formulasData {
		r := 7 + idx
		_ = f.SetRowHeight(sh, r, 22)
		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), fData.no)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), fData.komp)
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), fData.rumus)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), fData.logika)
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), fData.contoh)

		rowStyle := st.cellLeft
		if idx%2 == 1 {
			rowStyle = st.fillZebra
		}

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("F%d", r), rowStyle)
	}

	// Section 2: Glosarium
	sec2R := 7 + len(formulasData) + 2
	_ = f.MergeCell(sh, fmt.Sprintf("B%d", sec2R), fmt.Sprintf("F%d", sec2R))
	_ = f.SetCellValue(sh, fmt.Sprintf("B%d", sec2R), "2. GLOSARIUM & KAMUS ISTILAH TEKNIS ARSITEKTUR ENTERPRISE")
	_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", sec2R), fmt.Sprintf("F%d", sec2R), st.fontSecTitle)
	_ = f.SetRowHeight(sh, sec2R, 26)

	headersGlos := []string{"No", "Istilah Teknis", "Komponen Terkait", "Pengertian Arsitektur", "Fungsi Dalam Sistem SCM"}
	_ = f.SetRowHeight(sh, sec2R+1, 26)
	for i, h := range headersGlos {
		c := fmt.Sprintf("%s%d", cols[i], sec2R+1)
		_ = f.SetCellValue(sh, c, h)
		_ = f.SetCellStyle(sh, c, c, st.fillBlueHdr)
	}

	glosData := []struct {
		no    int
		term  string
		komp  string
		arti  string
		funcS string
	}{
		{1, "Active-Active RAC", "Oracle Database", "Dua node database berjalan bersamaan mengakses storage shared ASM", "Mencegah downtime transaksi SCM saat salah satu server database fail"},
		{2, "ASM Storage", "Storage Shared SAN", "Automatic Storage Management untuk striping LUN raw device database", "Optimasi throughput IOPS baca/tulis data transaksi dan redo log"},
		{3, "Sharding In-Memory", "Redis Enterprise", "Teknik pemecahan dataset cache ke dalam beberapa partisi memori", "Mempercepat response query hingga sub-millisecond"},
		{4, "Worker Pods", "Red Hat OpenShift", "Node komputasi Kubernetes tempat container microservices berjalan", "Autoscaling container aplikasi backend dan frontend"},
		{5, "Raw Device Mapping (RDM)", "Storage SAN Tier 2", "Pemetaan langsung storage LUN SAN ke mesin virtual", "Penyimpanan dokumen PO, invoice PDF & surat jalan transaksi"},
	}

	for idx, gData := range glosData {
		r := sec2R + 2 + idx
		_ = f.SetRowHeight(sh, r, 22)
		_ = f.SetCellValue(sh, fmt.Sprintf("B%d", r), gData.no)
		_ = f.SetCellValue(sh, fmt.Sprintf("C%d", r), gData.term)
		_ = f.SetCellValue(sh, fmt.Sprintf("D%d", r), gData.komp)
		_ = f.SetCellValue(sh, fmt.Sprintf("E%d", r), gData.arti)
		_ = f.SetCellValue(sh, fmt.Sprintf("F%d", r), gData.funcS)

		rowStyle := st.cellLeft
		if idx%2 == 1 {
			rowStyle = st.fillZebra
		}

		_ = f.SetCellStyle(sh, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), st.cellCenter)
		_ = f.SetCellStyle(sh, fmt.Sprintf("C%d", r), fmt.Sprintf("F%d", r), rowStyle)
	}
}
