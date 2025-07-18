package core

import (
	"fmt"
	"github.com/pkg/sftp"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ModifyEnv(env []byte, newHost, envFieldName string) []byte {
	lines := strings.Split(string(env), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, envFieldName) {
			lines[i] = fmt.Sprintf("%s\"%s\"", envFieldName, newHost)
		}
	}
	return []byte(strings.Join(lines, "\n"))
}

func prepareRemoteDirectory(client *sftp.Client, path string) error {
	if err := client.RemoveDirectory(path); err == nil {
		log.Printf("Removed existing directory: %s", path)
	}

	if err := client.MkdirAll(path); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}
	log.Printf("Created directory: %s", path)
	return nil
}

func CopyFiles(client *sftp.Client, localPath, remotePath string) error {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("local directory %s does not exist", localPath)
	}

	var filesToCopy []string
	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filesToCopy = append(filesToCopy, path)
			log.Printf("Found file to copy: %s (%d bytes)", path, info.Size())
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to scan local directory: %v", err)
	}

	if len(filesToCopy) == 0 {
		return fmt.Errorf("no files found in %s", localPath)
	}

	for _, srcPath := range filesToCopy {
		relPath, err := filepath.Rel(localPath, srcPath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %v", srcPath, err)
		}

		dstPath := filepath.ToSlash(filepath.Join(remotePath, relPath))
		dstDir := filepath.ToSlash(filepath.Dir(dstPath))

		if err := client.MkdirAll(dstDir); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dstDir, err)
		}

		if err := copySingleFileWithRetry(client, srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %v", srcPath, dstPath, err)
		}

		log.Printf("Successfully copied: %s -> %s", srcPath, dstPath)
	}

	return nil
}

func copySingleFileWithRetry(client *sftp.Client, srcPath, dstPath string) error {
	const maxAttempts = 3
	var lastError error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := tryCopyFile(client, srcPath, dstPath); err == nil {
			return nil
		} else {
			lastError = err
			log.Printf("Attempt %d failed for %s: %v", attempt, srcPath, err)
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	return fmt.Errorf("after %d attempts, last error: %v", maxAttempts, lastError)
}

func tryCopyFile(client *sftp.Client, srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := client.Create(dstPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy content failed: %v", err)
	}

	if err := client.Chmod(dstPath, 0644); err != nil {
		log.Printf("Warning: cannot set permissions for %s: %v", dstPath, err)
	}

	return nil
}
