// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"sync"

	"github.com/ruixueyun/ruixuego/bufferpool"
)

var _gzip = &gzipPool{}

type gzipPool struct {
	readers sync.Pool
	writers sync.Pool
}

func (pool *gzipPool) GetReader(src io.Reader) (reader *gzip.Reader) {
	if r := pool.readers.Get(); r != nil {
		reader = r.(*gzip.Reader)
		reader.Reset(src)
	} else {
		reader, _ = gzip.NewReader(src)
	}
	return reader
}

func (pool *gzipPool) PutReader(reader *gzip.Reader) {
	reader.Close()
	pool.readers.Put(reader)
}

func (pool *gzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	if w := pool.writers.Get(); w != nil {
		writer = w.(*gzip.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = gzip.NewWriterLevel(dst, gzip.BestCompression)
	}
	return writer
}

func (pool *gzipPool) PutWriter(writer *gzip.Writer) {
	writer.Close()
	pool.writers.Put(writer)
}

func gZIPCompress(b []byte) (*bytes.Buffer, error) {
	buf := bufferpool.Get()
	gw := _gzip.GetWriter(buf)
	defer _gzip.PutWriter(gw)

	_, err := gw.Write(b)
	if err != nil {
		bufferpool.Put(buf)
		return nil, err
	}
	return buf, nil
}

func gZIPUncompress(b []byte) ([]byte, error) {
	buf := bufferpool.Get()
	buf.Write(b)
	gr := _gzip.GetReader(buf)
	defer _gzip.PutReader(gr)
	defer bufferpool.Put(buf)

	ret, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GzipCompressV2 对输入数据进行 gzip 压缩。
func GzipCompressV2(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GzipDecompressV2 对 gzip 压缩数据进行解压。
func GzipDecompressV2(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
