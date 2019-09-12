// Package querytostruct is an extremely tiny utility for unmarshalling and scanning
// querystrings into structs. It supports scanning into int*, float*, string,
// and []byte types, along with []int*, []float*, and []*string slices.
//
// ```go
// qs := "name=John+Doe&yes=true&count=42&tag=x&tag=y"
// q, _ := url.ParseQuery(qs)

// type test struct {
// 	Name  string   `q:"name"`
// 	Yes   bool     `q:"yes"`
// 	Count int      `q:"count"`
// 	Tags  []string `q:"tag"`
// }

// var t test
// fields, err := querytostruct.Unmarshal(q, &t, "q")
// fmt.Println("fields=", fields, "error=", err)
// fmt.Println(t)
// ```
package querytostruct

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal takes a url.Values object, takes its keys and values
// and applies them to a given struct using reflection. The field names
// are mapped to the struct fields based on a given tag tag. The field
// names that have been mapped are also return as a list. Supports string,
// bool, number types and their slices.
//
// eg:
// type Order struct {
// 	Symbol string `url:"symbol"`
// 	Tags []string `url:"tag"`
// }
func Unmarshal(q url.Values, obj interface{}, fieldTag string) ([]string, error) {
	ob := reflect.ValueOf(obj)
	if ob.Kind() == reflect.Ptr {
		ob = ob.Elem()
	}

	if ob.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Failed to encode form values to struct. Non struct type: %T", ob)
	}

	// Go through every field in the struct and look for it in the query map.
	var fields []string
	for i := 0; i < ob.NumField(); i++ {
		f := ob.Field(i)
		if f.IsValid() && f.CanSet() {
			tag := ob.Type().Field(i).Tag.Get(fieldTag)
			if tag == "" || tag == "-" {
				continue
			}

			// Got a struct field with a tag.
			// If that field exists in the arg and convert its type.
			// Tags are of the type `tagname,attribute`
			tag = strings.Split(tag, ",")[0]
			if _, ok := q[tag]; !ok {
				continue
			}

			scanned := false
			// The struct field is a slice type.
			if f.Kind() == reflect.Slice {
				var (
					vals    = q[tag]
					numVals = len(vals)
				)

				// Make a slice.
				sl := reflect.MakeSlice(f.Type(), numVals, numVals)

				// If it's a []byte slice (=[]uint8), assign here.
				if f.Type().Elem().Kind() == reflect.Uint8 {
					br := q.Get(tag)
					b := make([]byte, len(br))
					copy(b, br)
					f.SetBytes(b)
					continue
				}

				// Iterate through fasthttp's multiple args and assign values
				// to each item in the slice.
				for i, v := range vals {
					scanned = setVal(sl.Index(i), string(v))
				}
				f.Set(sl)
			} else {
				scanned = setVal(f, string(q.Get(tag)))
			}

			if scanned {
				fields = append(fields, tag)
			}
		}
	}
	return fields, nil
}

func setVal(f reflect.Value, val string) bool {
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := strconv.ParseInt(val, 10, 0); err == nil {
			f.SetInt(v)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := strconv.ParseUint(val, 10, 0); err == nil {
			f.SetUint(v)
		}
	case reflect.Float32, reflect.Float64:
		if v, err := strconv.ParseFloat(val, 0); err == nil {
			f.SetFloat(v)
		}
	case reflect.String:
		f.SetString(val)
	case reflect.Bool:
		b, _ := strconv.ParseBool(val)
		f.SetBool(b)
	default:
		return false
	}
	return true
}
