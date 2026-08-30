package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// Simulasi Database PostgreSQL Asli (Lambat, baca dari disk 50ms)
type PostgresDB struct {
	data map[string]Product
}

func (db *PostgresDB) QueryProductFromDisk(id string) (Product, bool) {
	// Simulasi I/O latency disk database (50ms)
	time.Sleep(50 * time.Millisecond)
	p, ok := db.data[id]
	return p, ok
}

// Simulasi Redis In-Memory Cache (Super Cepat, baca dari RAM 1ms)
type RedisCache struct {
	mu      sync.RWMutex
	storage map[string]string // Key -> JSON String
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		storage: make(map[string]string),
	}
}

func (r *RedisCache) Get(key string) (string, bool) {
	// Simulasi latensi memori RAM (1ms)
	time.Sleep(1 * time.Millisecond)

	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.storage[key]
	return val, ok
}

func (r *RedisCache) Set(key string, jsonVal string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[key] = jsonVal
}

func (r *RedisCache) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, key)
}

// Service Layer yang Menerapkan Cache-Aside Pattern
type ProductService struct {
	db    *PostgresDB
	cache *RedisCache
}

func (s *ProductService) GetProductDetail(productID string) (*Product, string, time.Duration) {
	startTime := time.Now()
	cacheKey := "cache:product:" + productID

	// 1. CEK KE REDIS DULU (Langkah Pertama)
	if cachedJSON, found := s.cache.Get(cacheKey); found {
		var prod Product
		_ = json.Unmarshal([]byte(cachedJSON), &prod)
		duration := time.Since(startTime)
		return &prod, "🟢 CACHE HIT (Dari Redis RAM)", duration
	}

	// 2. JIKA TIDAK ADA DI REDIS (CACHE MISS), AMBIL DARI POSTGRESQL
	prod, exists := s.db.QueryProductFromDisk(productID)
	if !exists {
		return nil, "NOT FOUND", time.Since(startTime)
	}

	// 3. SIMPAN SALINANNYA KE REDIS UNTUK REQUEST BERIKUTNYA
	jsonBytes, _ := json.Marshal(prod)
	s.cache.Set(cacheKey, string(jsonBytes))

	duration := time.Since(startTime)
	return &prod, "🟡 CACHE MISS (Dari PostgreSQL Disk -> Lalu disimpan ke Redis)", duration
}

func (s *ProductService) UpdateProductPrice(productID string, newPrice int) {
	fmt.Printf("\n✍️ [ADMIN] Mengubah harga %s di PostgreSQL menjadi Rp%d...\n", productID, newPrice)

	// 1. Update di Database Utama
	if prod, ok := s.db.data[productID]; ok {
		prod.Price = newPrice
		s.db.data[productID] = prod
	}

	// 2. INVALIDASI CACHE (Hapus cache lama dari Redis agar data tidak basi!)
	cacheKey := "cache:product:" + productID
	s.cache.Delete(cacheKey)
	fmt.Println("🗑️ [CACHE INVALIDATION] Cache lama di Redis berhasil dihapus!")
}

func main() {
	// Inisialisasi Database & Redis
	pgDB := &PostgresDB{
		data: map[string]Product{
			"PRD-101": {ID: "PRD-101", Name: "Keychron K2 Keyboard", Price: 1250000, Stock: 15},
		},
	}
	redis := NewRedisCache()
	service := &ProductService{db: pgDB, cache: redis}

	fmt.Println("==================================================================")
	fmt.Println("   ⚡ DEMONSTRASI CACHE-ASIDE PATTERN (POSTGRESQL + REDIS)       ")
	fmt.Println("==================================================================\n")

	// Request 1: Pengunjung Pertama Membuka Produk PRD-101
	fmt.Println(">> [Request 1] User A membuka halaman produk PRD-101...")
	p1, source1, time1 := service.GetProductDetail("PRD-101")
	fmt.Printf("   Sumber : %s\n", source1)
	fmt.Printf("   Data   : %s (Harga: Rp%d)\n", p1.Name, p1.Price)
	fmt.Printf("   Waktu  : %v\n\n", time1.Round(time.Millisecond))

	// Request 2: Pengunjung Kedua Membuka Produk yang Sama
	fmt.Println(">> [Request 2] User B membuka halaman produk yang SAMA (PRD-101)...")
	p2, source2, time2 := service.GetProductDetail("PRD-101")
	fmt.Printf("   Sumber : %s\n", source2)
	fmt.Printf("   Data   : %s (Harga: Rp%d)\n", p2.Name, p2.Price)
	fmt.Printf("   Waktu  : %v (50x LEBIH CEPAT!)\n\n", time2.Round(time.Millisecond))

	// Request 3: Pengunjung Ketiga Membuka Produk yang Sama
	fmt.Println(">> [Request 3] User C membuka halaman produk yang SAMA (PRD-101)...")
	_, source3, time3 := service.GetProductDetail("PRD-101")
	fmt.Printf("   Sumber : %s\n", source3)
	fmt.Printf("   Waktu  : %v\n", time3.Round(time.Millisecond))

	// Admin Mengubah Harga & Invalidasi Cache
	service.UpdateProductPrice("PRD-101", 1100000)

	// Request 4: Pengunjung Keempat Membuka Produk Setelah Harganya Diubah
	fmt.Println("\n>> [Request 4] User D membuka produk setelah harga diubah...")
	p4, source4, time4 := service.GetProductDetail("PRD-101")
	fmt.Printf("   Sumber : %s\n", source4)
	fmt.Printf("   Data   : %s (Harga Baru: Rp%d)\n", p4.Name, p4.Price)
	fmt.Printf("   Waktu  : %v\n", time4.Round(time.Millisecond))

	fmt.Println("==================================================================")
}
