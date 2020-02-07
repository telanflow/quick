package requests

import (
	"fmt"
	"reflect"
	"testing"
)

func TestServer(t *testing.T) {
	list := []string{"1", "2"}
	listVal := reflect.ValueOf(list)
	fmt.Println(listVal.Kind())
}
