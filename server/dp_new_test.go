package server

import "testing"

func Test_pack(t *testing.T) {
	items := []int{9, 1, 2, 3}
	weights := []int{4, 4, 3, 2}
	values := []int{2, 2, 3, 4}
	//values := []int{2, 2, 3, 4}
	Pack(items, values, weights, 12)
}
