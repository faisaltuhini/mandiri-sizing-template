# Panduan Komprehensif Rumus, Teori & Metodologi Capacity Planning Enterprise

Dokumen ini merupakan panduan teoritis dan praktis yang menjelaskan seluruh rumus matematika, model ekstrapolasi, parameter arsitektur, serta petunjuk operasional untuk **Master Superpower Capacity Plan**. Panduan ini dirancang berdasarkan hasil sintesis mendalam dari arsitektur **RM Tools (RMT)** dan **Livin' Merchant (LM)**.

---

## 1. Konsep & Metodologi Capacity Planning Enterprise

Dalam skala enterprise, *Capacity Planning* bertujuan untuk menjamin bahwa infrastruktur teknologi (vCPU, Memory, Storage, Bandwidth, dan Lisensi) mampu melayani lonjakan transaksi bisnis tanpa mengorbankan SLA (*Service Level Agreement*), *latency*, maupun ketersediaan sistem (*High Availability*).

Kapasitas direncanakan menggunakan kombinasi **Dua Pendekatan Utama**:
1. **Business-Driven Capacity Planning**: Menghitung estimasi beban masa depan berdasarkan proyeksi metrik bisnis (jumlah user aktif, transaksi harian puncak, dan rasio jam sibuk).
2. **Performance Test-Driven (Extrapolation) Capacity Planning**: Menghitung konsumsi sumber daya aktual berbasis hasil uji beban (*load test/stress test*) pada lingkungan *Performance Test* (PT), lalu memproyeksikan konsumsi tersebut secara non-linear menggunakan **Model Power-Law ($a \cdot TPS^b$)**.

---

## 2. Derivasi Matematika & Rumus-Rumus Inti

### 2.1 Model Pertumbuhan Bisnis & Kalkulasi Peak TPS

Untuk menghitung proyeksi *Transactions Per Second* (TPS) puncak pada jam sibuk:

$$TPS_{peak} = \frac{Trx_{daily} \times Factor_{concurrency}}{PeakHours \times 3600}$$

Di mana:
- $Trx_{daily}$: Total estimasi transaksi harian puncak.
- $Factor_{concurrency}$: Rasio lonjakan transaksi pada menit/detik teramai di jam sibuk (tipikal $1.5 \times$ hingga $2.0 \times$).
- $PeakHours$: Durasi window jam sibuk operasional (misal: 8 jam).

Proyeksi transaksi tahunan menggunakan model pertumbuhan eksponensial compound:

$$Trx_{t} = Trx_{base} \times (1 + g)^{t - t_{base}}$$

Di mana $g$ adalah tingkat pertumbuhan tahunan (misal $g = 15\% = 0.15$).

---

### 2.2 Model Ekstrapolasi Non-Linear Power-Law ($a \cdot TPS^b$)

Penggunaan CPU dan Memory oleh aplikasi/microservice tidak selalu bersifat linier $1:1$ terhadap kenaikan TPS karena adanya *overhead* sinkronisasi thread, *garbage collection*, dan manajemen *connection pool*. Oleh karena itu, digunakan **Model Power-Law**:

$$Resource(TPS) = a \cdot TPS^b$$

#### Langkah Derivasi Parameter $b$ dan $a$:
Diberikan 2 titik data mentah hasil uji beban Performance Test (PT):
- Titik PT 1: $TPS_1$, dengan total konsumsi core $C_1 = \frac{\%CPU_1 \times SpecNode \times NodeQty}{100}$
- Titik PT 2: $TPS_2$, dengan total konsumsi core $C_2 = \frac{\%CPU_2 \times SpecNode \times NodeQty}{100}$

1. **Mencari Eksponen Scaling ($b$)**:
   $$\frac{C_2}{C_1} = \frac{a \cdot TPS_2^b}{a \cdot TPS_1^b} = \left(\frac{TPS_2}{TPS_1}\right)^b$$

   Mengambil logaritma natural ($\ln$) di kedua sisi:
   $$\ln\left(\frac{C_2}{C_1}\right) = b \cdot \ln\left(\frac{TPS_2}{TPS_1}\right)$$

   $$b = \frac{\ln(C_2 / C_1)}{\ln(TPS_2 / TPS_1)}$$

2. **Mencari Konstanta Scaling ($a$)**:
   $$a = \frac{C_1}{TPS_1^b}$$

3. **Menghitung Kebutuhan Core untuk Target TPS Masa Depan ($TPS_t$)**:
   $$Cores_{predicted, t} = a \cdot (TPS_t)^b$$

4. **Kebutuhan Core dengan Buffer Utilisasi Aman ($U_{max} = 40\%$ atau $0.40$)**:
   $$Cores_{required, t} = \left\lceil \frac{Cores_{predicted, t}}{U_{max}} \right\rceil$$

5. **Menghitung Jumlah Instance / Node yang Dibutuhkan**:
   $$NodeQty_t = \left\lceil \frac{Cores_{required, t}}{Cores\_per\_Node} \right\rceil$$

---

### 2.3 Microservice Container Pod Sizing & OpenShift Worker Node Math

Untuk aplikasi berarsitektur microservice yang berjalan di atas OpenShift / Kubernetes:

1. **Kalkulasi Pod Replica Count**:
   $$Pods_i = \max\left( \left\lceil \frac{TPS_t}{TPS_{per\_pod}} \right\rceil, MinReplicas \right) \times HA\_Factor$$

2. **Total Cluster CPU & RAM Request**:
   $$CPU_{cluster} = \sum_{i=1}^{N} (Pods_i \times RequestCPU_i) \times (1 + Headroom)$$
   $$RAM_{cluster} = \sum_{i=1}^{N} (Pods_i \times RequestRAM_i) \times (1 + Headroom)$$

3. **Kebutuhan OpenShift Worker Node (Baremetal/VM)**:
   $$Nodes_{worker} = \left\lceil \frac{CPU_{cluster}}{NodeCPU \times MaxUtil_{worker}} \right\rceil + N_{HA}$$
   Di mana $N_{HA} = +1$ atau $+2$ node cadangan untuk toleransi kegagalan node (*High Availability*).

---

### 2.4 Database Capacity & Proyeksi Storage 5 Tahun

Database Relational (Oracle RAC / SQL Server / PostgreSQL) dipisahkan menjadi **Data DB** dan **System DB**.

1. **Memory Sizing Database (Buffer Pool / Shared Buffers)**:
   $$RAM_{db} = Shared\_Buffers + (Work\_Mem \times MaxConnections) + OS\_Cache$$
   Aturan umum: RAM DB diset sebesar $25\% - 50\%$ dari total ukuran dataset aktif yang sering diakses (*working set*).

2. **Proyeksi Pertumbuhan Storage 5 Tahun**:
   - **Data Harian Baru (GB)**:
     $$Data_{daily} = \frac{Trx_{daily} \times AvgRowSize_{KB} \times (1 + IndexOverhead)}{1024 \times 1024}$$
   - **Pertumbuhan Storage Tahunan (GB)**:
     $$Storage_{annual} = Data_{daily} \times 365 \times (1 + LogOverhead)$$
     Di mana $IndexOverhead \approx 30\%$ dan $LogOverhead \approx 20\%$ untuk WAL/REDO Logs dan Temp Space.
   - **Total Disk Disediakan Tahun ke-$t$**:
     $$Storage_{total, t} = Storage_{base} + \sum_{y=1}^{t} Storage_{annual, y} + Storage_{buffer}$$

---

### 2.5 Redis In-Memory Cache Sizing

Kebutuhan RAM untuk Cluster Redis Enterprise / On-Premise:

$$RAM_{redis} = \frac{Keys_{total} \times AvgKeySize_{KB} \times (1 + DictOverhead)}{1024 \times 1024} \times Factor_{repl} \times \frac{1}{Util_{target}}$$

Di mana:
- $Keys_{total}$: Jumlah key sesi / data aktif dalam cache.
- $AvgKeySize_{KB}$: Ukuran rata-rata payload JSON per key (misal 4.5 KB).
- $DictOverhead$: Overhead struktur data internal Redis ($20\% = 0.20$).
- $Factor_{repl}$: Pengali replikasi ($2$ untuk 1 Primary + 1 Replica).
- $Util_{target}$: Max memory limit eviction target ($70\% = 0.70$).

---

### 2.6 Bandwidth & Throughput Jaringan

Estimasi kebutuhan lebar pita (*bandwidth*) Ingress/Egress:

$$Bandwidth_{Mbps} = \frac{TPS_{peak} \times AvgPayload_{KB} \times 8}{1024} \times Headroom_{network}$$

Di mana $Headroom_{network} = 1.5 \times$ (buffer $50\%$ lonjakan traffic HTTP/gRPC).

---

### 2.7 Aturan Aturan Lisensi Hardware & Compliance

1. **Oracle Database Core Factor**:
   Total Core License Oracle dihitung menggunakan Core Factor $0.5$ untuk prosessor Intel/AMD x86:
   $$OracleLicenses = TotalPhysicalCores \times 0.5$$
2. **Microsoft Windows Server Core License Rule**:
   Setiap physical host wajib dilisensikan minimal **16 Core License**, terlepas dari jumlah vCPU VM.

---

## 3. Struktur File Deliverable & Panduan Operasional

### 3.1 Template Excel Interaktif ([`Capacity_Plan_Superpower_Template.xlsx`](file:///Users/mcomacbook/.gemini/antigravity-ide/brain/7fce6081-d43a-4541-b2ae-00fa4c16054b/Capacity_Plan_Superpower_Template.xlsx))

File Excel ini berisi 10 tab terinterkoneksi:
1. **`00_Cover_&_Control`**: Tempat mengatur parameter batas utilisasi aman (CPU 40%, RAM 70%, Storage 80%) dan log versi.
2. **`01_Business_Growth_Plan`**: Mengisi asumsi pertumbuhan user dan transaksi 5 tahun (2025–2029).
3. **`02_SLA_&_Buffer_Parameters`**: Definisi faktor pengali HA dan buffer.
4. **`03_PT_Baseline_Data`**: Tempat menginput data mentah uji beban PT 1 & PT 2.
5. **`04_Extrapolation_Engine`**: Tab kalkulasi otomatis rumus Power-Law ($a \cdot TPS^b$).
6. **`05_Microservices_&_Containers`**: Perhitungan pod microservice & worker node OpenShift.
7. **`06_Database_Capacity`**: Perhitungan spesifikasi & storage DB 5 tahun.
8. **`07_Redis_&_Middleware`**: Perhitungan cluster Redis, Kafka, dan Elasticsearch.
9. **`08_Network_&_Bandwidth`**: Perhitungan throughput jaringan.
10. **`09_Hardware_BOM_Summary`**: Ringkasan akhir Bill of Materials (vCPU, RAM, Storage, Host, Lisensi).

---

### 3.2 Dokumen Word Resmi ([`Capacity_Plan_Superpower_Document.docx`](file:///Users/mcomacbook/.gemini/antigravity-ide/brain/7fce6081-d43a-4541-b2ae-00fa4c16054b/Capacity_Plan_Superpower_Document.docx))

Dokumen Word ini merupakan format standar bank/enterprise yang siap diprint atau di-export ke PDF untuk pengajuan pengadaan hardware dan audit arsitektur:
- **Judul Resmi & Document Change Control**.
- **Bab 1. System Configuration**: Arsitektur 8 tier & inventaris server eksisting.
- **Bab 2. Capacity Planning & Extrapolation**: Tabel pertumbuhan bisnis, rumus matematika, dan tabel proyeksi 5 tahun.
- **Bab 3. Infrastructure Procurement & BOM Summary**: Tabel ringkasan BOM vCPU/RAM/Disk/Lisensi.
- **Bab 4. Deliverable Acceptance Certificate (DAC)**: Lembar persetujuan resmi (Prepared, Reviewed, Accepted).

---

## 4. Kesimpulan

Dengan menggabungkan seluruh komponen dari **RMT** dan **Livin' Merchant** ke dalam 3 deliverable terintegrasi ini, proses perencanaan kapasitas kini bersifat presisi, *auditable*, otomatis, serta memiliki acuan matematis yang kokoh.
