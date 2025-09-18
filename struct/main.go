package main

import "fmt"

type User struct {
	name string
}

func main() {
	u1 := User{}
	u1Ptr := &u1
	var u2 User = *u1Ptr
	fmt.Printf("%v %p \n", u1, &u1)
	fmt.Printf("%v %p \n", u1Ptr, &u1Ptr)
	fmt.Printf("%v %p \n", u2, &u2)
}
