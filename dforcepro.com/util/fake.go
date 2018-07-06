package util

import (
	"math/rand"
	"reflect"
	"time"
)

func Fuzz(e interface{}) {
	ty := reflect.TypeOf(e)

	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}

	if ty.Kind() == reflect.Struct {
		value := reflect.ValueOf(e).Elem()
		for i := 0; i < ty.NumField(); i++ {
			field := value.Field(i)
			if "ObjectId" == field.Type().Name() {
				continue
			}

			if field.CanSet() {
				field.Set(fuzzValueFor(field.Kind()))
			}
		}

	}
}

// fuzzValueFor Generates random values for the following types:
// string, bool, int, int32, int64, float32, float64
func fuzzValueFor(kind reflect.Kind) reflect.Value {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	switch kind {
	case reflect.String:
		return reflect.ValueOf(randomString(25))
	case reflect.Int:
		return reflect.ValueOf(r.Int())
	case reflect.Int32:
		return reflect.ValueOf(r.Int31())
	case reflect.Int64:
		return reflect.ValueOf(r.Int63())
	case reflect.Float32:
		return reflect.ValueOf(r.Float32())
	case reflect.Float64:
		return reflect.ValueOf(r.Float64() * 100)
	case reflect.Bool:
		val := r.Intn(2) > 0
		return reflect.ValueOf(val)
	}

	return reflect.ValueOf("")
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
