package store

import (
	"github.com/jhump/protoreflect/desc"
)

type Storer interface {
	GetFileDescriptor(filepath string) (*desc.FileDescriptor, error)
}
