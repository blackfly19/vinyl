package qwe

import (
	"errors"
	"os"
	"time"
)

type FileMetaData struct {
	FileContentMD5Hash string
	LastModifiedTime   time.Time
}

func IsModified(projectFileHashes map[string]FileMetaData, path string) (FileMetaData, error) {

	// Check if the file exists and getting info
	info, err := os.Stat(path)
	if err != nil {
		return FileMetaData{}, err
	}

	file, err := os.ReadFile(path)
	if err != nil {

		return FileMetaData{}, err
	}

	if _, exists := projectFileHashes[path]; exists {

		if info.ModTime().After(projectFileHashes[path].LastModifiedTime) {
			updatedHash := CalculateMD5(file)

			if updatedHash != projectFileHashes[path].FileContentMD5Hash {
				return FileMetaData{LastModifiedTime: info.ModTime(), FileContentMD5Hash: updatedHash}, nil
			}
		}
	} else {
		return FileMetaData{LastModifiedTime: info.ModTime(), FileContentMD5Hash: CalculateMD5(file)}, nil
	}

	return FileMetaData{}, nil
}

func DeletedFiles(projectFileHashes map[string]FileMetaData) []string {

	var deletedFiles []string
	for file, _ := range projectFileHashes {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			deletedFiles = append(deletedFiles, file)
		}
	}
	return deletedFiles
}
