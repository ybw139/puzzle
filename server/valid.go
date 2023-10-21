package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func valid(c *gin.Context) {
	var req *ToolStruct
	ret := make(map[string]interface{})
	if err := c.ShouldBindJSON(&req); err != nil {
		ret["results"] = []struct{}{}
		ret["errMsg"] = "请求参数错误:" + err.Error()
		c.JSON(http.StatusOK, ret)
		return
	}
	if req.Container == nil || len(req.Objects) == 0 {
		ret["results"] = []struct{}{}
		c.JSON(http.StatusOK, ret)
		return
	}

	for _, v := range req.Objects {
		v.Shape = handleObjShape(v.Shape)
	}

	one := req.Results[0]
	nameMap := map[string]bool{}
	for _, v := range one {
		for _, v1 := range v {
			if v1 != "Z" && v1 != "" {
				nameMap[v1] = true
			}
		}
	}

	nameShapeMap := map[string]*ToolStruct_sub2{}
	indexShapeId := map[string]int{} // index id 映射
	for k, v := range req.Objects {
		nameShapeMap[v.Index] = v
		indexShapeId[v.Index] = k + 1
	}
	objectWeight := 0
	for k := range nameMap {
		objectWeight += nameShapeMap[k].Weight

	}
	// 物品数量
	ret["objectCount"] = len(nameMap)
	ret["objectWeight"] = objectWeight
	ret["objectCheck"] = true
	ret["blockCheck"] = true

	for _, v := range req.Container.Blocks {
		if one[v[0]][v[1]] != "Z" {
			ret["blockCheck"] = false
		}
	}

	//
	// 查找1个值
	maxValue, _ := DpParse4(req)
	if maxValue != objectWeight {
		ret["objectCheck"] = false
		c.JSON(http.StatusOK, ret)
		return
	} else {
		// 每个图形是不是合法的图形
		// 还原图形
		var r = map[string]Shape{}
		for index := range nameMap {
			w, h, s := genShape(one, index, indexShapeId[index])
			r[index] = NewShape(h, w, s)
		}
		var puzzles = map[string]*Puzzle{}
		//req.puzzles = make([]Puzzle, len(req.Objects))
		// 图形编号
		for _, v := range req.Objects {
			//shapeIdName[k+1] = v.Index
			h := len(v.Shape)    // height
			w := len(v.Shape[0]) // 宽
			puzzles[v.Index] = &Puzzle{}
			puzzles[v.Index].InitShape(NewShape(h, w, v.Shape), v.Index)

		}
		for index, _ := range nameMap {
			tempShape := r[index]
			p := puzzles[index]
			shapeNum := *p.ShapeNum
			succ := false
			for j := 0; j < shapeNum; j++ {
				if tempShape.Equal(p.allShapes[j]) {
					fmt.Println("对比成功=========", index)
					succ = true
				}
			}
			if !succ {
				ret["objectCheck"] = false
				c.JSON(http.StatusOK, ret)
				return
			}
		}
	}
	c.JSON(http.StatusOK, ret)
	return

}

func genShape(l [][]string, index string, id int) (int, int, [][]int) {
	m := make([][]int, len(l))
	var x []int
	var y []int
	for k, v := range l {
		m[k] = make([]int, len(v))
		for k1, v1 := range v {
			if v1 == index {
				m[k][k1] = id
				y = append(y, k)
				x = append(x, k1)
				fmt.Println(k, k1)
			}
		}
	}
	sort.Sort(sort.IntSlice(x))
	sort.Sort(sort.IntSlice(y))
	minX := x[0]
	maxX := x[len(x)-1]
	minY := y[0]
	maxY := y[len(y)-1]
	width := maxX - minX + 1
	height := maxY - minY + 1

	var point = make([][]int, height)
	indexy := minY
	for i := 0; i < height; i++ {
		point[i] = m[indexy][minX : minX+width]
		indexy++
	}
	return width, height, point
}

func handleObjShape(l [][]int) [][]int {
	var x []int
	var y []int
	for k, v := range l {
		//m[k] = make([]int, len(v))
		for k1, v1 := range v {
			if v1 == 1 {
				//m[k][k1] = 1
				y = append(y, k)
				x = append(x, k1)
				fmt.Println(k, k1)
			}
		}
	}
	sort.Sort(sort.IntSlice(x))
	sort.Sort(sort.IntSlice(y))
	minX := x[0]
	maxX := x[len(x)-1]
	minY := y[0]
	maxY := y[len(y)-1]
	width := maxX - minX + 1
	height := maxY - minY + 1
	fmt.Println(minX, maxX, minY, maxY, width, height)
	if len(l) != height || len(l[0]) != width {
		var point = make([][]int, height)
		indexy := minY
		for i := 0; i < height; i++ {
			point[i] = l[indexy][minX : minX+width]
			indexy++
		}
		return point
	}
	return l
}
