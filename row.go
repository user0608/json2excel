package json2excel

import (
	"errors"

	"github.com/tidwall/gjson"
)

var ErrInconsistentJSON = errors.New("all JSON objects in the array must be of the same type and have consistent field names")
var ErrInvalidJSON = errors.New("please check the JSON data for errors")
var ErrInvalidJSONType = errors.New("only type array or object are allowed")

type exceldata struct {
	Columns    []string
	RowsValues [][]any
}

func (*exceldata) getKeys(results []gjson.Result) []string {
	var keys = make([]string, len(results))
	for i, r := range results {
		keys[i] = r.String()
	}
	return keys
}
func (excel *exceldata) ensureEquals(keys []string) error {

	if len(excel.Columns) != len(keys) {
		return ErrInconsistentJSON
	}
	for i := range excel.Columns {
		if keys[i] != excel.Columns[i] {
			return ErrInconsistentJSON
		}
	}
	return nil
}
func (excel *exceldata) UnmarshalJSON(data []byte) error {
	if !gjson.ValidBytes(data) {
		return ErrInvalidJSON
	}
	result := gjson.ParseBytes(data)
	if !result.IsArray() && !result.IsObject() {
		return ErrInvalidJSONType
	}
	for i, r := range result.Array() {
		if !r.IsObject() {
			return ErrInvalidJSONType
		}
		keys := excel.getKeys(r.Get("@keys").Array())
		if i == 0 {
			if len(keys) == 0 {
				return nil
			}
			excel.Columns = append(excel.Columns, keys...)
		}
		if err := excel.ensureEquals(keys); err != nil {
			return err
		}
		var row = make([]any, len(keys))
		for i, k := range keys {
			row[i] = r.Get(k).Value()
		}
		excel.RowsValues = append(excel.RowsValues, row)
	}
	return nil
}
