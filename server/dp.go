package server

import (
	"fmt"
)

var (
//weight = [5]int{0, 2, 3, 4, 5}
//value  = [5]int{0, 3, 4, 5, 6}
//dp     = [5][9]int{}
//object = [5]int{}
)

func Dynamic(p *DpStruct) { //动态规划找到背包所能装下的最大价值
	dp := p.dp
	var value []int //价值
	value = p.weight
	var vol []int
	vol = p.vol
	for i := 1; i < p.num; i++ {
		for j := 1; j <= p.totalPoint; j++ {
			if vol[i] > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-vol[i]]+value[i])
			}
		}
	}
	//fmt.Println(p, dp)
	// DP表
	for k := 0; k < p.num; k++ {
		if k == 0 {
			for l := 0; l <= p.totalPoint; l++ {
				fmt.Printf("%d	", l)
			}
			fmt.Println()
			for l := 0; l <= p.totalPoint; l++ {
				fmt.Printf("------")
			}
			fmt.Println()
		}
		for l := 0; l <= p.totalPoint; l++ {
			fmt.Printf("%d	", dp[k][l])
		}
		fmt.Println()
	}
}

func Find(p *DpStruct, i, j int) { //回溯找到到底装了哪些物品
	dp := p.dp
	if i == 0 {
		for _, v := range p.object {
			fmt.Println(v)
		}
		return
	}
	if dp[i][j] == dp[i-1][j] {
		p.object[i] = 0
		Find(p, i-1, j)
	} else if dp[i][j] == dp[i-1][j-p.vol[i]]+p.weight[i] {
		p.object[i] = 1
		Find(p, i-1, j-p.vol[i])
	}

}

type DpStruct struct {
	weight     []int
	vol        []int
	dp         [][]int
	totalPoint int
	num        int
	object     []int
}

func DpParse(req *ToolStruct) *ToolStruct {
	//1.计算总共有多少个点
	totalPoint := req.Container.Column*req.Container.Row - len(req.Container.Blocks) // 总共有多少个点

	volList := []int{0}
	weightList := []int{0}
	//2.记录每个物体有几个点
	for k, v := range req.Objects {
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
	}

	// 物品个数
	num := len(req.Objects)
	// 3.01背包计算最大能放多少
	p := new(DpStruct)
	p.vol = volList
	p.weight = weightList
	p.num = num + 1
	p.totalPoint = totalPoint
	for i := num; i >= 0; i-- {
		p.dp = append(p.dp, make([]int, totalPoint+2))
	}
	p.object = make([]int, p.num)
	Dynamic(p)
	Find(p, num, totalPoint)
	if show {
		fmt.Println("物品数量：", num, "共有多少点：", totalPoint)
		fmt.Println("最大重量", p.dp[num][totalPoint])
		//p.object = make([]int, p.num)
		//
		//Find(p, 4, 12)
		//p.object = make([]int, p.num)
		//
		//fmt.Println("----")
		//Find(p, 4, 11)
		//p.object = make([]int, p.num)
		//
		//fmt.Println("----")
		//Find(p, 4, 10)
		p.object = make([]int, p.num)

		fmt.Println("----")
		Find(p, 4, 9)
		fmt.Println("----")
		p.object = make([]int, p.num)

		Find(p, num, totalPoint)

	}

	// 选中了哪些物品
	var Objects []*ToolStruct_sub2
	for k, v := range p.object {
		if v == 1 {
			Objects = append(Objects, req.Objects[k-1])
		}
	}

	req.MAP_WIDTH = req.Container.Column
	req.MAP_HEIGHT = req.Container.Row
	req.PIECE_NUM = len(Objects)
	// 4.选出有哪些图形
	req2 := *req
	req2.Objects = Objects
	InitShape(&req2)
	return &req2
}

func DpParse2(req *ToolStruct) *ToolStruct {
	//1.计算总共有多少个点
	totalPoint := req.Container.Column*req.Container.Row - len(req.Container.Blocks) // 总共有多少个点

	volList := []int{}
	weightList := []int{}
	indexList := []int{}
	//2.记录每个物体有几个点
	for k, v := range req.Objects {
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

	// 选中了哪些物品
	var Objects []*ToolStruct_sub2
	for k, v := range selected[0] {
		if v {
			Objects = append(Objects, req.Objects[k])
		}
	}
	if show {
		fmt.Println("物品数量：", num, "共有多少点：", totalPoint)
		fmt.Println("最大重量", sol.res[len(sol.res)-1].value)
		fmt.Println("选中的节点", Objects)
		fmt.Println("下次要求的值", maxValue)
	}

	req.MAP_WIDTH = req.Container.Column
	req.MAP_HEIGHT = req.Container.Row
	req.PIECE_NUM = len(Objects)

	var shapeIdName = map[int]string{}
	req.puzzles = make([]Puzzle, len(req.Objects))
	// 图形编号
	for k, v := range req.Objects {
		shapeIdName[k+1] = v.Index
	}
	req.ShapeIdName = shapeIdName

	// 4.选出有哪些图形
	req2 := *req
	req2.Objects = Objects
	InitShape(&req2)
	return &req2
}
