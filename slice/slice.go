package main

func MakeMap() {
	s := make(map[string]int)

	s["a"] = 1
	s["b"] = 2

	for idx, val := range s {
		println(idx, "   ", val)
	}
}
