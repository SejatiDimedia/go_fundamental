package calculator

import (
	"errors"
	"testing"
)

func TestCalculateTotalDiscount(t *testing.T) {
	// Mendefinisikan Tabel Kasus Uji (Test Cases Table)
	testCases := []struct {
		name            string
		price           float64
		discountPercent int
		expectedTotal   float64
		expectedErr     error
	}{
		{
			name:            "Diskon Normal 10%",
			price:           100000,
			discountPercent: 10,
			expectedTotal:   90000,
			expectedErr:     nil,
		},
		{
			name:            "Diskon 0% (Tanpa Diskon)",
			price:           50000,
			discountPercent: 0,
			expectedTotal:   50000,
			expectedErr:     nil,
		},
		{
			name:            "Diskon 100% (Gratis)",
			price:           75000,
			discountPercent: 100,
			expectedTotal:   0,
			expectedErr:     nil,
		},
		{
			name:            "Error: Harga Negatif",
			price:           -5000,
			discountPercent: 10,
			expectedTotal:   0,
			expectedErr:     ErrNegativePrice,
		},
		{
			name:            "Error: Diskon Melebihi 100%",
			price:           100000,
			discountPercent: 150,
			expectedTotal:   0,
			expectedErr:     errors.New("invalid"), // Diskon tidak valid
		},
	}

	// Menjalankan seluruh skenario di dalam tabel
	for _, tc := range testCases {
		// t.Run menjalankan setiap baris sebagai Sub-Test mandiri
		t.Run(tc.name, func(t *testing.T) {
			result, err := CalculateTotalDiscount(tc.price, tc.discountPercent)

			// 1. Validasi Error
			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Ekspektasi error, tapi fungsi mengembalikan nil!")
				}
				if errors.Is(tc.expectedErr, ErrNegativePrice) && !errors.Is(err, ErrNegativePrice) {
					t.Errorf("Ekspektasi error %v, tapi dapat %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Tidak diekspektasikan error, tapi terjadi error: %v", err)
				}
			}

			// 2. Validasi Hasil Nilai
			if result != tc.expectedTotal {
				t.Errorf("Hasil salah! Ekspektasi: %.2f, Didapat: %.2f", tc.expectedTotal, result)
			}
		})
	}
}

// 1. Unit Test untuk Fungsi Divide
func TestDivide(t *testing.T) {
	t.Run("Pembagian Normal", func(t *testing.T) {
		res, err := Divide(10, 2)
		if err != nil || res != 5 {
			t.Errorf("Ekspektasi 5 tanpa error, tapi dapat %.2f (err: %v)", res, err)
		}
	})

	t.Run("Pembagian dengan Nol (Error)", func(t *testing.T) {
		_, err := Divide(10, 0)
		if !errors.Is(err, ErrDivisionByZero) {
			t.Errorf("Ekspektasi ErrDivisionByZero, tapi dapat: %v", err)
		}
	})
}

// 2. Benchmark Fungsi Diskon
func BenchmarkCalculateTotalDiscount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateTotalDiscount(100000, 15)
	}
}
