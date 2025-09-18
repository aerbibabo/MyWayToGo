package main

import "fmt"

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Printf("地址 %p, 值 %d \n", &i, i)
			//println(i)
		}()
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			fmt.Printf("地址 %p, 值 %v\n", &i, i)
			//println(i)
		}(i)
	}
}

func DeferClosureLoopV3() {
	//var j int
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			fmt.Printf("地址 %p, 值 %v\n", &j, j)
			//println(j)
		}()
	}
}

func main() {
	//DeferClosureLoopV1()
	println("\n")
	//DeferClosureLoopV2()
	println("\n")
	DeferClosureLoopV3()

}
