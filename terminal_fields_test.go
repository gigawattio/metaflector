package metaflector

import (
	"reflect"
	"testing"

	"github.com/deckarep/golang-set"
)

// Sample structs for testing.
// Try not to think too hard about them.

type Foo struct {
	Bar       Bar
	StructPtr *Bar
	Contents  []Content
}

type Bar struct {
	Baz   Baz
	Stock string
}

type Baz struct {
	Name              string
	Multiplier        float64
	Active            bool
	Contents          []Content
	ContentPtrs       []*Content
	PtrContentPtrPtrs *[]**Content
	PtrA              *uint8
	PtrB              *int64
	hiddenString      string
	hiddenStruct      struct {
		StillHidden string
	}
}

type Content struct {
	Key     string
	Value   string
	Version int
}

var (
	threeve      = int64(3)
	uEight       = uint8(8)
	notHotdogPtr = &Content{
		Key:   "not",
		Value: "hotdog",
	}
)

func TestTerminalFields(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected []string
	}{
		{
			in: Foo{
				Bar: Bar{
					Baz: Baz{
						Name:       "hotdog",
						Multiplier: 10.5,
						Active:     true,
						Contents: []Content{
							{
								Key:     "99",
								Value:   "problems",
								Version: -1,
							},
						},
						PtrB: &threeve,
					},
					Stock: "max",
				},
				Contents: []Content{
					{
						Key:     "hotdog",
						Value:   "yes",
						Version: 1,
					},
					{
						Key:     "not hotdog",
						Value:   "no",
						Version: 2,
					},
				},
			},
			expected: []string{
				"Bar.Baz.Active",
				"Bar.Baz.Contents.Key",
				"Bar.Baz.Contents.Value",
				"Bar.Baz.Contents.Version",
				"Bar.Baz.Multiplier",
				"Bar.Baz.Name",
				"Bar.Baz.PtrA",
				"Bar.Baz.PtrB",
				"Bar.Stock",
				"Contents.Key",
				"Contents.Value",
				"Contents.Version",
			},
		},
		{
			in: &Foo{
				Bar: Bar{
					Baz: Baz{
						ContentPtrs: []*Content{
							nil,
							nil,
							nil,
							notHotdogPtr,
						},
					},
					Stock: "pile",
				},
				StructPtr: &Bar{},
			},
			expected: []string{
				"Bar.Baz.Active",
				"Bar.Baz.ContentPtrs.Key",
				"Bar.Baz.ContentPtrs.Value",
				"Bar.Baz.ContentPtrs.Version",
				"Bar.Baz.Multiplier",
				"Bar.Baz.Name",
				"Bar.Baz.PtrA",
				"Bar.Baz.PtrB",
				"Bar.Stock",
				"StructPtr.Baz.Active",
				"StructPtr.Baz.Multiplier",
				"StructPtr.Baz.Name",
				"StructPtr.Baz.PtrA",
				"StructPtr.Baz.PtrB",
				"StructPtr.Stock",
			},
		},
		{
			in: &Foo{
				Bar: Bar{
					Baz: Baz{
						PtrContentPtrPtrs: &[]**Content{
							nil,
							nil,
							nil,
							&notHotdogPtr,
						},
					},
					Stock: "pile",
				},
				StructPtr: &Bar{},
			},
			expected: []string{
				"Bar.Baz.Active",
				"Bar.Baz.Multiplier",
				"Bar.Baz.Name",
				"Bar.Baz.PtrA",
				"Bar.Baz.PtrB",
				"Bar.Baz.PtrContentPtrPtrs.Key",
				"Bar.Baz.PtrContentPtrPtrs.Value",
				"Bar.Baz.PtrContentPtrPtrs.Version",
				"Bar.Stock",
				"StructPtr.Baz.Active",
				"StructPtr.Baz.Multiplier",
				"StructPtr.Baz.Name",
				"StructPtr.Baz.PtrA",
				"StructPtr.Baz.PtrB",
				"StructPtr.Stock",
			},
		},
	}

	for i, test := range tests {
		if expected, actual := test.expected, TerminalFields(test.in); !reflect.DeepEqual(actual, expected) {
			var (
				setExpected = mapset.NewSetFromSlice(toIfaces(test.expected))
				setActual   = mapset.NewSetFromSlice(toIfaces(actual))
				missing     = setExpected.Difference(setActual)
				extra       = setActual.Difference(setExpected)
			)
			if len(missing.ToSlice()) == 0 && len(extra.ToSlice()) == 0 {
				t.Errorf("[i=%v] Improper field ordering detected", i)
				t.Errorf("[i=%v] expected=%v", i, expected)
				t.Errorf("[i=%v]   actual=%v", i, actual)
			} else {
				t.Errorf("[i=%v] Expected fields mismatch; missing=%v extra=%v", i, missing, extra)
				t.Errorf("[i=%v] expected=%v", i, expected)
				t.Errorf("[i=%v]   actual=%v", i, actual)
			}
		}
	}
}

func TestEachField(t *testing.T) {
	tests := []interface{}{
		false,
		true,
		"hotdog",
		"not hotdog",
		3,
		3.3,
		&threeve,
		'c',
	}

	for i, test := range tests {
		if ok := EachField(test, func(_ interface{}, _ string, _ reflect.Kind) {}); ok {
			t.Errorf("[i=%v] 'ok' should have been false but actual=%v", i, ok)
		}
	}
}

func TestResolveUnderlying(t *testing.T) {
	var (
		structPtr = &Content{
			Key:     "structPtr",
			Version: 98,
		}
		structPtrPtr = &structPtr
		slicePtr     = &[]**Content{
			nil,
			structPtrPtr,
		}
		slicePtrPtr = &slicePtr
	)

	tests := []struct {
		obj      interface{}
		resolved interface{}
		ok       bool
	}{
		{}, // Empty test case.
		{
			obj:      3,
			resolved: nil,
			ok:       false,
		},
		{
			obj: Content{
				Key:     "test1",
				Version: 99,
			},
			resolved: Content{
				Key:     "test1",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: &Content{
				Key:     "test2",
				Version: 99,
			},
			resolved: Content{
				Key:     "test2",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: []Content{
				{
					Key:     "test3",
					Version: 99,
				},
			},
			resolved: Content{
				Key:     "test3",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: []Content{
				{
					Key:     "test5.1",
					Version: 99,
				},
				{
					Key:     "test5.2",
					Version: 100,
				},
			},
			resolved: Content{
				Key:     "test5.1",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: []*Content{
				{
					Key:     "test6",
					Version: 99,
				},
			},
			resolved: Content{
				Key:     "test6",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: []*Content{
				nil,
				{
					Key:     "test7",
					Version: 99,
				},
			},
			resolved: Content{
				Key:     "test7",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: []*Content{
				nil,
				{
					Key:     "test8.1",
					Version: 99,
				},
				{
					Key:     "test8.2",
					Version: 100,
				},
				nil,
			},
			resolved: Content{
				Key:     "test8.1",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: &[]Content{
				{
					Key:     "test9.1",
					Version: 99,
				},
				{
					Key:     "test9.2",
					Version: 100,
				},
			},
			resolved: Content{
				Key:     "test9.1",
				Version: 99,
			},
			ok: true,
		},
		{
			obj: &[]*Content{
				nil,
				{
					Key:     "test10.1",
					Version: 99,
				},
				{
					Key:     "test10.2",
					Version: 100,
				},
				nil,
			},
			resolved: Content{
				Key:     "test10.1",
				Version: 99,
			},
			ok: true,
		},
		{
			obj:      structPtrPtr,
			resolved: *structPtr,
			ok:       true,
		},
		{
			obj:      slicePtrPtr,
			resolved: *structPtr,
			ok:       true,
		},
	}

	for i, test := range tests {
		obj, ok := ResolveUnderlying(test.obj)
		if expected, actual := test.resolved, obj; !reflect.DeepEqual(actual, expected) {
			t.Errorf("[i=%v] Expected obj=%v but actual=%v for test=%+v", i, expected, actual, test)
		}
		if expected, actual := test.ok, ok; actual != expected {
			t.Errorf("[i=%v] Expected ok=%v but actual=%v for test=%+v", i, expected, actual, test)
		}
	}
}

// toIfaces converts a slice of string to a slice of interface.
func toIfaces(src []string) []interface{} {
	ifaces := make([]interface{}, len(src))
	for i, s := range src {
		ifaces[i] = s
	}
	return ifaces
}

func TestGet(t *testing.T) {
	tests := []struct {
		obj      interface{}
		path     string
		expected interface{}
	}{
		{},
		{
			obj:      "",
			expected: "",
		},
		{
			obj:      "foo",
			expected: "foo",
		},
		{
			obj:      true,
			expected: true,
		},
		{
			obj:      5,
			expected: 5,
		},
		{
			obj:      5.59,
			expected: 5.59,
		},
		{
			obj:      map[string]string{},
			expected: map[string]string{},
		},
		{
			obj:      map[string]string{},
			path:     "Foo",
			expected: nil,
		},
		{
			obj: Content{
				Key:     "test1",
				Value:   "v1",
				Version: 99,
			},
			path:     "Key",
			expected: "test1",
		},
		{
			obj: &Content{
				Key:     "test1",
				Value:   "v1",
				Version: 99,
			},
			path:     "Key",
			expected: "test1",
		},
		{
			obj: []Content{
				{
					Key:     "test1",
					Value:   "v1",
					Version: 99,
				},
			},
			path:     "Key",
			expected: []interface{}{"test1"},
		},
		{
			obj: []*Content{
				&Content{
					Key:     "test1",
					Value:   "v1",
					Version: 99,
				},
			},
			path:     "Key",
			expected: []interface{}{"test1"},
		},
		{
			obj: []*Content{
				nil,
				&Content{
					Key:     "test1",
					Value:   "v1",
					Version: 99,
				},
			},
			path:     "Key",
			expected: []interface{}{"test1"},
		},
		{
			obj: []*Content{
				nil,
				&Content{
					Key:     "test1",
					Value:   "v1",
					Version: 99,
				},
				&Content{
					Key:     "test2",
					Value:   "v2",
					Version: 100,
				},
			},
			path:     "Key",
			expected: []interface{}{"test1", "test2"},
		},
		{
			obj: [][]*Content{
				[]*Content{
					nil,
					&Content{
						Key:   "test1",
						Value: "v1",
					},
					&Content{
						Key:   "test2",
						Value: "v2",
					},
				},
				[]*Content{
					nil,
					&Content{
						Key:   "test3",
						Value: "v3",
					},
					&Content{
						Key:   "test4",
						Value: "v4",
					},
				},
			},
			path: "Key",
			expected: []interface{}{
				[]interface{}{"test1", "test2"},
				[]interface{}{"test3", "test4"},
			},
		},
		{
			obj: [][]*Content{
				[]*Content{
					nil,
					&Content{
						Key:   "test1",
						Value: "v1",
					},
					&Content{
						Key:   "test2",
						Value: "v2",
					},
				},
				[]*Content{
					nil,
					&Content{
						Key:   "test3",
						Value: "v3",
					},
					&Content{
						Key:   "test4",
						Value: "v4",
					},
				},
			},
			path: "Key",
			expected: []interface{}{
				[]interface{}{"test1", "test2"},
				[]interface{}{"test3", "test4"},
			},
		},
		{
			obj: Foo{
				Bar: Bar{
					Baz: Baz{
						Name: "Nested",
					},
				},
			},
			path:     "Bar.Baz.Name",
			expected: "Nested",
		},
		{
			obj: Foo{
				StructPtr: &Bar{
					Stock: "back",
				},
			},
			path:     "StructPtr.Stock",
			expected: "back",
		},
		{
			obj: Foo{
				Bar: Bar{
					Baz: Baz{
						PtrA: &uEight,
					},
				},
			},
			path:     "Bar.Baz.PtrA",
			expected: &uEight,
		},
		{
			obj: Foo{
				Bar: Bar{
					Baz: Baz{
						PtrB: &threeve,
					},
				},
			},
			path:     "Bar.Baz.PtrB",
			expected: &threeve,
		},
		{
			obj:  Foo{},
			path: "Bar.Baz.PtrB",
			expected: func() interface{} {
				// NB: This is to return the "right type" of nil (*int64).
				// See also: https://stackoverflow.com/a/43059772/293064
				var x *int64
				return x
			}(),
		},
		{
			obj: Foo{
				StructPtr: &Bar{
					Stock: "back",
				},
			},
			path:     "StructPtr.Stock",
			expected: "back",
		},
		{
			obj: struct {
				Foo *Foo
			}{
				Foo: &Foo{
					Bar: Bar{
						Baz: Baz{
							ContentPtrs: []*Content{
								&Content{Version: 1},
								&Content{Version: 5},
								&Content{Version: 3},
								&Content{Version: 2},
								&Content{Version: 4},
							},
						},
					},
				},
			},
			path: "Foo.Bar.Baz.ContentPtrs.Version",
			expected: []interface{}{
				int64(1),
				int64(5),
				int64(3),
				int64(2),
				int64(4),
			},
		},
	}

	for i, test := range tests {
		t.Logf("Starting %v", i)
		if expected, actual := test.expected, Get(test.obj, test.path); !reflect.DeepEqual(actual, expected) {
			t.Errorf("[i=%v] Expected value=%[2]T/%[2]v but actual=%[3]T/%[3]v for test=%# v", i, expected, actual, test)
		}
	}
}
