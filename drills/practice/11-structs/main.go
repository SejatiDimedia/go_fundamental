package main

import (
	"encoding/json"
	"fmt"
)

// TODO 1: Buat struct 'AuditLog'
// Field:
// - CreatedAt (string, tag json: "created_at")
// - UpdatedAt (string, tag json: "updated_at")
type AuditLog struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TODO 2: Buat struct 'Transaction' yang meng-embed AuditLog
// Field:
// - ID (string, tag json: "id")
// - Amount (int, tag json: "amount")
// - AuditLog (Anonymous field / Struct Embedding)
// - SecretNote (string, tag json: "-" agar rahasia & tidak pernah keluar di JSON)
type Transaction struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	AuditLog
	SecretNote string `json:"-"`
}

func main() {
	// TODO 3:
	// Buat instance Transaction:
	// - ID: "TRX-999"
	// - Amount: 750000
	// - CreatedAt: "2026-11-18" (isi langsung via field promotion: trx.CreatedAt = ...)
	// - UpdatedAt: "2026-11-18"
	// - SecretNote: "Catatan internal fraud analysis"
	var trx Transaction
	trx.ID = "TRX-999"
	trx.Amount = 750000
	trx.CreatedAt = "2026-11-18"
	trx.UpdatedAt = "2026-11-18"
	trx.SecretNote = "Catatan internal fraud analysis"

	// TODO 4:
	// Marshal struct Transaction tersebut ke JSON:
	// jsonBytes, err := json.Marshal(trx)
	// Cetak string JSON-nya ke terminal.
	//
	// Ekspektasi Output JSON:
	// {"id":"TRX-999","amount":750000,"created_at":"2026-11-18","updated_at":"2026-11-18"}
	// (Perhatikan: SecretNote TIDAK BOLEH muncul di dalam JSON!)

	jsonBytes, err := json.Marshal(trx)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
