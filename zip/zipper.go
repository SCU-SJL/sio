package zip

import (
	"archive/zip"
	"io"

	"github.com/SCU-SJL/sio/stream"
)

type Zipper struct {
	handler *stream.SafeStreamHandler
	pr      *io.PipeReader
	c       *streamDataConsumer
}

func NewZipper(inputStream stream.Stream, inputErrCh <-chan error) *Zipper {

	zipper := &Zipper{}

	pr, pw := io.Pipe()

	zipper.pr = pr
	zipper.c = newStreamDataConsumer(pw)

	zipper.handler = stream.NewSafeStreamHandler(
		inputStream,
		inputErrCh,
		zipper.c,
	)

	return zipper

}

func (z *Zipper) WithBuf(buf []byte) *Zipper {

	z.c.buf = buf

	return z

}

func (z *Zipper) StartZipping() (io.ReadCloser, <-chan error) {
	return z.pr, z.handler.StartHandling()
}

type streamDataConsumer struct {
	zw  *zip.Writer
	pw  *io.PipeWriter
	buf []byte
}

func newStreamDataConsumer(pw *io.PipeWriter) *streamDataConsumer {

	return &streamDataConsumer{
		pw: pw,
		zw: zip.NewWriter(pw),
	}

}

func (c *streamDataConsumer) Consume(streamData stream.StreamData) error {

	zipStreamData, err := stream.ToZipStreamData(streamData)
	if err != nil {
		return err
	}

	streamR := zipStreamData.GetReadCloser()
	defer streamR.Close()
	filename := zipStreamData.GetFileName()

	if streamR == nil || len(filename) == 0 {
		return nil
	}

	fw, err := c.zw.Create(filename)
	if err != nil {
		return err
	}

	if len(c.buf) == 0 {
		_, err = io.Copy(fw, streamR)
	} else {
		_, err = io.CopyBuffer(fw, streamR, c.buf)
	}

	return err

}

func (c *streamDataConsumer) Finalize() {

	defer c.pw.Close()

	defer c.zw.Close()

}
