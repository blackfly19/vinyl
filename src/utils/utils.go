package utils

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"os"
)

type md5Type interface {
	~string | ~[]byte
}

func CalculateMD5[T md5Type](input T) string {
	switch v := any(input).(type) {
	case []byte:
		newHash := md5.Sum(v)
		return hex.EncodeToString(newHash[:])
	case string:
		newHash := md5.Sum([]byte(v))
		return hex.EncodeToString(newHash[:])
	default:
		panic("Invalid type")
	}
}

func GenerateCommitID(rootHash string, commitMessage string, dateTime string) string {
	return CalculateMD5(rootHash + commitMessage + dateTime)
}

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
