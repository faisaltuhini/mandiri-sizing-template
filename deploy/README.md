# 🚀 Panduan Deployment OpenShift (Enterprise Architecture)
## Enterprise Architecture Sizing & Capacity Hub

Panduan resmi untuk melakukan containerization dan deployment aplikasi Sizing Suite (Go Backend + Nuxt 3 Frontend) ke cluster **Red Hat OpenShift Container Platform (OCP) 4.x** menggunakan **Kustomize**.

---

### 📁 Struktur Direktori Deployment

```text
00_SIZING_APPS/
├── backend/
│   ├── Dockerfile             # Multi-stage Golang image (OpenShift Restricted-v2 SCC ready)
│   └── .dockerignore
├── frontend/
│   ├── Dockerfile             # Multi-stage Nuxt 3 image (OpenShift Restricted-v2 SCC ready)
│   └── .dockerignore
└── deploy/
    ├── README.md              # Dokumentasi ini
    └── kustomize/
        ├── base/              # Base Kubernetes & OpenShift resources
        │   ├── backend/
        │   │   ├── deployment.yaml
        │   │   ├── service.yaml
        │   │   ├── configmap.yaml
        │   │   └── secret.example.yaml
        │   ├── frontend/
        │   │   ├── deployment.yaml
        │   │   ├── service.yaml
        │   │   └── configmap.yaml
        │   ├── openshift/
        │   │   ├── route-frontend.yaml   # Route dengan TLS Edge Termination & Redirect HTTPS
        │   │   └── route-backend.yaml    # Route API Backend
        │   └── kustomization.yaml
        └── overlays/          # Konfigurasi per Environment
            ├── dev/           # 1 replica, namespace sizing-dev
            ├── sit/           # 2 replicas HA, namespace sizing-sit
            ├── uat/           # 2 replicas HA, namespace sizing-uat
            └── prod/          # 3 replicas HA + HorizontalPodAutoscaler (HPA), namespace sizing-prod
```

---

### 🛡️ Kepatuhan Keamanan OpenShift (Security Context Constraints)

Kedua image container telah dikonfigurasi secara ketat mematuhi kebijakan **Restricted-v2 SCC** Red Hat OpenShift:
1. **Arbitrary Non-Root UID**: Container berjalan menggunakan `USER 1001` dengan direktori milik `group 0 (root)` dan izin `g=u` (`chmod -R g=u /app`).
2. **Security Context**:
   - `runAsNonRoot: true`
   - `allowPrivilegeEscalation: false`
   - `seccompProfile.type: RuntimeDefault`
   - `capabilities.drop: ["ALL"]`
3. **Health Probes**:
   - Backend: HTTP GET `/api/health` (Port 8089)
   - Frontend: HTTP GET `/` (Port 3000)

---

### 🛠️ Langkah 1: Build & Push Container Image

Anda dapat membangun image menggunakan Docker/Podman lokal atau OpenShift BuildConfig internal:

#### Opsi A: Menggunakan Docker / Podman Lokal

```bash
# Masuk ke direktori 00_SIZING_APPS
cd 00_SIZING_APPS

# 1. Build Backend Image
docker build -t sizing-backend:1.0.0 -f backend/Dockerfile backend/

# 2. Build Frontend Image
docker build -t sizing-frontend:1.0.0 -f frontend/Dockerfile frontend/

# 3. Tag & Push ke Image Registry OpenShift (atau Nexus / Harbor Enterprise)
# Login ke registry OpenShift:
podman login -u $(oc whoami) -p $(oc whoami -t) default-route-openshift-image-registry.apps.enterprise.local

# Tag ke namespace target (contoh: sizing-dev):
docker tag sizing-backend:1.0.0 default-route-openshift-image-registry.apps.enterprise.local/sizing-dev/sizing-backend:dev-latest
docker tag sizing-frontend:1.0.0 default-route-openshift-image-registry.apps.enterprise.local/sizing-dev/sizing-frontend:dev-latest

# Push:
docker push default-route-openshift-image-registry.apps.enterprise.local/sizing-dev/sizing-backend:dev-latest
docker push default-route-openshift-image-registry.apps.enterprise.local/sizing-dev/sizing-frontend:dev-latest
```

#### Opsi B: Menggunakan OpenShift Binary Build (Langsung dari source)

```bash
# 1. Buat BuildConfig di OpenShift
oc new-build --name=sizing-backend --binary --strategy=docker -n sizing-dev
oc new-build --name=sizing-frontend --binary --strategy=docker -n sizing-dev

# 2. Trigger build langsung dari direktori lokal
oc start-build sizing-backend --from-dir=backend -n sizing-dev --follow
oc start-build sizing-frontend --from-dir=frontend -n sizing-dev --follow
```

---

### 🔐 Langkah 2: Setup Secret Database

Aplikasi membutuhkan Secret untuk string koneksi database (`DATABASE_URL`). Buat secret di namespace tujuan sebelum melakukan apply:

```bash
# Untuk Microsoft SQL Server:
oc create secret generic sizing-db-secret \
  --from-literal=DATABASE_URL="sqlserver://sa:YourPassword123@sqlserver.database.svc:1433?database=sizing&trustServerCertificate=true" \
  -n sizing-dev

# Atau untuk PostgreSQL:
oc create secret generic sizing-db-secret \
  --from-literal=DATABASE_URL="postgres://sizing_user:YourPassword123@postgres.database.svc:5432/sizing?sslmode=disable" \
  -n sizing-dev
```

---

### 🚀 Langkah 3: Deploy Menggunakan Kustomize

Pilih overlay sesuai environment target:

#### 1. Deploy ke DEV (`sizing-dev`):
```bash
oc apply -k deploy/kustomize/overlays/dev
```

#### 2. Deploy ke SIT (`sizing-sit`):
```bash
oc apply -k deploy/kustomize/overlays/sit
```

#### 3. Deploy ke UAT (`sizing-uat`):
```bash
oc apply -k deploy/kustomize/overlays/uat
```

#### 4. Deploy ke PRODUCTION (`sizing-prod`):
```bash
oc apply -k deploy/kustomize/overlays/prod
```

> **Tips Preview Dry-Run:**
> Jalankan `kubectl kustomize deploy/kustomize/overlays/dev` untuk melihat output YAML lengkap sebelum di-apply.

---

### 🔍 Langkah 4: Verifikasi & Troubleshooting

```bash
# 1. Cek status Pod
oc get pods -n sizing-dev -l app.kubernetes.io/part-of=sizing-apps

# 2. Cek status OpenShift Route (URL publik)
oc get routes -n sizing-dev

# 3. Cek log backend
oc logs -f deployment/sizing-backend -n sizing-dev

# 4. Cek log frontend
oc logs -f deployment/sizing-frontend -n sizing-dev

# 5. Uji endpoint Health Check dari CLI
curl -k https://$(oc get route sizing-backend -n sizing-dev -o jsonpath='{.spec.host}')/api/health
```
Output yang diharapkan:
```json
{
  "service": "Enterprise Capacity Planning & Sizing Engine (Enterprise Architecture)",
  "status": "ok",
  "version": "1.0.0"
}
```
