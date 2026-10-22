// ─────────────────────────────────────────────
// SIZING GLOSSARY: Penjelasan Istilah untuk Non-Tech Savvy User
// ─────────────────────────────────────────────

export interface GlossaryEntry {
  term: string
  friendlyName: string
  shortDesc: string
  businessImpact: string
}

export const SIZING_GLOSSARY: Record<string, GlossaryEntry> = {
  nodeCount: {
    term: 'Jumlah Node / Server',
    friendlyName: 'Unit Server Aktif',
    shortDesc: 'Berapa banyak server mandiri yang dialokasikan untuk menjalankan komponen ini.',
    businessImpact: 'Semakin banyak server (misal 2N = 2 server aktif), semakin minim risiko sistem mati mendadak saat salah satu server gangguan (High Availability).'
  },
  vcpu: {
    term: 'vCPU (Virtual CPU)',
    friendlyName: 'Kapasitas Pemroses (Otak Komputer)',
    shortDesc: 'Kekuatan pemrosesan komputasi untuk menjalankan logika transaksi per detik.',
    businessImpact: 'Mempengaruhi kecepatan respon aplikasi saat banyak nasabah mengakses secara bersamaan.'
  },
  ram: {
    term: 'RAM (Memory)',
    friendlyName: 'Memori Kerja Cepat',
    shortDesc: 'Ruang simpan sementara berkecepatan ultra-tinggi saat data sedang diproses.',
    businessImpact: 'Mencegah aplikasi lag atau crash saat beban transaksi puncak (peak hour).'
  },
  storage: {
    term: 'SAN Storage (Disk)',
    friendlyName: 'Penyimpanan Data Utama',
    shortDesc: 'Kapasitas hard drive / flash storage jaringan untuk menyimpan data aplikasi & database.',
    businessImpact: 'Menentukan seberapa lama data historis dan volume file nasabah dapat disimpan tanpa kehabisan ruang.'
  },
  storageTier2: {
    term: 'SAN Tier 2 (Standard)',
    friendlyName: 'Penyimpanan Standar Andal',
    shortDesc: 'Penyimpanan SAN standar enterprise, ideal untuk web, container, dan log.',
    businessImpact: 'Biaya optimal (Rp 33.000 / GB / thn) untuk kebutuhan workload umum non-database utama.'
  },
  storageTier3: {
    term: 'SAN Tier 3 (High-Perf ASM)',
    friendlyName: 'Penyimpanan Ultra Cepat (High-IOPS)',
    shortDesc: 'Storage premium latensi mikro-detik khusus klaster Oracle ASM.',
    businessImpact: 'Kecepatan read-write maksimal untuk database Core Banking & transaksi finansial intensif.'
  },
  backupNBU: {
    term: 'NetBackup (NBU)',
    friendlyName: 'Proteksi Cadangan Harian',
    shortDesc: 'Sistem backup otomatis harian enterprise terpusat Bank Mandiri.',
    businessImpact: 'Menjamin data dapat dipulihkan utuh jika terjadi bencana atau kerusakan data (RPO/RTO ketat).'
  },
  computeVM: {
    term: 'Virtual Machine (VM)',
    friendlyName: 'Server Virtual Fleksibel',
    shortDesc: 'Server berbasis virtualisasi yang efisien dan mudah di-scale up/down.',
    businessImpact: 'Pilihan paling fleksibel dan hemat biaya untuk sebagian besar microservices.'
  },
  computeBaremetal: {
    term: 'Baremetal Dedicated',
    friendlyName: 'Server Fisik Khusus (Tanpa Berbagi)',
    shortDesc: 'Mesin server fisik seutuhnya tanpa lapisan virtualisasi tambahan.',
    businessImpact: 'Performa puncak murni tanpa overhead, direkomendasikan untuk Core Database & OpenShift Worker berbeban raksasa.'
  },
  computeShared: {
    term: 'Shared / SAN LUN',
    friendlyName: 'Penyimpanan Berbagi (Storage-only)',
    shortDesc: 'Alokasi ruang penyimpanan bersama yang dikelola terpusat tanpa server compute terpisah.',
    businessImpact: 'Biaya server Rp 0, hanya membayar kapasitas disk yang dialokasikan.'
  },
  licenseFree: {
    term: 'Community OSS',
    friendlyName: 'Software Open Source Gratis',
    shortDesc: 'Software open source tanpa biaya lisensi tahunan vendor.',
    businessImpact: 'Biaya lisensi Rp 0, cocok untuk environment staging, dev, atau service internal dengan support tim internal.'
  },
  licenseEnterprise: {
    term: 'Enterprise Subscription',
    friendlyName: 'Lisensi Resmi Vendor Bergaransi',
    shortDesc: 'Software berlisensi resmi dengan jaminan SLA perbaikan 24/7 dan patch keamanan dari vendor resmi (Red Hat, Oracle, F5).',
    businessImpact: 'Wajib untuk sistem produksi kritis kepatuhan OJK / regulasi perbankan.'
  }
}
