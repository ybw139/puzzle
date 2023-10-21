package server

import (
	"fmt"
	"sort"
	"sync"
)

type DpStruct struct {
	weight     []int
	vol        []int
	dp         [][]int
	totalPoint int
	num        int
	object     []int
}

func DpParse2(req *ToolStruct) *ToolStruct {
	//1.计算总共有多少个点
	totalPoint := req.Container.Column*req.Container.Row - len(req.Container.Blocks) // 总共有多少个点

	volList := []int{}
	weightList := []int{}
	indexList := []int{}
	var shapeIdName = map[int]string{}
	//2.记录每个物体有几个点
	for k, v := range req.Objects {
		shapeIdName[k+1] = v.Index
		total := 0
		for _, item := range v.Shape {
			for k1, val := range item {
				if val == 1 {
					total++
					item[k1] = k + 1
				}
			}
		}
		v.Vol = total
		volList = append(volList, total)
		weightList = append(weightList, v.Weight)
		indexList = append(indexList, k)
	}
	req.ShapeIdName = shapeIdName

	req.MAP_WIDTH = req.Container.Column
	req.MAP_HEIGHT = req.Container.Row

	// 物品个数
	num := len(req.Objects)
	sol := Pack(indexList, weightList, volList, totalPoint)

	maxValue := sol.res[len(sol.res)-1].value
	var selected [][]bool
	for i := len(sol.res) - 1; i >= 0; i-- {
		if sol.res[i].value == maxValue {
			maxValue = sol.res[i].value
			selected = append(selected, sol.res[i].selected)
		} else {
			maxValue = sol.res[i].value
			break
		}
	}

	if show {
		fmt.Println("物品数量：", num, "共有多少点：", totalPoint)
		fmt.Println("最大重量", sol.res[len(sol.res)-1].value)
		//fmt.Println("选中的节点", Objects)
		fmt.Println("选中组合", len(selected))
		fmt.Println("下次要求的值", maxValue)
	}

	// 选中了哪些物品
	for _, v := range selected {
		var Objects []*ToolStruct_sub2
		for k1, v1 := range v {
			if v1 {
				Objects = append(Objects, req.Objects[k1])
			}
		}
		req.PIECE_NUM = len(Objects)
		// 4.选出有哪些图形
		req2 := *req
		req2.Objects = Objects
		InitShape(&req2)
		return &req2
	}

	return nil

}

func DpParse3(req *ToolStruct) [][][]string {
	//1.计算总共有多少个点
	totalPoint := req.Container.Column*req.Container.Row - len(req.Container.Blocks) // 总共有多少个点

	volList := []int{}
	weightList := []int{}
	indexList := []int{}
	var shapeIdName = map[int]string{}
	//2.记录每个物体有几个点
	for k, v := range req.Objects {
		shapeIdName[k+1] = v.Index
		total := 0
		for _, item := range v.Shape {
			for k1, val := range item {
				if val == 1 {
					total++
					item[k1] = k + 1
				}
			}
		}
		v.Vol = total
		volList = append(volList, total)
		weightList = append(weightList, v.Weight)
		indexList = append(indexList, k)
	}
	req.ShapeIdName = shapeIdName

	req.MAP_WIDTH = req.Container.Column
	req.MAP_HEIGHT = req.Container.Row

	// 物品个数
	num := len(req.Objects)
	sol := Pack(indexList, weightList, volList, totalPoint)

	var ret [][][]string
	var selected [][]bool
	maxValue := sol.res[len(sol.res)-1].value
	req.lock = &sync.Mutex{}
	for i := len(sol.res) - 1; i >= 0; i-- {
		if sol.res[i].value == maxValue {
			maxValue = sol.res[i].value
			selected = append(selected, sol.res[i].selected)
		} else {
			maxValue = sol.res[i].value
			if show {
				fmt.Println("物品数量：", num, "共有多少点：", totalPoint)
				fmt.Println("最大重量", sol.res[len(sol.res)-1].value)
				//fmt.Println("选中的节点", Objects)
				fmt.Println("选中组合", len(selected))
				fmt.Println("下次要求的值", maxValue)
			}
			wg := &sync.WaitGroup{}
			// 选中了哪些物品
			for _, v := range selected {
				wg.Add(1)
				v0 := v
				go func() {
					defer wg.Done()
					var Objects []*ToolStruct_sub2
					for k1, v1 := range v0 {
						if v1 {
							Objects = append(Objects, req.Objects[k1])
						}
					}
					req.PIECE_NUM = len(Objects)
					// 4.选出有哪些图形
					req2 := *req
					req2.Objects = Objects
					InitShape(&req2)
					_, _, rs := searchAllRes(&req2, true, true)
					if len(rs) > 0 {
						ret = append(ret, rs...)
					}

				}()
			}
			wg.Wait()
			if len(ret) > 0 {
				// 找到了最优解
				return ret
			}
			selected = [][]bool{}
		}
	}
	return nil
}

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
	if show {
		for _, v := range solution.res {
			fmt.Printf("%v\n", v)
		}
	}
	return solution
}
