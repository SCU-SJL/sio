package stream

import (
	"fmt"
)

type StreamDataConsumer interface {
	Consume(StreamData) error
	Finalize()
}

type SafeStreamHandler struct {
	stream      Stream
	inputErrCh  <-chan error
	outputErrCh chan error
	consumer    StreamDataConsumer
}

func NewSafeStreamHandler(stream Stream, errCh <-chan error, c StreamDataConsumer) *SafeStreamHandler {

	sr := &SafeStreamHandler{
		stream:     stream,
		inputErrCh: errCh,
		consumer:   c,
	}

	if sr.inputErrCh == nil {
		ch := make(chan error)
		close(ch)
		sr.inputErrCh = ch
	}

	sr.outputErrCh = make(chan error, 1)

	return sr

}

func (sr *SafeStreamHandler) StartHandling() <-chan error {

	if sr.inputErrCh == nil || sr.outputErrCh == nil || sr.stream == nil || sr.consumer == nil {
		panic("stream handler is not initialized")
	}

	go func() {

		defer func() {

			if r := recover(); r != nil {
				sr.outputErrCh <- fmt.Errorf("stream handler panicked: %v", r)
			}

			if err := sr.safeFinalize(); err != nil {
				sr.outputErrCh <- err
			}

		}()

		for {

			streamData, ok := <-sr.stream

			if !ok {
				break
			}

			if err := sr.consumer.Consume(streamData); err != nil {
				sr.outputErrCh <- err
				break
			}

		}

		if err, ok := <-sr.inputErrCh; ok && err != nil {
			sr.outputErrCh <- err
		}

	}()

	return sr.outputErrCh

}

func (sr *SafeStreamHandler) safeFinalize() error {

	defer func() {
		if r := recover(); r != nil {
			sr.outputErrCh <- fmt.Errorf("stream handler finalize failed: %v", r)
		}
	}()

	close(sr.outputErrCh)

	sr.consumer.Finalize()

}
