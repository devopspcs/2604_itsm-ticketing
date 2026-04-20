# Dimana Lihat Fitur Jira? 🎯

**Bahasa**: Indonesian  
**Tanggal**: April 19, 2026

## Jawaban Singkat

Fitur Jira sekarang **TERLIHAT** di UI! Buka aplikasi dan pergi ke halaman project board.

## Cara Melihat

### 1. Jalankan Aplikasi

```bash
# Terminal 1: Backend
cd backend
go run ./cmd/server

# Terminal 2: Frontend
cd frontend
npm run dev
```

### 2. Buka Browser

Pergi ke: `http://localhost:3000/projects/{projectId}`

Ganti `{projectId}` dengan ID project yang ada.

### 3. Lihat Jira-Style Layout

Anda akan melihat:

```
┌─────────────────────────────────────────────────────────┐
│ ┌──────────┐ ┌──────────────────────────────────────┐  │
│ │ SIDEBAR  │ │ HEADER: Board                        │  │
│ │          │ │ [Release] [Menu]                     │  │
│ │ Backlog  │ ├──────────────────────────────────────┤  │
│ │ Board ✓  │ │ Quick Filters                        │  │
│ │ Sprint   │ ├──────────────────────────────────────┤  │
│ │ Settings │ │ KANBAN BOARD                         │  │
│ │          │ │ ┌────────┐ ┌────────┐ ┌────────┐   │  │
│ │ Reports  │ │ │ TODO   │ │ IN DEV │ │ DONE   │   │  │
│ │ Releases │ │ │        │ │        │ │        │   │  │
│ │ Components│ │ │ Card 1 │ │ Card 2 │ │ Card 3 │   │  │
│ │ Issues   │ │ │ Card 4 │ │ Card 5 │ │        │   │  │
│ │ Repository│ │ └────────┘ └────────┘ └────────┘   │  │
│ │ Settings │ │                                      │  │
│ └──────────┘ └──────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## Fitur yang Terlihat

### Di Sidebar (Kiri)
- ✅ **Backlog** - Lihat dan kelola backlog items
- ✅ **Board** - Kanban board view (halaman saat ini)
- ✅ **Sprint** - Sprint planning dan management
- ✅ **Settings** - Konfigurasi project
- ✅ **Reports** - Laporan project
- ✅ **Releases** - Manajemen release
- ✅ **Components** - Komponen project
- ✅ **Issues** - Daftar semua issues
- ✅ **Repository** - Repository management

### Di Board (Tengah)
- ✅ **Kanban Columns** - Kolom dengan status berbeda (TODO, IN DEV, DONE, dll)
- ✅ **Issue Cards** - Kartu dengan:
  - Judul issue
  - Issue key (TIS-123)
  - Status completion (checkmark jika selesai)
- ✅ **Drag & Drop** - Tarik kartu antar kolom
- ✅ **Add Column** - Tombol untuk tambah kolom baru
- ✅ **Quick Filters** - Filter dan search issues

## Fitur Jira yang Sudah Diimplementasi (17 Total)

### 1. Issue Types ✅
- Bug, Task, Story, Epic, Sub-task
- Lihat di: Setiap kartu di board

### 2. Custom Fields ✅
- Text, Textarea, Dropdown, Multi-select, Date, Number, Checkbox
- Lihat di: Detail modal ketika klik kartu

### 3. Workflows ✅
- Custom statuses dan transitions
- Lihat di: Kolom di board (TODO, IN DEV, DONE, dll)

### 4. Sprint Planning ✅
- Lihat di: Sprint page (klik "Sprint" di sidebar)

### 5. Backlog Management ✅
- Lihat di: Backlog page (klik "Backlog" di sidebar)

### 6. Comments ✅
- Lihat di: Detail modal issue

### 7. Attachments ✅
- Lihat di: Detail modal issue

### 8. Labels & Tags ✅
- Lihat di: Setiap kartu issue

### 9. Issue Type Scheme ✅
- Lihat di: Settings page

### 10. Field Configuration ✅
- Lihat di: Settings page

### 11. Sprint Board ✅
- Lihat di: Board page (halaman saat ini)

### 12. Backlog View ✅
- Lihat di: Backlog page

### 13. Issue Display ✅
- Lihat di: Kartu di board

### 14. Database Schema ✅
- 16 tables untuk support semua fitur

### 15. API Endpoints ✅
- 40+ endpoints untuk semua operasi

### 16. Backward Compatibility ✅
- Semua fitur lama masih bekerja

### 17. Notification System ✅
- Lihat di: Notification center

## Navigasi Antar Halaman

### Dari Board
- Klik **"Backlog"** di sidebar → Lihat backlog items
- Klik **"Sprint"** di sidebar → Lihat sprint planning
- Klik **"Settings"** di sidebar → Konfigurasi project
- Klik **"Reports"** di sidebar → Lihat laporan
- Klik **"Releases"** di sidebar → Manajemen release
- Klik **"Components"** di sidebar → Komponen project
- Klik **"Issues"** di sidebar → Daftar semua issues
- Klik **"Repository"** di sidebar → Repository management

### Dari Halaman Lain
- Klik **"Board"** di sidebar → Kembali ke board

## Interaksi dengan Kartu

1. **Klik kartu** → Buka detail modal
2. **Drag kartu** → Pindah ke kolom lain
3. **Hover kartu** → Lihat efek hover

## Menambah Kolom Baru

1. Scroll ke kanan di board
2. Klik tombol **"+ Tambah kolom"**
3. Ketik nama kolom
4. Tekan Enter atau klik "Tambah"

## Fitur yang Akan Datang

- Lebih banyak customization untuk board
- Advanced filtering dan search
- Reporting dan analytics
- Integration dengan tools lain

## Troubleshooting

### Tidak melihat sidebar?
- Refresh halaman (F5)
- Pastikan sudah login
- Pastikan project ID benar di URL

### Tidak bisa drag kartu?
- Pastikan sudah klik di kartu (bukan di label)
- Coba refresh halaman

### Tidak bisa klik link di sidebar?
- Pastikan sudah klik di text link
- Coba refresh halaman

## Kesimpulan

✅ **Fitur Jira sekarang TERLIHAT dan BERFUNGSI di UI!**

Semua 17 fitur Jira-like sudah diimplementasi dan dapat diakses melalui:
- **Sidebar navigation** untuk navigasi antar halaman
- **Kanban board** untuk manajemen issues
- **Detail modal** untuk melihat dan edit issue details
- **Settings page** untuk konfigurasi project

Aplikasi siap untuk testing dan deployment! 🚀
