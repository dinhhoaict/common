package utils

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func SprettyPrint(data interface{}) string {
	if data == nil {
		return "nil"
	}
	var sb strings.Builder
	sprettyPrint(reflect.ValueOf(data), 0, &sb)
	return sb.String()
}

func sprettyPrint(v reflect.Value, depth int, sb *strings.Builder) {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		sb.WriteString(fmt.Sprintf("%s{", t))
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Interface) && f.IsNil() {
				continue
			}
			if !f.CanInterface() {
				if !f.CanAddr() {
					continue
				}
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()

			}
			sb.WriteString(fmt.Sprintf("%s:", t.Field(i).Name))
			sprettyPrint(f, depth + 1, sb)
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%s}", tab(depth)))
	case reflect.Ptr:
		if v.IsNil() {
			sb.WriteString("nil")
		}else{
			sb.WriteString("&")
			sprettyPrint(v.Elem(), depth, sb)
		}
	case reflect.Interface:
		if v.IsNil() {
			sb.WriteString("nil")
		}else{
			sprettyPrint(v.Elem(), depth, sb)
		}
	default:
		sb.WriteString(fmt.Sprintf("%v", v.Interface()))
	}
}

func tab(n int) string{
	return strings.Repeat("", n)
}