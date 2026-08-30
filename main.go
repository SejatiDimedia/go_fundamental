package main

import (
	"log/slog"
	"os"
	"time"
)

func main() {
	// ==========================================================
	// 1. Text Logger (Format Manusiawi untuk Development Lokal)
	// ==========================================================
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Tangkap semua level dari Debug sampai Error
	}))

	slog.SetDefault(textLogger)

	slog.Info("Aplikasi E-Commerce Service dimulai",
		"version", "1.0.0",
		"port", 8080,
		"env", "development",
	)

	slog.Debug("Inisialisasi koneksi pool database", "max_conns", 20)

	// ==========================================================
	// 2. JSON Logger (Standar Wajib Production / Cloud Kubernetes)
	// ==========================================================
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Simulasi Transaksi Pembayaran
	jsonLogger.Info("Transaksi berhasil diproses",
		slog.String("trx_id", "TRX-998877"),
		slog.Int("amount", 250000),
		slog.String("customer_id", "USR-101"),
		slog.Duration("latency", 45*time.Millisecond),
		slog.Group("payment_details",
			slog.String("method", "BCA_VA"),
			slog.String("status", "PAID"),
		),
	)

	// Simulasi Error di Server
	jsonLogger.Error("Gagal menghubungi Payment Gateway pihak ketiga",
		slog.String("trx_id", "TRX-998878"),
		slog.String("gateway", "Midtrans"),
		slog.String("error_reason", "connection refused / timeout"),
		slog.Int("http_status", 504),
	)
}
