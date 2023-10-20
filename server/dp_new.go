package server

import (
	"fmt"
	"sort"
)

// 背包问题的解结构体
type Solution struct {
	items     []int
	weights   []int
	values    []int
	selected  []bool
	maxWeight int
	res       ResultList
}

type result struct {
	value    int
	selected []bool
}

type ResultList []*result

func (r ResultList) Len() int {
	return len(r)
}
func (r ResultList) Less(i, j int) bool {
	return r[i].value < r[j].value
}
func (r ResultList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// 初始化背包问题的解构体
func NewSolution(items []int, weights []int, values []int, maxWeight int) *Solution {
	selected := make([]bool, len(items))
	return &Solution{items, weights, values, selected, maxWeight, nil}
}

// 输出可满足解
func (s *Solution) PrintSolutions() {
	s.printHelper(len(s.items)-1, 0, 0)
}

// 递归辅助函数输出可满足解
func (s *Solution) printHelper(itemIdx, Weight, currValue int) { // 判断是否达到终止条件
	if itemIdx < 0 {
		it := &result{
			value: currValue,
		}
		it.selected = append(it.selected, s.selected...)
		//fmt.Printf("Weight: %d, Value %d, Items: [", Weight, currValue)
		//for i, item := range s.items {
		//	if s.selected[i] {
		//		fmt.Printf("%s", item)
		//	}
		//}
		s.res = append(s.res, it)
		//fmt.Printf("]\n")
		return
	}

	currWeight := Weight
	// 不当前物品
	s.printHelper(itemIdx-1, Weight, currValue)

	//x选
	if currWeight+s.weights[itemIdx] <= s.maxWeight {
		s.selected[itemIdx] = true
		s.printHelper(itemIdx-1, currWeight+s.weights[itemIdx], currValue+s.values[itemIdx])
		s.selected[itemIdx] = false
	}
}

// 背包的最大重量限制

func Pack(items, values, weights []int, maxWeight int) *Solution {
	solution := NewSolution(items, weights, values, maxWeight)
	solution.PrintSolutions()
	sort.Sort(solution.res)
	for _, v := range solution.res {
		fmt.Printf("%v\n", v)
	}
	return solution
}
