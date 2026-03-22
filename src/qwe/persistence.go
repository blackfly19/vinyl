package qwe

type Persistence interface {
	WriteToDisk() error
	ReadFromDisk() error
}
