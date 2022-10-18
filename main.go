package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}
	nums = append(nums[:2], nums[3:]...)

	fmt.Println(nums
}