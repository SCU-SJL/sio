package stream

import "io"

type ZipStreamData interface {
	StreamData
	GetFileName() string
}

type DefaultZipStreamData struct {
	baseData DefaultStreamData
	filename string
}

func NewDefaultZipStreamData(r io.ReadCloser, filename string) ZipStreamData {
	return DefaultZipStreamData{
		baseData: NewDefaultStreamData(r),
		filename: filename,
	}
}

func (d DefaultZipStreamData) GetReadCloser() io.ReadCloser {
	return d.baseData.GetReadCloser()
}

func (d DefaultZipStreamData) GetFileName() string {
	return d.filename
}
