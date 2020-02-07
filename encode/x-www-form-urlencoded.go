package encode

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"reflect"
	"strings"
)

type XWwwFormUrlencoded struct {
	v interface{}
}

func (x *XWwwFormUrlencoded) SetValue(v interface{}) {
	x.v = v
}

func (x *XWwwFormUrlencoded) Encode(w io.Writer) error {
	switch t := x.v.(type) {
	case string:
		if _, err := w.Write(StringToBytes(t)); err != nil {
			return err
		}
	case []string:
		var buffer bytes.Buffer
		for i := 0; i < len(t); i++ {
			buffer.WriteString(t[i] + "&")
		}
		if _, err := io.WriteString(w, strings.Join(t, "&")); err != nil {
			return err
		}
	case []byte:
		if _, err := w.Write(t); err != nil {
			return err
		}
	case map[string]string:
		values := make(url.Values)
		for key, value := range t {
			values.Set(key, value)
		}
		if _, err := io.WriteString(w, values.Encode()); err != nil {
			return err
		}
	case url.Values:
		if _, err := io.WriteString(w, t.Encode()); err != nil {
			return err
		}
	case io.Reader:
		if _, err := io.Copy(w, t); err != nil {
			return err
		}
	case map[string]interface{}:
		values := make(url.Values)
		for k, v := range t {
			// reflect value
			val := reflect.ValueOf(v)

			// reflect loopElem pointer dereference
			val = LoopElem(val)

			switch t := val.Kind(); t {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			case reflect.Float32, reflect.Float64:
			case reflect.String:
			default:
				if _, ok := val.Interface().([]byte); !ok {
					return errors.New("unknown type")
				}
			}

			values.Add(k, valToStr(val, emptyField))
		}
		if _, err := io.WriteString(w, values.Encode()); err != nil {
			return err
		}
	default:
		return errors.New("unknown type")
	}
	return nil
}
