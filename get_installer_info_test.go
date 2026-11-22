package hush

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"os/exec"
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
		// ROMs
		{"test.nes", InstallerTypeNaked},
		{"test.GBA", InstallerTypeNaked},
		{"test.nds", InstallerTypeNaked},
		// Unregistered extensions install as naked
		{"test.unknown", InstallerTypeNaked},
		{"test.UNKNOWN", InstallerTypeNaked},
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

func TestGetInstallerInfo_UnregisteredIsNaked(t *testing.T) {
	consumer := &state.Consumer{}

	zipBytes := testZipBytes(t)

	tests := []struct {
		filename     string
		content      []byte
		expectedType InstallerType
	}{
		{"tfitch-07212016-bin", []byte("\x7fELF\x01\x01\x01\x00rest of binary"), InstallerTypeNaked},
		{"game-macos", []byte("\xcf\xfa\xed\xfe\x07\x00\x00\x01"), InstallerTypeNaked},
		{"README", []byte("just some text"), InstallerTypeNaked},
		{"empty", []byte{}, InstallerTypeNaked},
		{"Terri-Fried.nds", []byte("\x2e\x00\x00\xea not really a rom"), InstallerTypeNaked},
		{"Game v1.2", []byte("version suffix parsed as an extension"), InstallerTypeNaked},
		// zips inside, but with no registered extension we don't extract
		{"map.pk3", zipBytes, InstallerTypeNaked},
		{"zipped-no-ext", zipBytes, InstallerTypeNaked},
		// registered zip-shaped formats stay whole
		{"game.love", zipBytes, InstallerTypeNaked},
		{"game.jar", zipBytes, InstallerTypeNaked},
		{"library.yymps", zipBytes, InstallerTypeNaked},
		// registered archives are still probed and extracted
		{"game.zip", zipBytes, InstallerTypeArchive},
		{"game.ZIP", zipBytes, InstallerTypeArchive},
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

			if tt.expectedType == InstallerTypeArchive {
				require.NotNil(t, info.ArchiveInfo)
				assert.Len(t, info.Entries, 1)
				return
			}

			assert.Nil(t, info.ArchiveInfo)

			// a naked source must not have been read from
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
	require.NoError(t, os.WriteFile(path, testZipBytes(t), 0644))
}

func testZipBytes(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer

	w := zip.NewWriter(&buf)
	entry, err := w.Create("hello.txt")
	require.NoError(t, err)
	_, err = entry.Write([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func TestGetInstallerInfo_ShortTarExtensions(t *testing.T) {
	if _, err := exec.LookPath("tar"); err != nil {
		t.Skip("tar not available")
	}
	consumer := &state.Consumer{}

	tests := []struct {
		filename string
		tarFlag  string
		format   string
	}{
		{"game.tgz", "-z", "tar.gz"},
		{"game.TGZ", "-z", "tar.gz"},
		{"game.tbz2", "-j", "tar.bz2"},
	}

	// .txz is registered too, but probing xz loads the 7-zip library, so
	// it is left to boar's own tests.
	assert.Equal(t, InstallerTypeArchive, installerForExt[".txz"])

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello"), 0644))

			filePath := filepath.Join(dir, tt.filename)
			out, err := exec.Command("tar", "-c", tt.tarFlag, "-f", filePath, "-C", dir, "hello.txt").CombinedOutput()
			require.NoError(t, err, string(out))

			file, err := os.Open(filePath)
			require.NoError(t, err)
			defer file.Close()

			info, err := GetInstallerInfo(consumer, file)
			require.NoError(t, err)
			assert.Equal(t, InstallerTypeArchive, info.Type)
			require.NotNil(t, info.ArchiveInfo)
			assert.Equal(t, tt.format, info.ArchiveInfo.Format)
		})
	}
}
