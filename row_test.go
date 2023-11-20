package json2excel

import (
	"errors"

	"reflect"
	"testing"
)

func TestStringArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected exceldata
		err      error
	}{
		{
			name:     "Empty Array",
			input:    []byte("[]"),
			expected: exceldata{},
		},
		{
			name:     "Empty Object",
			input:    []byte(`{}`),
			expected: exceldata{},
		},
		{
			name:     "Object",
			input:    []byte(`{"name":"kevin","surname":"saucedo"}`),
			expected: exceldata{Columns: []string{"name", "surname"}, RowsValues: [][]any{{"kevin", "saucedo"}}},
		},
		{
			name:  "Object 2",
			input: []byte(`{"name":"kevin","surname":"saucedo","age":25,"state":true}`),
			expected: exceldata{
				Columns:    []string{"name", "surname", "age", "state"},
				RowsValues: [][]any{{"kevin", "saucedo", float64(25), true}},
			},
		},
		{
			name:  "Array",
			input: []byte(`[{"name":"kevin","surname":"saucedo"},{"name":"jose","surname":"perez"}]`),
			expected: exceldata{
				Columns:    []string{"name", "surname"},
				RowsValues: [][]any{{"kevin", "saucedo"}, {"jose", "perez"}},
			},
		},
		{
			name:  "Array 2",
			input: []byte(`[{"name":"kevin","surname":"saucedo","state":true},{"name":"jose","surname":"perez","state":false}]`),
			expected: exceldata{
				Columns:    []string{"name", "surname", "state"},
				RowsValues: [][]any{{"kevin", "saucedo", true}, {"jose", "perez", false}},
			},
		},
		{
			name:     "Invalid Array",
			input:    []byte(`[[]`),
			expected: exceldata{},
			err:      ErrInvalidJSON,
		},
		{
			name:     "Invalid Object",
			input:    []byte(`{`),
			expected: exceldata{},
			err:      ErrInvalidJSON,
		},
		{
			name:     "Invalid Type",
			input:    []byte(`[1,2,3,4]`),
			expected: exceldata{},
			err:      ErrInvalidJSONType,
		},
		{
			name:     "Invalid Type 2",
			input:    []byte(`3`),
			expected: exceldata{},
			err:      ErrInvalidJSONType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result exceldata
			err := result.UnmarshalJSON(test.input)
			if err != nil && !errors.Is(err, test.err) {
				t.Errorf("UnmarshalJSON returned an error: %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("UnmarshalJSON result %v does not match expected %v", result, test.expected)
			}
		})
	}
}
