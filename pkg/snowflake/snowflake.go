package snowflake

import (
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// -- 每台机器相当于一个结点
var node *sf.Node

func InitSnowID(startTime string, machineID int64) {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(fmt.Errorf("init SnowFlake fail, err:%s", err))
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	if err != nil {
		panic(fmt.Errorf("init SnowFlake fail, err:%s", err))
	}
}
func GenID() string {
	return node.Generate().String()
}
