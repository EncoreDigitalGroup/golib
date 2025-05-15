/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Right Reserved.
 */

package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type CopyDirectory struct {
	BufferSize int
}

// CountFiles recursively counts the total number of files in a directory tree
func (cd *CopyDirectory) CountFiles(directory string) (int, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			subCount, err := cd.CountFiles(filepath.Join(directory, entry.Name()))
			if err != nil {
				return 0, err
			}
			count += subCount
		} else {
			count++
		}
	}
	return count, nil
}

// CopyFilesAndDirectory copies files and directories, sending progress updates via a channel
func (cd *CopyDirectory) CopyFilesAndDirectory(sourceDirectory, destinationDirectory string, progressChan chan struct{}) (int, error) {
	// Check if the destination exists
	if _, err := os.Stat(destinationDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(destinationDirectory, 0755)
		if err != nil {
			return 0, err
		}
	}

	// Get all files and subdirectories from the source directory
	entries, err := os.ReadDir(sourceDirectory)
	if err != nil {
		return 0, err
	}

	// Use a WaitGroup to track directory copying operations
	var waitGroup sync.WaitGroup

	// Track errors from goroutines
	var errorMutex sync.Mutex
	var firstErr error

	// Track file counts
	var fileCount int
	var countMutex sync.Mutex

	// Copy all files and subdirectories from the source directory
	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDirectory, entry.Name())
		destinationPath := filepath.Join(destinationDirectory, entry.Name())

		if entry.IsDir() {
			// Handle subdirectories concurrently
			waitGroup.Add(1)
			go func(src, dst string) {
				defer waitGroup.Done()
				count, err := cd.CopyFilesAndDirectory(src, dst, progressChan)
				countMutex.Lock()
				fileCount += count
				countMutex.Unlock()

				if err != nil {
					errorMutex.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errorMutex.Unlock()
				}
			}(sourcePath, destinationPath)
		} else {
			// Copy file with optimized buffer
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				return fileCount, err
			}

			// Create destination file (overwrite if exists)
			destinationFile, err := os.Create(destinationPath)
			if err != nil {
				_ = sourceFile.Close()
				return fileCount, err
			}

			// Use buffered copy for better performance
			bufSize := 1024 * 1024 // 1MB buffer
			if cd.BufferSize > 0 {
				bufSize = cd.BufferSize
			}
			buf := make([]byte, bufSize)
			_, err = io.CopyBuffer(destinationFile, sourceFile, buf)

			_ = sourceFile.Close()
			_ = destinationFile.Close()

			if err != nil {
				return fileCount, err
			}

			// Send progress update instead of logging
			progressChan <- struct{}{}
			countMutex.Lock()
			fileCount++
			countMutex.Unlock()
		}
	}

	// Wait for all directory copies to complete
	waitGroup.Wait()
	return fileCount, firstErr
}

// Copy orchestrates the copying process with a progress bar
func (cd *CopyDirectory) Copy(source, destination string) (int, error) {
	// Count total files for progress bar initialization
	totalFiles, err := cd.CountFiles(source)
	if err != nil {
		return 0, err
	}

	// Initialize progress bar
	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Copying files: "),
	)

	// Channel for progress updates
	progressChan := make(chan struct{})

	// Goroutine to update progress bar
	go func() {
		for range progressChan {
			_ = bar.Add(1)
		}
	}()

	// Perform the copy operation
	fileCount, err := cd.CopyFilesAndDirectory(source, destination, progressChan)
	close(progressChan) // Close channel after copying is done

	if err != nil {
		return fileCount, err
	}

	// Finish the progress bar
	_ = bar.Finish()
	return fileCount, nil
}

// CopyMultipleSources copies multiple sources to a destination with progress
func (cd *CopyDirectory) CopyMultipleSources(sources []string, destination string) (int, error) {
	// Count total files across all sources first
	totalFiles := 0
	for _, source := range sources {
		count, err := cd.CountFiles(source)
		if err != nil {
			return 0, err
		}
		totalFiles += count
	}

	// Create a single progress bar for all sources
	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Copying files: "),
	)

	// Channel for progress updates from all copy operations
	progressChan := make(chan struct{})

	// Goroutine to update progress bar
	go func() {
		for range progressChan {
			_ = bar.Add(1)
		}
	}()

	// Ensure destination exists
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		err = os.MkdirAll(destination, 0755)
		if err != nil {
			close(progressChan)
			return 0, err
		}
	}

	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(sources))
	countChan := make(chan int, len(sources))

	// Process sources in parallel
	for _, source := range sources {
		waitGroup.Add(1)
		go func(src string) {
			defer waitGroup.Done()
			// Use copyFilesAndDirectory instead of Copy to avoid creating multiple progress bars
			count, err := cd.CopyFilesAndDirectory(src, destination, progressChan)
			countChan <- count
			if err != nil {
				errChan <- err
			}
		}(source)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()
	close(progressChan)
	close(errChan)
	close(countChan)

	// Calculate total files copied
	totalFilesCopied := 0
	for count := range countChan {
		totalFilesCopied += count
	}

	_ = bar.Finish()

	// Return first error if any
	if len(errChan) > 0 {
		return totalFilesCopied, <-errChan
	}
	return totalFilesCopied, nil
}
