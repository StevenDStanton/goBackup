package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	fmt.Println("File Monitoring Service Starting...")

	// Initialize AWS Session
	awsConfig := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}
	sess := session.Must(session.NewSession(awsConfig))
	s3Client := s3.New(sess)

	// Setup local Database
	db, err := sql.Open("sqlite", "filehashes.db")
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	// Create the database table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS filehashes (
        Path TEXT PRIMARY KEY,
        Hash TEXT,
        Uploaded BOOLEAN
    );`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// File system watcher setup for specified path
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	// Add path to watcher
	err = watcher.Add(os.Getenv("LOCAL_FILE_PATH"))
	if err != nil {
		fmt.Println("Error adding path to watcher:", err)
		return
	}

	// Timer to upload every 15 minutes
	timer := time.NewTicker(15 * time.Minute)
	defer timer.Stop()

	// Event handling loop
	for {
		select {
		case event := <-watcher.Events:
			go handleFileEvent(event, db)
		case <-timer.C:
			go handleUploads(db, s3Client)
		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)
		}
	}
}

func handleFileEvent(event fsnotify.Event, db *sql.DB) {
	fmt.Println("Detected file change:", event.Name)
	// Implement logic to update DB and handle file change
}

func handleUploads(db *sql.DB, s3Client *s3.S3) {
	fmt.Println("Handling scheduled uploads...")
	// Implement upload logic
}
