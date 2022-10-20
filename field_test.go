package bear_log

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
)

type errorMarshaling struct{}

func (e errorMarshaling) MarshalJSON() ([]byte, error) {
	return nil, ErrTest
}

func TestLogField_Key(t *testing.T) {
	field := LogField{
		key: randomText(),
	}
	assert.Equal(t, field.key, field.Key())
}

func TestLogField_Type(t *testing.T) {
	field := LogField{
		valueType: FieldType(rand.Intn(256)),
	}
	assert.Equal(t, field.valueType, field.Type())
}

func TestLogField_StringValue(t *testing.T) {
	field := LogField{
		str: randomText(),
	}
	assert.Equal(t, field.str, field.StringValue())
}

func TestLogField_Int(t *testing.T) {
	field := LogField{
		i64: rand.Int63(),
	}
	assert.Equal(t, field.i64, field.Int())
}

func TestLogField_Uint(t *testing.T) {
	field := LogField{
		ui64: rand.Uint64(),
	}
	assert.Equal(t, field.ui64, field.UInt())
}

func TestLogField_Float(t *testing.T) {
	field := LogField{
		f64: rand.Float64(),
	}
	assert.Equal(t, field.f64, field.Float())
}

func TestLogField_Addressable(t *testing.T) {
	field := LogField{
		addr: make([]string, rand.Intn(12)),
	}
	assert.Equal(t, field.addr, field.Addressable())
}

func TestLogField_Errors(t *testing.T) {
	field := LogField{
		valueType: TypeBinary,
		addr:      make([]int, 10),
	}
	_, err := field.String()
	assert.Error(t, err)

	field = LogField{
		valueType: TypeArray,
		addr:      errorMarshaling{},
	}
	_, err = field.String()
	assert.Error(t, err)
}

func TestLogField_String(t *testing.T) {
	intNeg := rand.Int() * -1
	intPos := rand.Int()
	uint64 := rand.Uint64()
	float := rand.Float64()
	str := randomText()
	keys := []string{randomText(), randomText(), randomText()}

	fields := []Field{
		Int(randomText(), intNeg),
		Int(randomText(), intPos),
		Uint(randomText(), uint64),
		Float(randomText(), float),
		String(randomText(), str),
		Binary(randomText(), []byte(str)),
		Array(randomText(), []any{intNeg, intPos, float, str}),
		Map(randomText(), map[string]any{
			keys[0]: intNeg,
			keys[1]: float,
			keys[2]: str,
		}),
	}
	values := []string{
		strconv.Itoa(intNeg),
		strconv.Itoa(intPos),
		strconv.FormatUint(uint64, 10),
		strconv.FormatFloat(float, 'f', -1, 64),
		str,
		base64.StdEncoding.EncodeToString([]byte(str)),
		fmt.Sprintf(`[%d,%d,%s,"%s"]`, intNeg, intPos, strconv.FormatFloat(float, 'f', -1, 64), str),
		fmt.Sprintf(`{"%s":%d,"%s":%s,"%s":"%s"}`, keys[0], intNeg, keys[1], strconv.FormatFloat(float, 'f', -1, 64), keys[2], str),
	}

	for i := range fields {
		s, err := fields[i].String()
		require.NoError(t, err)

		if fields[i].Type() == TypeMap {
			assert.JSONEq(t, values[i], s)
		} else {
			assert.Equal(t, values[i], s)
		}
	}
}

func TestLogField_Value(t *testing.T) {
	fields := []Field{
		Int(randomText(), rand.Int()),
		Uint(randomText(), rand.Uint64()),
		Float(randomText(), rand.Float64()),
		Binary(randomText(), []byte(randomText())),
		String(randomText(), randomText()),
		Array(randomText(), []int{rand.Int(), rand.Int(), rand.Int()}),
	}
	values := []any{
		int64(0),
		uint64(0),
		float64(0),
		"",
		"",
		[]int{},
	}

	for i := range fields {
		v, err := fields[i].Value()
		require.NoError(t, err)

		assert.IsType(t, values[i], v)
	}
}

func TestInt(t *testing.T) {
	names := []string{"negative", "zero", "positive"}
	multi := []int{-1, 0, 1}

	for i := range names {
		t.Run(names[i], func(t *testing.T) {
			key := randomText()
			value := rand.Int() * multi[i]

			field := Int(key, value)
			assert.Equal(t, TypeInt, field.Type())
			assert.Equal(t, key, field.Key())
			assert.Equal(t, int64(value), field.Int())
		})
	}
}

func TestUint(t *testing.T) {
	key := randomText()
	value := rand.Uint32()

	field := Uint(key, value)
	assert.Equal(t, TypeUInt, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, uint64(value), field.UInt())
}

func TestFloat(t *testing.T) {
	key := randomText()
	value := rand.Float32()

	field := Float(key, value)
	assert.Equal(t, TypeFloat, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, float64(value), field.Float())
}

func TestString(t *testing.T) {
	names := []string{"empty", "value"}
	values := []string{"", randomText()}

	for i := range names {
		t.Run(names[i], func(t *testing.T) {
			key := randomText()

			field := String(key, values[i])
			assert.Equal(t, TypeString, field.Type())
			assert.Equal(t, key, field.Key())
			assert.Equal(t, values[i], field.StringValue())
		})
	}
}

func TestBinary(t *testing.T) {
	key := randomText()
	value := []byte(randomText())

	field := Binary(key, value)
	assert.Equal(t, TypeBinary, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, value, field.Addressable())
}

func TestArray(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		key := randomText()
		value := []int{rand.Int(), rand.Int(), rand.Int()}
		assertArray(t, key, value, Array(key, value))
	})
	t.Run("string", func(t *testing.T) {
		key := randomText()
		value := []string{randomText(), randomText(), randomText()}
		assertArray(t, key, value, Array(key, value))
	})
	t.Run("any", func(t *testing.T) {
		key := randomText()
		value := []any{rand.Int(), rand.Uint32(), rand.Float64(), randomText(), nil}
		assertArray(t, key, value, Array(key, value))
	})
}

func TestMap(t *testing.T) {
	t.Run("string:int", func(t *testing.T) {
		key := randomText()
		value := map[string]int{
			randomText(): rand.Int(),
			randomText(): rand.Int(),
			randomText(): rand.Int(),
		}

		assertMap(t, key, value, Map(key, value))
	})
	t.Run("int:string", func(t *testing.T) {
		key := randomText()
		value := map[int]string{
			rand.Int(): randomText(),
			rand.Int(): randomText(),
			rand.Int(): randomText(),
		}

		assertMap(t, key, value, Map(key, value))
	})
}

func TestObject(t *testing.T) {
	key := randomText()
	value := struct {
		id string
	}{randomText()}

	field := Object(key, value)
	assert.Equal(t, TypeAny, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, value, field.Addressable())
}

func assertArray[T any](t *testing.T, key string, expected []T, field Field) {
	assert.Equal(t, TypeArray, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, expected, field.Addressable())
}

func assertMap[K comparable, V any](t *testing.T, key string, expected map[K]V, field Field) {
	assert.Equal(t, TypeMap, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, expected, field.Addressable())
}
