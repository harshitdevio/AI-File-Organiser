package mover

import (
	"os"
)

type MockFileSystemOps struct {
	MkdirAllFunc func(path string, perm os.FileMode) error
	RenameFunc   func(oldpath, newpath string) error

	MkdirAllCalls []MkdirAllCall
	RenameCalls   []RenameCall
}

type MkdirAllCall struct {
	Path string
	Perm os.FileMode
}

type RenameCall struct {
	OldPath string
	NewPath string
}

func (m *MockFileSystemOps) MkdirAll(path string, perm os.FileMode) error {
	m.MkdirAllCalls = append(m.MkdirAllCalls, MkdirAllCall{Path: path, Perm: perm})
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockFileSystemOps) Rename(oldpath, newpath string) error {
	m.RenameCalls = append(m.RenameCalls, RenameCall{OldPath: oldpath, NewPath: newpath})
	if m.RenameFunc != nil {
		return m.RenameFunc(oldpath, newpath)
	}
	return nil
}