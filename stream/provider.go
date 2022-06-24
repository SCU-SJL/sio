package stream

import (
	"fmt"
)

type StreamDataProvider interface {
	Next() (StreamData, error)
}

type SafeStreamWriter struct {
	stream       Stream
	errCh        chan error
	dataProvider StreamDataProvider
}

func NewSafeStreamWriter(dataProvider StreamDataProvider) *SafeStreamWriter {
	return &SafeStreamWriter{
		stream:       make(Stream),
		errCh:        make(chan error, 1),
		dataProvider: dataProvider,
	}
}

func (w *SafeStreamWriter) StartWriting() (Stream, <-chan error) {

	go func() {

		defer func() {

			if r := recover(); r != nil {
				err := fmt.Errorf("stream writer panicked: %v", r)
				w.errCh <- err
			}

			close(w.stream)
			close(w.errCh)

		}()

		for {

			// read
			streamData, err := w.dataProvider.Next()

			if err != nil {
				w.errCh <- err
				break
			}

			if streamData == nil {
				break
			}

			// block write
			w.stream <- streamData

		}

	}()

	return w.stream, w.errCh

}
