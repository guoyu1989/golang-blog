package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4
	// First way to get the type
	fmt.Println("type: ", reflect.TypeOf(x))

	// Second way to get the type
	v := reflect.ValueOf(x)
	fmt.Println("type: ", reflect.TypeOf(v))
	fmt.Println("type: ", v.Type())
	fmt.Println("kind: ", v.Kind())
	fmt.Println("value: ", v.Float())

	// Get back the original interface
	y := v.Interface()
	fmt.Println(y)

	VarCanSet()

	SetStructField()
}

func VarCanSet() {
	var x float64 = 3.4
	v := reflect.ValueOf(&x)
	fmt.Println("Type: ", v.Type())
	fmt.Println("settability: ", v.CanSet())
}

func SetStructField() {
	t := struct {
		A int
		B string
	}{
		23,
		"skidoo",
	}

	s := reflect.ValueOf(&t).Elem()
	fmt.Println("Type is: ", s.Type())

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf(
			"%d: %s %s = %v\n", i, s.Type().Field(i).Name, f.Type(), f.Interface())
	}
}
