package main

import "fmt"

// ==========================================
// KASUS: Mengapa Return Struct Jauh Lebih Baik?
// ==========================================

// 1. Dependency (Database): MENGGUNAKAN INTERFACE (Accept Interface)
type UserRepository interface {
	SaveUser(name string) error
}

// 2. Service Utama: Struct Konkret
type UserService struct {
	repo UserRepository
}

func (s *UserService) Register(name string) {
	fmt.Println("1. Registrasi user:", name)
	s.repo.SaveUser(name)
}

func (s *UserService) SendWelcomeEmail(name string) {
	fmt.Println("2. Mengirim email selamat datang ke:", name)
}

// 3. Constructor: MENGEMBALIKAN POINTER STRUCT (*UserService)
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ==========================================
// Implementasi A: Mock DB (Untuk Unit Test)
// ==========================================
type MockRepo struct{}

func (m *MockRepo) SaveUser(name string) error {
	fmt.Println("   [Mock DB] Data berhasil disimpan ke memori dummy.")
	return nil
}

// ==========================================
// Implementasi B: Real DB (Untuk Production!)
// ==========================================
type PostgresRepo struct {
	ConnectionString string
}

func (p *PostgresRepo) SaveUser(name string) error {
	fmt.Printf("   [Postgres Production DB] INSERT INTO users VALUES ('%s') via %s\n", name, p.ConnectionString)
	return nil
}

func main() {
	fmt.Println("=== 1. SAAT UNIT TEST ===")
	mockDB := &MockRepo{}
	testService := NewUserService(mockDB)
	testService.Register("User Dummy Test")

	fmt.Println("\n=== 2. SAAT PRODUCTION DI SERVER ===")
	// Cukup ganti mockDB dengan struct DB asli:
	realDB := &PostgresRepo{ConnectionString: "postgres://admin:secret@localhost:5432/production_db"}
	prodService := NewUserService(realDB)
	prodService.Register("Sultan Andara")
	prodService.SendWelcomeEmail("Sultan Andara")
}
