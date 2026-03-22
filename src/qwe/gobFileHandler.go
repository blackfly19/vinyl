package qwe

import (
	"encoding/gob"
	"os"
)

func GobEncoder[T any](file *os.File, object T) error {
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(object)
	if err != nil {
		return err
	}
	return nil
}

func GobDecoder[T any](file *os.File, object *T) error {
	decoder := gob.NewDecoder(file)
	err := decoder.Decode(object)
	if err != nil {
		return err
	}
	return nil
}
