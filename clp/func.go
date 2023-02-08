package clp

import (
	"fmt"
	"math"
	"time"
)

const (
	S_FUNC    = "func"
	S_PRINT   = "print"
	S_POW     = "pow"
	S_TIME    = "time"
	S_TIME2TS = "time2ts"
	S_TS2TIME = "ts2time"
	S_DATE2TS = "date2ts"
	S_TS2DATE = "ts2date"
	S_MAIN    = "main"
)

type BuiltinFunc func(...IValue) IValue

type Function struct {
	name   string      //函数名称
	ptype  []int       //函数参数类型
	ftype  int         //函数返回类型
	sysF   BuiltinFunc //内建函数指针
	indexF int         //函数块索引
}

//函数表
var g_funcs []*Function = []*Function{
	&Function{S_PRINT, nil, VT_NONE, SysPrint, 0},
	&Function{S_POW, []int{VT_FLOAT, VT_FLOAT}, VT_FLOAT, SysPow, 0},
	&Function{S_TIME, nil, VT_INT, SysTime, 0},
	&Function{S_TIME2TS, []int{VT_STRING}, VT_INT, SysTime2TS, 0},
	&Function{S_TS2TIME, []int{VT_INT}, VT_STRING, SysTS2Time, 0},
	&Function{S_DATE2TS, []int{VT_STRING}, VT_INT, SysDate2TS, 0},
	&Function{S_TS2DATE, []int{VT_INT}, VT_STRING, SysTS2Date, 0},
}

func FindFuncIndex(name string) (int, bool) {
	for i, v := range g_funcs {
		if v.name == name {
			if v.sysF == nil {
				return i, false
			}
			return i, true
		}
	}
	return -1, false
}

func GetSysFuncNum() int {
	for i, v := range g_funcs {
		if v.sysF == nil {
			return i
		}
	}
	return 0
}

func SysPrint(values ...IValue) IValue {
	iv := []interface{}{}
	for _, v := range values {
		iv = append(iv, v.GetValue())
	}
	fmt.Println(iv...)
	return &ValNone{}
}

func SysPow(values ...IValue) IValue {
	if len(values) != 2 {
		return &ValFloat{0}
	}
	return &ValFloat{math.Pow(values[0].GetValue().(float64), values[1].GetValue().(float64))}
}

func SysTime(values ...IValue) IValue {
	return &ValInt{time.Now().UnixNano()}
}

func SysTime2TS(values ...IValue) IValue {
	strTime := values[0].GetValue().(string)
	fmtTime, err := time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
	if err != nil {
		return &ValInt{0}
	}
	return &ValInt{fmtTime.Unix()}
}

func SysTS2Time(values ...IValue) IValue {
	ts := values[0].GetValue().(int64)
	return &ValString{time.Unix(ts, 0).Format("2006-01-02 15:04:05")}
}

func SysDate2TS(values ...IValue) IValue {
	strTime := values[0].GetValue().(string)
	fmtTime, err := time.ParseInLocation("2006-01-02", strTime, time.Local)
	if err != nil {
		return &ValInt{0}
	}
	return &ValInt{fmtTime.Unix()}
}

func SysTS2Date(values ...IValue) IValue {
	ts := values[0].GetValue().(int64)
	return &ValString{time.Unix(ts, 0).Format("2006-01-02")}
}
