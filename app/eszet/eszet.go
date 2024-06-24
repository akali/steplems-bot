package eszet

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"sync"
)

const (
	escape = '/'
)

type Eszet struct {
	file    io.ReadWriteCloser
	content map[string][]byte
	mu      sync.Mutex
}

func New(file io.ReadWriteCloser) *Eszet {
	return &Eszet{
		file: file,
	}
}

func readValue(r *bufio.Reader) ([]byte, error) {
	b, err := r.ReadBytes(escape)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read count: %w", err)
	}
	cnt, err := strconv.ParseInt(string(b[:len(b)-1]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse count: %w", err)
	}
	keyBuf := make([]byte, cnt)
	keyBufCnt, err := io.ReadFull(r, keyBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	if keyBufCnt != int(cnt) {
		return nil, fmt.Errorf("failed to read: expected %d bytes, got %d", cnt, keyBufCnt)
	}
	return keyBuf, nil
}

func (e *Eszet) Init() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	r := bufio.NewReader(e.file)
	e.content = make(map[string][]byte)
	for {
		keyBuf, err := readValue(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read key: %w", err)
		}
		value, err := readValue(r)
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}

		key := string(keyBuf)

		e.content[key] = value
	}

	return nil
}

func (e *Eszet) Close() error {
	return e.file.Close()
}

func (e *Eszet) Get(key string) ([]byte, bool) {
	val, ok := e.content[key]
	return val, ok
}

func (e *Eszet) Write(key string, val []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	w := bufio.NewWriter(e.file)
	fmt.Fprintf(w, "%d%c%s%d%c%s", len(key), escape, key, len(val), escape, val)
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	e.content[key] = val
	return nil
}
