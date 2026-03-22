package qwe

import (
	"crypto/md5"
	"encoding/hex"
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

func GenerateCommitID(rootHash string, commitMessage string) string {
	return CalculateMD5(rootHash + commitMessage)
}
