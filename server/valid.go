package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	one := req.results[0]
	nameMap := map[string]bool{}
	for _, v := range one {
		for _, v1 := range v {
			if v1 != "Z" && v1 != "" {
				nameMap[v1] = true
			}
		}
	}

	nameShapeMap := map[string]*ToolStruct_sub2{}

	for _, v := range req.Objects {
		nameShapeMap[v.Index] = v
	}
	objectWeight := 0
	for k := range nameMap {
		objectWeight += nameShapeMap[k].Weight
	}
	// 物品数量
	ret["objectCount"] = len(nameMap)
	ret["objectWeight"] = objectWeight
	ret["objectCheck"] = false
	ret["blockCheck"] = true

	for _, v := range req.Container.Blocks {
		if one[v[0]][v[1]] != "Z" {
			ret["blockCheck"] = false
		}
	}

	//
	// 查找1个值
	maxValue := 0 //todo
	if maxValue != objectWeight {
		ret["objectCheck"] = false
		c.JSON(http.StatusOK, ret)
		return
	} else {
		// 每个图形是不是合法的图形
	}

	c.JSON(http.StatusOK, ret)
	return

}
