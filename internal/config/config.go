package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AllowedOrigins                  []string
	DbFilePath                      string
	GuestDbFilePath                 string
	AccessTokenSecret               string
	RefreshTokenSecret              string
	DeployFilePath                  string
	DbBackupDir                     string
	RedisAddress                    string
	Timezone                        string
	AwsRegion                       string
	AwsEndpointUrl                  string
	S3Bucket                        string
	S3DbBackupPath                  string
	GoogleServiceAccountKeyFilePath string
	BrevoApiKey                     string
}

func NewConfig() *Config {
	exePath, _ := os.Executable()
	deployFilePath := filepath.Join(filepath.Dir(exePath), "deploy.json")

	return &Config{
		AllowedOrigins:                  strings.Split(os.Getenv("ALLOWED_ORIGINS"), ";"),
		DbFilePath:                      os.Getenv("DB_FILE_PATH"),
		GuestDbFilePath:                 os.Getenv("GUEST_DB_FILE_PATH"),
		AccessTokenSecret:               os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret:              os.Getenv("REFRESH_TOKEN_SECRET"),
		DeployFilePath:                  deployFilePath,
		DbBackupDir:                     os.Getenv("DB_BACKUP_DIR"),
		RedisAddress:                    os.Getenv("REDIS_ADDRESS"),
		Timezone:                        os.Getenv("TIMEZONE"),
		AwsRegion:                       os.Getenv("AWS_REGION"),
		AwsEndpointUrl:                  os.Getenv("AWS_ENDPOINT_URL"),
		S3Bucket:                        os.Getenv("S3_BUCKET"),
		S3DbBackupPath:                  os.Getenv("S3_DB_BACKUP_PATH"),
		GoogleServiceAccountKeyFilePath: os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_FILE_PATH"),
		BrevoApiKey:                     os.Getenv("BREVO_API_KEY"),
	}
}
