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

func CreateAsyncStream(dataList []StreamData) Stream {

	stream := make(Stream)

	go func() {

		for _, data := range dataList {
			stream <- data
		}

	}()

	return stream

}

func CreateAsyncStreamWithReadCloser(rList []io.ReadCloser) Stream {

	stream := make(Stream)

	go func() {

		for _, r := range rList {
			stream <- NewDefaultStreamData(r)
		}

	}()

	return stream

}
