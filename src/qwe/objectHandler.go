package qwe

import (
	"github.com/blackfly19/vcs/src/constants"
	"os"
)

func CreateDataObjects(fileContent []byte, fileName string) error {

	file, err := os.Create(constants.DIR_OBJECTS + fileName)
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

	err = GobEncoder(file, fileContent)
	if err != nil {
		return err
	}

	return nil

}
