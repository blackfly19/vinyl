package qwe

import (
	"os"
)

type MapHandler[V any] struct {
	FileMap  map[string]V
	FilePath string
}

func NewMapHandler[V any](filePath string) *MapHandler[V] {
	return &MapHandler[V]{FilePath: filePath, FileMap: make(map[string]V)}
}

func (m *MapHandler[V]) WriteToDisk() error {
	file, err := os.OpenFile(m.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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

	err = GobEncoder(file, m.FileMap)
	if err != nil {
		return err
	}

	return nil
}

func (m *MapHandler[V]) ReadFromDisk() error {
	gobfile, err := os.Open(m.FilePath)
	if err != nil {
		return err
	}

	err = GobDecoder(gobfile, &m.FileMap)
	if err != nil {
		return err
	}

	return nil
}
