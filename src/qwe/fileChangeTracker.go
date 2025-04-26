package qwe

import (
	"os"
	"time"
)

type FileMetaData struct {
	FileContentMD5Hash string
	LastModifiedTime   time.Time
}

func WriteFileHashToDisk(filePath string, filePathMetaDataMap map[string]FileMetaData) error {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(file *os.File) error {
		err := file.Close()
		if err != nil {
			return err
		}
		return nil
	}(file)

	err = GobEncoder(file, filePathMetaDataMap)
	if err != nil {
		return err
	}

	return nil
}

func ReadFileHashFromDisk(fileName string) (map[string]FileMetaData, error) {

	gobfile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	filePathMetaDataMap := make(map[string]FileMetaData)

	err = GobDecoder(gobfile, &filePathMetaDataMap)
	if err != nil {
		return nil, err
	}

	return filePathMetaDataMap, nil
}
