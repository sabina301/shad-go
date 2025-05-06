//go:build !solution

package jsonlist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(w io.Writer, slice interface{}) error {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	sliceValue := reflect.ValueOf(slice)
	builder := &strings.Builder{}
	for i := 0; i < sliceValue.Len(); i++ {
		if i != 0 {
			builder.WriteByte(' ')
		}
		elem := sliceValue.Index(i)
		if elem.Kind() == reflect.Struct {
			builder.WriteString("{")
			numField := elem.NumField()
			firstFieldFlag := true
			for j := 0; j < numField; j++ {
				field := elem.Type().Field(j)
				value := elem.Field(j)
				if value.IsZero() {
					continue
				}
				if !firstFieldFlag {
					builder.WriteString(", ")
				} else {
					firstFieldFlag = false
				}

				_, err := fmt.Fprintf(builder, `"%s": `, field.Name)
				if err != nil {
					return err
				}

				switch value.Kind() {
				case reflect.String:
					_, err = fmt.Fprintf(builder, `"%v"`, value)
					if err != nil {
						return err
					}
				default:
					_, err = fmt.Fprintf(builder, "%v", value)
					if err != nil {
						return err
					}
				}
			}
			builder.WriteString("}")
		} else {
			_, err := fmt.Fprintf(builder, "%v", elem)
			if err != nil {
				return err
			}
		}
	}

	_, err := w.Write([]byte(builder.String()))
	return err
}

func splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' && (i > 0 && data[i-1] != ':') {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func deleteQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, "")
}

func scanUnitOfSlice(r io.Reader, slice any) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitter)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		if scanner.Text() == "" {
			return nil
		}
		text := deleteQuotes(scanner.Text())
		err := decode(slice, text)
		if err != nil {
			return err
		}
	}
	return nil
}

func valueWithTrueType(s string) any {
	s = strings.TrimSpace(s)
	if boolValue, err := strconv.ParseBool(s); err == nil {
		return boolValue
	}
	if floatValue, err := strconv.ParseFloat(s, 64); err == nil {
		return floatValue
	}
	if intValue, err := strconv.Atoi(s); err == nil {
		return intValue
	}
	return s
}

func getFieldNameForStruct(str string) string {
	builder := &strings.Builder{}
	for _, e := range str[1:] {
		if e == ':' {
			break
		}
		builder.WriteRune(e)
	}
	return builder.String()
}

func getFieldValueForStruct(str string) string {
	builder := &strings.Builder{}
	valFlag := false
	for i := 1; i < len(str); i++ {
		if str[i] == ':' {
			i += 2
			valFlag = true
		}
		if str[i] == '}' {
			break
		}
		if valFlag == true {
			builder.WriteByte(str[i])
		}
	}
	return builder.String()
}

func parseValue(data string, targetType reflect.Type) (reflect.Value, error) {
	value := valueWithTrueType(data)
	resultValue := reflect.ValueOf(value)

	if !resultValue.Type().ConvertibleTo(targetType) {
		return reflect.Value{}, fmt.Errorf("error")
	}

	return resultValue.Convert(targetType), nil
}

func decode(slice any, data string) error {
	sliceValue := reflect.ValueOf(slice).Elem()
	switch sliceValue.Type().Elem().Kind() {
	case reflect.Int:
		intData, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(int(intData))))
	case reflect.String:
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(data)))
	case reflect.Struct:
		fieldName := getFieldNameForStruct(data)
		fieldValue := getFieldValueForStruct(data)

		elemType := sliceValue.Type().Elem()
		newStruct := reflect.New(elemType).Elem()

		field := newStruct.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("error")
		}
		if !field.CanSet() {
			return fmt.Errorf("error")
		}

		parsedValue, err := parseValue(fieldValue, field.Type())
		if err != nil {
			return fmt.Errorf("error")
		}

		field.Set(parsedValue)
		sliceValue.Set(reflect.Append(sliceValue, newStruct))
	case reflect.Interface:
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(valueWithTrueType(data))))
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func Unmarshal(r io.Reader, slice interface{}) error {
	if reflect.TypeOf(slice).Kind() != reflect.Ptr {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	if reflect.TypeOf(slice).Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	return scanUnitOfSlice(r, slice)
}
