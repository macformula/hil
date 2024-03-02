package canlink

import "os"

type FileType interface {
	dumpToFile(*os.File) error
	getFile() (*os.File, error)
}
