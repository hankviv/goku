package _struct

import (
	"fmt"
	"testing"
)

type From struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type To struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func Test_structtomap(t *testing.T) {
	from := From{
		A: 110,
		B: "B",
	}
	to := To{}

	res := CopyStruct(from, to)

	fmt.Printf("%+v", res)
}

func Test_map2st(t *testing.T) {
	type A struct {
		Name  string
		Id    int
		state int
	}
	type B struct {
		Name  string
		Id    int
		state int
	}
	a := A{Name: "aaa", Id: 11, state: 1}
	b := B{}
	//结构
	bb := CopyStruct(a, &b)
	fmt.Println(bb)
}

