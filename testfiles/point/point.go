package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
}

func add(num *int) {
	*num++
}

func main() {
	//p1 := Person{Name: "Udin"}
	p2 := &Person{Name: "Nigga"}

	//fmt.Println(p1.Name)
	fmt.Println(reflect.TypeOf(p2))
	//fmt.Println(reflect.Ptr)

	t := reflect.TypeOf(p2)
	t = t.Elem()
	fmt.Println(t.Field(0).Type)

	num := 5
	add(&num)
	fmt.Println(num)
}
