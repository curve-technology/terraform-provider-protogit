package store

import (
	"fmt"
	"os"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

type FileSystem struct {
	folderPath string
}

func NewFileSystem(folderPath string) (*FileSystem, error) {
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("provided path must be a location to a directory")
	}

	return &FileSystem{folderPath: folderPath}, nil
}

func (fs *FileSystem) GetFileDescriptor(filepath string) (*desc.FileDescriptor, error) {
	parser := protoparse.Parser{ImportPaths: []string{fs.folderPath}}

	descriptors, err := parser.ParseFiles(filepath)
	if err != nil {
		return nil, err
	}

	fileDescriptor := descriptors[0]
	return fileDescriptor, nil
}
