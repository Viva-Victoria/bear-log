package log

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorMarshaling struct{}

func (e errorMarshaling) MarshalJSON() ([]byte, error) {
	return nil, ErrTest
}

func TestLogField_Key(t *testing.T) {
	field := FieldImpl{
		key: randomText(),
	}
	assert.Equal(t, field.key, field.Key())
}

func TestLogField_Type(t *testing.T) {
	field := FieldImpl{
		valueType: FieldType(randomInt(128, 256)),
	}
	assert.Equal(t, field.valueType, field.Type())
}

func TestLogField_StringValue(t *testing.T) {
	field := FieldImpl{
		str: randomText(),
	}
	assert.Equal(t, field.str, field.StringValue())
}

func TestLogField_Int(t *testing.T) {
	field := FieldImpl{
		i64: int64(randomInt(-1_000_000, 1_000_000)),
	}
	assert.Equal(t, field.i64, field.Int())
}

func TestLogField_Uint(t *testing.T) {
	field := FieldImpl{
		ui64: uint64(randomUInt()),
	}
	assert.Equal(t, field.ui64, field.UInt())
}

func TestLogField_Float(t *testing.T) {
	field := FieldImpl{
		f64: randomFloat(),
	}
	assert.Equal(t, field.f64, field.Float())
}

func TestLogField_Addressable(t *testing.T) {
	field := FieldImpl{
		addr: make([]string, randomInt(11, 12)),
	}
	assert.Equal(t, field.addr, field.Addressable())
}

func TestLogField_Errors(t *testing.T) {
	field := FieldImpl{
		valueType: TypeBinary,
		addr:      make([]int, 10),
	}
	_, err := field.String()
	assert.Error(t, err)

	field = FieldImpl{
		valueType: TypeArray,
		addr:      errorMarshaling{},
	}
	_, err = field.String()
	assert.Error(t, err)
}

func TestLogField_String(t *testing.T) {
	intNeg := randomInt(-100_000, -1)
	intPos := randomInt(1, 100_000)
	uint64 := uint64(randomUInt())
	float := randomFloat()
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
		Int(randomText(), randomInt(-100_000, 100_000)),
		Uint(randomText(), randomUInt()),
		Float(randomText(), randomFloat()),
		Binary(randomText(), []byte(randomText())),
		String(randomText(), randomText()),
		Array(randomText(), []int{
			randomInt(-1000, 1000),
			randomInt(-1000, 1000),
			randomInt(-1000, 1000),
		}),
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
			value := randomInt(0, 100_000) * multi[i]

			field := Int(key, value)
			assert.Equal(t, TypeInt, field.Type())
			assert.Equal(t, key, field.Key())
			assert.Equal(t, int64(value), field.Int())
		})
	}
}

func TestUint(t *testing.T) {
	key := randomText()
	value := randomUInt()

	field := Uint(key, value)
	assert.Equal(t, TypeUInt, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, uint64(value), field.UInt())
}

func TestFloat(t *testing.T) {
	key := randomText()
	value := randomFloat()

	field := Float(key, value)
	assert.Equal(t, TypeFloat, field.Type())
	assert.Equal(t, key, field.Key())
	assert.Equal(t, value, field.Float())
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
		value := []int{
			randomInt(-1000, 1000),
			randomInt(-1000, 1000),
			randomInt(-1000, 1000),
		}
		assertArray(t, key, value, Array(key, value))
	})
	t.Run("string", func(t *testing.T) {
		key := randomText()
		value := []string{randomText(), randomText(), randomText()}
		assertArray(t, key, value, Array(key, value))
	})
	t.Run("any", func(t *testing.T) {
		key := randomText()
		value := []any{
			randomInt(-1000, 1000),
			randomUInt(),
			randomFloat(),
			randomText(),
			nil,
		}
		assertArray(t, key, value, Array(key, value))
	})
}

func TestMap(t *testing.T) {
	t.Run("string:int", func(t *testing.T) {
		key := randomText()
		value := map[string]int{
			randomText(): randomInt(-1000, 1000),
			randomText(): randomInt(-1000, 1000),
			randomText(): randomInt(-1000, 1000),
		}

		assertMap(t, key, value, Map(key, value))
	})
	t.Run("int:string", func(t *testing.T) {
		key := randomText()
		value := map[int]string{
			randomInt(-1000, 1000): randomText(),
			randomInt(-1000, 1000): randomText(),
			randomInt(-1000, 1000): randomText(),
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
