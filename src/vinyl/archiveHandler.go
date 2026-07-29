package vinyl

import (
	"archive/zip"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadFileFromArchives(archiveName string, fileName string) *zip.File {

	archiveReader, err := zip.OpenReader(archiveName)
	if err != nil {
	}
	defer func(archive *zip.ReadCloser) {
		err := archive.Close()
		if err != nil {
		} // Handle error
	}(archiveReader)

	for _, file := range archiveReader.File {
		if file.Name == fileName {
			return file
		}
	}
	return nil
}

func CreateArchivesFromRoot(archiveName string, rootPath string, skipFileNames ...string) {

	archive, err := os.Create(archiveName)
	if err != nil {
	} // Add error handling
	defer func(archive *os.File) {
		err := archive.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(archive)

	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(zipWriter)

	err = filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			//fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == ".vinyl" {
			return filepath.SkipDir
		}

		file, err := os.ReadFile(path)

		if !info.IsDir() {

			for _, fileName := range skipFileNames {
				if strings.Contains(path, fileName) {
					return nil
				}
			}

			writer, err := zipWriter.Create(strings.TrimPrefix(path, "./"))
			if err != nil {
				return err
			}

			_, err = writer.Write(file)
		}
		return nil
	})
}
