package main

import (
	"encoding/json" // JSON verilerini işlemek için kullanılan paket
	"fmt"
	"github.com/shirou/gopsutil/cpu" // CPU bilgilerini almak için kullanılan  kütüphane
	"github.com/shirou/gopsutil/mem" // Bellek bilgilerini almak için kullanılan  kütüphane
	"html/template"
	"net/http"
)

// Sistem bilgilerini alır ve döndürür
func getSystemStats() (float64, float64, float64, float64, float64, string, float64, error) {
	// CPU bilgilerini alır
	cpuInfo, err := cpu.Info()
	if err != nil {
		return 0, 0, 0, 0, 0, "", 0, err // Hata durumunda sıfır değerler ve hata döndürür
	}
	cpuSpeed := cpuInfo[0].Mhz      // CPU'nun çalışma hızını MHz olarak alır
	cpuName := cpuInfo[0].ModelName // CPU model ismini alır

	// CPU kullanım yüzdesini alır
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, 0, 0, 0, 0, "", 0, err // Hata durumunda sıfır değerler ve hata döndürür
	}

	// Bellek bilgilerini alır
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, 0, 0, "", 0, err // Hata durumunda sıfır değerler ve hata döndürür
	}
	usedMemoryMB := float64(virtualMem.Used) / 1024 / 1024          // Kullanılan belleği MB olarak hesaplar
	totalMemoryGB := float64(virtualMem.Total) / 1024 / 1024 / 1024 // Toplam belleği GB olarak hesaplar
	memoryPercent := (usedMemoryMB / (totalMemoryGB * 1024)) * 100  // Bellek kullanım yüzdesini hesaplar

	// Kullanılan toplam MHz miktarını hesaplar
	usedCpuMHz := (cpuSpeed * cpuPercent[0]) / 100

	return cpuSpeed, cpuPercent[0], usedMemoryMB, totalMemoryGB, memoryPercent, cpuName, usedCpuMHz, nil
}

// Bilgileri JSON formatında döndürür
func statsHandler(w http.ResponseWriter, r *http.Request) {
	cpuSpeed, cpuPercent, memUsage, totalMemory, memPercent, cpuName, usedCpuMHz, err := getSystemStats()
	if err != nil {
		// Hata durumunda 500 kodu döndürür
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Verileri bir maple
	stats := map[string]interface{}{
		"cpuSpeed":      cpuSpeed,
		"cpuPercent":    cpuPercent,
		"memory":        memUsage,
		"memoryPercent": memPercent,
		"totalMemory":   totalMemory,
		"cpuName":       cpuName,
		"usedCpuMHz":    usedCpuMHz,
	}
	json.NewEncoder(w).Encode(stats) // Json Çıktı al
}

// HTML şablonunu yükler - anasayfa
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)      // Anasayfa
	http.HandleFunc("/stats", statsHandler) // /Stats handler
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
