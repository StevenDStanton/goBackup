// func main() {

// 	// Initialize AWS Session
// 	awsConfig := &aws.Config{
// 		Region: aws.String(os.Getenv("AWS_REGION")),
// 		Credentials: credentials.NewStaticCredentials(
// 			os.Getenv("AWS_ACCESS_KEY_ID"),
// 			os.Getenv("AWS_SECRET_ACCESS_KEY"),
// 			"",
// 		),
// 	}
// 	sess := session.Must(session.NewSession(awsConfig))
// 	s3Client := s3.New(sess)

// 	// Timer to upload every 15 minutes
// 	timer := time.NewTicker(15 * time.Minute)
// 	defer timer.Stop()

// 	// Event handling loop
// 	for {
// 		select {
// 		case event := <-watcher.Events:
// 			go handleFileEvent(event, db)
// 		case <-timer.C:
// 			go handleUploads(db, s3Client)
// 		case err := <-watcher.Errors:
// 			fmt.Println("Watcher error:", err)
// 		}
// 	}
// }

// func handleFileEvent(event fsnotify.Event, db *sql.DB) {
// 	fmt.Println("Detected file change:", event.Name)

// 	// Only handle write and create operations
// 	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
// 		newHash, err := computeFileHash(event.Name, db)
// 		if err != nil {
// 			fmt.Println("Error hashing file:", err)
// 			return
// 		}

// 		// Insert or update the database
// 		_, err = db.Exec("REPLACE INTO filehashes (Path, Hash, Uploaded) VALUES (?, ?, FALSE)", event.Name, newHash)
// 		if err != nil {
// 			fmt.Println("Error updating database:", err)
// 			return
// 		}
// 		fmt.Println("Database updated for:", event.Name)
// 	}
// }

// func handleUploads(db *sql.DB, s3Client *s3.S3) {
// 	fmt.Println("Handling scheduled uploads...")

// 	// Query for all entries where 'Uploaded' is false
// 	rows, err := db.Query("SELECT Path, Hash FROM filehashes WHERE Uploaded = FALSE")
// 	if err != nil {
// 		fmt.Println("Error querying for unuploaded files:", err)
// 		return
// 	}
// 	defer rows.Close()

// 	// Iterate over the query results
// 	for rows.Next() {
// 		var path, hash string
// 		err := rows.Scan(&path, &hash)
// 		if err != nil {
// 			fmt.Println("Error reading row:", err)
// 			continue
// 		}

// 		// Simulate the upload by logging (replace this with actual S3 upload logic when ready)
// 		fmt.Printf("Simulating upload for file: %s, Hash: %s\n", path, hash)

// 		// Update the 'Uploaded' status in the database
// 		_, err = db.Exec("UPDATE filehashes SET Uploaded = TRUE WHERE Path = ?", path)
// 		if err != nil {
// 			fmt.Println("Error updating upload status:", err)
// 			continue
// 		}
// 	}

// 	// Check for errors from iterating over rows
// 	if err = rows.Err(); err != nil {
// 		fmt.Println("Error during rows iteration:", err)
// 	}
// }

// func setupWatcher(watcher *fsnotify.Watcher, path string) error {
// 	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			if info.Name() == ".git" {
// 				return filepath.SkipDir
// 			}
// 			// Add the directory to the watcher
// 			err = watcher.Add(path)
// 			if err != nil {
// 				fmt.Println("Error adding path to watcher:", err)
// 			}
// 		}
// 		return nil
// 	})
// }

// func retryFiles(db *sql.DB) {
// 	rows, err := db.Query("SELECT Path FROM retry_files")
// 	if err != nil {
// 		fmt.Println("Error querying retry files:", err)
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var path string
// 		if err := rows.Scan(&path); err != nil {
// 			fmt.Println("Error reading retry row:", err)
// 			continue
// 		}

// 		hash, err := computeFileHash(path, db)
// 		if err != nil {
// 			fmt.Println("Still failing to access file:", path, err)
// 			continue
// 		}

// 		// If successful, update the filehashes table and remove from retry_files
// 		_, err = db.Exec("REPLACE INTO filehashes (Path, Hash, Uploaded) VALUES (?, ?, FALSE)", path, hash)
// 		if err != nil {
// 			fmt.Println("Error updating file hash:", err)
// 			continue
// 		}
// 		_, err = db.Exec("DELETE FROM retry_files WHERE Path = ?", path)
// 		if err != nil {
// 			fmt.Println("Error removing file from retry table:", err)
// 			continue
// 		}
// 		fmt.Println("Successfully retried file:", path)
// 	}
// }

