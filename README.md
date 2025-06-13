# LocalDNS dengan CoreDNS dan Pi-hole (mode internal)

Sistem ini adalah implementasi DNS lokal berbasis **CoreDNS**, di mana semua perangkat tetap menggunakan **Pi-hole** (192.168.2.3) sebagai DNS utama. Pi-hole melakukan **penerusan (forwarding)** domain lokal tertentu seperti `.test` ke **CoreDNS** (192.168.2.5), dan sisanya ke resolver publik melalui **dnscrypt-proxy**.

---

## 📡 Topologi Jaringan

```plaintext
+-------------+        +------------------+        +------------------+        +-------------------+
|    Klien    |  DNS→  |     Pi-hole      |  →DNS  |     CoreDNS      |  →DNS  |   dnscrypt-proxy  |
| 192.168.X.X |        | 192.168.2.3:53   |        | 192.168.2.5:53   |        |   127.0.0.1:54    |
+-------------+        +------------------+        +------------------+        +-------------------+
                             │                          ▲
                             └──── forward domain .test─┘
```

---

## 🖥️ Komponen dan IP

| Komponen         | IP              | Fungsi                              |
|------------------|-----------------|-------------------------------------|
| Pi-hole          | 192.168.2.3     | DNS utama seluruh klien             |
| CoreDNS          | 192.168.2.5     | Resolver domain `.test`, rDNS lokal |
| dnscrypt-proxy   | 127.0.0.1:54    | Resolver DNS publik terenkripsi     |

---

## 📁 Struktur Proyek

```
.
├── docker-compose.yml
├── Corefile
└── zones/
    ├── domain.test.zone
    ├── contoh.test.zone
    ├── reverse.192.168.X.zone
    └── reverse.10.0.0.zone
```

---

## ⚙️ Langkah Konfigurasi

### 1. Jalankan CoreDNS
```bash
docker-compose up -d
```

### 2. Konfigurasi Forward `.test` di Pi-hole
Tambahkan ke `/etc/dnsmasq.d/99-local.conf`:
```ini
server=/test/192.168.2.5
server=/domain.test/192.168.2.5
server=/contoh.test/192.168.2.5
```
Lalu restart:
```bash
pihole restartdns
```

> ❗ Tidak perlu mengubah konfigurasi DNS pada perangkat klien. Cukup arahkan semuanya ke Pi-hole.

---

## 🧪 Pengujian

### Tes domain lokal dari klien:
```bash
dig @192.168.2.3 api.domain.test
```

### Tes PTR/rDNS lokal:
```bash
dig -x 192.168.2.20 @192.168.2.3
```

### Tes DNS publik:
```bash
dig google.com @192.168.2.3
```

---

## 📌 Catatan Penting
- Pi-hole tetap menjadi entry point DNS utama
- CoreDNS hanya aktif untuk domain lokal (zone file)
- dnscrypt tetap bekerja sebagai resolver utama untuk domain global melalui `127.0.0.1:54`
- Tidak ada konfigurasi DNS yang diubah di klien

---

## 🧰 Lisensi
MIT License
