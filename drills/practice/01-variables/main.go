package main

import "fmt"

// TODO 1: Definisikan custom type 'UserRole' dengan tipe dasar 'int'
// type UserRole ...
type UserRole int

// TODO 2: Buat konstanta enum menggunakan iota
// - RoleUnknown = 0 (Zero Value Defense)
// - RoleGuest   = 1
// - RoleMember  = 2
// - RoleAdmin   = 3
// const ( ... )
const (
	RoleUnknown UserRole = iota
	RoleGuest
	RoleMember
	RoleAdmin
)

// TODO 3: Buat fungsi IsAuthorized(role UserRole) bool
// Return true HANYA jika RoleMember atau RoleAdmin. Selain itu return false.
// func IsAuthorized ...
func IsAuthorized(role UserRole) bool {
	if role == RoleMember || role == RoleAdmin {
		return true
	}
	return false
}

func main() {
	// TODO 4:
	// a. Deklarasikan variabel 'var unassignedRole UserRole' (manfaatkan zero value)
	// b. Panggil IsAuthorized(unassignedRole) dan cetak hasilnya
	var unassignedRole UserRole
	fmt.Println(IsAuthorized(unassignedRole))

	// TODO 5:
	// a. Buat variabel dengan nilai RoleAdmin
	// b. Panggil IsAuthorized(...) dan cetak hasilnya
	adminRole := RoleAdmin
	fmt.Println(IsAuthorized(adminRole))

	fmt.Println("Silakan lengkapi kode di atas!")
}
