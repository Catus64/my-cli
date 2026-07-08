package main

import "fmt"

type Parent struct {
	Name string
}

func (p *Parent) sharedMethod() string {
	return "Shared logic: " + p.Name
}

type Child interface {
	DoSomething() string
	getName() string
}

type ChildA struct {
	Parent
	SpecificA string
}

func (c *ChildA) DoSomething() string {
	return "ChildA doing smtg" + c.SpecificA
}

func (c *ChildA) getName() string {
	return c.Name
}

type ChildB struct {
	Parent
	SpecificB string
}

func (b *ChildB) DoSomething() string {
	return "ChildA doing smtg " + b.SpecificB
}

func (c *ChildB) getName() string {
	return c.Name
}

func NewChild(childType string, name string) Child {
	switch childType {
	case "A":
		return &ChildA{
			Parent:    Parent{Name: name},
			SpecificA: "Special A",
		}
	case "B":
		return &ChildB{
			Parent:    Parent{Name: name},
			SpecificB: "Special B",
		}
	default:
		return nil
	}
}

func main() {
	c := NewChild("A", "Ali")
	fmt.Println(c.getName())
}
