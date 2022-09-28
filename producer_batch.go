// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"net/http"
	"sync"
	"time"
)

func newBatchWriter(c *Client, conf *BigDataConfig) *batchWriter {
	return &batchWriter{
		conf:        conf,
		client:      c,
		bufferMutex: new(sync.RWMutex),
		buffer:      make([]*BigDataLog, 0, conf.BatchSize),
		cacheMutex:  new(sync.RWMutex),
		cache:       make([]*BigDataLog, 0, conf.BatchSize*2),
		closed:      make(chan struct{}, 1),
	}
}

type batchWriter struct {
	conf        *BigDataConfig
	client      *Client
	bufferMutex *sync.RWMutex
	buffer      []*BigDataLog
	cacheMutex  *sync.RWMutex
	cache       []*BigDataLog
	gzipPool    *gzipPool
	closed      chan struct{}
}

func (bw *batchWriter) Init() error {
	if !bw.conf.AutoFlush {
		return nil
	}
	go func() {
		ticker := time.NewTicker(bw.conf.AutoFlushInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := bw.Flush()
				if err != nil {
					logger.Errorf(err.Error())
				}
			case <-bw.closed:
				return
			}
		}
	}()
	return nil
}

func (bw *batchWriter) Write(logData *BigDataLog) error {
	bw.bufferMutex.Lock()
	bw.buffer = append(bw.buffer, logData)
	bw.bufferMutex.Unlock()

	if bw.bufferLength() >= bw.conf.BatchSize || bw.cacheLength() > 0 {
		return bw.Flush()
	}
	return nil
}

func (bw *batchWriter) Flush() error {
	bw.bufferMutex.Lock()
	bw.cacheMutex.Lock()
	defer bw.cacheMutex.Unlock()

	if len(bw.cache) == 0 && len(bw.buffer) == 0 {
		bw.bufferMutex.Unlock()
		return nil
	}

	defer func() {
		if len(bw.cache) > bw.conf.CacheCapacity {
			bw.cache = append(bw.cache[:0], bw.cache[len(bw.cache)-bw.conf.CacheCapacity:]...)
		}
	}()

	if len(bw.cache) == 0 || len(bw.buffer) >= bw.conf.BatchSize {
		bw.cache = append(bw.cache, bw.buffer...)
		bw.buffer = bw.buffer[:0]
	}
	bw.bufferMutex.Unlock()

	n := len(bw.cache)
	if n > bw.conf.BatchSize {
		n = bw.conf.BatchSize
	}
	b, err := MarshalJSON(bw.cache[:n])
	if err != nil {
		return err
	}

	code := 0
	for i := 0; i < 3; i++ {
		code, err = bw.client.track(b, n, !bw.conf.DisableCompress)
		if err != nil {
			logger.Errorf("failed to send track log: [%d] %s, data: %s", code, err.Error(), b)
			if code != http.StatusOK {
				continue
			}
			return err
		} else {
			bw.cache = append(bw.cache[:0], bw.cache[n:]...)
			break
		}
	}
	return err
}

func (bw *batchWriter) Close() error {
	close(bw.closed)
	err := bw.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (bw *batchWriter) bufferLength() int {
	bw.bufferMutex.RLock()
	defer bw.bufferMutex.RUnlock()
	return len(bw.buffer)
}

func (bw *batchWriter) cacheLength() int {
	bw.cacheMutex.RLock()
	defer bw.cacheMutex.RUnlock()
	return len(bw.cache)
}
