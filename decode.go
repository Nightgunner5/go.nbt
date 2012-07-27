package nbt

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Unmarshal(compression Compression, in io.Reader, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	new(decodeState).init(compression, in).unmarshal(v)
	return
}

type decodeState struct {
	in    io.Reader
	stack []Tag
}

func (d *decodeState) init(compression Compression, in io.Reader) *decodeState {
	if in == nil {
		panic(fmt.Errorf("nbt: Input stream is nil"))
	}

	switch compression {
	case Uncompressed:
		d.in = in
	case GZip:
		r, err := gzip.NewReader(in)
		if err != nil {
			panic(err)
		}
		d.in = r
	case ZLib:
		r, err := zlib.NewReader(in)
		if err != nil {
			panic(err)
		}
		d.in = r
	default:
		panic(fmt.Errorf("nbt: Unknown compression type: %d", compression))
	}

	d.stack = make([]Tag, 0, 32)

	return d
}

func (d *decodeState) unmarshal(v interface{}) {
	d.readTag(reflect.ValueOf(v).Elem())
}

func (d *decodeState) r(i interface{}) {
	err := binary.Read(d.in, binary.BigEndian, i)
	if err != nil {
		panic(err)
	}
}

// Returns true if TAG_End was NOT the tag read.
func (d *decodeState) readTag(v reflect.Value) bool {
	var tag Tag
	d.r(&tag)

	switch tag {
	case TAG_End:
		return false
	case TAG_Byte:
		var value uint8
		d.r(&value)
		switch v.Kind() {
		case reflect.Bool:
			v.SetBool(value != 0)
		case reflect.Int8:
			v.SetInt(int64(int8(value)))
		case reflect.Uint8:
			v.SetUint(uint64(value))
		default:
			panic(fmt.Errorf("nbt: Tag is %s, but I don't know how to put that in a %s!", tag, v.Kind()))
		}
	default:
		panic(fmt.Errorf("nbt: Unhandled tag: %s", tag))
	}
	return true
}
