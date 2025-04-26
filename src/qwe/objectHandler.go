package qwe

import (
	"os"
)

func CreateDataObjects(fileContent []byte, fileName string) error {

	file, err := os.Create(".qwe/objects/" + fileName + ".gob")
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
