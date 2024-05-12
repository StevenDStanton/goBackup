# File Monitor and Backup System

## Overview

As Chrome moved to manifest version 3 and adblock became less effective I started to move out of the Google ecosystem. While I have been a Google fanboy for a number of years and still love many of their products, I am tired of being one of their products. Making the switch to Proton mail earlier this year I have become less reliant on google in general. After reading a wired article about people being locked out of their google drive accounts I decided it was time for me to start moving some of my critical files from Google Drive to Amazon. This is part of that effort.

This system monitors specified directories for file changes, logs these events, and synchronizes changed files to an AWS S3 bucket. It uses an SQLite database to store file hashes and configuration settings, ensuring data integrity and enabling recovery in case of system failures.

## Features

- Real-time file monitoring.
- File changes logged in an SQLite database.
- Automatic file backup to AWS S3.
- Encryption using AWS on-disk encryption.
- UI for managing configurations and monitoring system status.
- Backup and recovery of the SQLite database to/from AWS S3.

## Development Checklist

### Setup and Initial Configuration

- [x] Set up a Go project environment.
- [x] Install necessary Go packages (`fsnotify`, `aws-sdk-go`, etc.).
- [x] Create an initial project structure.

### Database Design

- [x] Design and set up an SQLite database schema for storing file hashes and configurations.
- [x] Implement database interaction utilities (connect, read, write).

### File Monitoring

- [x] Implement file monitoring using `fsnotify`.
- [ ] Develop logic to calculate and compare file hashes.
- [ ] Store and retrieve file hashes from the SQLite database.

### AWS Integration

- [x] Set up AWS CLI and ensure it's configured on the development machine.
- [ ] Implement file upload functionality using AWS SDK for Go.
- [ ] Ensure encryption is enabled for file uploads.

### User Interface

- [ ] Design a basic UI layout using a suitable Go library or framework.
- [ ] Implement system tray integration for minimization and background running.
- [ ] Develop configuration management sections in the UI.

### Backup and Recovery

- [ ] Implement automatic backup of the SQLite database to S3.
- [ ] Develop a recovery process to restore the SQLite database and files from S3.
- [ ] Allow the user to specify a local folder for recovery downloads.

### Notifications and Logging

- [ ] Set up a logging system to track file changes and system errors.
- [ ] Implement user notifications for critical events and statuses.

### Security and Performance

- [ ] Ensure secure storage of sensitive configuration data.
- [ ] Optimize performance for handling large numbers of files or very large files.

### Testing and Deployment

- [ ] Thoroughly test file monitoring, backup, and recovery functionalities.
- [ ] Test UI interactions and configuration management.
- [ ] Prepare deployment guidelines and scripts.

### Documentation

- [ ] Write comprehensive user documentation covering setup, use, and troubleshooting.
- [ ] Document the code and API endpoints (if any).

## Future Enhancements

- [ ] Consider adding support for monitoring multiple directories.
- [ ] Explore the integration of additional cloud storage providers.
