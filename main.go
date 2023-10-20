package main

import (
	"log"
	"puzzle/server"
)

var (
	showPic = true
	port    = "8888"
)

//var object = []int{1, 2, 3, 4}
//var vol = []int{1, 2, 3, 4}
//var value = []int{2, 3, 4, 5}
//var limitVol int = 13

func main() {
	log.Println("running as a server ...")
	server.Init(showPic, port)
	server.Run()
}

//func test() {
//
//	// 第一层 2中选和不选
//	dfs("0", 0, 1, 0, 0, 0, 0)
//	dfs("0", 1, 1, 0, 0, 0, 0)
//
//}

//// 返回最大值和路径
//func dfs(trace string, chosen int, index int, vol int, currVol int, value int, max int) (string, int, int) {
//	if chosen > 0 {
//
//	} else {
//		return trace + ",-" + strconv.Itoa(chosen), chosen, max
//	}
//	if vol+currVol > limitVol {
//		// 终止
//		// 返回路径和最大值
//		return trace + ",-" + strconv.Itoa(chosen), chosen, max
//	}
//}
