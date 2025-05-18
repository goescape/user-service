package middlewares

import (
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// salt adalah nilai tetap (static salt) yang digunakan untuk hashing password.
// Perlu dipertimbangkan penggunaan salt yang unik per user untuk keamanan lebih baik.
var salt = []byte("secret-password")

// GeneratePassword menghasilkan hash dari password plaintext menggunakan algoritma Argon2.
// Hasil hash dikembalikan dalam bentuk string base64 tanpa padding.
func GeneratePassword(password string) string {
	timeCost := uint32(1)        // jumlah iterasi (semakin tinggi, semakin aman tapi lambat)
	memCost := uint32(64 * 1024) // memori dalam KB yang digunakan
	parallelism := uint8(1)      // jumlah thread yang digunakan
	hashLength := uint32(32)     // panjang hasil hash dalam byte

	hash := argon2.IDKey([]byte(password), salt, timeCost, memCost, parallelism, hashLength)

	return base64.RawStdEncoding.EncodeToString(hash)
}

// VerifyPassword membandingkan password plaintext dengan hash yang sudah disimpan.
// Hash ulang password lalu cocokkan hasilnya dengan hash input.
func VerifyPassword(hash, password string) bool {
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 1, 32)
	encodedHash := base64.RawStdEncoding.EncodeToString(newHash)

	return hash == encodedHash
}
