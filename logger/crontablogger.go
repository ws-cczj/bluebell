package logger

import (
	"go.uber.org/zap"
)

type Log struct {
	*zap.Logger
}

func (l Log) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Info(msg, zapKVHandler(keysAndValues)...)
}

func (l Log) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error", err)
	l.Logger.Error(msg, zapKVHandler(keysAndValues)...)
}

//// formatTimes formats any time.Time values as "2006-01-02 15:04:05".
//func formatTimes(keysAndValues []interface{}) []interface{} {
//	var formattedArgs []interface{}
//	for _, arg := range keysAndValues {
//		if t, ok := arg.(time.Time); ok {
//			arg = t.Format("2006-01-02 15:04:05")
//		}
//		formattedArgs = append(formattedArgs, arg)
//	}
//	return formattedArgs
//}

// zapKVHandler 将key 和 values 处理为zapAny
func zapKVHandler(keysAndValues ...interface{}) []zap.Field {
	var kvs []zap.Field
	keyVals := len(keysAndValues)
	// 如果长度为奇数，说明不匹配，这里做一下处理
	if keyVals%2 != 0 {
		keyVals--
	}
	for i := 0; i < keyVals; i += 2 {
		kv := zap.Any(keysAndValues[i].(string), keysAndValues[i+1])
		kvs = append(kvs, kv)
	}
	return kvs
}
