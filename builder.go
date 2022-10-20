package log

import (
	"bytes"
	"sync"
)

var (
	_pool = sync.Pool{
		New: func() any {
			buf := bytes.NewBuffer(make([]byte, 4096))
			buf.Reset()
			return buf
		},
	}
)

type StringBuilder struct {
	buffer *bytes.Buffer
}

func NewStringBuilder() StringBuilder {
	return StringBuilder{
		buffer: _pool.Get().(*bytes.Buffer),
	}
}

func (s StringBuilder) Append(values ...string) StringBuilder {
	for _, value := range values {
		s.buffer.WriteString(value)
	}
	return s
}

func (s StringBuilder) AppendBytes(values ...[]byte) StringBuilder {
	for _, value := range values {
		s.buffer.Write(value)
	}
	return s
}

func (s StringBuilder) String() string {
	return s.buffer.String()
}

func (s StringBuilder) Bytes() []byte {
	return []byte(s.String())
}

func (s StringBuilder) Dispose() {
	_pool.Put(s.buffer)
}
