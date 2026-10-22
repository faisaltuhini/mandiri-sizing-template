# Enterprise Sizing & Capacity Planning Hub — Bank Mandiri

Aplikasi terdedikasi untuk **Enterprise Sizing & Capacity Planning**, arsitektur baseline 13-layer, ekstrapolasi non-linear Power-Law ($a \cdot \text{TPS}^b$), analitik telemetri SolarWinds database multi-node, dan deliverable EAWG.

---

## 🏗️ Architecture

Aplikasi ini menggunakan arsitektur **Decoupled**:

```text
00_SIZING_APPS/
├── backend/                  # Golang (Go Fiber v2)
│   ├── cmd/api/main.go       # Entry point HTTP server (:8089)
│   ├── internal/
│   │   ├── handler/          # REST controllers
│   │   ├── service/          # Calculation engine, power-law, EAWG, telemetry
│   │   ├── data/             # 13 architecture layers & enterprise baseline
│   │   └── model/            # Go structs, request/response DTOs
│   ├── Makefile              # make run, make build, make test
│   └── go.mod
│
├── frontend/                 # Nuxt 3 / Vue 3
│   ├── components/template/  # Sizing Hub, Capacity Calculator, EAWG, SolarWinds, Guide
│   ├── composables/          # useSizingApi.ts (integrasi ke Go Backend)
│   ├── layouts/default.vue   # Executive layout & live Go backend health status
│   ├── pages/index.vue       # Dedicated Sizing Hub view
│   ├── server/api/           # Local file openers & file streaming
│   └── package.json
│
├── templates/                # Master Deliverables (.xlsx, .docx, .md, .csv)
│   ├── CAPLAN/               # Capacity planning documents & templates
│   ├── TECH_STACK/           # Tech stack architecture sizing xlsx
│   ├── SOLARWIND/            # Sample raw CPU, memory, storage CSVs
│   └── Template Sizing.xlsx  # Official master enterprise sizing spreadsheet
│
└── Makefile                  # Unified root runner
```

---

## 🚀 Cara Menjalankan

### 1. Menjalankan Sekaligus (Backend & Frontend)
```bash
make dev
```
- **Backend Go**: berjalan di `http://localhost:8089`
- **Frontend Nuxt**: berjalan di `http://localhost:5454`

### 2. Menjalankan Satuan

**Backend (Go Fiber):**
```bash
make dev-backend
# atau:
cd backend && make run
```

**Frontend (Nuxt 3):**
```bash
make dev-frontend
# atau:
cd frontend && pnpm dev
```

### 3. Build Production
```bash
make build
```
- Menghasilkan binary binary Go di `backend/bin/sizing-backend`
- Menghasilkan bundle Nuxt di `frontend/.output`

---

## ⚡ API Endpoints (Backend Go - `:8089`)

| Method | Endpoint | Keterangan |
| :--- | :--- | :--- |
| `GET` | `/api/health` | Healthcheck server & version |
| `GET` | `/api/v1/tech-stack/baseline` | 13 layers enterprise architecture baseline |
| `GET` | `/api/v1/tech-stack/kpis` | Executive summary KPIs (Moderate vs Critical) |
| `POST` | `/api/v1/sizing/calculate-power-law` | Ekstrapolasi non-linear formula $a \cdot \text{TPS}^b$ |
| `POST` | `/api/v1/sizing/calculate-eawg` | Kalkulasi 5-year projection EAWG |
| `POST` | `/api/v1/solarwinds/upload` | Parser & statistical analytics CSV telemetri DB |
| `GET` | `/api/v1/templates` | Daftar berkas master template |
| `GET` | `/api/v1/templates/download/:filename` | Download berkas template resmi |
