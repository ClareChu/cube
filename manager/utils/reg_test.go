package utils

import (
	"fmt"
	"gotest.tools/assert"
	"reflect"
	"regexp"
	"testing"
)

type Foo struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int    `tag_name:"tag 3"`
	Bar       Bar    `json:"bar" tag_name:"bar"`
}

type Bar struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int    `tag_name:"tag 3"`
}

type Foo1 struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Bar       Bar1   `json:"bar"`
}

type Bar1 struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

func (f *Foo) reflect() {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		switch val.Kind() {
		case reflect.Int:
			//valueField := val.Field(i)
			//typeField := val.Type().Field(i)
		case reflect.String:
			val.SetString("Foo")
		case reflect.Bool:
			val.SetBool(true)
		case reflect.Struct:
			// Iterate over the struct fields
			for i := 0; i < val.NumField(); i++ {
				err := setDefaultValue(val.Field(i).Addr())
				if err != nil {
					//return err
				}
			}
		default:
			//return errors.New("Unsupported kind: " + val.Kind().String())

		}
		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
}

func TestBuild(t *testing.T) {
	f := &Foo{
		FirstName: "Drew",
		LastName:  "Olson",
		Age:       30,
	}

	f.reflect()
	//err := SetDefault(&f)
}

func TestRegx(t *testing.T) {
	f := &Foo1{
		FirstName: "Drew",
		LastName:  "Olson",
		Age:       30,
	}

	reg := regexp.MustCompile(`\$\{(.*?)\}`)


	st := Regx("${firstName}-${age}/${firstName}.....${lastName}", "${", "}", f)
	match := reg.MatchString("${firstName}-${age}/${firstName}.....${lastName}")
	assert.Equal(t, true, match)
	assert.Equal(t, "Drew-30/Drew.....Olson", st)

	st1 := Regx(".....${firstName}-${age}/${firstName}.....${lastName}adsadasdasd", "${", "}", f)
	assert.Equal(t, ".....Drew-30/Drew.....Olsonadsadasdasd", st1)

	st2 := "3123123"
	match2 := reg.MatchString(st2)
	assert.Equal(t, false, match2)

	st3 := "312${3123"
	match3 := reg.MatchString(st3)
	assert.Equal(t, false, match3)

	st3 = "312${}3123"
	match3 = reg.MatchString(st3)
	assert.Equal(t, true, match3)

	st3 = "312${}3123"
	match3 = reg.MatchString(st3)
	assert.Equal(t, true, match3)
}
