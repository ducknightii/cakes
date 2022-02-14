package types

import "io"

type Reader interface {
	io.Reader
	io.ReaderAt
}
