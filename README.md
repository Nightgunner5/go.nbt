All the other Go NBT libraries I tried fell short of my expectations. Either they had their own type system
which I had to convert into maps or parse a second time into my structs or they had cryptic error messages
or custom "error handling" that omitted important details like where exactly the problem occurred.

I built go.nbt to solve two problems: NBT parsing (NBT is Mojang's "Named Binary Tag" format used heavily
in Minecraft) and meaningful error handling.

But that's enough of the trivial details. You're here to use the library, not hear a story about how it was
made.

Setup
=====

To set up go.nbt, simply run:

    go get github.com/Nightgunner5/go.nbt

Or use `go get` with no parameters in the package you use it in and it will be automagically downloaded and
compiled.

But you already knew that.

Example 1
=========

```go
package example

import (
	"github.com/Nightgunner5/go.nbt"
	"io"
)

// Just define your struct the way you normally would. Most data types can be used with no modifications.
type Example1 struct {
	Name string `nbt:"name"` // If you need a lowercase first letter for an NBT field name or
	                         // your field name is invalid as an identifier in Go, you can
	                         // use tags similar to encoding/json and encoding/xml.

	Data [256]byte // go.nbt supports both arrays and slices for TAG_Byte_Array and TAG_Int_Array.

	Children []Example1 // Any type that can be used as a TAG_Compound can also be used as an element
	                    // in a TAG_List.
}

func ReadExample1(in io.Reader) (Example1, error) {
	var out Example1

	err := nbt.Unmarshal(nbt.Uncompressed, in, &out)

	return out, err
}
```
