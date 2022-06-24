package stream

import "io"

type Stream chan StreamData

type StreamData interface {
	GetReadCloser() io.ReadCloser
}

type DefaultStreamData struct {
	r io.ReadCloser
}

func NewDefaultStreamData(r io.ReadCloser) DefaultStreamData {
	return DefaultStreamData{
		r: r,
	}
}

func (d DefaultStreamData) GetReadCloser() io.ReadCloser {
	return d.r
}
