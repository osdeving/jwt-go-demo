package main

import (
	"encoding/base64"
	"testing"
)

// Casos de teste
var testCases = []struct {
	input    string
	expected string
}{
	{"Manaed", "TWFuYWVk"},
	{"1234567890", "MTIzNDU2Nzg5MA=="},
	{"Hello World", "SGVsbG8gV29ybGQ="},
	{"Go", "R28="},
	{"A", "QQ=="},
	{"abcdef", "YWJjZGVm"},
	{"Base64 encoding", "QmFzZTY0IGVuY29kaW5n"},
	{"", ""},
}

// Testa se a nossa implementação gera a mesma saída do Go
func TestEncodeBase64(t *testing.T) {
	for _, tc := range testCases {
		got := Encb64([]byte(tc.input)) // ⚠️ Use o nome correto da função
		want := base64.StdEncoding.EncodeToString([]byte(tc.input))

		if got != want {
			t.Errorf("Erro na conversão de %q. Obtido: %q, Esperado: %q", tc.input, got, want)
		}
	}
}

// Benchmark da nossa implementação
func BenchmarkEncodeBase64(b *testing.B) {
	data := []byte("Benchmarking Base64 encoding performance!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encb64(data) // ⚠️ Use o nome correto da função
	}
}

// Benchmark da implementação oficial do Go
func BenchmarkEncodeBase64Go(b *testing.B) {
	data := []byte("Benchmarking Base64 encoding performance!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = base64.StdEncoding.EncodeToString(data)
	}
}
