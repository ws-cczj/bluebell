package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// -- 每台机器相当于一个结点
var node *sf.Node

func InitSnowID(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}
