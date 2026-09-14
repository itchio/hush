package hush

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/itchio/headway/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInstallerInfo_CaseInsensitiveExtension(t *testing.T) {
	consumer := &state.Consumer{}

	tests := []struct {
		filename     string
		expectedType InstallerType
	}{
		// Lowercase
		{"test.mp3", InstallerTypeNaked},
		{"test.pdf", InstallerTypeNaked},
		{"test.exe", InstallerTypeNaked},
		// Uppercase
		{"test.MP3", InstallerTypeNaked},
		{"test.PDF", InstallerTypeNaked},
		{"test.EXE", InstallerTypeNaked},
		// Mixed case
		{"test.Mp3", InstallerTypeNaked},
		{"test.Pdf", InstallerTypeNaked},
		{"test.ExE", InstallerTypeNaked},
		// AppImage (Linux portable apps)
		{"test.appimage", InstallerTypeNaked},
		{"test.AppImage", InstallerTypeNaked},
		// Unknown extensions should remain unknown regardless of case
		{"test.unknown", InstallerTypeUnknown},
		{"test.UNKNOWN", InstallerTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			filePath := filepath.Join(dir, tt.filename)

			err := os.WriteFile(filePath, []byte("test content"), 0644)
			require.NoError(t, err)

			file, err := os.Open(filePath)
			require.NoError(t, err)
			defer file.Close()

			info, err := GetInstallerInfo(consumer, file)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedType, info.Type)
		})
	}
}

func TestGetInstallerInfo_SniffExtensionlessExecutable(t *testing.T) {
	consumer := &state.Consumer{}

	tests := []struct {
		filename     string
		content      []byte
		expectedType InstallerType
	}{
		{"tfitch-07212016-bin", []byte("\x7fELF\x01\x01\x01\x00rest of binary"), InstallerTypeNaked},
		{"game-macos", []byte("\xcf\xfa\xed\xfe\x07\x00\x00\x01"), InstallerTypeNaked},
		{"game-fat", []byte("\xca\xfe\xba\xbe\x00\x00\x00\x02"), InstallerTypeNaked},
		{"README", []byte("just some text"), InstallerTypeUnknown},
		{"tiny", []byte("\x7fE"), InstallerTypeUnknown},
		{"empty", []byte{}, InstallerTypeUnknown},
		{"binary.unknown", []byte("\x7fELF\x01\x01\x01\x00"), InstallerTypeNaked},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			filePath := filepath.Join(dir, tt.filename)

			err := os.WriteFile(filePath, tt.content, 0644)
			require.NoError(t, err)

			file, err := os.Open(filePath)
			require.NoError(t, err)
			defer file.Close()

			info, err := GetInstallerInfo(consumer, file)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedType, info.Type)

			// sniffing must leave the file at the start
			pos, err := file.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(0), pos)
		})
	}
}

func TestGetInstallerInfo_ExtensionOverrides(t *testing.T) {
	consumer := &state.Consumer{}
	dir := t.TempDir()

	// a real zip so the archive probe succeeds
	filePath := filepath.Join(dir, "game.dmg")
	writeTestZip(t, filePath)

	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	// default registry: .dmg is naked
	info, err := GetInstallerInfoWithParams(GetInstallerInfoParams{
		Consumer: consumer,
		File:     file,
	})
	require.NoError(t, err)
	assert.Equal(t, InstallerTypeNaked, info.Type)

	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)

	// override: .dmg routed to archive, and probed as one
	info, err = GetInstallerInfoWithParams(GetInstallerInfoParams{
		Consumer: consumer,
		File:     file,
		ExtensionOverrides: map[string]InstallerType{
			".dmg": InstallerTypeArchive,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, InstallerTypeArchive, info.Type)
	require.NotNil(t, info.ArchiveInfo)
	assert.Len(t, info.Entries, 1)
}

func writeTestZip(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	w := zip.NewWriter(f)
	entry, err := w.Create("hello.txt")
	require.NoError(t, err)
	_, err = entry.Write([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, w.Close())
}
