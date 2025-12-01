package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// Benchmark untuk endpoint /recommend
func BenchmarkRecommendHandler(b *testing.B) {
    // Buat request baru
    req, _ := http.NewRequest("GET", "/recommend?month=Januari&region=Utara", nil)
    
    // Buat recorder untuk menangkap response
    w := httptest.NewRecorder()
    
    // Reset timer sebelum loop
    b.ResetTimer()
    
    // Jalankan test sebanyak b.N kali
    for i := 0; i < b.N; i++ {
        recommendHandler(w, req)
    }
}

// Benchmark untuk endpoint /care
func BenchmarkCareHandler(b *testing.B) {
    req, _ := http.NewRequest("GET", "/care?plant=padi", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        careHandler(w, req)
    }
}

// Benchmark untuk endpoint /production
func BenchmarkProductionHandler(b *testing.B) {
    req, _ := http.NewRequest("GET", "/production", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        productionHandler(w, req)
    }
}

// Benchmark untuk endpoint /fertilizer
func BenchmarkFertilizerHandler(b *testing.B) {
    req, _ := http.NewRequest("GET", "/fertilizer?plant=jagung", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fertilizerHandler(w, req)
    }
}

// Benchmark untuk endpoint /weather
func BenchmarkWeatherHandler(b *testing.B) {
    req, _ := http.NewRequest("GET", "/weather?month=Juli", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        weatherHandler(w, req)
    }
}

// Benchmark untuk endpoint /analysis
func BenchmarkAnalysisHandler(b *testing.B) {
    req, _ := http.NewRequest("GET", "/analysis", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        analysisHandler(w, req)
    }
}