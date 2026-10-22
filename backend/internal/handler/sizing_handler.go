package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"mandiri/sizing-apps-backend/internal/excel"
	"mandiri/sizing-apps-backend/internal/model"
	"mandiri/sizing-apps-backend/internal/quote"
	"mandiri/sizing-apps-backend/internal/service"
)

type SizingHandler struct {
	svc          *service.SizingService
	db           *gorm.DB
	templatesDir string

	// High-performance thread-safe in-memory cache (sub-millisecond catalog & quote calculations)
	cacheMu            sync.RWMutex
	cachedComponents   []model.CatalogComponent
	cachedPricingItems []model.PricingItem
	cachedRules        []model.LicenseRule
	cacheExpiresAt     time.Time
}

func NewSizingHandler(svc *service.SizingService, database *gorm.DB, templatesDir string) *SizingHandler {
	return &SizingHandler{
		svc:          svc,
		db:           database,
		templatesDir: templatesDir,
	}
}

func (h *SizingHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// System / Health
	app.Get("/api/health", h.HealthCheck)

	// Catalog & Pricing (Database-driven CRUD)
	api.Get("/catalog/components", h.GetCatalogComponents)
	api.Get("/pricing/items", h.GetPricingItems)
	api.Post("/pricing/items", h.CreatePricingItem)
	api.Put("/pricing/items/:id", h.UpdatePricingItem)
	api.Delete("/pricing/items/:id", h.DeletePricingItem)

	// License Rules (Database-driven CRUD)
	api.Get("/license/rules", h.GetLicenseRules)
	api.Post("/license/rules", h.CreateLicenseRule)
	api.Put("/license/rules/:id", h.UpdateLicenseRule)
	api.Delete("/license/rules/:id", h.DeleteLicenseRule)

	// Live Quote Calculation & Excel Export
	api.Post("/quotes/calculate", h.CalculateQuote)
	api.Post("/quotes/export-excel", h.ExportQuoteExcel)

	// Baseline & KPI
	api.Get("/tech-stack/baseline", h.GetTechStackBaseline)
	api.Get("/tech-stack/kpis", h.GetKPISummaries)

	// Sizing Engine Calculations
	api.Post("/sizing/calculate-power-law", h.CalculatePowerLaw)
	api.Post("/sizing/calculate-eawg", h.CalculateEAWG)

	// SolarWinds Analytics
	api.Post("/solarwinds/upload", h.UploadSolarwindsCSV)

	// Templates
	api.Get("/templates", h.ListTemplates)
	api.Get("/templates/download/:filename", h.DownloadTemplate)
}

func (h *SizingHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "Enterprise Capacity Planning & Sizing Engine (Bank Mandiri)",
		"version": "1.0.0",
	})
}

func (h *SizingHandler) GetTechStackBaseline(c *fiber.Ctx) error {
	items := h.svc.GetTechStackItems()
	return c.JSON(fiber.Map{
		"total": len(items),
		"data":  items,
	})
}

func (h *SizingHandler) GetKPISummaries(c *fiber.Ctx) error {
	kpis := h.svc.GetKPISummaries()
	return c.JSON(fiber.Map{
		"data": kpis,
	})
}

func (h *SizingHandler) CalculatePowerLaw(c *fiber.Ctx) error {
	var req model.PowerLawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	result := h.svc.CalculatePowerLaw(req)
	return c.JSON(result)
}

func (h *SizingHandler) CalculateEAWG(c *fiber.Ctx) error {
	var req model.CalculateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	result := h.svc.CalculateEAWG(req)
	return c.JSON(result)
}

func (h *SizingHandler) UploadSolarwindsCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "CSV file is required in 'file' field",
		})
	}

	metricName := c.FormValue("metric_name", "CPU Utilization")
	unit := c.FormValue("unit", "%")

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer f.Close()

	summary, err := h.svc.AnalyzeSolarwindsCSV(metricName, unit, f)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

func (h *SizingHandler) ListTemplates(c *fiber.Ctx) error {
	var templates []model.TemplateFile

	err := filepath.Walk(h.templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".xlsx" || ext == ".docx" || ext == ".md" || ext == ".csv" {
			rel, _ := filepath.Rel(h.templatesDir, path)
			category := "GENERAL"
			if len(filepath.Dir(rel)) > 0 && filepath.Dir(rel) != "." {
				category = filepath.Dir(rel)
			}
			templates = append(templates, model.TemplateFile{
				ID:          rel,
				Name:        info.Name(),
				Category:    category,
				FileName:    info.Name(),
				SizeBytes:   info.Size(),
				Format:      ext,
				Description: "Official deliverable template: " + info.Name(),
			})
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to scan templates directory",
		})
	}

	return c.JSON(fiber.Map{
		"total": len(templates),
		"data":  templates,
	})
}

func (h *SizingHandler) DownloadTemplate(c *fiber.Ctx) error {
	filename := c.Params("filename")
	targetPath := filepath.Join(h.templatesDir, filename)

	// Also check subfolders
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		var foundPath string
		_ = filepath.Walk(h.templatesDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info != nil && !info.IsDir() && info.Name() == filename {
				foundPath = path
			}
			return nil
		})
		if foundPath != "" {
			targetPath = foundPath
		}
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Template file not found",
		})
	}

	return c.Download(targetPath)
}

// --- IN-MEMORY CACHE HELPERS (Sub-millisecond read access) ---

func (h *SizingHandler) invalidateCache() {
	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()
	h.cachedComponents = nil
	h.cachedPricingItems = nil
	h.cachedRules = nil
	h.cacheExpiresAt = time.Time{}
}

func (h *SizingHandler) getCachedPricingItems() ([]model.PricingItem, error) {
	h.cacheMu.RLock()
	if len(h.cachedPricingItems) > 0 && time.Now().Before(h.cacheExpiresAt) {
		items := make([]model.PricingItem, len(h.cachedPricingItems))
		copy(items, h.cachedPricingItems)
		h.cacheMu.RUnlock()
		return items, nil
	}
	h.cacheMu.RUnlock()

	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()

	// Double check after write lock
	if len(h.cachedPricingItems) > 0 && time.Now().Before(h.cacheExpiresAt) {
		items := make([]model.PricingItem, len(h.cachedPricingItems))
		copy(items, h.cachedPricingItems)
		return items, nil
	}

	var items []model.PricingItem
	if err := h.db.Order("category asc, sku asc").Find(&items).Error; err != nil {
		return nil, err
	}

	h.cachedPricingItems = items
	h.cacheExpiresAt = time.Now().Add(5 * time.Minute)

	res := make([]model.PricingItem, len(items))
	copy(res, items)
	return res, nil
}

func (h *SizingHandler) getCachedComponents() ([]model.CatalogComponent, error) {
	h.cacheMu.RLock()
	if len(h.cachedComponents) > 0 && time.Now().Before(h.cacheExpiresAt) {
		comps := h.cachedComponents
		h.cacheMu.RUnlock()
		return comps, nil
	}
	h.cacheMu.RUnlock()

	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()

	if len(h.cachedComponents) > 0 && time.Now().Before(h.cacheExpiresAt) {
		return h.cachedComponents, nil
	}

	var components []model.CatalogComponent
	if err := h.db.Preload("Baselines").Preload("Licenses").Order("layer_no asc").Find(&components).Error; err != nil {
		return nil, err
	}

	h.cachedComponents = components
	h.cacheExpiresAt = time.Now().Add(5 * time.Minute)
	return components, nil
}

func (h *SizingHandler) getCachedLicenseRules() ([]model.LicenseRule, error) {
	h.cacheMu.RLock()
	if len(h.cachedRules) > 0 && time.Now().Before(h.cacheExpiresAt) {
		rules := make([]model.LicenseRule, len(h.cachedRules))
		copy(rules, h.cachedRules)
		h.cacheMu.RUnlock()
		return rules, nil
	}
	h.cacheMu.RUnlock()

	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()

	if len(h.cachedRules) > 0 && time.Now().Before(h.cacheExpiresAt) {
		rules := make([]model.LicenseRule, len(h.cachedRules))
		copy(rules, h.cachedRules)
		return rules, nil
	}

	var rules []model.LicenseRule
	if err := h.db.Order("category asc, rule_name asc").Find(&rules).Error; err != nil {
		return nil, err
	}

	h.cachedRules = rules
	h.cacheExpiresAt = time.Now().Add(5 * time.Minute)

	res := make([]model.LicenseRule, len(rules))
	copy(res, rules)
	return res, nil
}

func (h *SizingHandler) GetCatalogComponents(c *fiber.Ctx) error {
	components, err := h.getCachedComponents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch components: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total": len(components),
		"data":  components,
	})
}

func (h *SizingHandler) GetPricingItems(c *fiber.Ctx) error {
	items, err := h.getCachedPricingItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pricing items: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total": len(items),
		"data":  items,
	})
}

type CreatePricingItemRequest struct {
	SKU          string `json:"sku"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	Environment  string `json:"environment"`
	Unit         string `json:"unit"`
	UnitPriceIDR *int64 `json:"unit_price_idr"`
	IsConfirmed  bool   `json:"is_confirmed"`
	Source       string `json:"source"`
}

func (h *SizingHandler) CreatePricingItem(c *fiber.Ctx) error {
	var req CreatePricingItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.SKU == "" || req.Name == "" || req.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SKU, Nama Item, dan Kategori wajib diisi",
		})
	}

	// Find active price list version
	var version model.PriceListVersion
	if err := h.db.Where("status = ?", "active").First(&version).Error; err != nil {
		if err := h.db.First(&version).Error; err != nil {
			version = model.PriceListVersion{
				ID:            uuid.New(),
				Name:          "Katalog Harga Master Mandiri",
				EffectiveDate: time.Now(),
				Status:        "active",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			h.db.Create(&version)
		}
	}

	if req.Environment == "" {
		req.Environment = "Production"
	}
	if req.Unit == "" {
		req.Unit = "Unit"
	}
	if req.Source == "" {
		req.Source = "Katalog Manual"
	}

	item := model.PricingItem{
		ID:                 uuid.New(),
		PriceListVersionID: version.ID,
		SKU:                req.SKU,
		Category:           req.Category,
		Name:               req.Name,
		Environment:        req.Environment,
		Unit:               req.Unit,
		UnitPriceIDR:       req.UnitPriceIDR,
		IsConfirmed:        req.IsConfirmed,
		Source:             req.Source,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := h.db.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pricing item: " + err.Error(),
		})
	}

	h.invalidateCache()
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *SizingHandler) UpdatePricingItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	itemUUID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID format",
		})
	}

	var item model.PricingItem
	if err := h.db.First(&item, "id = ?", itemUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pricing item not found",
		})
	}

	var req CreatePricingItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.SKU != "" {
		item.SKU = req.SKU
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Environment != "" {
		item.Environment = req.Environment
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.Source != "" {
		item.Source = req.Source
	}

	item.UnitPriceIDR = req.UnitPriceIDR
	item.IsConfirmed = req.IsConfirmed
	item.UpdatedAt = time.Now()

	if err := h.db.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update pricing item: " + err.Error(),
		})
	}

	h.invalidateCache()
	return c.JSON(item)
}

func (h *SizingHandler) DeletePricingItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	itemUUID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID format",
		})
	}

	if err := h.db.Delete(&model.PricingItem{}, "id = ?", itemUUID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete pricing item: " + err.Error(),
		})
	}

	h.invalidateCache()
	return c.JSON(fiber.Map{
		"message": "Pricing item deleted successfully",
		"id":      idParam,
	})
}

// --- LICENSE RULES CRUD ---

type LicenseRuleRequest struct {
	RuleName       string `json:"rule_name"`
	Vendor         string `json:"vendor"`
	Category       string `json:"category"`
	MetricType     string `json:"metric_type"`
	PackSize       int    `json:"pack_size"`
	MinClusterNode int    `json:"min_cluster_node"`
	TargetLayer    string `json:"target_layer"`
	PricingSKU     string `json:"pricing_sku"`
	UnitPriceIDR   *int64 `json:"unit_price_idr"`
	Description    string `json:"description"`
	IsActive       *bool  `json:"is_active"`
}

func (h *SizingHandler) GetLicenseRules(c *fiber.Ctx) error {
	rules, err := h.getCachedLicenseRules()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch license rules: " + err.Error(),
		})
	}
	return c.JSON(rules)
}

func (h *SizingHandler) CreateLicenseRule(c *fiber.Ctx) error {
	var req LicenseRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.RuleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rule_name is required",
		})
	}

	packSize := req.PackSize
	if packSize <= 0 {
		packSize = 1
	}

	minNodes := req.MinClusterNode
	if minNodes <= 0 {
		minNodes = 1
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rule := model.LicenseRule{
		ID:             uuid.New(),
		RuleName:       req.RuleName,
		Vendor:         req.Vendor,
		Category:       req.Category,
		MetricType:     req.MetricType,
		PackSize:       packSize,
		MinClusterNode: minNodes,
		TargetLayer:    req.TargetLayer,
		PricingSKU:     req.PricingSKU,
		UnitPriceIDR:   req.UnitPriceIDR,
		Description:    req.Description,
		IsActive:       isActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(&rule).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create license rule: " + err.Error(),
		})
	}

	h.invalidateCache()
	return c.Status(fiber.StatusCreated).JSON(rule)
}

func (h *SizingHandler) UpdateLicenseRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	ruleUUID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID format",
		})
	}

	var rule model.LicenseRule
	if err := h.db.First(&rule, "id = ?", ruleUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "License rule not found",
		})
	}

	var req LicenseRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.RuleName != "" {
		rule.RuleName = req.RuleName
	}
	if req.Vendor != "" {
		rule.Vendor = req.Vendor
	}
	if req.Category != "" {
		rule.Category = req.Category
	}
	if req.MetricType != "" {
		rule.MetricType = req.MetricType
	}
	if req.PackSize > 0 {
		rule.PackSize = req.PackSize
	}
	if req.MinClusterNode > 0 {
		rule.MinClusterNode = req.MinClusterNode
	}
	if req.TargetLayer != "" {
		rule.TargetLayer = req.TargetLayer
	}
	if req.PricingSKU != "" {
		rule.PricingSKU = req.PricingSKU
	}
	if req.UnitPriceIDR != nil {
		rule.UnitPriceIDR = req.UnitPriceIDR
	}
	if req.Description != "" {
		rule.Description = req.Description
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	rule.UpdatedAt = time.Now()

	if err := h.db.Save(&rule).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update license rule: " + err.Error(),
		})
	}

	h.invalidateCache()
	return c.JSON(rule)
}

func (h *SizingHandler) DeleteLicenseRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	ruleUUID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID format",
		})
	}

	if err := h.db.Delete(&model.LicenseRule{}, "id = ?", ruleUUID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete license rule: " + err.Error(),
		})
	}

	h.invalidateCache()
	return c.JSON(fiber.Map{
		"message": "License rule deleted successfully",
		"id":      idParam,
	})
}

func (h *SizingHandler) CalculateQuote(c *fiber.Ctx) error {
	var req model.QuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid quote request body: " + err.Error(),
		})
	}

	pricing, err := h.getCachedPricingItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pricing: " + err.Error(),
		})
	}

	components, err := h.getCachedComponents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch catalog: " + err.Error(),
		})
	}

	engine := quote.NewEngine(pricing, components)
	quoteResp := engine.CalculateQuote(req)
	return c.JSON(quoteResp)
}

type ExportQuoteRequest struct {
	Quote       model.QuoteRequest `json:"quote"`
	ProjectName string             `json:"project_name"`
}

func (h *SizingHandler) ExportQuoteExcel(c *fiber.Ctx) error {
	var req ExportQuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid export request: " + err.Error(),
		})
	}

	if req.ProjectName == "" {
		req.ProjectName = "Sizing Arsitektur Bank Mandiri"
	}

	pricing, err := h.getCachedPricingItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pricing: " + err.Error(),
		})
	}

	components, err := h.getCachedComponents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch components: " + err.Error(),
		})
	}

	engine := quote.NewEngine(pricing, components)
	quoteResp := engine.CalculateQuote(req.Quote)

	buf, err := excel.GenerateQuoteWorkbook(quoteResp, req.Quote, req.ProjectName, components, pricing)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to build Excel spreadsheet: " + err.Error(),
		})
	}

	fileName := fmt.Sprintf("Mandiri_Sizing_BOQ_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return c.Send(buf.Bytes())
}

