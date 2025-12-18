package dbbackup

import (
	"context"
	"database/sql"
	"log"
	"murky_api/internal/clock"
	"murky_api/internal/config"
	"murky_api/internal/s3"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/hibiken/asynq"
)

const (
	NumOfDbBackupsToKeep = 7
)

func Handle(db *sql.DB, conf *config.Config, s3Client s3.Client, cl clock.Clock) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		backupFilePath, err := createDbBackup(db, conf.DbBackupDir, cl)
		if err != nil {
			log.Println(err)
			return err
		}

		err = uploadBackupToS3(s3Client, backupFilePath, conf.S3DbBackupPath)
		if err != nil {
			log.Println(err)
			return err
		}

		err = deleteOldBackupsFromLocal(conf.DbBackupDir, NumOfDbBackupsToKeep)
		if err != nil {
			log.Println(err)
			return err
		}

		err = deleteOldBackupsFromS3(s3Client, conf.S3DbBackupPath, NumOfDbBackupsToKeep)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	}
}

var backupMu sync.Mutex

func createDbBackup(db *sql.DB, dbBackupDir string, cl clock.Clock) (backupFilePath string, err error) {
	backupMu.Lock()
	defer backupMu.Unlock()

	timestamp := cl.Now().Format("20060102T150405")
	backupFilePath = filepath.Join(dbBackupDir, timestamp+".sqlite3")

	_, err = db.Exec("VACUUM INTO ?", backupFilePath)
	if err != nil {
		if strings.Contains(err.Error(), "output file already exists") {
			log.Println("Backup already exists, skipping:", backupFilePath)
			return backupFilePath, nil
		}
		return backupFilePath, err
	}

	log.Println("Backup created at:", backupFilePath)

	return backupFilePath, nil
}

func uploadBackupToS3(s3Client s3.Client, backupFilePath string, s3DbBackupPath string) error {
	f, err := os.Open(backupFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	key := strings.TrimSuffix(s3DbBackupPath, "/") + "/" + filepath.Base(backupFilePath)

	err = s3Client.Upload(key, f)
	if err != nil {
		return err
	}

	log.Println("Backup uploaded to s3")

	return nil
}

func deleteOldBackupsFromLocal(dbBackupDir string, keep int) error {
	entries, err := os.ReadDir(dbBackupDir)
	if err != nil {
		return err
	}

	var backups []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sqlite3") {
			backups = append(backups, e.Name())
		}
	}

	if len(backups) <= keep {
		return nil
	}

	sort.Strings(backups)

	toDelete := backups[:len(backups)-keep]

	for _, f := range toDelete {
		_ = os.Remove(filepath.Join(dbBackupDir, f))
	}

	log.Println("Deleted old backups from local")

	return nil
}

func deleteOldBackupsFromS3(client s3.Client, s3DbBackupPath string, keep int) error {
	objectKeys, err := client.ListObjectKeys(s3DbBackupPath)
	if err != nil {
		return err
	}

	var keys []string
	for _, key := range objectKeys {
		if strings.HasSuffix(key, ".sqlite3") {
			keys = append(keys, key)
		}
	}

	if len(keys) <= keep {
		return nil
	}

	sort.Strings(keys)

	for _, k := range keys[:len(keys)-keep] {
		_ = client.DeleteObject(k)
	}

	log.Println("Deleted old backups from s3")

	return nil
}
