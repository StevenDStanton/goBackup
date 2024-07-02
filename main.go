package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/unix"
)

var (
	rootDir = "~/Dropbox/Writing"
)

func init() {
	localOS := runtime.GOOS

	if localOS != "linux" {
		panic("Sorry, I only work on linux...")
	}
}

func main() {
	fd, err := initInotifyWatches()
	if err != nil {
		log.Fatalf("Error 1: %v", err)
	}
	defer unix.Close(fd)

	buffer := make([]byte, unix.SizeofInotifyEvent*4096)

	for {
		n, err := unix.Read(fd, buffer)
		if err != nil {
			log.Fatalf("Error 2: %v", err)
		}
	}

}

func initInotifyWatches() (int, error) {
	fd, err := unix.InotifyInit()
	if err != nil {
		return -1, err
	}

	err = addWacthesRecusively(fd)
	if err != nil {
		unix.Close(fd)
		return -1, err
	}

	return fd, nil
}

func addWacthesRecusively(fd int) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			_, err := unix.InotifyAddWatch(fd, path, unix.IN_CREATE|unix.IN_MODIFY|unix.IN_DELETE)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
