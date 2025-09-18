package main

func SumInt(nums []int) int {
	var res int
	for _, num := range nums {
		res += num
	}
	return res
}

func SumInt64(nums []int64) int64 {
	var res int64
	for _, num := range nums {
		res += num
	}
	return res
}
