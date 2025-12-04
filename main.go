package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

type Plant struct {
	Nama      string `json:"nama"`
	Musim     string `json:"musim"`
	Wilayah   string `json:"wilayah"`
	Produksi  int    `json:"produksi"`
	Deskripsi string `json:"deskripsi"`
	Tips      string `json:"tips,omitempty"`
	Pupuk     string `json:"pupuk,omitempty"`
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

var plants = []Plant{
	{"Padi", "Hujan", "Utara", 95, "Tanaman utama di wilayah utara, cocok di musim hujan dengan curah hujan tinggi.", "Gunakan sistem irigasi yang baik dan pupuk organik.Gunakan sistem irigasi yang baik dan pupuk organik.Gunakan sistem irigasi yang baik dan pupuk organik.", "Urea 100kg/ha, NPK 150kg/ha, pupuk organik 2 ton/ha"},
	{"Kedelai", "Hujan", "Utara", 85, "Ditanam setelah padi di musim hujan dengan drainase baik.", "Tanam di tanah gembur dan hindari genangan air.", "NPK 100kg/ha, Pupuk kandang 1 ton/ha"},
	{"Jagung", "Peralihan", "Tengah", 90, "Tahan terhadap cuaca tidak menentu dan cocok di lahan sedang.", "Pastikan sinar matahari cukup dan pengairan teratur.", "Urea 120kg/ha, KCl 50kg/ha, NPK 100kg/ha"},
	{"Cabai", "Peralihan", "Tengah", 80, "Cocok di tanah gembur dengan sinar matahari cukup.", "Gunakan mulsa plastik hitam perak untuk menjaga kelembapan.", "Kompos 2 ton/ha, NPK 200kg/ha, dolomit 100kg/ha"},
	{"Tembakau", "Kemarau", "Selatan", 98, "Unggulan Jember bagian selatan pada musim kemarau.", "Cocok di musim kemarau, hindari curah hujan tinggi.", "ZA 100kg/ha, SP36 75kg/ha, pupuk organik 1,5 ton/ha"},
	{"Jagung", "Kemarau", "Selatan", 85, "Tahan panas dan minim curah hujan.", "Pastikan sinar matahari cukup dan pengairan teratur.", "Urea 120kg/ha, KCl 50kg/ha, NPK 100kg/ha"},
	{"Padi", "Peralihan", "Utara", 80, "Masih cocok ditanam di awal musim peralihan.", "Gunakan sistem irigasi yang baik dan pupuk organik.Gunakan sistem irigasi yang baik dan pupuk organik.Gunakan sistem irigasi yang baik dan pupuk organik.", "Urea 100kg/ha, NPK 150kg/ha"},
	{"Kedelai", "Kemarau", "Utara", 70, "Masih bisa tumbuh di akhir kemarau dengan irigasi cukup.", "Tanam di tanah gembur dan hindari genangan air.", "NPK 100kg/ha, Pupuk kandang 1 ton/ha"},
	{"Cabai", "Kemarau", "Tengah", 85, "Hasil baik di tanah gembur saat panas tidak ekstrem.", "Gunakan mulsa plastik hitam perak untuk menjaga kelembapan.", "Kompos 2 ton/ha, NPK 200kg/ha"},
}

func getProduksiPerWilayah() map[string]map[string]int {
	hasil := make(map[string]map[string]int)
	for _, p := range plants {
		region := strings.Title(strings.ToLower(strings.TrimSpace(p.Wilayah)))
		if hasil[region] == nil {
			hasil[region] = make(map[string]int)
		}
		existing, ok := hasil[region][p.Nama]
		if !ok || p.Produksi > existing {
			hasil[region][p.Nama] = p.Produksi
		}
	}
	return hasil
}

func getTopProduksiGlobal() (string, int, string) {
	data := getProduksiPerWilayah()
	topName := "-"
	topVal := 0
	topWilayah := "-"
	for wilayah, m := range data {
		for nama, val := range m {
			if val > topVal {
				topVal = val
				topName = nama
				topWilayah = wilayah
			}
		}
	}
	return topName, topVal, topWilayah
}

func findPlantByName(name string) *Plant {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, p := range plants {
		if strings.ToLower(p.Nama) == name {
			cp := p
			return &cp
		}
	}
	return nil
}

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
        
        if region != "" {
            region = strings.ToLower(strings.TrimSpace(region))
            
            var regionPart string
            if strings.Contains(region, "utara") {
                regionPart = "utara"
            } else if strings.Contains(region, "tengah") {
                regionPart = "tengah"
            } else if strings.Contains(region, "selatan") {
                regionPart = "selatan"
            }
            
            if regionPart != "" && strings.EqualFold(regionPart, strings.ToLower(plant.Wilayah)) {
                skor += 40
            }
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

    sort.Slice(recs, func(i, j int) bool {
        return recs[i].Skor > recs[j].Skor
    })

    resp := Response{
        Bulan:       month,
        Musim:       season,
        Wilayah:     region,
        Rekomendasi: recs,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func plantsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plants)
}

func careHandler(w http.ResponseWriter, r *http.Request) {
	plantName := r.URL.Query().Get("plant")
	if plantName == "" {
		http.Error(w, "missing plant parameter", http.StatusBadRequest)
		return
	}
	p := findPlantByName(plantName)
	if p == nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Tanaman tidak ditemukan"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"tanaman": p.Nama,
		"panduan": p.Tips,
	})
}

func productionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getProduksiPerWilayah())
}

func fertilizerHandler(w http.ResponseWriter, r *http.Request) {
	plantName := r.URL.Query().Get("plant")
	if plantName == "" {
		http.Error(w, "missing plant parameter", http.StatusBadRequest)
		return
	}
	p := findPlantByName(plantName)
	if p == nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Tanaman tidak ditemukan"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"tanaman":    p.Nama,
		"pupuk":      p.Pupuk,
		"keterangan": "Rekomendasi dosis berdasarkan rata-rata hasil panen terbaik.",
	})
}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"bulan": strings.Title(month),
		"musim": season,
		"info":  info,
	})
}

func analysisHandler(w http.ResponseWriter, r *http.Request) {
	prod := getProduksiPerWilayah()

	wilUtara := []string{}
	wilTengah := []string{}
	wilSelatan := []string{}
	for wilayah, m := range prod {
		names := make([]string, 0, len(m))
		for nama := range m {
			names = append(names, nama)
		}
		sort.Strings(names)
		switch strings.ToLower(wilayah) {
		case "utara":
			wilUtara = names
		case "tengah":
			wilTengah = names
		case "selatan":
			wilSelatan = names
		}
	}

	topTanaman, topNilai, _ := getTopProduksiGlobal()

	freq := map[string]int{}
	for _, p := range plants {
		freq[p.Musim]++
	}
	type mf struct {
		M string
		F int
	}
	mfs := []mf{}
	for k, v := range freq {
		mfs = append(mfs, mf{M: k, F: v})
	}
	sort.Slice(mfs, func(i, j int) bool { return mfs[i].F > mfs[j].F })
	musimTerbaik := []string{}
	for _, x := range mfs {
		musimTerbaik = append(musimTerbaik, x.M)
	}

	result := map[string]interface{}{
		"wilayah_utara":       wilUtara,
		"wilayah_tengah":      wilTengah,
		"wilayah_selatan":     wilSelatan,
		"produksi_tertinggi":  map[string]interface{}{"tanaman": topTanaman, "nilai": topNilai},
		"musim_berurut":       musimTerbaik, 
		"total_tanaman":       len(plants),
		"produksi_per_wilayah": prod,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

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

	log.Println("Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", r)
}