package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ===== Struktur Data =====
type Plant struct {
	Nama      string `json:"nama"`
	Musim     string `json:"musim"`
	Wilayah   string `json:"wilayah"`
	Produksi  int    `json:"produksi"`
	Deskripsi string `json:"deskripsi"`
}

type Response struct {
	Bulan       string        `json:"bulan"`
	Musim       string        `json:"musim"`
	Wilayah     string        `json:"wilayah"`
	Rekomendasi []Rekomendasi `json:"rekomendasi"`
}

type Rekomendasi struct {
	Tanaman   string `json:"tanaman"`
	Skor      int    `json:"skor"`
	Deskripsi string `json:"deskripsi"`
}

// ====== DATA UTAMA ======
var plants = []Plant{
	{"Padi", "Hujan", "Utara", 95, "Tanaman utama di wilayah utara, cocok di musim hujan dengan curah hujan tinggi."},
	{"Kedelai", "Hujan", "Utara", 85, "Ditanam setelah padi di musim hujan dengan drainase baik."},
	{"Jagung", "Peralihan", "Tengah", 90, "Tahan terhadap cuaca tidak menentu dan cocok di lahan sedang."},
	{"Cabai", "Peralihan", "Tengah", 80, "Cocok di tanah gembur dengan sinar matahari cukup."},
	{"Tembakau", "Kemarau", "Selatan", 98, "Unggulan Jember bagian selatan pada musim kemarau."},
	{"Jagung", "Kemarau", "Selatan", 85, "Tahan panas dan minim curah hujan."},
	{"Padi", "Peralihan", "Utara", 80, "Masih cocok ditanam di awal musim peralihan."},
	{"Kedelai", "Kemarau", "Utara", 70, "Masih bisa tumbuh di akhir kemarau dengan irigasi cukup."},
	{"Cabai", "Kemarau", "Tengah", 85, "Hasil baik di tanah gembur saat panas tidak ekstrem."},
}

// ====== DATA TAMBAHAN ======
var produksi = map[string]int{
	"Padi":     95,
	"Kedelai":  85,
	"Jagung":   90,
	"Cabai":    80,
	"Tembakau": 98,
}

var tips = map[string]string{
	"padi":     "Gunakan sistem irigasi yang baik dan pupuk organik.",
	"kedelai":  "Tanam di tanah gembur dan hindari genangan air.",
	"jagung":   "Pastikan sinar matahari cukup dan pengairan teratur.",
	"cabai":    "Gunakan mulsa plastik hitam perak untuk menjaga kelembapan.",
	"tembakau": "Cocok di musim kemarau, hindari curah hujan tinggi.",
}

var pupuk = map[string]string{
	"padi":     "Urea 100kg/ha, NPK 150kg/ha, pupuk organik 2 ton/ha",
	"kedelai":  "NPK 100kg/ha, Pupuk kandang 1 ton/ha",
	"jagung":   "Urea 120kg/ha, KCl 50kg/ha, NPK 100kg/ha",
	"cabai":    "Kompos 2 ton/ha, NPK 200kg/ha, dolomit 100kg/ha",
	"tembakau": "ZA 100kg/ha, SP36 75kg/ha, pupuk organik 1,5 ton/ha",
}

// ===== FUNGSI PENDUKUNG =====
func getMusimFromMonth(month string) string {
	month = strings.ToLower(strings.TrimSpace(month))
	switch month {
	case "desember", "januari", "februari", "maret":
		return "Hujan"
	case "april", "mei", "oktober", "november":
		return "Peralihan"
	case "juni", "juli", "agustus", "september":
		return "Kemarau"
	default:
		return "Tidak diketahui"
	}
}

// ===== HANDLER API =====

// 1️⃣ Rekomendasi tanaman
func recommendHandler(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	region := r.URL.Query().Get("region")
	season := r.URL.Query().Get("season")

	if season == "" {
		season = getMusimFromMonth(month)
	}

	var recs []Rekomendasi

	for _, plant := range plants {
		skor := 0
		if strings.EqualFold(plant.Musim, season) {
			skor += 40
		}
		if strings.Contains(strings.ToLower(region), strings.ToLower(plant.Wilayah)) {
			skor += 40
		}
		skor += plant.Produksi / 3

		if skor > 60 {
			recs = append(recs, Rekomendasi{
				Tanaman:   plant.Nama,
				Skor:      skor,
				Deskripsi: plant.Deskripsi,
			})
		}
	}

	resp := Response{
		Bulan:       month,
		Musim:       season,
		Wilayah:     region,
		Rekomendasi: recs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 2️⃣ Menampilkan semua tanaman
func plantsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(plants)
}

// 3️⃣ Panduan perawatan
func careHandler(w http.ResponseWriter, r *http.Request) {
	plant := strings.ToLower(r.URL.Query().Get("plant"))
	if val, ok := tips[plant]; ok {
		json.NewEncoder(w).Encode(map[string]string{
			"tanaman": plant,
			"panduan": val,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Tanaman tidak ditemukan",
		})
	}
}

// 4️⃣ Produksi
func productionHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(produksi)
}

// 5️⃣ Rekomendasi pupuk
func fertilizerHandler(w http.ResponseWriter, r *http.Request) {
	plant := strings.ToLower(r.URL.Query().Get("plant"))
	if val, ok := pupuk[plant]; ok {
		json.NewEncoder(w).Encode(map[string]string{
			"tanaman":   plant,
			"pupuk":     val,
			"keterangan": "Rekomendasi dosis berdasarkan rata-rata hasil panen terbaik.",
		})
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Tanaman tidak ditemukan",
		})
	}
}

// 6️⃣ Simulasi cuaca
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	month := strings.ToLower(r.URL.Query().Get("month"))
	season := getMusimFromMonth(month)

	cuaca := map[string]string{
		"Hujan":     "Curah hujan tinggi, suhu 23-28°C, kelembapan 80-90%.",
		"Peralihan": "Hujan tidak menentu, suhu 26-30°C, kelembapan sedang.",
		"Kemarau":   "Curah hujan rendah, suhu 30-34°C, kelembapan rendah.",
	}

	info := cuaca[season]
	if info == "" {
		info = "Data cuaca tidak tersedia untuk bulan tersebut."
	}

	json.NewEncoder(w).Encode(map[string]string{
		"bulan": month,
		"musim": season,
		"info":  info,
	})
}

// 7️⃣ Analisis keseluruhan wilayah
func analysisHandler(w http.ResponseWriter, r *http.Request) {
	hasil := map[string]any{
		"wilayah_utara":   []string{"Padi", "Kedelai"},
		"wilayah_tengah":  []string{"Jagung", "Cabai"},
		"wilayah_selatan": []string{"Tembakau", "Jagung"},
		"musim_terbaik":   "Kemarau dan Peralihan",
		"produksi_tertinggi": map[string]any{
			"tanaman": "Tembakau",
			"nilai":   98,
		},
	}
	json.NewEncoder(w).Encode(hasil)
}

// ===== MAIN =====
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/recommend", recommendHandler).Methods("GET")
	r.HandleFunc("/plants", plantsHandler).Methods("GET")
	r.HandleFunc("/care", careHandler).Methods("GET")
	r.HandleFunc("/production", productionHandler).Methods("GET")
	r.HandleFunc("/fertilizer", fertilizerHandler).Methods("GET")
	r.HandleFunc("/weather", weatherHandler).Methods("GET")
	r.HandleFunc("/analysis", analysisHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	log.Println("🌾 Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
