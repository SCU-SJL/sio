package multipart

import (
	"io"
	"mime/multipart"

	"github.com/SCU-SJL/sio/stream"
)

type Multiparter struct {
	handler *stream.SafeStreamHandler
	pr      *io.PipeReader
	c       *streamDataConsumerForMultipart
}

func NewMultiparter(inputStream stream.Stream, inputErrCh <-chan error) *Multiparter {

	multiparter := &Multiparter{}

	pr, pw := io.Pipe()

	multiparter.pr = pr
	multiparter.c = newStreamDataConsumerForMultipart(pw)

	multiparter.handler = stream.NewSafeStreamHandler(
		inputStream,
		inputErrCh,
		multiparter.c,
	)

	return multiparter

}

func (m *Multiparter) WithBuf(buf []byte) *Multiparter {

	m.c.buf = buf

	return m

}

func (m *Multiparter) StartMultiparting() (io.ReadCloser, <-chan error) {
	return m.pr, m.handler.StartHandling()
}

type streamDataConsumerForMultipart struct {
	mpw *multipart.Writer
	pw  *io.PipeWriter
	buf []byte
}

func newStreamDataConsumerForMultipart(pw *io.PipeWriter) *streamDataConsumerForMultipart {

	return &streamDataConsumerForMultipart{
		pw:  pw,
		mpw: multipart.NewWriter(pw),
	}

}

func (c *streamDataConsumerForMultipart) Consume(streamData stream.StreamData) error {

	multipartStreamData, err := stream.ToMultipartStreamData(streamData)
	if err != nil {
		return err
	}

	streamR := multipartStreamData.GetReadCloser()
	defer streamR.Close()
	fieldname := multipartStreamData.GetFieldName()
	filename := multipartStreamData.GetFileName()

	var part io.Writer

	if part, err = c.mpw.CreateFormFile(fieldname, filename); err != nil {
		return err
	}

	if len(c.buf) == 0 {
		_, err = io.Copy(part, streamR)
	} else {
		_, err = io.CopyBuffer(part, streamR, c.buf)
	}

	return err

}

func (c *streamDataConsumerForMultipart) Finalize() {

	defer c.pw.Close()

	defer c.mpw.Close()

}
