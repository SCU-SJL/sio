package stream

import (
	"fmt"
)

type StreamDataConsumer interface {
	Consume(StreamData) error
	Finalize()
}

type SafeStreamConsumer struct {
	stream      *Stream
	inputErrCh  <-chan error
	outputErrCh chan error
	consumer    StreamDataConsumer
}

func NewSafeStreamConsumer(stream *Stream, errCh <-chan error, c StreamDataConsumer) *SafeStreamConsumer {

	sr := &SafeStreamConsumer{
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

func (sr *SafeStreamConsumer) StartHandling() <-chan error {

	if sr.inputErrCh == nil || sr.outputErrCh == nil || sr.stream == nil || sr.consumer == nil {
		panic("stream handler is not initialized")
	}

	go func() {

		defer func() {

			if r := recover(); r != nil {
				sr.outputErrCh <- fmt.Errorf("stream handler panicked: %v", r)
			}

			sr.safeFinalize()

		}()

		for {

			streamData, ok := sr.stream.GetData()

			if !ok {
				break
			}

			if err := sr.consumer.Consume(streamData); err != nil {
				sr.outputErrCh <- err
				break
			}

		}

		if err, ok := ReadErr(sr.inputErrCh); ok && err != nil {
			sr.outputErrCh <- err
		}

	}()

	return sr.outputErrCh

}

func (sr *SafeStreamConsumer) safeFinalize() {

	defer func() {
		if r := recover(); r != nil {
			sr.outputErrCh <- fmt.Errorf("stream handler finalize failed: %v", r)
		}
	}()

	close(sr.outputErrCh)

	sr.consumer.Finalize()

}
