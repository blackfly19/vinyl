package vinyl

type Persistence interface {
	WriteToDisk(fileName string) error
	ReadFromDisk(fileName string) error
}
