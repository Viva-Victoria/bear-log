package log

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
)

type FieldType uint8

const (
	TypeAny FieldType = iota
	TypeInt
	TypeUInt
	TypeFloat
	TypeBinary
	TypeString
	TypeArray
	TypeMap
)

type Field interface {
	Type() FieldType
	StringValue() string
	Int() int64
	UInt() uint64
	Float() float64
	Addressable() any

	Key() string
	Value() (any, error)
	String() (string, error)
}

type FieldImpl struct {
	addr      any
	key       string
	str       string
	i64       int64
	ui64      uint64
	f64       float64
	valueType FieldType
}

func (l FieldImpl) Key() string {
	return l.key
}

func (l FieldImpl) Type() FieldType {
	return l.valueType
}

func (l FieldImpl) StringValue() string {
	return l.str
}

func (l FieldImpl) Int() int64 {
	return l.i64
}

func (l FieldImpl) UInt() uint64 {
	return l.ui64
}

func (l FieldImpl) Float() float64 {
	return l.f64
}

func (l FieldImpl) Addressable() any {
	return l.addr
}

func (l FieldImpl) String() (string, error) {
	var value string
	switch l.Type() {
	case TypeInt:
		value = strconv.FormatInt(l.Int(), 10)

	case TypeUInt:
		value = strconv.FormatUint(l.UInt(), 10)

	case TypeFloat:
		value = strconv.FormatFloat(l.Float(), 'f', -1, 64)

	case TypeBinary:
		array, ok := l.Addressable().([]byte)
		if !ok {
			return "", fmt.Errorf("%v is not []byte", l.Addressable())
		}

		value = base64.StdEncoding.EncodeToString(array)

	case TypeString:
		value = l.StringValue()

	case TypeArray, TypeMap, TypeAny:
		data, err := json.Marshal(l.Addressable())
		if err != nil {
			return "", err
		}

		value = string(data)
	}

	return value, nil
}

func (l FieldImpl) Value() (any, error) {
	switch l.Type() {
	case TypeInt:
		return l.Int(), nil

	case TypeUInt:
		return l.UInt(), nil

	case TypeFloat:
		return l.Float(), nil

	case TypeBinary:
		array, ok := l.Addressable().([]byte)
		if !ok {
			return nil, fmt.Errorf("%v is not []byte", l.Addressable())
		}

		return base64.StdEncoding.EncodeToString(array), nil

	case TypeString:
		return l.StringValue(), nil

	case TypeArray, TypeMap, TypeAny:
		return l.Addressable(), nil
	}

	return nil, nil
}

func Int[I int | int8 | int16 | int32 | int64](k string, v I) Field {
	return FieldImpl{
		valueType: TypeInt,
		key:       k,
		i64:       int64(v),
	}
}

func Uint[U uint | uint8 | uint16 | uint32 | uint64](k string, v U) Field {
	return FieldImpl{
		valueType: TypeUInt,
		key:       k,
		ui64:      uint64(v),
	}
}

func Float[F float32 | float64](k string, f F) Field {
	return FieldImpl{
		valueType: TypeFloat,
		key:       k,
		f64:       float64(f),
	}
}

func String(k, s string) Field {
	return FieldImpl{
		valueType: TypeString,
		key:       k,
		str:       s,
	}
}

func Binary(k string, v []byte) Field {
	return FieldImpl{
		valueType: TypeBinary,
		key:       k,
		addr:      v,
	}
}

func Array[T any](k string, v []T) Field {
	return FieldImpl{
		valueType: TypeArray,
		key:       k,
		addr:      v,
	}
}

func Map[K comparable, V any](k string, v map[K]V) Field {
	return FieldImpl{
		valueType: TypeMap,
		key:       k,
		addr:      v,
	}
}

func Object(k string, v any) Field {
	return FieldImpl{
		valueType: TypeAny,
		key:       k,
		addr:      v,
	}
}
