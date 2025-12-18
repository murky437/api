package dbbackup

import (
	"context"
	"io"
	"murky_api/internal/app"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockClient struct {
	UploadedKeys  []string
	ListKeysCount int
	DeletedKeys   []string
}

func newMockClient() *mockClient {
	return &mockClient{
		UploadedKeys:  []string{},
		ListKeysCount: 0,
		DeletedKeys:   []string{},
	}
}

func (m *mockClient) Upload(key string, bodyReader io.Reader) error {
	m.UploadedKeys = append(m.UploadedKeys, key)
	return nil
}

func (m *mockClient) ListObjectKeys(path string) (objectKeys []string, err error) {
	m.ListKeysCount++
	return []string{
		"20251111T000000.sqlite3",
		"20251111T000001.sqlite3",
		"20251111T000002.sqlite3",
		"20251111T000003.sqlite3",
		"20251111T000004.sqlite3",
		"20251111T000008.sqlite3",
		"20251111T000007.sqlite3",
		"20251111T000006.sqlite3",
		"20251111T000005.sqlite3",
		"20251111T000009.sqlite3",
	}, nil
}

func (m *mockClient) DeleteObject(key string) error {
	m.DeletedKeys = append(m.DeletedKeys, key)
	return nil
}

type mockClock struct{}

func (f mockClock) Now() time.Time {
	return time.Date(2025, 11, 22, 12, 0, 0, 0, time.UTC)
}

func TestHandleSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	c.Config.DbBackupDir = "test_backup_dir"
	c.Config.S3DbBackupPath = "test_s3_backup_path"

	err := os.MkdirAll(c.Config.DbBackupDir, 0o755)
	require.NoError(t, err)

	mockS3Client := newMockClient()

	err = Handle(c.Db, c.Config, mockS3Client, mockClock{})(context.Background(), nil)
	require.NoError(t, err)

	_, err = os.Stat("test_backup_dir/20251122T120000.sqlite3")
	require.NoError(t, err)

	err = os.RemoveAll(c.Config.DbBackupDir)
	require.NoError(t, err)

	require.Equal(t, []string{"test_s3_backup_path/20251122T120000.sqlite3"}, mockS3Client.UploadedKeys)
	require.Equal(t, 1, mockS3Client.ListKeysCount)
	require.Equal(t, []string{
		"20251111T000000.sqlite3",
		"20251111T000001.sqlite3",
		"20251111T000002.sqlite3",
	}, mockS3Client.DeletedKeys)

}
