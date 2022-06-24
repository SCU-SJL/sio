package zip

import (
	"archive/zip"
	"io"

	"github.com/SCU-SJL/sio/stream"
)

type Zipper struct {
	handler *stream.SafeStreamHandler
	pr      *io.PipeReader
}

func NewZipper(inputStream stream.Stream, inputErrCh <-chan error) *Zipper {

	zipper := &Zipper{}

	pr, pw := io.Pipe()

	zipper.pr = pr

	zipper.handler = stream.NewSafeStreamHandler(
		inputStream,
		inputErrCh,
		newStreamDataConsumer(pw),
	)

	return zipper

}

func (z *Zipper) StartZipping() (io.ReadCloser, <-chan error) {
	return z.pr, z.handler.StartHandling()
}

type streamDataConsumer struct {
	zw *zip.Writer
	pw *io.PipeWriter
}

func newStreamDataConsumer(pw *io.PipeWriter) *streamDataConsumer {
	zw := zip.NewWriter(pw)
	return &streamDataConsumer{pw: pw, zw: zw}
}

func (h *streamDataConsumer) Consume(streamData stream.StreamData) error {

	zipStreamData, err := stream.ToZipStreamData(streamData)
	if err != nil {
		return err
	}

	streamR := zipStreamData.GetReadCloser()
	filename := zipStreamData.GetFileName()

	if streamR == nil || len(filename) == 0 {
		return nil
	}

	fw, err := h.zw.Create(filename)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, streamR); err != nil {
		return err
	}

	return streamR.Close()

}

func (h *streamDataConsumer) Finalize() {

	defer h.pw.Close()

	defer h.zw.Close()

}
