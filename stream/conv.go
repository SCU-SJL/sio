package stream

import "errors"

func ToZipStreamData(sd StreamData) (ZipStreamData, error) {
	zd, ok := sd.(ZipStreamData)
	if !ok {
		return nil, errors.New("StreamData cannot be converted to ZipStreamData")
	}
	return zd, nil
}

func ToMultipartStreamData(sd StreamData) (MultipartStreamData, error) {
	hd, ok := sd.(MultipartStreamData)
	if !ok {
		return nil, errors.New("StreamData cannot be converted to MultipartStreamData")
	}
	return hd, nil
}
