package nbt

import (
	"fmt"
	"reflect"
)

func parseStruct(v reflect.Value) map[string]reflect.Value {
	parsed := make(map[string]reflect.Value)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			continue
		}

		name := f.Name
		if tag := f.Tag.Get("nbt"); tag != "" {
			name = tag
		}

		if _, exists := parsed[name]; exists {
			panic(fmt.Errorf("Multiple fields with name %#v", name))
		}

		parsed[name] = v.Field(i)
	}

	return parsed
}
