package hush

import (
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
