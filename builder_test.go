package bear_log

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
)

func TestStringBuilder_Append(t *testing.T) {
	builder := NewStringBuilder()

	values := []string{randomText(), randomText(), randomText()}
	builder.Append(values...)

	assert.Equal(t, strings.Join(values, ""), builder.buffer.String())
}

func TestStringBuilder_AppendBytes(t *testing.T) {
	builder := NewStringBuilder()

	values := []string{randomText(), randomText(), randomText()}
	builder.AppendBytes([]byte(values[0]), []byte(values[1]), []byte(values[2]))

	assert.Equal(t, strings.Join(values, ""), builder.buffer.String())
}

func TestStringBuilder_String(t *testing.T) {
	builder := NewStringBuilder()

	count := rand.Intn(10)
	values := make([]string, 0, count)
	for i := 0; i < count; i++ {
		values = append(values, randomText())
		builder.Append(values[i])
	}

	assert.Equal(t, strings.Join(values, ""), builder.String())
}

func TestStringBuilder_Bytes(t *testing.T) {
	builder := NewStringBuilder()

	count := rand.Intn(10)
	values := make([]string, 0, count)
	for i := 0; i < count; i++ {
		values = append(values, randomText())
		builder.Append(values[i])
	}

	assert.Equal(t, []byte(strings.Join(values, "")), builder.Bytes())
}
