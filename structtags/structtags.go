//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var cache sync.Map

type fieldUnit struct {
	index    int
	isSlice  bool
	elemType reflect.Type
}

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	valuePtr := reflect.ValueOf(ptr).Elem()

	cVal, ok := cache.Load(valuePtr.Type())
	if !ok {
		fields := make(map[string]fieldUnit)
		v := reflect.ValueOf(ptr).Elem()
		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i)
			tag := fieldInfo.Tag
			name := tag.Get("http")
			if name == "" {
				name = strings.ToLower(fieldInfo.Name)
			}
			isSlice := fieldInfo.Type.Kind() == reflect.Slice
			var elemType reflect.Type
			if isSlice {
				elemType = fieldInfo.Type.Elem()
			}
			fields[name] = fieldUnit{
				index:    i,
				isSlice:  isSlice,
				elemType: elemType,
			}
		}
		cache.Store(valuePtr.Type(), fields)
		cVal = fields
	}

	fields := cVal.(map[string]fieldUnit)

	for name, values := range req.Form {
		f, ok := fields[name]
		if !ok {
			continue
		}
		field := valuePtr.Field(f.index)
		if f.isSlice {
			slice := reflect.MakeSlice(field.Type(), len(values), len(values))
			for i, value := range values {
				elem := reflect.New(f.elemType).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				slice.Index(i).Set(elem)
			}
			field.Set(slice)
		} else if len(values) > 0 {
			if err := populate(field, values[len(values)-1]); err != nil {
				return fmt.Errorf("%s: %v", name, err)
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
