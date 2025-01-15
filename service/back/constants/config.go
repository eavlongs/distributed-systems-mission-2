package constants

import "os"

const (
	MAIN_SERVER    = 1
	BACKUP_SERVER  = 2
	ARCHIVE_SERVER = 3
)

type Config struct {
	ServerMapping      map[int]string
	FileStoragePrefix  string
	SyncFileScriptPath string
}

func NewConfig() *Config {
	return &Config{
		ServerMapping: map[int]string{
			MAIN_SERVER:    os.Getenv("MAIN_SERVER_IP"),
			BACKUP_SERVER:  os.Getenv("BACKUP_SERVER_IP"),
			ARCHIVE_SERVER: os.Getenv("ARCHIVE_SERVER_IP"),
		},
		FileStoragePrefix:  os.Getenv("FILE_STORAGE_PATH_PREFIX"),
		SyncFileScriptPath: os.Getenv("SYNC_FILE_SCRIPT_PATH"),
	}
}
