package stream

import "errors"

func ToZipStreamData(sd StreamData) (ZipStreamData, error) {
	zd, ok := sd.(ZipStreamData)
	if !ok {
		return nil, errors.New("StreamData cannot be converted to ZipStreamData")
	}
	return zd, nil
}
