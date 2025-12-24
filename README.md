# 📚 Analisis Kompleksitas: Binary Search Iteratif vs Rekursif

**Judul Proyek:** tubes-AKA-kompleksitas-binary-search-iteratif-rekursif-untuk-perpustakaan

Develob by : 
- Salman Baihaqi 
- Syafiq Yusuf Ikhsan W.S
Proyek ini adalah implementasi dan analisis perbandingan kinerja algoritma **Binary Search** menggunakan pendekatan **Iteratif** dan **Rekursif**. Dibuat sebagai pemenuhan Tugas Besar mata kuliah Analisis Kompleksitas Algoritma (AKA), simulasi ini mengambil studi kasus pencarian buku dalam sistem perpustakaan digital berskala besar.

## 🌟 Fitur Utama

- **Pembangkit Data Unik (Lehmer Code Logic):** Mampu menghasilkan hingga **4.020.960** judul buku unik yang terurut secara otomatis tanpa duplikasi, menggunakan kombinasi permutasi dari 160 kata baku.
- **Dual Algorithm Implementation:** Membandingkan Binary Search Iteratif dan Rekursif secara _apple-to-apple_.
- **High-Precision Benchmarking:** Menggunakan strategi _warm-up_ dan _multiple batch runs_ untuk mendapatkan rata-rata waktu eksekusi dalam nanosecond (ns) yang stabil.
- **Web-based Visualization:** Antarmuka pengguna (UI) berbasis web yang interaktif untuk memvisualisasikan perbandingan langkah dan waktu.
- **Analisis Kompleksitas:** Mendemonstrasikan efisiensi \$O(\\log n)\$ pada dataset besar (hingga jutaan data).

## 🚀 Prasyarat

Sebelum menjalankan proyek ini, pastikan komputer Anda telah terinstal:

- **Go (Golang)**: Versi 1.18 atau yang lebih baru. [Download di sini](https://go.dev/dl/).
- **Web Browser**: Chrome, Firefox, Safari, atau Edge untuk melihat hasil visualisasi.

## 📥 Cara Download & Instalasi

Anda bisa mengunduh proyek ini menggunakan Git atau download manual sebagai ZIP.

### Opsi 1: Via Git (Direkomendasikan)

git clone \[<https://github.com/username-anda/tubes-AKA-kompleksitas-binary-search-iteratif-rekursif-untuk-perpustakaan.git\>](<https://github.com/username-anda/tubes-AKA-kompleksitas-binary-search-iteratif-rekursif-untuk-perpustakaan.git>)  
cd tubes-AKA-kompleksitas-binary-search-iteratif-rekursif-untuk-perpustakaan  

### Opsi 2: Download ZIP

- Klik tombol **Code** berwarna hijau di halaman atas repositori ini.
- Pilih **Download ZIP**.
- Ekstrak file ZIP ke folder tujuan Anda.

## 🛠️ Cara Menjalankan Aplikasi

Setelah berhasil diunduh, ikuti langkah berikut:

- Buka Terminal  
    Arahkan terminal ke dalam folder proyek yang sudah diekstrak/clone.  

- **Jalankan Perintah Go**  
    go run tubesAKA.go  

- Akses Web GUI  
    Tunggu hingga muncul pesan Server Binary Search (Unik) Berjalan, lalu buka browser dan akses:  
    <http://localhost:8080>  

## 📊 Metodologi Pengujian

### 1\. Data Generation

Data judul buku tidak di-hardcode, melainkan dibangkitkan secara prosedural menggunakan algoritma permutasi parsial.

- **Kamus:** 160 kata (Administrasi, Advokasi, ..., Zoologi).
- **Format:** 3 kata per judul (misal: "Data Bisnis Akuntansi").
- **Sifat:** Data otomatis terurut (Sorted), prasyarat mutlak untuk Binary Search.

### 2\. Skenario Benchmark

- Pengguna memasukkan jumlah data (N), minimal 500.
- Sistem memilih target pencarian secara acak maupun spesifik.
- Algoritma dijalankan ribuan kali (Warm-up + Real run) untuk menghilangkan _noise_ sistem operasi.

### 3\. Perbandingan

| **Aspek** | **Iteratif** | **Rekursif** |
| --- | --- | --- |
| **Pendekatan** | Loop (for) | Pemanggilan Fungsi Diri Sendiri |
| --- | --- | --- |
| **Memori (Space)** | \$O(1)\$ | \$O(\\log n)\$ (Stack Overhead) |
| --- | --- | --- |
| **Waktu (Time)** | \$O(\\log n)\$ | \$O(\\log n)\$ |
| --- | --- | --- |
| **Overhead** | Rendah | Sedang (Context Switching) |
| --- | --- | --- |

## 📝 Struktur Kode

- tubesAKA.go: File utama yang berisi semua logika (Backend, Algoritma, HTML Template).
  - GenerateJudul(): Logika matematika untuk permutasi kata.
  - BinarySearchIteratif(): Implementasi iteratif.
  - BinarySearchRekursif(): Implementasi rekursif.
  - runStableBenchmark(): Logika pengujian performa.
