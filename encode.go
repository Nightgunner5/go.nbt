package nbt

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Marshal(compression Compression, out io.Writer, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok {
				err = fmt.Errorf(s)
			} else {
				err = r.(error)
			}
		}
	}()

	if out == nil {
		panic(fmt.Errorf("nbt: Output stream is nil"))
	}

	switch compression {
	case Uncompressed:
		break
	case GZip:
		w := gzip.NewWriter(out)
		defer w.Close()
		out = w
	case ZLib:
		w := zlib.NewWriter(out)
		defer w.Close()
		out = w
	default:
		panic(fmt.Errorf("nbt: Unknown compression type: %d", compression))
	}

	writeRootTag(out, reflect.ValueOf(v))

	return
}

func writeRootTag(out io.Writer, v reflect.Value) {
	writeTag(out, "", v)
}

func w(out io.Writer, v interface{}) {
	err := binary.Write(out, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
}

func writeTag(out io.Writer, name string, v reflect.Value) {
	v = reflect.Indirect(v)
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("%v\n\t\tat struct field %#v", r, name))
		}
	}()
	switch v.Kind() {
	case reflect.Bool:
		w(out, TAG_Byte)
		writeValue(out, TAG_String, name)
		if v.Bool() {
			writeValue(out, TAG_Byte, byte(1))
		} else {
			writeValue(out, TAG_Byte, byte(0))
		}

	case reflect.Int8:
		w(out, TAG_Byte)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Byte, int8(v.Int()))

	case reflect.Uint8:
		w(out, TAG_Byte)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Byte, uint8(v.Uint()))

	case reflect.Int16:
		w(out, TAG_Short)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Short, int16(v.Int()))

	case reflect.Uint16:
		w(out, TAG_Short)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Short, uint16(v.Uint()))

	case reflect.Int32:
		w(out, TAG_Int)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Int, int32(v.Int()))

	case reflect.Uint32:
		w(out, TAG_Int)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Int, uint32(v.Uint()))

	case reflect.Int64:
		w(out, TAG_Long)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Long, v.Int())

	case reflect.Uint64:
		w(out, TAG_Long)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Long, v.Uint())

	case reflect.Float32:
		w(out, TAG_Float)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Float, float32(v.Float()))

	case reflect.Float64:
		w(out, TAG_Double)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_Double, v.Float())

	case reflect.String:
		w(out, TAG_String)
		writeValue(out, TAG_String, name)
		writeValue(out, TAG_String, v.String())

	case reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			w(out, TAG_Byte_Array)
			writeValue(out, TAG_String, name)
			writeValue(out, TAG_Byte_Array, v.Slice(0, v.Len()).Bytes())

		case reflect.Int32, reflect.Uint32:
			w(out, TAG_Int_Array)
			writeValue(out, TAG_String, name)
			for i := 0; i < v.Len(); i++ {
				writeValue(out, TAG_Int, v.Index(i).Interface())
			}

		default:
			panic(fmt.Errorf("nbt: Unhandled array type: %v", v.Type().Elem()))
		}

	case reflect.Slice:
		w(out, TAG_List)
		writeValue(out, TAG_String, name)
		writeList(out, v)

	case reflect.Map:
		w(out, TAG_Compound)
		writeValue(out, TAG_String, name)
		writeMap(out, v)

	case reflect.Struct:
		w(out, TAG_Compound)
		writeValue(out, TAG_String, name)
		writeCompound(out, v)

	default:
		panic(fmt.Errorf("nbt: Unhandled type: %v (%v)", v.Type(), v.Interface()))
	}
}

func writeValue(out io.Writer, tag Tag, v interface{}) {
	switch tag {
	case TAG_Byte, TAG_Short, TAG_Int, TAG_Long, TAG_Float, TAG_Double:
		w(out, v)

	case TAG_String:
		w(out, uint16(len(v.(string))))
		_, err := out.Write([]byte(v.(string)))
		if err != nil {
			panic(err)
		}

	case TAG_Byte_Array:
		w(out, uint32(len(v.([]byte))))
		_, err := out.Write(v.([]byte))
		if err != nil {
			panic(err)
		}

	default:
		panic(fmt.Errorf("nbt: Unhandled tag: %s (%v)", tag, v))
	}
}

func writeList(out io.Writer, v reflect.Value) {
	v = reflect.Indirect(v)
	var tag Tag
	mustConvertBool := false
	mustConvertMap := false
	switch v.Type().Elem().Kind() {
	case reflect.Bool:
		mustConvertBool = true
		fallthrough
	case reflect.Int8, reflect.Uint8:
		tag = TAG_Byte

	case reflect.Int16, reflect.Uint16:
		tag = TAG_Short

	case reflect.Int32, reflect.Uint32:
		tag = TAG_Int

	case reflect.Int64, reflect.Uint64:
		tag = TAG_Long

	case reflect.Float32:
		tag = TAG_Float

	case reflect.Float64:
		tag = TAG_Double

	case reflect.String:
		tag = TAG_String

	case reflect.Array:
		switch v.Type().Elem().Elem().Kind() {
		case reflect.Uint8:
			tag = TAG_Byte_Array

		case reflect.Int32, reflect.Uint32:
			tag = TAG_Int_Array

		default:
			panic(fmt.Errorf("nbt: Unhandled array type: %v", v.Type().Elem().Elem()))
		}

	case reflect.Slice:
		tag = TAG_List

	case reflect.Map:
		mustConvertMap = true
		fallthrough
	case reflect.Struct:
		tag = TAG_Compound

	default:
		panic(fmt.Errorf("nbt: Unhandled list element type: %v", v.Type().Elem()))
	}
	w(out, tag)
	w(out, uint32(v.Len()))

	var i int
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("%v\n\t\tat list index %d", r, i))
		}
	}()
	for i = 0; i < v.Len(); i++ {
		if mustConvertBool {
			if v.Index(i).Bool() {
				writeValue(out, TAG_Byte, uint8(1))
			} else {
				writeValue(out, TAG_Byte, uint8(0))
			}
		} else if tag == TAG_Compound {
			if mustConvertMap {
				writeMap(out, v.Index(i))
			} else {
				writeCompound(out, v.Index(i))
			}
		} else if tag == TAG_List {
			writeList(out, v.Index(i))
		} else if tag == TAG_Byte_Array {
			writeValue(out, tag, v.Index(i).Bytes())
		} else if tag == TAG_Int_Array {
			for j := 0; j < v.Index(i).Len(); j++ {
				writeValue(out, TAG_Int, v.Index(i).Index(j).Interface())
			}
		} else {
			writeValue(out, tag, v.Index(i).Interface())
		}
	}
}

func writeMap(out io.Writer, v reflect.Value) {
	for _, name := range v.MapKeys() {
		writeTag(out, name.String(), reflect.Indirect(v.MapIndex(name)))
	}
	w(out, TAG_End)
}

func writeCompound(out io.Writer, v reflect.Value) {
	v = reflect.Indirect(v)
	fields := parseStruct(v)

	for name, value := range fields {
		writeTag(out, name, value)
	}
	w(out, TAG_End)
}
