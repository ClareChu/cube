package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Regx(value, firstKey, lastKey string, o interface{}) (newValue string) {
	for _, k := range strings.Split(value, firstKey) {
		if !strings.Contains(k, lastKey) {
			newValue = newValue + k
		} else {
			s := strings.Split(k, lastKey)
			s1 := s[0]
			tag := fmt.Sprintf("%v", GetFileName(o, s1))
			newValue = newValue + tag
			s2 := s[1]
			newValue = newValue + s2
		}
	}
	return
}

func Reflect(o interface{}) {
	reflectType := reflect.TypeOf(o)
	switch reflectType.Kind() {
	case reflect.Struct:
		for i := 0; i < reflectType.NumField(); i++ {
			f := reflectType.Field(i)
			if f.Tag != "" {
				f.Tag.Get("json")
			}
		}
	}
}

func GetFileName(f interface{}, tagName string) (o interface{}) {
	val := reflect.ValueOf(f).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
		jsonName := tag.Get("json")
		if jsonName == tagName {
			return valueField.Interface()
		}
	}
	return ""
}

func SetDefault(s interface{}) error {
	return setDefaultValue(reflect.ValueOf(s))
}

func setDefaultValue(v reflect.Value) error {

	if v.Kind() != reflect.Ptr {
		return errors.New("Not a pointer value")
	}

	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Int:
		v.Elem().NumField()

	case reflect.String:
		v.SetString("Foo")
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Struct:
		// Iterate over the struct fields
		for i := 0; i < v.NumField(); i++ {
			err := setDefaultValue(v.Field(i).Addr())
			if err != nil {
				return err
			}
		}

	default:
		return errors.New("Unsupported kind: " + v.Kind().String())

	}

	return nil
}
