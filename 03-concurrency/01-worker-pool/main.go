package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TransactionJob merepresentasikan satu tugas transaksi yang masuk ke antrean
type TransactionJob struct {
	ID            string
	CustomerName  string
	Amount        int
	PaymentMethod string
}

// TransactionResult merepresentasikan laporan hasil pemrosesan dari worker
type TransactionResult struct {
	Job       TransactionJob
	WorkerID  int
	Success   bool
	Duration  time.Duration
	Reference string
}

// Worker function: Mengambil job dari antrean, memproses, dan mengirim hasil
func transactionWorker(id int, jobs <-chan TransactionJob, results chan<- TransactionResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		startTime := time.Now()
		fmt.Printf("👷 [Worker %d] Mulai memproses transaksi %s (%s - Rp%d)...\n",
			id, job.ID, job.PaymentMethod, job.Amount)

		// Simulasi latensi pemrosesan payment gateway (300ms - 800ms)
		processDelay := time.Duration(300+rand.Intn(500)) * time.Millisecond
		time.Sleep(processDelay)

		// Simulasi hasil transaksi
		refCode := fmt.Sprintf("REF-%s-%04d", job.PaymentMethod, rand.Intn(9999))
		duration := time.Since(startTime)

		fmt.Printf("✅ [Worker %d] Selesai %s dalam %v\n", id, job.ID, duration.Round(time.Millisecond))

		// Kirim hasil ke channel results
		results <- TransactionResult{
			Job:       job,
			WorkerID:  id,
			Success:   true,
			Duration:  duration,
			Reference: refCode,
		}
	}
}

func main() {
	const totalWorkers = 3
	const totalJobs = 10

	// 1. Data Dummy 10 Transaksi
	transactions := []TransactionJob{
		{ID: "TRX-101", CustomerName: "Budi Santoso", Amount: 500000, PaymentMethod: "BCA_VA"},
		{ID: "TRX-102", CustomerName: "Siti Rahma", Amount: 125000, PaymentMethod: "GOPAY"},
		{ID: "TRX-103", CustomerName: "Andi Wijaya", Amount: 2500000, PaymentMethod: "CREDIT_CARD"},
		{ID: "TRX-104", CustomerName: "Dewi Lestari", Amount: 75000, PaymentMethod: "OVO"},
		{ID: "TRX-105", CustomerName: "Eko Prasetyo", Amount: 350000, PaymentMethod: "MANDIRI_VA"},
		{ID: "TRX-106", CustomerName: "Fajar Nugraha", Amount: 890000, PaymentMethod: "QRIS"},
		{ID: "TRX-107", CustomerName: "Gita Savitri", Amount: 1500000, PaymentMethod: "BCA_VA"},
		{ID: "TRX-108", CustomerName: "Hendra Gunawan", Amount: 420000, PaymentMethod: "SHOPEEPAY"},
		{ID: "TRX-109", CustomerName: "Indah Permata", Amount: 650000, PaymentMethod: "BNI_VA"},
		{ID: "TRX-110", CustomerName: "Joko Susilo", Amount: 2100000, PaymentMethod: "CREDIT_CARD"},
	}

	// 2. Siapkan Channels
	jobsChan := make(chan TransactionJob, totalJobs)
	resultsChan := make(chan TransactionResult, totalJobs)

	var workerWg sync.WaitGroup

	overallStart := time.Now()
	fmt.Println("==================================================================")
	fmt.Printf("   🚀 MEMULAI WORKER POOL (%d Workers untuk memproses %d Jobs)   \n", totalWorkers, totalJobs)
	fmt.Println("==================================================================")

	// 3. Menyalakan Tim Worker
	for w := 1; w <= totalWorkers; w++ {
		workerWg.Add(1)
		go transactionWorker(w, jobsChan, resultsChan, &workerWg)
	}

	// 4. Masukkan seluruh 10 pekerjaan ke dalam jobs channel
	for _, job := range transactions {
		jobsChan <- job
	}
	close(jobsChan) // Tutup jobs channel karena semua data sudah dimasukkan

	// 5. Goroutine Pengawas: Menutup resultsChan hanya setelah SEMUA worker selesai
	go func() {
		workerWg.Wait()
		close(resultsChan)
	}()

	// 6. Mengumpulkan & Menampilkan Laporan Hasil Pemrosesan
	fmt.Println("\n--- LAPORAN HASIL TRANSAKSI REALTIME ---")
	totalSuccessAmount := 0
	processedCount := 0

	for res := range resultsChan {
		processedCount++
		totalSuccessAmount += res.Job.Amount
		fmt.Printf("[%d/%d] Ref: %-18s | %-8s | Rp%-9d | Oleh Worker %d\n",
			processedCount, totalJobs, res.Reference, res.Job.ID, res.Job.Amount, res.WorkerID)
	}

	totalDuration := time.Since(overallStart)

	fmt.Println("==================================================================")
	fmt.Printf("🎉 SELURUH BATCH SELESAI!\n")
	fmt.Printf("   Total Transaksi Berhasil : %d dari %d\n", processedCount, totalJobs)
	fmt.Printf("   Total Perputaran Uang    : Rp%d\n", totalSuccessAmount)
	fmt.Printf("   Total Waktu Eksekusi     : %v\n", totalDuration.Round(time.Millisecond))
	fmt.Println("==================================================================")
}
