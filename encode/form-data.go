package encode

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"reflect"
)

// the post request upload file object
type File struct {
	Name string
	Data []byte
}

type FormData struct {
	v interface{}
}

func (x *FormData) SetValue(v interface{}) {
	x.v = v
}

func (f *FormData) Encode(w io.Writer) error {
	multipartWrite := multipart.NewWriter(w)
	defer func() {
		if err := multipartWrite.Close(); err != nil {
			panic(err)
		}
	}()

	switch t := f.v.(type) {
	case map[string]interface{}:
		for k, v := range t {
			// reflect value
			val := reflect.ValueOf(v)
			// reflect loopElem pointer dereference
			val = LoopElem(val)

			switch t2 := val.Kind(); t2 {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Float32, reflect.Float64,
				reflect.String:
				if err := multipartWrite.WriteField(k, valToStr(val, emptyField)); err != nil {
					return err
				}
			case reflect.Struct:
				file := v.(File) // TODO

				formFile, err := multipartWrite.CreateFormFile(k, file.Name)
				if err != nil {
					return err
				}
				if _, err := io.Copy(formFile, bytes.NewReader(file.Data)); err != nil {
					return err
				}
			default:
				if _, ok := val.Interface().([]byte); !ok {
					return errors.New("unknown type")
				}
				if err := multipartWrite.WriteField(k, valToStr(val, emptyField)); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("unknown type")
	}
	return nil
}
