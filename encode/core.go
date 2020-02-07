package encode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	timeType   = reflect.TypeOf(time.Time{})
	emptyField = reflect.StructField{}
)

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// reflect loopElem pointer dereference
func LoopElem(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}

func timeToStr(v reflect.Value, sf reflect.StructField) string {
	t := v.Interface().(time.Time)
	timeFormat := sf.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		var tv int64
		if tf == "unix" {
			tv = t.Unix()
		} else {
			tv = t.UnixNano()
		}

		return strconv.FormatInt(tv, 10)
	}

	return t.Format(timeFormat)
}

func valToStr(v reflect.Value, sf reflect.StructField) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Type() == timeType {
		return timeToStr(v, sf)
	}

	if b, ok := v.Interface().([]byte); ok {
		return *(*string)(unsafe.Pointer(&b))
	}

	return fmt.Sprint(v.Interface())
}
