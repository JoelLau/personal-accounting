package testutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TemporaryFile(t testing.TB, content string) *os.File {
	t.Helper()

	tempFile, err := os.CreateTemp(t.TempDir(), "test-file-*")
	require.NoError(t, err, "failed to create temporary file")

	t.Cleanup(func() {
		_ = tempFile.Close()
	})

	_, err = tempFile.WriteString(content)
	require.NoError(t, err, "failed to write content to temp file")

	_, err = tempFile.Seek(0, 0)
	require.NoError(t, err, "failed to seek to start of file")

	return tempFile
}
