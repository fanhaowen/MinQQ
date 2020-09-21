package main

import "fmt"

func main() {
	array := [5]int{1, 2, 3, 4, 5}
	for index, val := range array {
		fmt.Printf("index array[%d] = %d\n", index, val)
	}
}
