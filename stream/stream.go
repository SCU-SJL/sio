package stream

import (
	"io"
	"sync"
)

type Stream struct {
	dataCh  chan StreamData
	mu      *sync.RWMutex
	closeCh chan struct{}
}

func NewStream() *Stream {
	return &Stream{
		dataCh:  make(chan StreamData),
		mu:      &sync.RWMutex{},
		closeCh: make(chan struct{}),
	}
}

func (s *Stream) PutData(data StreamData) {
	s.dataCh <- data
}

func (s *Stream) GetData() (StreamData, bool) {
	data, ok := <-s.dataCh
	return data, ok
}

func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return
	}
	close(s.dataCh)
	close(s.closeCh)
}

func (s *Stream) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isClosed()
}

func (s *Stream) isClosed() bool {
	select {
	case <-s.closeCh:
		return true
	default:
		return false
	}
}

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

func CreateAsyncStream(dataList []StreamData) *Stream {

	stream := NewStream()

	go func() {

		for _, data := range dataList {
			stream.PutData(data)
		}

	}()

	return stream

}

func CreateAsyncStreamWithReadCloser(rList []io.ReadCloser) *Stream {

	stream := NewStream()

	go func() {

		for _, r := range rList {
			stream.PutData(NewDefaultStreamData(r))
		}

	}()

	return stream

}

func ReadErr(errCh <-chan error) (err error, ok bool) {
	select {
	case err, ok = <-errCh:
		return
	default:
		return
	}
}
