package nbt

import "reflect"

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

		parsed[name] = v.Field(i)
	}

	return parsed
}
