package scanner

import (
	"io"
	"strings"
	"testing"
)

func BenchmarkGetDirectory_WithFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetDirectory("/home/user/test", "")
	}
}
 
func BenchmarkGetDirectory_WithInput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetDirectory("", "/home/user/test")
	}
}
 
func BenchmarkScanWithDeps_WithFlag(b *testing.B) {
	writer := io.Discard
	reader := strings.NewReader("")
 
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ScanWithDeps("/home/user/test", reader, writer)
	}
}
 
func BenchmarkScanWithDeps_WithInput(b *testing.B) {
	writer := io.Discard
 
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader("/home/user/test\n")
		ScanWithDeps("", reader, writer)
	}
}