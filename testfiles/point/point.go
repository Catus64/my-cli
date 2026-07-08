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

func add2(num *int, num2 int) *int {
	temp := *num + num2
	return &temp
}

func main() {
	//p1 := Person{Name: "Udin"}
	p2 := &Person{Name: "Nigga"}

	//fmt.Println(p1.Name)
	fmt.Println(reflect.TypeOf(p2))
	//fmt.Println(reflect.Ptr)

	//t := reflect.TypeOf(p2)
	//t = t.Elem()
	//fmt.Println(t.Field(0).Type)

	num := 10
	num2 := 10
	add(&num)
	//fmt.Println(num)
	numres := add2(&num, num2)
	println(*numres)

}
