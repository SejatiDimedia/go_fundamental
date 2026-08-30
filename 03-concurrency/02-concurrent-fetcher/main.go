package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Vendor merepresentasikan endpoint payment gateway yang akan dicek
type Vendor struct {
	Name     string
	Endpoint string
	// Waktu simulasi latensi server vendor (dalam milidetik)
	SimulatedLatency time.Duration
}

// VendorResponse merepresentasikan hasil pengecekan status
type VendorResponse struct {
	VendorName string
	Latency    time.Duration
	RateUSD    float64
	IsOnline   bool
	Error      error
}

// fetchVendorStatus bertugas memanggil API vendor dengan memperhatikan Context Cancellation
func fetchVendorStatus(ctx context.Context, v Vendor, resultChan chan<- VendorResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	// Membuat channel internal untuk simulasi respons dari server vendor
	internalDone := make(chan VendorResponse, 1)

	// Goroutine anak untuk simulasi request jaringan
	go func() {
		time.Sleep(v.SimulatedLatency)

		// Simulasi kurs acak
		rate := 15500.0 + (rand.Float64() * 200.0)

		internalDone <- VendorResponse{
			VendorName: v.Name,
			Latency:    v.SimulatedLatency,
			RateUSD:    rate,
			IsOnline:   true,
			Error:      nil,
		}
	}()

	// Menunggu: Mana yang lebih dulu terjadi? Server vendor membalas ATAU Context Timeout?
	select {
	case <-ctx.Done():
		// Skenario A: Waktu Timeout Habis sebelum vendor membalas!
		resultChan <- VendorResponse{
			VendorName: v.Name,
			Latency:    time.Since(startTime),
			IsOnline:   false,
			Error:      ctx.Err(), // ctx.Err() akan berisi: "context deadline exceeded"
		}

	case resp := <-internalDone:
		// Skenario B: Server vendor berhasil merespons tepat waktu
		resultChan <- resp
	}
}

func main() {
	// 1. Daftar 5 Payment Gateway yang akan dicek serentak
	// Ada beberapa vendor yang sengaja disimulasikan lambat (> 1200ms) untuk menguji timeout
	vendors := []Vendor{
		{Name: "Midtrans Gateway", Endpoint: "https://api.midtrans.com/v2/status", SimulatedLatency: 350 * time.Millisecond},
		{Name: "Xendit Payment", Endpoint: "https://api.xendit.co/rates", SimulatedLatency: 600 * time.Millisecond},
		{Name: "DOKU Core", Endpoint: "https://api.doku.com/health", SimulatedLatency: 1500 * time.Millisecond}, // LAMBAT (Akan kena timeout!)
		{Name: "Stripe International", Endpoint: "https://api.stripe.com/v1/ping", SimulatedLatency: 450 * time.Millisecond},
		{Name: "Faspay Engine", Endpoint: "https://api.faspay.co.id/check", SimulatedLatency: 1800 * time.Millisecond}, // LAMBAT (Akan kena timeout!)
	}

	// 2. Menentukan Batas Waktu Toleransi (Timeout = 1.2 Detik)
	timeoutDuration := 1200 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel() // Best practice: Selalu panggil cancel() untuk membersihkan resource context

	resultChan := make(chan VendorResponse, len(vendors))
	var wg sync.WaitGroup

	fmt.Println("==========================================================================")
	fmt.Printf("   🌐 MULTI-API CONCURRENT FETCHER (Batas Waktu Timeout: %v)\n", timeoutDuration)
	fmt.Println("==========================================================================")
	fmt.Println(">> Mengirim request ke 5 Payment Gateway secara bersamaan...\n")

	fetchStart := time.Now()

	// 3. Luncurkan pengecekan ke 5 vendor secara paralel di background
	for _, v := range vendors {
		wg.Add(1)
		go fetchVendorStatus(ctx, v, resultChan, &wg)
	}

	// 4. Goroutine pengawas untuk menutup channel hasil
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 5. Menampilkan Laporan Status Vendor
	fmt.Printf("%-24s | %-12s | %-12s | %-20s\n", "Vendor Gateway", "Status", "Latensi", "Info / Kurs USD")
	fmt.Println("--------------------------------------------------------------------------")

	onlineCount := 0
	timeoutCount := 0

	for res := range resultChan {
		if res.IsOnline {
			onlineCount++
			fmt.Printf("%-24s | 🟢 %-10s | %-12v | Rp%.2f / USD\n",
				res.VendorName, "ONLINE", res.Latency, res.RateUSD)
		} else {
			timeoutCount++
			fmt.Printf("%-24s | 🔴 %-10s | %-12v | ❌ %v\n",
				res.VendorName, "TIMEOUT", res.Latency.Round(time.Millisecond), res.Error)
		}
	}

	totalElapsed := time.Since(fetchStart)

	fmt.Println("==========================================================================")
	fmt.Printf("📊 RINGKASAN HASIL:\n")
	fmt.Printf("   - Gateway Online  : %d dari %d vendor\n", onlineCount, len(vendors))
	fmt.Printf("   - Gateway Timeout : %d dari %d vendor\n", timeoutCount, len(vendors))
	fmt.Printf("   - Total Waktu     : %v (Dibatasi ketat oleh batas timeout 1.2s)\n", totalElapsed.Round(time.Millisecond))
	fmt.Println("==========================================================================")
}
