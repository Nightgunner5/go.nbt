package nbt

import "fmt"

// All tags are big endian.

type Tag byte

const (
	TAG_End        Tag = 0  // No payload, no name.
	TAG_Byte       Tag = 1  // Signed 8 bit integer.
	TAG_Short      Tag = 2  // Signed 16 bit integer.
	TAG_Int        Tag = 3  // Signed 32 bit integer.
	TAG_Long       Tag = 4  // Signed 64 bit integer.
	TAG_Float      Tag = 5  // IEEE 754-2008 32 bit floating point number.
	TAG_Double     Tag = 6  // IEEE 754-2008 64 bit floating point number.
	TAG_Byte_Array Tag = 7  // size TAG_Int, then payload [size]byte.
	TAG_String     Tag = 8  // length TAG_Short, then payload (utf-8) string (of length length).
	TAG_List       Tag = 9  // tagID TAG_Byte, length TAG_Int, then payload [length]tagID.
	TAG_Compound   Tag = 10 // { tagID TAG_Byte, name TAG_String, payload tagID }... TAG_End
	TAG_Int_Array  Tag = 11 // size TAG_Int, then payload [size]TAG_Int
	TAG_Long_Array Tag = 12
)

func (tag Tag) String() string {
	name := "Unknown"
	switch tag {
	case TAG_End:
		name = "TAG_End"
	case TAG_Byte:
		name = "TAG_Byte"
	case TAG_Short:
		name = "TAG_Short"
	case TAG_Int:
		name = "TAG_Int"
	case TAG_Long:
		name = "TAG_Long"
	case TAG_Float:
		name = "TAG_Float"
	case TAG_Double:
		name = "TAG_Double"
	case TAG_Byte_Array:
		name = "TAG_Byte_Array"
	case TAG_String:
		name = "TAG_String"
	case TAG_List:
		name = "TAG_List"
	case TAG_Compound:
		name = "TAG_Compound"
	case TAG_Int_Array:
		name = "TAG_Int_Array"
	case TAG_Long_Array:
		name = "TAG_Long_Array"
	}
	return fmt.Sprintf("%s (0x%02x)", name, byte(tag))
}

type Compression byte

const (
	Uncompressed Compression = 0
	GZip         Compression = 1
	ZLib         Compression = 2
)
