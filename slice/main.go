package main

func main() {
	//s1 := [7]int{2, 4, 6, 8, 10, 12, 14}
	//s2 := s1[1:3]
	//
	//fmt.Printf("s2: %v len %d cap %d\n", s2, len(s2), cap(s2))

	s1 := []int{1, 2, 3}
	res1 := SumInt(s1)
	println(res1)

	s2 := make([]int64, 0, len(s1))
	for _, num := range s1 {
		s2 = append(s2, int64(num))
	}

	res2 := SumInt64(s2)
	println(res2)

}
