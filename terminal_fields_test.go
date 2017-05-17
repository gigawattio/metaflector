package metaflect

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
	Name         string
	Multiplier   float64
	Active       bool
	Contents     []Content
	PtrA         *uint8
	PtrB         *int64
	hiddenString string
	hiddenStruct struct {
		StillHidden string
	}
}

type Content struct {
	Key     string
	Value   string
	Version int
}

var threeve = int64(3)

func TestTerminalFields(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected []string
	}{
		{
			in: Foo{
				Bar: Bar{
					Baz: Baz{
						Name:       "shellz",
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
						Key:     "name-o",
						Value:   "bingo",
						Version: 1,
					},
					{
						Key:     "name-o",
						Value:   "dingo",
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
			// if missing.Equal(mapset.NewSet()) && extra.Equal(mapset.NewSet())
			if len(missing.ToSlice()) == 0 && len(extra.ToSlice()) == 0 {
				t.Errorf("[i=%v] Improper field ordering detected", i)
				t.Errorf("[i=%v] expected=%v", i, expected)
				t.Errorf("[i=%v]  actual=%v", i, actual)
			} else {
				t.Errorf("[i=%v] Expected fields mismatch; missing=%v extra=%v", i, missing, extra)
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

// toIfaces converts a slice of string to a slice of interface.
func toIfaces(src []string) []interface{} {
	ifaces := make([]interface{}, len(src))
	for i, s := range src {
		ifaces[i] = s
	}
	return ifaces
}
