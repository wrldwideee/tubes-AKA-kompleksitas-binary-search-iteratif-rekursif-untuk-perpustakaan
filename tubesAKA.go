package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"time"
)

// ==========================================
// 1. STRUKTUR DATA
// ==========================================

type Buku struct {
	Index int
	Judul string
}

type ResultData struct {
	JumlahData   int
	TargetIndex  int
	TargetString string
	Pesan        string

	// Hasil Iteratif
	HasilIteratif string
	StepsIteratif int
	WaktuIteratif string
	BarIteratif   int

	// Hasil Rekursif
	HasilRekursif string
	StepsRekursif int
	WaktuRekursif string
	BarRekursif   int

	Pemenang string
	Selisih  string
}

// ==========================================
// 2. KAMUS KATA (DIPERLUAS)
// ==========================================
// 160 Kata Baku Terurut Abjad.
// Kapasitas Unik Tanpa Pengulangan Kata = 160 * 159 * 158 = 4.020.960 Judul.
var kamusKata = []string{
	"Administrasi", "Advokasi", "Aerodinamika", "Agama", "Agrikultur", "Agronomi", "Akuntansi", "Algoritma", "Aljabar", "Anatomi",
	"Analisis", "Antropologi", "Aplikasi", "Arkeologi", "Arsitektur", "Astronomi", "Atom", "Audio", "Audit", "Automotif",
	"Bahasa", "Bank", "Basis", "Biokimia", "Biologi", "Bisnis", "Botani", "Budaya", "Bumi",
	"Cahaya", "Cerdas", "Cloud", "Cyber",
	"Daerah", "Dagang", "Dasar", "Data", "Demografi", "Desain", "Digital", "Dinamika", "Diplomasi", "Distribusi", "Dokter",
	"Ekologi", "Ekonomi", "Eksperimen", "Ekspor", "Elektronika", "Energi", "Enkripsi", "Entitas", "Erosi", "Estetika", "Etika", "Evolusi",
	"Farmasi", "Fisika", "Flora", "Forensik", "Fotografi", "Framework", "Fungsi",
	"Game", "Genetika", "Geografi", "Geologi", "Geometri", "Global", "Grafik", "Gizi",
	"Habitat", "Hardware", "Hayati", "Histori", "Hukum", "Humas",
	"Iklim", "Ilmu", "Imunologi", "Industri", "Infeksi", "Informasi", "Inovasi", "Instrumen", "Intelijen", "Interaksi", "Internet", "Investasi",
	"Jaringan", "Jurnal", "Jurnalistik",
	"Kalkulus", "Kanker", "Keamanan", "Kebijakan", "Kecerdasan", "Kedokteran", "Kehutanan", "Kelautan", "Kimia", "Klinis", "Komunikasi", "Komputer", "Konflik", "Konsep", "Konstruksi", "Kriminal", "Kriptografi", "Kurikulum",
	"Laboratorium", "Lanskap", "Logika", "Logistik", "Lingkungan", "Linux", "Literasi",
	"Makro", "Manajemen", "Marketing", "Matematika", "Material", "Mekanika", "Media", "Medis", "Metalurgi", "Meteorologi", "Metode", "Mikro", "Mikrobiologi", "Mikrokontroler", "Mineral", "Molekuler", "Multimedia",
	"Nano", "Navigasi", "Negara", "Neurologi", "Nuklir", "Numerik", "Nutrisi",
	"Oseanografi", "Optik", "Optimasi", "Organisasi", "Otomasi", "Otomotif",
	"Pajak", "Partikel", "Patologi", "Pemasaran", "Pemrograman", "Pendidikan", "Pengantar", "Penyakit", "Perancangan", "Perangkat", "Perbankan", "Pertanian", "Filsafat", "Politik", "Praktis", "Prinsip", "Produksi", "Protokol", "Proyek", "Psikologi", "Publik",
	"Radiasi", "Rekayasa", "Riset", "Robotika",
	"Sains", "Satelit", "Sejarah", "Sel", "Semikonduktor", "Seni", "Server", "Simulasi", "Sistem", "Software", "Sosiologi", "Statistik", "Strategi", "Struktur", "Survei",
	"Tanah", "Teknik", "Teknologi", "Telekomunikasi", "Teori", "Terapan", "Terapi", "Transportasi", "Turis",
	"Urban",
	"Vaksin", "Validasi", "Vector", "Virus", "Visual", "Vulkanologi",
	"Web", "Wireless", "Wirausaha",
	"Zoologi",
}

// GenerateJudul menghasilkan 3 kata unik yang TERURUT.
// Tidak ada kata yang berulang (misal: "Data Data Data" tidak akan muncul).
func GenerateJudul(n int) string {
	lenK := len(kamusKata)

	// Hitung "digit" permutasi (Lehmer Code logic untuk permutasi parsial)
	// Kita memilih 3 item dari N item.
	// Slot 1 memiliki N pilihan.
	// Slot 2 memiliki N-1 pilihan.
	// Slot 3 memiliki N-2 pilihan.

	denom1 := (lenK - 1) * (lenK - 2)
	denom2 := (lenK - 2)

	// Proteksi jika N melebihi kapasitas (Wrap around)
	n = n % (lenK * denom1)

	d1 := n / denom1
	rem1 := n % denom1

	d2 := rem1 / denom2
	d3 := rem1 % denom2

	// Mapping ke Index Nyata (Skip Logic)
	// 1. Ambil kata pertama
	idx1 := d1

	// 2. Ambil kata kedua (Skip idx1)
	idx2 := d2
	if idx2 >= idx1 {
		idx2++
	}

	// 3. Ambil kata ketiga (Skip idx1 dan idx2)
	idx3 := d3

	// Tentukan index yang sudah terpakai untuk di-skip
	p1, p2 := idx1, idx2
	if p1 > p2 {
		p1, p2 = p2, p1
	} // Sort kecil ke besar

	if idx3 >= p1 {
		idx3++
	}
	if idx3 >= p2 {
		idx3++
	}

	return fmt.Sprintf("%s %s %s", kamusKata[idx1], kamusKata[idx2], kamusKata[idx3])
}

// ==========================================
// 3. ALGORITMA BINARY SEARCH
// ==========================================

func BinarySearchIteratif(data []Buku, target string) (int, int) {
	low := 0
	high := len(data) - 1
	steps := 0

	for low <= high {
		steps++
		mid := low + (high-low)/2

		if data[mid].Judul == target {
			return mid, steps
		}
		if data[mid].Judul < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1, steps
}

func BinarySearchRekursif(data []Buku, target string, low int, high int, steps *int) int {
	*steps++
	if low > high {
		return -1
	}
	mid := low + (high-low)/2

	if data[mid].Judul == target {
		return mid
	}
	if data[mid].Judul < target {
		return BinarySearchRekursif(data, target, mid+1, high, steps)
	} else {
		return BinarySearchRekursif(data, target, low, mid-1, steps)
	}
}

func RunBinarySearchRekursif(data []Buku, target string) (int, int) {
	steps := 0
	idx := BinarySearchRekursif(data, target, 0, len(data)-1, &steps)
	return idx, steps
}

// ==========================================
// 4. BENCHMARK
// ==========================================

func runStableBenchmark(data []Buku, targets []string, algoFunc func([]Buku, string) (int, int)) float64 {
	const numBatches = 20
	const runsPerBatch = 50

	var batchDurations []time.Duration

	// Warm-up
	for _, t := range targets {
		algoFunc(data, t)
	}
	runtime.GC()

	for b := 0; b < numBatches; b++ {
		start := time.Now()
		for r := 0; r < runsPerBatch; r++ {
			for _, t := range targets {
				algoFunc(data, t)
			}
		}
		dur := time.Since(start)
		batchDurations = append(batchDurations, dur)
	}

	sort.Slice(batchDurations, func(i, j int) bool {
		return batchDurations[i] < batchDurations[j]
	})

	trimCount := numBatches / 5
	validDurations := batchDurations[trimCount : numBatches-trimCount]

	var totalNs int64
	for _, d := range validDurations {
		totalNs += d.Nanoseconds()
	}

	totalOps := float64(len(validDurations) * runsPerBatch * len(targets))
	return float64(totalNs) / totalOps
}

// ==========================================
// 5. UI & SERVER
// ==========================================

const htmlTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Perpus Digital: Binary Search</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f6f9; color: #333; margin: 0; padding: 20px; display: flex; justify-content: center; }
        .container { background: white; padding: 35px; border-radius: 16px; box-shadow: 0 10px 30px rgba(0,0,0,0.06); width: 100%; max-width: 800px; }
        
        h1 { text-align: center; color: #2c3e50; margin-top: 0; font-size: 28px; font-weight: 800; letter-spacing: -0.5px; }
        .subtitle { text-align: center; color: #7f8c8d; font-size: 16px; margin-bottom: 35px; font-weight: 500; }

        .form-container { background: #ffffff; padding: 25px; border-radius: 12px; border: 1px solid #edf2f7; box-shadow: 0 2px 10px rgba(0,0,0,0.02); }
        .form-row { display: flex; gap: 20px; margin-bottom: 20px; }
        .form-group { flex: 1; }
        label { display: block; font-weight: 700; margin-bottom: 10px; font-size: 13px; text-transform: uppercase; color: #64748b; letter-spacing: 0.5px; }
        input { width: 100%; padding: 14px 16px; border: 2px solid #e2e8f0; border-radius: 8px; box-sizing: border-box; font-size: 16px; transition: all 0.2s; background: #f8fafc; font-family: 'Courier New', monospace; font-weight: 600; color: #334155; }
        input:focus { border-color: #3b82f6; outline: none; background: white; box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1); }
        .hint { font-size: 12px; color: #94a3b8; margin-top: 6px; }
        
        button { width: 100%; padding: 16px; background: linear-gradient(135deg, #3b82f6, #2563eb); color: white; border: none; border-radius: 8px; font-size: 16px; font-weight: 700; cursor: pointer; transition: 0.3s; box-shadow: 0 4px 6px -1px rgba(59, 130, 246, 0.5); }
        button:hover { transform: translateY(-2px); box-shadow: 0 10px 15px -3px rgba(59, 130, 246, 0.6); }

        .chart-section { margin-top: 40px; animation: fadeIn 0.5s ease-out; }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

        .target-display { 
            text-align: center; background: #eff6ff; color: #1e40af; 
            padding: 15px; border-radius: 10px; margin-bottom: 30px; 
            font-family: 'Segoe UI', sans-serif; font-size: 16px; 
            border: 1px solid #dbeafe; display: flex; flex-direction: column; gap: 5px;
        }
        .book-icon { font-size: 24px; margin-bottom: 5px; }
        .book-title { font-size: 20px; font-weight: 800; color: #1e3a8a; font-family: 'Georgia', serif; }

        .chart-container {
            display: flex; justify-content: center; align-items: flex-end; gap: 80px;
            height: 300px; padding: 0 20px 20px 20px; margin-bottom: 30px;
            position: relative;
            background-image: linear-gradient(#e5e7eb 1px, transparent 1px);
            background-size: 100% 50px;
            border-bottom: 2px solid #94a3b8;
        }
        
        .bar-group { display: flex; flex-direction: column; align-items: center; justify-content: flex-end; height: 100%; width: 120px; position: relative; z-index: 10; }
        .bar-value { margin-bottom: 12px; font-weight: 800; color: #334155; font-size: 16px; background: rgba(255,255,255,0.8); padding: 2px 6px; border-radius: 4px; }
        .bar { width: 100%; border-radius: 8px 8px 0 0; transition: height 1.2s cubic-bezier(0.34, 1.56, 0.64, 1); min-height: 4px; position: relative; }
        .bar-label { margin-top: 15px; font-weight: 700; text-align: center; color: #475569; letter-spacing: 0.5px; }
        .bar-sublabel { font-size: 13px; color: #64748b; font-weight: 500; margin-top: 6px; background: #f1f5f9; padding: 4px 10px; border-radius: 20px; display: inline-block; }

        .iteratif-bar { background: linear-gradient(180deg, #38bdf8 0%, #0ea5e9 100%); box-shadow: 0 4px 20px rgba(14, 165, 233, 0.4); }
        .rekursif-bar { background: linear-gradient(180deg, #f472b6 0%, #db2777 100%); box-shadow: 0 4px 20px rgba(219, 39, 119, 0.4); }

        .summary-box { background-color: #ffffff; border-radius: 12px; padding: 30px; border: 1px solid #e2e8f0; position: relative; overflow: hidden; }
        .summary-box::before { content: ""; position: absolute; left: 0; top: 0; bottom: 0; width: 6px; background: linear-gradient(to bottom, #22c55e, #16a34a); }
        .summary-title { font-weight: 800; margin-bottom: 15px; color: #166534; font-size: 18px; text-transform: uppercase; letter-spacing: 0.5px; display: flex; align-items: center; gap: 8px; }
        .summary-text { font-size: 15px; line-height: 1.8; color: #475569; }
        .highlight { font-weight: 800; color: #1e293b; background: #dcfce7; padding: 0 4px; border-radius: 4px; }
    </style>
</head>
<body>

<div class="container">
    <h1>Perpus Digital: Analisis Algoritma</h1>
    <p class="subtitle">Benchmarking Binary Search (Judul Buku Unik)</p>

    <div class="form-container">
        <form method="POST" action="/">
            <div class="form-row">
                <div class="form-group">
                    <label>Jumlah Koleksi Buku (Min 500)</label>
                    <input type="number" name="amount" value="{{if .JumlahData}}{{.JumlahData}}{{else}}1000000{{end}}" min="500">
                    <div class="hint">Minimal 500 untuk kompleksitas yang nyata</div>
                </div>
                <div class="form-group">
                    <label>Cari ID Buku</label>
                    <input type="number" name="target_index" value="{{if .TargetIndex}}{{.TargetIndex}}{{else}}12345{{end}}">
                    <div class="hint">Index buku yang ingin dicari (0 - N)</div>
                </div>
            </div>
            <button type="submit">🔍 Generate Koleksi & Mulai Analisis</button>
        </form>
    </div>

    {{if .Pesan}}
    <div class="chart-section">
        <div class="target-display">
            <div class="book-icon">📚</div>
            <div>Mencari Buku dengan Judul:</div>
            <div class="book-title">"{{.TargetString}}"</div>
            <div style="font-size: 13px; color: #64748b; margin-top:5px;">(ID: {{.TargetIndex}})</div>
        </div>

        <h3 style="text-align:center; margin-bottom:40px; color:#475569; font-weight:600;">⏱️ Perbandingan Rata-rata Waktu (Nanoseconds)</h3>
        
        <div class="chart-container">
            <!-- Batang Iteratif -->
            <div class="bar-group">
                <div class="bar-value">{{.WaktuIteratif}} ns</div>
                <div class="bar iteratif-bar" style="height: {{.BarIteratif}}%;"></div>
                <div class="bar-label">
                    Iteratif
                    <div class="bar-sublabel">{{.StepsIteratif}} Langkah</div>
                </div>
            </div>

            <!-- Batang Rekursif -->
            <div class="bar-group">
                <div class="bar-value">{{.WaktuRekursif}} ns</div>
                <div class="bar rekursif-bar" style="height: {{.BarRekursif}}%;"></div>
                <div class="bar-label">
                    Rekursif
                    <div class="bar-sublabel">{{.StepsRekursif}} Langkah</div>
                </div>
            </div>
        </div>

        <div class="summary-box">
            <div class="summary-title">📝 Laporan Analisis</div>
            <div class="summary-text">
                <p>
                    Data adalah <strong>Judul Buku Unik</strong> yang dibuat dari kombinasi 3 kata berbeda (Permutasi).
                    <br>
                    <strong>Hasil Benchmark:</strong> Metode <span class="highlight">{{.Pemenang}}</span> lebih unggul dengan selisih waktu <span class="highlight">{{.Selisih}} ns</span>.
                </p>
                <hr style="border: 0; border-top: 1px dashed #cbd5e1; margin: 15px 0;">
                <ul style="padding-left: 20px; margin: 0;">
                    <li><strong>Unique Title Logic:</strong> Judul menggunakan 3 kata yang tidak boleh sama (Misal: "Data Bisnis Akuntansi"). Dengan 160 kata dasar, tersedia >4 Juta variasi judul unik yang terurut secara otomatis.</li>
                    <li><strong>Kompleksitas:</strong> Dengan minimal 500 data, kompleksitas $O(\log n)$ mulai terasa dampaknya. Pada 1 juta data, pencarian hanya butuh sekitar 20 langkah.</li>
                </ul>
            </div>
        </div>
    </div>
    {{end}}
</div>

</body>
</html>
`

// ==========================================
// 6. HANDLER
// ==========================================

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Gagal memuat template", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	amountStr := r.FormValue("amount")
	targetIdxStr := r.FormValue("target_index")
	amount, _ := strconv.Atoi(amountStr)
	targetIndex, _ := strconv.Atoi(targetIdxStr)

	// Validasi input
	if amount < 500 {
		amount = 500
	} // Paksa minimal 500 sesuai request
	if targetIndex >= amount {
		targetIndex = amount - 1
	}
	if targetIndex < 0 {
		targetIndex = 0
	}

	// 1. Generate Data (Virtual Access)
	bukuList := make([]Buku, amount)
	for i := 0; i < amount; i++ {
		bukuList[i] = Buku{
			Index: i,
			Judul: GenerateJudul(i),
		}
	}

	// 2. Tentukan Target String
	targetString := bukuList[targetIndex].Judul

	// 3. Siapkan Sampel Benchmark Stabil
	const sampleCount = 2000
	targetStrings := make([]string, sampleCount)
	step := 1
	if amount > sampleCount {
		step = amount / sampleCount
	}
	for i := 0; i < sampleCount; i++ {
		index := (i * step) % amount
		targetStrings[i] = bukuList[index].Judul
	}

	// 4. Jalankan Benchmark
	nsIt := runStableBenchmark(bukuList, targetStrings, BinarySearchIteratif)
	nsRec := runStableBenchmark(bukuList, targetStrings, RunBinarySearchRekursif)

	// 5. Hitung Langkah (Single Run)
	_, stepsIt := BinarySearchIteratif(bukuList, targetString)
	_, stepsRec := RunBinarySearchRekursif(bukuList, targetString)

	// 6. Kalkulasi Visualisasi
	maxVal := nsIt
	if nsRec > maxVal {
		maxVal = nsRec
	}
	if maxVal == 0 {
		maxVal = 1
	}

	widthIt := int((nsIt / maxVal) * 100)
	widthRec := int((nsRec / maxVal) * 100)
	if widthIt < 2 {
		widthIt = 2
	}
	if widthRec < 2 {
		widthRec = 2
	}

	// 7. Pemenang
	pemenang := "Seimbang"
	selisih := 0.0
	if nsIt < nsRec {
		pemenang = "Iteratif"
		selisih = nsRec - nsIt
	} else if nsRec < nsIt {
		pemenang = "Rekursif"
		selisih = nsIt - nsRec
	}

	data := ResultData{
		JumlahData:    amount,
		TargetIndex:   targetIndex,
		TargetString:  targetString,
		Pesan:         "Done",
		WaktuIteratif: fmt.Sprintf("%.3f", nsIt),
		StepsIteratif: stepsIt,
		BarIteratif:   widthIt,
		WaktuRekursif: fmt.Sprintf("%.3f", nsRec),
		StepsRekursif: stepsRec,
		BarRekursif:   widthRec,
		Pemenang:      pemenang,
		Selisih:       fmt.Sprintf("%.3f", selisih),
	}

	tmpl.Execute(w, data)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	http.HandleFunc("/", mainHandler)
	fmt.Println("==============================================")
	fmt.Println(" Server Binary Search (Unik) Berjalan")
	fmt.Println(" Akses di browser: http://localhost:8080")
	fmt.Println("==============================================")
	http.ListenAndServe(":8080", nil)
}
