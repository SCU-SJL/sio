package stream

import (
	"io"
)

type MultipartStreamData interface {
	StreamData
	GetFileName() string
	GetFieldName() string
}

type DefaultMultipartStreamData struct {
	DefaultStreamData
	filename, fieldName string
}

func NewDefaultMultipartStreamData(r io.ReadCloser, filename, fieldName string) MultipartStreamData {
	return DefaultMultipartStreamData{
		DefaultStreamData: NewDefaultStreamData(r),
		filename:          filename,
		fieldName:         fieldName,
	}
}

func (d DefaultMultipartStreamData) GetFileName() string {
	return d.filename
}

func (d DefaultMultipartStreamData) GetFieldName() string {
	return d.fieldName
}
