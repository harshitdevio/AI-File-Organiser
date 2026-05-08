package mover

import (
	"bytes"
	"testing"
)

func BenchmarkProcessMoveWithDeps_Mock(b *testing.B) {
	mockFS := &MockFileSystemOps{}
	writer := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessMoveWithDeps("/source/file.txt", "documents", 0.8, mockFS, writer)
		writer.Reset()
		mockFS.MkdirAllCalls = nil
		mockFS.RenameCalls = nil
	}
}

func BenchmarkProcessMoveWithDeps_ConfidenceCheck(b *testing.B) {
	mockFS := &MockFileSystemOps{}
	writer := &bytes.Buffer{}

	confidences := []float64{0.1, 0.3, 0.5, 0.6, 0.7, 0.9}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conf := confidences[i%len(confidences)]
		ProcessMoveWithDeps("/source/file.txt", "topic", conf, mockFS, writer)
		writer.Reset()
		mockFS.MkdirAllCalls = nil
		mockFS.RenameCalls = nil
	}
}