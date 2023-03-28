package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/golang/glog"
	"os"
	"sync"
)

var node *snowflake.Node
var once sync.Once

// GetSnowflake 获取雪花算法ID
func GetSnowflake() int64 {
	if nil == node {
		once.Do(func() {
			var err error
			var id int64
			//由机器名称生成不同的nodeID
			h, _ := os.Hostname()
			for _, v := range h {
				id += int64(v)
			}
			node, err = snowflake.NewNode(id % 1024)
			if nil != err {
				glog.Errorln(err, "snowflake")
			}
		})
	}
	return node.Generate().Int64()
}
