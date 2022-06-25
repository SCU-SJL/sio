package stream

import "io"

type ZipStreamData interface {
	StreamData
	GetFileName() string
}

type DefaultZipStreamData struct {
	DefaultStreamData
	filename string
}

func NewDefaultZipStreamData(r io.ReadCloser, filename string) ZipStreamData {
	return DefaultZipStreamData{
		DefaultStreamData: NewDefaultStreamData(r),
		filename:          filename,
	}
}

func (d DefaultZipStreamData) GetFileName() string {
	return d.filename
}
