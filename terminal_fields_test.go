package metaflector

import (
	"reflect"
	"testing"

	"github.com/deckarep/golang-set"
)

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
