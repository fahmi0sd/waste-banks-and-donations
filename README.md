# Bank Sampah & Donasi

REST API untuk platform bank sampah multi-lokasi yang terintegrasi dengan donasi — user menukar sampah jadi saldo (rupiah) di wallet digital, saldo itu bisa ditarik atau didonasikan ke kampanye penggalangan dana. Dibangun sebagai final project bootcamp Hacktiv8.

## Tech Stack

- **Bahasa & Framework**: Go 1.25.5, [Echo](https://echo.labstack.com/) (REST API)
- **Database**: PostgreSQL (data bisnis + log, satu database untuk semuanya)
- **ORM**: [GORM](https://gorm.io/)
- **Auth**: JWT (`golang-jwt/jwt`)
- **Email**: [Mailjet](https://www.mailjet.com/) API
- **Scheduler**: [robfig/cron](https://github.com/robfig/cron) (backup otomatis harian)
- **Validasi**: [go-playground/validator](https://github.com/go-playground/validator)

## Fitur

| Modul | Deskripsi |
|---|---|
| Auth & Profil | Registrasi, login, kelola profil sendiri |
| Master Data | CRUD lokasi bank sampah, kategori sampah, harga (dengan histori), kelola akun admin |
| Antrian & Kalkulator | Simulasi nilai tukar sampah, request & verifikasi antrian kunjungan |
| Transaksi Penukaran Sampah | Input sampah multi-kategori, kredit wallet otomatis, atomic dengan update stok lokasi |
| Wallet & Withdraw | Cek saldo, riwayat mutasi, tarik saldo |
| Donasi | Kampanye donasi, donasi dari saldo wallet, update penyaluran dana + notifikasi ke donor |
| Notifikasi | Email otomatis saat deposit masuk & update donasi (via Mailjet) |
| Laporan | Agregasi transaksi per lokasi/tanggal (location-scoped untuk admin) |
| Backup | Backup database manual/terjadwal (tiap hari jam 2 pagi), retensi otomatis |

## Database: Supabase

Project ini pakai [Supabase](https://supabase.com) (Postgres terkelola) sebagai database.

### Load Skema ke Supabase 

**Opsi A — lewat Dashboard** : buka Supabase Dashboard → SQL Editor → paste isi `sql/ddl.sql` → Run.

**Opsi B — lewat psql** (kalau sudah install `postgresql-client`):
```bash
psql "postgresql://postgres.xxxxxxxxxxxx:[PASSWORD]@aws-0-[region].pooler.supabase.com:5432/postgres" -f sql/ddl.sql
```
Ambil connection string persis dari Supabase Dashboard → Connect → Session pooler.

## Menjalankan dengan Docker

### Prasyarat
- Docker & Docker Compose terpasang
- Sudah buat project Supabase dan sudah load skema (lihat di atas)

### Langkah

```bash
cp .env.example .env
# edit .env — isi DB_HOST/DB_USER/DB_PASS/DB_NAME dari Supabase Dashboard (Connect -> Session pooler),
# plus MAILJET_* dan JWT_SECRET

docker compose up --build
```

API akan jalan di `http://localhost:8080`.

Cek API sudah hidup:
```bash
curl http://localhost:8080/ping
# {"message":"pong"}
```

## Menjalankan Tanpa Docker

### Prasyarat
- Go 1.25.5+
- Project Supabase yang skemanya sudah di-load (lihat di atas)

### Langkah

```bash
cp .env.example .env
# edit .env sesuai connection string Supabase kalian

go mod tidy
go run ./app/echo-server
```

## Environment Variables

| Variable | Wajib? | Default | Keterangan |
|---|---|---|---|
| `PORT` | Tidak | `8080` | Port API |
| `DB_HOST` | Ya | - | Host pooler Supabase (Dashboard -> Connect -> Session pooler) |
| `DB_PORT` | Ya | - | Port pooler Supabase (biasanya `5432`) |
| `DB_USER` | Ya | - | Format `postgres.xxxxxxxxxxxx`, dari Supabase Dashboard |
| `DB_PASS` | Ya | - | Password database Supabase |
| `DB_NAME` | Ya | - | Biasanya `postgres` untuk Supabase |
| `DB_SSLMODE` | Tidak | `require` | Supabase selalu dukung SSL secara default, jadi `require` langsung jalan tanpa setup tambahan |
| `JWT_SECRET` | Ya | - | Secret key untuk tanda tangan JWT — pakai string acak yang panjang |
| `JWT_EXPIRY_HOURS` | Tidak | `24` | Masa berlaku token login (jam) |
| `MAILJET_API_KEY` | Ya (untuk notifikasi email) | - | Lihat panduan setup Mailjet |
| `MAILJET_SECRET_KEY` | Ya (untuk notifikasi email) | - | - |
| `MAILJET_SENDER_EMAIL` | Ya (untuk notifikasi email) | - | Harus email yang sudah diverifikasi di dashboard Mailjet |
| `MAILJET_SENDER_NAME` | Tidak | `Bank Sampah & Donasi` | Nama pengirim email |
| `MAILJET_BASE_URL` | Tidak | endpoint resmi Mailjet | Cuma untuk testing (arahkan ke mock server) |
| `BACKUP_DIR` | Tidak | `backup` | **Kalau tidak pakai Docker**, sebaiknya arahkan ke path di luar folder project (lihat catatan di bawah) |
| `BACKUP_KEEP` | Tidak | `30` | Jumlah file backup terakhir yang disimpan sebelum yang lebih lama dihapus otomatis |

## Struktur Project

```
app/echo-server/
├── main.go              # entry point, dependency injection
├── router/               # daftar route + struct Controllers
├── middleware/            # JWT middleware
└── controller/{domain}/   # HTTP handler per domain

service/{domain}/          # business logic + interface Repository
repository/{domain}/       # implementasi GORM dari interface Repository
pkg/                       # infrastruktur lintas-domain (Mailjet client, identity lookup)
util/                      # helper kecil (response format, ambil user ID dari context)
scheduler/                 # cron job (backup otomatis)
sql/ddl.sql                # skema database (satu-satunya sumber kebenaran skema)
```


## Testing

### Unit Test
```bash
go test ./...
```

## Autentikasi

Semua endpoint yang butuh login memakai `Authorization: Bearer <token>`. Dapatkan token lewat:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Role yang tersedia: `user`, `admin` (terikat ke satu lokasi), `master_admin` (akses penuh). Registrasi mandiri (`/auth/register`) selalu menghasilkan role `user` — akun `admin`/`master_admin` cuma bisa dibuat oleh `master_admin` yang sudah ada lewat `POST /admin/users`.

**Bootstrap master admin pertama**: karena tidak ada master admin di awal, registrasi manual satu akun lalu naikkan role-nya langsung lewat database:
```sql
UPDATE users SET role='master_admin' WHERE email='email_calon_master_admin@example.com';
```

## Tim

Final project Hacktiv8 — Bank Sampah & Donasi.
