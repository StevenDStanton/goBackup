package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type App struct {
	db      *sql.DB
	watcher *fsnotify.Watcher
	logger  *logrus.Logger
	config  Config
}

type Config struct {
	LocalFilePath string
	DatabasePath  string
}

func NewApp() *App {
	logger := logrus.New()
	config := loadConfig()
	db := setupDatabase(config, logger)
	watcher := setupWatcher(config.LocalFilePath, logger)

	return &App{
		db:      db,
		watcher: watcher,
		logger:  logger,
		config:  config,
	}
}

func loadConfig() Config {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	return Config{
		LocalFilePath: os.Getenv("LOCAL_FILE_PATH"),
		DatabasePath:  "filehashes.db",
	}
}

func setupDatabase(config Config, logger *logrus.Logger) *sql.DB {
	db, err := sql.Open("sqlite", config.DatabasePath)
	if err != nil {
		logger.Fatalf("Error opening database connection: %v", err)
	}

	createFilehashesTable := `
    CREATE TABLE IF NOT EXISTS filehashes (
        Path TEXT PRIMARY KEY,
        Hash TEXT,
        Uploaded BOOLEAN
    );`
	if _, err = db.Exec(createFilehashesTable); err != nil {
		db.Close()
		logger.Fatalf("Unable to create filehashes table: %v", err)
	}

	createRetryFilesTable := `
    CREATE TABLE IF NOT EXISTS retry_files (
        Path TEXT PRIMARY KEY
    );`
	if _, err = db.Exec(createRetryFilesTable); err != nil {
		db.Close()
		logger.Fatalf("Unable to create retry_files table: %v", err)
	}

	return db
}

func setupWatcher(pathToWatch string, logger *logrus.Logger) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatalf("Error creating watcher: %v", err)
	}

	if err := watcher.Add(pathToWatch); err != nil {
		logger.Fatalf("Error adding path to watcher: %v", err)
	}

	return watcher
}

func (app *App) Run() {
	defer app.db.Close()
	defer app.watcher.Close()

	app.initializeFileHashes()

	app.logger.Info("Starting to watch for file events...")
	for {
		select {
		case event := <-app.watcher.Events:
			go app.handleFileEvent(event)
		case err := <-app.watcher.Errors:
			app.logger.Errorf("Error watching directory: %v", err)
		}
	}
}

func (app *App) handleFileEvent(event fsnotify.Event) {
	app.logger.Infof("Detected file change: %s", event.Name)
	// Add your specific handling logic here, e.g., rehash the file
}

func (app *App) initializeFileHashes() {
	err := filepath.Walk(app.config.LocalFilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return app.addFileToDB(path)
	})
	if err != nil {
		app.logger.Error("Error walking the directory:", err)
	}
}

func (app *App) addFileToDB(path string) error {
	newHash, err := app.computeFileHash(path)
	if err != nil {
		return err
	}

	var dbHash string
	err = app.db.QueryRow("SELECT Hash FROM filehashes WHERE Path = ?", path).Scan(&dbHash)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows || dbHash != newHash {
		_, err = app.db.Exec("REPLACE INTO filehashes (Path, Hash, Uploaded) VALUES (?, ?, FALSE)", path, newHash)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) computeFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	app := NewApp()
	app.Run()
}
