package metaflector

import (
	"reflect"
	"sort"
)

// Separator is the string to use as the delimiter between field names.
var Separator = "."

// TerminalFields returns a slice of strings representing the full path in dot
// notation for each "terminal" field, where a terminal field is defined as a
// primitive type (and without additional sub-fields, e.g. an int).
//
// This implementation uses a BFS queue-based traversal to minimize stack
// depth.
//
// Important note: Circular references aren't supported yet and will blow up.
func TerminalFields(obj interface{}) []string {
	if obj == nil {
		return nil
	}

	type item struct {
		obj  interface{}
		path string
	}

	var (
		paths = []string{}
		queue = []item{
			{obj: obj},
		}
	)

	for len(queue) > 0 {
		EachField(queue[0].obj, func(child interface{}, name string, kind reflect.Kind) {
			if len(queue[0].path) > 0 {
				name = queue[0].path + Separator + name
			}

			// Filter and exclude non-terminal types.
			if isTerminal(kind) {
				paths = append(paths, name)
			} else {
				i := item{
					obj:  child,
					path: name,
				}
				queue = append(queue, i)
			}
		})
		queue = queue[1:]
	}

	sort.Strings(paths)

	return paths
}

// IterFunc is the type signature of callbacks sent to `EachField`.
type IterFunc func(child interface{}, name string, kind reflect.Kind)

// EachField invokes a callback with the value, name, and kind for each field in
// a struct.  The function returns false if the passed object cannot be
// resolved to a struct or non-empty slice / array (i.e. if must be a
// non-terminal type).
func EachField(obj interface{}, fn IterFunc) (ok bool) {
	if obj, ok = ResolveUnderlying(obj); !ok || obj == nil {
		ok = false
		return
	}

	v := reflect.ValueOf(obj)

	for i := 0; i < v.NumField(); i++ {
		// Skip unexported (signaled by non-mepty pkgpath) or anonymous fields.
		if sf := v.Type().Field(i); sf.PkgPath != "" || sf.Anonymous {
			continue
		}

		var (
			field = v.Field(i)
			name  = v.Type().Field(i).Name
			kind  = field.Kind()
		)

		for kind == reflect.Ptr {
			// Resolve underlying pointer type.
			kind = field.Type().Elem().Kind()
		}

		switch kind {
		case reflect.Struct:
			fn(field.Interface(), name, kind)

		case reflect.Slice, reflect.Array:
			if firstObj, ok := ResolveUnderlying(field.Interface()); ok {
				EachField(firstObj, func(child interface{}, childName string, childKind reflect.Kind) {
					fn(child, name+Separator+childName, childKind)
				})
			}

		case reflect.String, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fn(nil, name, kind)
		}
	}

	ok = true
	return
}

// ResolveUnderlying takes an interface{} (object) and resolves it to an
// instance of the underlying type through 3 varieties of resolution mutation:
//
// 1. Pointers are resolved to whatever they're referencing.
//
// 2. Slices and arrays, when not empty, are resolved to the type of the first
// element.
//
// 3. Test if the end result is a struct.
func ResolveUnderlying(obj interface{}) (resolved interface{}, ok bool) {
	if obj, ok = resolvePointer(obj); !ok {
		return
	}

	if hasType(obj, []reflect.Kind{reflect.Slice, reflect.Array}) {
		v := reflect.ValueOf(obj)
		if v.Len() == 0 {
			return
		}
		obj = nil
		// Find first non-nil element.
		for i := 0; i < v.Len(); i++ {
			if value := v.Index(i); isStruct(value.Interface()) || !value.IsNil() {
				obj = value.Interface()
				break
			}
		}
	}

	if obj, ok = resolvePointer(obj); !ok {
		return
	}

	if !isStruct(obj) {
		ok = false
		return
	}

	resolved = obj
	ok = true
	return
}

// resolvePointer keeps digging until it can't inspect any further or a
// non-pointer is unearthed.
func resolvePointer(obj interface{}) (interface{}, bool) {
	for isPointer(obj) {
		v := reflect.ValueOf(obj)
		if v.IsNil() {
			// Can't do further inspection on nil values.
			return nil, false
		}
		obj = reflect.Indirect(v).Interface()
	}
	if obj == nil {
		return nil, false
	}
	return obj, true
}

func isPointer(obj interface{}) bool {
	if obj == nil {
		return false
	}
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}

func isStruct(obj interface{}) bool {
	if obj == nil {
		return false
	}
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func hasType(obj interface{}, types []reflect.Kind) bool {
	if obj == nil {
		return false
	}
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}

// isTerminal returns true if the supplied reflect.Kind is a terminal (i.e.
// primitive) type with no additional sub-fields (e.g. an int, bool, string).
func isTerminal(kind reflect.Kind) bool {
	return kind != reflect.Struct && kind != reflect.Slice && kind != reflect.Array
}
