package decoders

import (
	"encoding/json"
	"io"
	"reflect"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	t := val.Type()

	if t.Kind() != reflect.Slice {
		return json.NewDecoder(dec.r).Decode(v)
	}

	elem := val.Type().Elem()
	decoder := json.NewDecoder(dec.r)

	for decoder.More() {
		obj := reflect.New(elem).Interface()

		if err := decoder.Decode(&obj); err != nil {
			return err
		}
		ptr := reflect.Indirect(reflect.ValueOf(obj))
		val.Set(reflect.Append(val, ptr))
	}

	return nil
}
