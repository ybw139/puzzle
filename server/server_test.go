package server

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

func Test_pack(t *testing.T) {
	items := []int{9, 1, 2, 3}
	weights := []int{4, 4, 3, 2}
	values := []int{2, 2, 3, 4}
	//values := []int{2, 2, 3, 4}
	Pack(items, values, weights, 12)
}

func Test_All(t *testing.T) {
	var req *ToolStruct
	err := json.Unmarshal([]byte(jsonObj), &req)
	if err != nil {
		fmt.Println("=====err", err)
	}
	DpParse3(req)
}

var jsonObj = `{
   "container": {
       "row": 4,
       "column": 4,
       "blocks": [

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

func Test_json(t *testing.T) {
	var req map[string][][][]string
	err := json.Unmarshal([]byte(resultJson), &req)
	if err != nil {
		fmt.Println("=====err", err)
	}
	fmt.Println(req["results"][0])

	start := time.Now()
	l := req["results"][0]
	m := make([][]int, len(l))
	var x []int
	var y []int
	for k, v := range l {
		m[k] = make([]int, len(v))
		for k1, v1 := range v {
			if v1 == "C" {
				m[k][k1] = 1
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
	fmt.Println(point)
	fmt.Println(time.Since(start))
}

func Test_obj_shape(t *testing.T) {
	l := [][]int{{0, 1}, {0, 1}}
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
	var point = make([][]int, height)
	indexy := minY
	for i := 0; i < height; i++ {
		point[i] = l[indexy][minX : minX+width]
		indexy++
	}
	fmt.Println(point)
}

func Test_abc(t *testing.T) {
	type Shape struct {
		Height  int
		Width   int
		MyShape *[][]int
	}
	a := Shape{1, 2, &[][]int{{1, 2}}}
	b := Shape{1, 2, &[][]int{{1, 2}}}
	fmt.Println(reflect.DeepEqual(a, b)) //false

}

var resultJson = `{"results":[[["A","A","Z","Z"],["A","A","C","Z"],["D","C","C","Z"],["D","Z","Z","Z"]],[["A","A","Z","Z"],["A","A","C","Z"],["Z","C","C","Z"],["D","D","Z","Z"]],[["A","A","Z","Z"],["A","A","C","Z"],["Z","C","C","Z"],["Z","D","D","Z"]],[["A","A","D","Z"],["A","A","D","Z"],["C","Z","Z","Z"],["C","C","Z","Z"]],[["A","A","D","Z"],["A","A","D","Z"],["C","Z","Z","Z"],["C","C","Z","Z"]],[["A","A","Z","Z"],["A","A","D","Z"],["C","Z","D","Z"],["C","C","Z","Z"]],[["B","B","C","Z"],["B","C","C","Z"],["B","D","D","Z"],["Z","Z","Z","Z"]],[["B","B","C","Z"],["B","C","C","Z"],["B","D","Z","Z"],["Z","D","Z","Z"]],[["B","B","C","Z"],["B","C","C","Z"],["B","Z","Z","Z"],["D","D","Z","Z"]],[["B","B","C","Z"],["B","C","C","Z"],["B","Z","D","Z"],["Z","Z","D","Z"]],[["B","B","C","Z"],["B","C","C","Z"],["B","D","Z","Z"],["Z","D","Z","Z"]],[["B","B","C","Z"],["B","C","C","Z"],["B","D","D","Z"],["Z","Z","Z","Z"]]]}`
