# File Monitor and Backup System

## Overview

As Chrome moved to manifest version 3 and adblock became less effective I started to move out of the Google ecosystem. While I have been a Google fanboy for a number of years and still love many of their products, I am tired of being one of their products. Making the switch to Proton mail earlier this year I have become less reliant on google in general. After reading a wired article about people being locked out of their google drive accounts I decided it was time for me to start moving some of my critical files from Google Drive to Amazon. This is part of that effort.

This system monitors specified directories for file changes, logs these events, and synchronizes changed files to an AWS S3 bucket. It uses an SQLite database to store file hashes and configuration settings, ensuring data integrity and enabling recovery in case of system failures.

## Interesting side-note and a move to only support Linux

Since I first started this project I have switched from Windows to Linux. I was not happy with how Windows was performing for this application since it was not detecting file updates in real time. I am not sure if this was due to the iNotify library I was using or if it was due to the Windows API's being weird, however, at the time I read about other people having the same sort of issues.

Since moving to linux I have circled back to this project. I would still like to have my writing backed up automatically to an S3 bucket. So I am making it work for linux only and using the unix package.

## Features

- Real-time file monitoring.
- File changes logged in an SQLite database.
- Automatic file backup to AWS S3.
- Encryption using AWS on-disk encryption.
- UI for managing configurations and monitoring system status.
- Backup and recovery of the SQLite database to/from AWS S3.

## Development Checklist

- [ ] Set up File Monitoring
- [ ] Scan existing files for Hash
- [ ] Save hash to SQLite DB
- [ ] When restarting scan files for changes since last run
- [ ] Look into some way to check file integrity after upload?
