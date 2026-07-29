package vinyl

import (
	"github.com/blackfly19/vcs/src/utils"
	"os"
)

type MapHandler[V any] struct {
	FileMap  map[string]V
	FilePath string
}

func NewMapHandler[V any](filePath string) *MapHandler[V] {
	return &MapHandler[V]{FilePath: filePath, FileMap: make(map[string]V)}
}

func (m *MapHandler[V]) WriteToDisk(fileName string) error {
	file, err := os.OpenFile(m.FilePath+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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

	err = utils.GobEncoder(file, m.FileMap)
	if err != nil {
		return err
	}

	return nil
}

func (m *MapHandler[V]) ReadFromDisk(fileName string) error {
	gobfile, err := os.Open(m.FilePath + fileName)
	if err != nil {
		return err
	}

	err = utils.GobDecoder(gobfile, &m.FileMap)
	if err != nil {
		return err
	}

	return nil
}
