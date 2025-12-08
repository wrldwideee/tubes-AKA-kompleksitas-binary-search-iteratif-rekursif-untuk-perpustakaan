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
	ID    int
	Judul string
}

type ResultData struct {
	JumlahData int
	TargetID   int
	Pesan      string

	// Hasil Iteratif
	HasilIteratif string
	StepsIteratif int
	WaktuIteratif string
	BarIteratif   int // Tinggi grafik (persen)

	// Hasil Rekursif
	HasilRekursif string
	StepsRekursif int
	WaktuRekursif string
	BarRekursif   int // Tinggi grafik (persen)

	// Kesimpulan
	Pemenang string
	Selisih  string
}

// ==========================================
// 2. ALGORITMA BINARY SEARCH
// ==========================================

// Binary Search Iteratif (Loop)
func BinarySearchIteratif(data []Buku, targetID int) (int, int) {
	low := 0
	high := len(data) - 1
	steps := 0

	for low <= high {
		steps++
		mid := low + (high-low)/2
		if data[mid].ID == targetID {
			return mid, steps
		}
		if data[mid].ID < targetID {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1, steps
}

// Binary Search Rekursif
func BinarySearchRekursif(data []Buku, targetID int, low int, high int, steps *int) int {
	*steps++
	if low > high {
		return -1
	}
	mid := low + (high-low)/2
	if data[mid].ID == targetID {
		return mid
	}
	if data[mid].ID < targetID {
		return BinarySearchRekursif(data, targetID, mid+1, high, steps)
	} else {
		return BinarySearchRekursif(data, targetID, low, mid-1, steps)
	}
}

// Wrapper Helper untuk Rekursif agar mudah dipanggil
func RunBinarySearchRekursif(data []Buku, targetID int) (int, int) {
	steps := 0
	idx := BinarySearchRekursif(data, targetID, 0, len(data)-1, &steps)
	return idx, steps
}

// ==========================================
// 3. LOGIKA BENCHMARK STABIL (Trimmed Mean)
// ==========================================
// Fungsi ini menjalankan algoritma berkali-kali dan membuang data outlier
// untuk mendapatkan hasil waktu yang sangat stabil dan akurat.

func runStableBenchmark(data []Buku, targets []int, algoFunc func([]Buku, int) (int, int)) float64 {
	const numBatches = 20
	const runsPerBatch = 50

	var batchDurations []time.Duration

	// 1. Warm-up (Pemanasan CPU)
	for _, tID := range targets {
		algoFunc(data, tID)
	}
	runtime.GC() // Bersihkan memori sebelum mulai agar GC tidak mengganggu

	// 2. Eksekusi Batch
	for b := 0; b < numBatches; b++ {
		start := time.Now()
		for r := 0; r < runsPerBatch; r++ {
			for _, tID := range targets {
				algoFunc(data, tID)
			}
		}
		dur := time.Since(start)
		batchDurations = append(batchDurations, dur)
	}

	// 3. Urutkan & Buang Outlier (Trimmed Mean)
	sort.Slice(batchDurations, func(i, j int) bool {
		return batchDurations[i] < batchDurations[j]
	})

	// Ambil 60% data tengah (buang 20% tercepat & 20% terlambat)
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
// 4. UI & SERVER (HTML CSS TERINTEGRASI)
// ==========================================

const htmlTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Analisis Binary Search</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f8f9fa; color: #333; margin: 0; padding: 20px; display: flex; justify-content: center; }
        .container { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.08); width: 100%; max-width: 700px; }
        
        h1 { text-align: center; color: #2c3e50; margin-top: 0; font-size: 26px; }
        .subtitle { text-align: center; color: #666; font-size: 14px; margin-bottom: 30px; }

        /* Form Styling */
        .form-row { display: flex; gap: 20px; margin-bottom: 20px; }
        .form-group { flex: 1; }
        label { display: block; font-weight: 600; margin-bottom: 8px; font-size: 14px; color: #444; }
        input { width: 100%; padding: 12px; border: 2px solid #e9ecef; border-radius: 8px; box-sizing: border-box; font-size: 16px; transition: 0.3s; }
        input:focus { border-color: #3498db; outline: none; }
        
        button { width: 100%; padding: 14px; background-color: #3498db; color: white; border: none; border-radius: 8px; font-size: 16px; font-weight: bold; cursor: pointer; transition: 0.2s; }
        button:hover { background-color: #2980b9; transform: translateY(-1px); }

        /* Grafik Vertikal Sederhana (CSS Only) */
        .chart-section { margin-top: 40px; padding-top: 20px; border-top: 2px solid #f1f1f1; }
        .chart-container {
            display: flex; justify-content: center; align-items: flex-end; gap: 60px;
            height: 250px; padding: 0 20px; margin-bottom: 20px;
            position: relative;
            background: linear-gradient(to top, #fff, #fcfcfc);
            border-bottom: 2px solid #ccc;
        }
        
        .bar-group { display: flex; flex-direction: column; align-items: center; justify-content: flex-end; height: 100%; width: 100px; }
        .bar-value { margin-bottom: 10px; font-weight: bold; color: #555; font-size: 14px; }
        .bar { width: 100%; border-radius: 6px 6px 0 0; transition: height 1s ease; min-height: 4px; position: relative; }
        .bar-label { margin-top: 15px; font-weight: bold; text-align: center; }
        .bar-sublabel { font-size: 12px; color: #888; font-weight: normal; margin-top: 4px; }

        .iteratif-bar { background-color: #3498db; box-shadow: 0 4px 10px rgba(52, 152, 219, 0.3); }
        .rekursif-bar { background-color: #e74c3c; box-shadow: 0 4px 10px rgba(231, 76, 60, 0.3); }

        /* Kesimpulan Box */
        .summary-box { background-color: #f8f9fa; border-radius: 10px; padding: 20px; margin-top: 30px; border-left: 5px solid #2ecc71; }
        .summary-title { font-weight: bold; margin-bottom: 10px; color: #27ae60; font-size: 16px; }
        .summary-text { font-size: 14px; line-height: 1.6; color: #555; }
    </style>
</head>
<body>

<div class="container">
    <h1>Analisis Pencarian Buku</h1>
    <p class="subtitle">Perbandingan Binary Search Iteratif vs Rekursif (Golang)</p>

    <form method="POST" action="/">
        <div class="form-row">
            <div class="form-group">
                <label>Jumlah Buku (N)</label>
                <input type="number" name="amount" value="{{if .JumlahData}}{{.JumlahData}}{{else}}1000000{{end}}" min="100">
            </div>
            <div class="form-group">
                <label>Cari ID Buku</label>
                <input type="number" name="target_id" value="{{if .TargetID}}{{.TargetID}}{{else}}100{{end}}">
            </div>
        </div>
        <button type="submit">Jalankan Analisis Stabil</button>
    </form>

    {{if .Pesan}}
    <div class="chart-section">
        <h3 style="text-align:center; margin-bottom:30px;">Rata-rata Waktu Eksekusi (Nanoseconds)</h3>
        
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

        <!-- Kotak Analisis & Kesimpulan -->
        <div class="summary-box">
            <div class="summary-title">📊 Analisis & Kesimpulan</div>
            <div class="summary-text">
                <p>
                    <strong>Hasil:</strong> {{.Pemenang}} lebih cepat dengan selisih <strong>{{.Selisih}} ns</strong>.
                </p>
                <ul style="padding-left: 20px; margin-top: 10px;">
                    <li><strong>Kompleksitas Waktu:</strong> Keduanya sama-sama <em>O(log n)</em>.</li>
                    <li><strong>Memori:</strong> Iteratif lebih efisien <em>O(1)</em>, sedangkan Rekursif butuh stack memory <em>O(log n)</em>.</li>
                    <li><strong>Saran:</strong> Gunakan <strong>Iteratif</strong> untuk data sangat besar (Production), dan <strong>Rekursif</strong> untuk pembelajaran karena kodenya lebih mudah dibaca.</li>
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
// 5. HANDLER
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
	targetStr := r.FormValue("target_id")
	amount, _ := strconv.Atoi(amountStr)
	targetID, _ := strconv.Atoi(targetStr)

	// 1. Generate Data
	bukuList := make([]Buku, amount)
	for i := 0; i < amount; i++ {
		bukuList[i] = Buku{ID: i * 2, Judul: "Buku"}
	}

	// 2. Siapkan Target Stabil (2000 titik sampel merata)
	const sampleCount = 2000
	targets := make([]int, sampleCount)
	step := 1
	if amount > sampleCount {
		step = amount / sampleCount
	}
	for i := 0; i < sampleCount; i++ {
		index := (i * step) % amount
		targets[i] = bukuList[index].ID
	}

	// 3. Jalankan Benchmark
	nsIt := runStableBenchmark(bukuList, targets, BinarySearchIteratif)
	nsRec := runStableBenchmark(bukuList, targets, RunBinarySearchRekursif)

	// 4. Hitung Langkah (Hanya untuk ID user untuk keperluan display)
	_, stepsIt := BinarySearchIteratif(bukuList, targetID)
	_, stepsRec := RunBinarySearchRekursif(bukuList, targetID)

	// 5. Persiapan Visualisasi Grafik
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
	} // Minimal tinggi agar terlihat
	if widthRec < 2 {
		widthRec = 2
	}

	// 6. Analisis Pemenang
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
		TargetID:      targetID,
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
	fmt.Println(" Server Binary Search Berjalan (Web Mode)")
	fmt.Println(" Akses di browser: http://localhost:8080")
	fmt.Println("==============================================")
	http.ListenAndServe(":8080", nil)
}
