package test

import "fmt"

func testdata() {
	a := []int{1, 2, 3}
	a = append(a, 4)

	b := append(a, 4) // want "should not assign from source variable a to different variable b"
	if len(b) == 0 {
		fmt.Errorf("Unreachable error")
	}
}
