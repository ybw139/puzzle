package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	show = true
	port string
)

func Init(isShow bool, p string) {
	show = isShow
	port = p
}

type ToolStruct struct {
	Container *ToolStruct_sub1   `fmt:"container"`
	Objects   []*ToolStruct_sub2 `fmt:"objects"`

	MAP_WIDTH  int
	MAP_HEIGHT int
	PIECE_NUM  int

	originMap   *Map
	puzzles     []Puzzle
	ShapeIdName map[int]string
}

type ToolStruct_sub1 struct {
	Blocks [][]int `fmt:"blocks"`
	Column int     `fmt:"column"`
	Row    int     `fmt:"row"`
}

type ToolStruct_sub2 struct {
	Index  string  `fmt:"index"`
	Shape  [][]int `fmt:"shape"`
	Weight int     `fmt:"weight"`
	Vol    int
}

func InitShape(req *ToolStruct) {
	// 初始化map
	//var shapeIdName = map[int]string{}
	req.puzzles = make([]Puzzle, len(req.Objects))
	// 图形编号
	for k, v := range req.Objects {
		//shapeIdName[k+1] = v.Index
		h := len(v.Shape)    // height
		w := len(v.Shape[0]) // 宽
		req.puzzles[k].InitShape(NewShape(h, w, v.Shape))
	}
	//req.ShapeIdName = shapeIdName
	req.originMap = NewMap(req, true)
	req.originMap.SetWall(req.Container.Blocks)
	return
}

func Run() {
	var r *gin.Engine
	if show {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	r.GET("/calcOne", resolve)
	r.POST("/calcOne", resolve)
	r.GET("/calc", resolveAll)
	r.POST("/calc", resolveAll)
	if err := r.Run(":" + port); err != nil {
		log.Fatalln(err)
	}
}

func resolve(c *gin.Context) {
	var req *ToolStruct
	//if err := c.ShouldBind(req); err != nil {
	//	c.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	err := json.Unmarshal([]byte(jsonObj), &req)
	if err != nil {
		fmt.Println("=====err", err)
	}

	ret := make(map[string]interface{})
	_, _, _, rs := resolveEasy()
	ret["results"] = [][][]string{rs}
	c.JSON(http.StatusOK, ret)
}

func resolveAll(c *gin.Context) {
	var req *ToolStruct
	//if err := c.ShouldBind(req); err != nil {
	//	c.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	err := json.Unmarshal([]byte(jsonObj), &req)
	if err != nil {
		fmt.Println("=====err", err)
	}

	ret := make(map[string]interface{})
	_, _, rs := searchAllRes(DpParse2(req), true, true)
	ret["results"] = rs
	c.JSON(http.StatusOK, ret)
}

//func resolveEasy(month, day int) ([][constant.MAP_WIDTH]int, int64, string) {
func resolveEasy() ([][]int, int64, string, [][]string) {
	var req *ToolStruct
	//if err := c.ShouldBind(req); err != nil {
	//	c.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	err := json.Unmarshal([]byte(jsonObj), &req)
	if err != nil {
		fmt.Println("=====err", err)
	}
	req2 := DpParse2(req)
	myMap := req2.originMap.DeepCopy(req.MAP_WIDTH)
	start := time.Now()
	m, count, rs := searchOneRes(req2, true, myMap, "")
	fmt.Println("=======", m)
	return m, count, time.Since(start).String(), rs
}

func searchOneRes(req *ToolStruct, modeEasy bool, calendar *Map, week string) ([][]int, int64, [][]string) {
	var (
		back       bool
		stackIndex int
		pieceNum   = req.PIECE_NUM
		height     = req.MAP_HEIGHT
		puzs       = req.puzzles
	)
	// 逐一为拼图块选好位置和形状，如果遇到无处安放的块，则回溯
	backCount := 0
	for stackIndex < pieceNum && stackIndex >= 0 {
		//	初始化
		var i, j, k int
		if back {
			backCount++
			//需要回溯，也就是当前拼图需要一个新的位置,要先从旧的位置删除掉
			puzs[stackIndex].Clear(calendar)
			i = *puzs[stackIndex].Y
			j = *puzs[stackIndex].X
			k = puzs[stackIndex].ShapeIndex + 1
		} else {
			i, j, k = 0, 0, 0
		}

		//为stack_index号拼图找一个位置
		success := false
		for ; i < height; i++ {
			for ; j < req.MAP_WIDTH; j++ {
				for ; k < *puzs[stackIndex].ShapeNum; k++ {
					if puzs[stackIndex].Check(calendar, j, i, k, height, req.MAP_WIDTH, modeEasy) {
						success = true
						break
					}
				}
				if success {
					break
				}
				k = 0
			}
			if success {
				break
			}
			j = 0
		}
		if success {
			stackIndex++
			back = false
		} else {
			stackIndex--
			back = true
		}
	}
	fmt.Printf("Down.Total search %d possibilities\n", backCount)
	if show {
		calendar.Show(req, height, week)
	}
	rs := calendar.Result(req, height, week)
	return *calendar, int64(backCount), rs
}

func searchAllRes(req *ToolStruct, modeEasy, inServer bool) ([]*Map, int64, [][][]string) {
	var calendars []*Map
	var (
		myMap      *Map
		back       bool //回溯标志
		stackIndex int  //当前待放置的拼图序号
		backCount  int  // 逐一为拼图块选好位置和形状，如果遇到无处安放的块，则回溯
		resCount   int  //已经找到的解的数量
		pieceNum   = req.PIECE_NUM
		height     = req.MAP_HEIGHT
		puzs       = req.puzzles
	)
	myMap = req.originMap.DeepCopy(req.MAP_WIDTH)

	for {
		for stackIndex < pieceNum && stackIndex >= 0 {
			//	初始化
			var i, j, k int
			if back {
				backCount++
				//需要回溯，也就是当前拼图需要一个新的位置,要先从旧的位置删除掉
				puzs[stackIndex].Clear(myMap)
				i = *puzs[stackIndex].Y
				j = *puzs[stackIndex].X
				k = puzs[stackIndex].ShapeIndex + 1
			} else {
				i, j, k = 0, 0, 0
			}

			//为stack_index号拼图找一个位置
			success := false
			for ; i < height; i++ {
				for ; j < req.MAP_WIDTH; j++ {
					for ; k < *puzs[stackIndex].ShapeNum; k++ {
						if puzs[stackIndex].Check(myMap, j, i, k, height, req.MAP_WIDTH, modeEasy) {
							success = true
							break
						}
					}
					if success {
						break
					}
					k = 0
				}
				if success {
					break
				}
				j = 0
			}
			if success {
				stackIndex++
				back = false
			} else {
				stackIndex--
				back = true
			}
		}
		if stackIndex == pieceNum {
			//循环因为找到解而中断
			//myMap.Show(height, week)
			back = true
			stackIndex--
			resCount++
			calendars = append(calendars, myMap.DeepCopy(req.MAP_WIDTH))
			if inServer && len(calendars) == 1000 {
				break
			}
		} else {
			//循环因为找不到解而中断
			break
		}
	}

	if show {
		showAllRes(req, calendars, height)
		fmt.Printf("Down.Total search %d possibilities\n", backCount)
	}
	rs := AllRes(req, calendars, height)
	return calendars, int64(backCount), rs
}

func showAllRes(req *ToolStruct, calendars []*Map, height int) {
	size := len(calendars)
	for k := 0; k < size; k += 6 {
		for i := 0; i < height; i++ {
			for l := k; l < k+6 && l < size; l++ {
				for j := 0; j < req.MAP_WIDTH; j++ {
					PrintBlock(req, (*calendars[l])[i][j])
				}
				PrintEmpty()
			}
			fmt.Println()
		}
		fmt.Println()
	}
	fmt.Printf("There are %d solutions\n", size)
}

func AllRes(req *ToolStruct, calendars []*Map, height int) [][][]string {
	ret := [][][]string{}
	size := len(calendars)
	for k := 0; k < size; k++ {
		retItem := [][]string{}
		for i := 0; i < height; i++ {
			item := []string{}
			for j := 0; j < req.MAP_WIDTH; j++ {
				if v, ok := req.ShapeIdName[(*calendars[k])[i][j]]; ok {
					item = append(item, v)
				} else {
					item = append(item, "Z")
				}
			}
			retItem = append(retItem, item)
		}
		ret = append(ret, retItem)
	}
	return ret
}

//var jsonObj = `{
//    "container": {
//        "row": 7,
//        "column": 7,
//        "blocks": [
//            [0, 6],
//            [1, 6],
//            [6, 3],
//            [6, 4],
//            [6, 5],
//            [6, 6]
//        ]
//    },
//    "objects": [
//        {
//            "index": "A",
//            "weight": 1,
//            "shape": [
//                [1, 1, 1],
//                [1, 0, 0],
//                [1, 0, 0]
//            ]
//        },
//        {
//            "index": "B",
//            "weight": 1,
//            "shape": [
//                [1, 0],
//                [1, 1],
//                [1, 0],
//                [1, 0]
//            ]
//        },
//        {
//            "index": "C",
//            "weight": 1,
//            "shape": [
//                [1, 1],
//                [1, 1],
//                [1, 1]
//            ]
//        },
//        {
//            "index": "D",
//            "weight": 1,
//            "shape": [
//                [1, 0, 1],
//                [1, 1, 1]
//            ]
//        },
//        {
//            "index": "E",
//            "weight": 1,
//            "shape": [
//                [0, 0, 0, 1],
//                [1, 1, 1, 1]
//            ]
//        },
//        {
//            "index": "F",
//            "weight": 1,
//            "shape": [
//                [0, 0, 1],
//                [1, 1, 1],
//                [1, 0, 0]
//            ]
//        },
//        {
//            "index": "G",
//            "weight": 1,
//            "shape": [
//                [0, 1],
//                [1, 1],
//                [1, 0],
//                [1, 0]
//            ]
//        },
//        {
//            "index": "H",
//            "weight": 1,
//            "shape": [
//                [0, 1],
//                [1, 1],
//                [1, 1]
//            ]
//        }
//    ]
//}`

var jsonObj = `{
   "container": {
       "row": 4,
       "column": 4,
       "blocks": [
[0, 3],
     [1, 3],
           [2, 3],
           [3, 3]
       ]
   },
   "objects": [
       {
           "index": "A",
           "weight": 2,
           "shape": [
               [1, 1],
               [1, 1]
           ]
       },
       {
           "index": "B",
           "weight": 2,
           "shape": [
               [1, 1],
               [1, 0],
               [1, 0]
           ]
       },
       {
           "index": "C",
           "weight": 3,
           "shape": [
               [1, 0],
               [1, 1]
           ]
       },
       {
           "index": "D",
           "weight": 4,
           "shape": [
               [0, 1],
               [0, 1]
           ]
       }
   ]
}`
