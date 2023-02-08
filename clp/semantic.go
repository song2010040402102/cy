package clp

//import "fmt"

//规范语义表，也可以类似C语言定义错误表，给出警告提示，但不利于代码规范化
var g_mapSemNorm map[int][][]int = map[int][][]int{
	OP_NOT:  [][]int{{VT_BOOL}},
	OP_NEG:  [][]int{{VT_INT}, {VT_FLOAT}},
	OP_LDI:  [][]int{{VT_INT}},
	OP_ADS1: [][]int{{VT_INT}},
	OP_ADS2: [][]int{{VT_INT}},
	OP_SBS1: [][]int{{VT_INT}},
	OP_SBS2: [][]int{{VT_INT}},
	OP_MUL:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}},
	OP_DIV:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}},
	OP_MOD:  [][]int{{VT_INT, VT_INT}},
	OP_ADD:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_STRING, VT_STRING}},
	OP_SUB:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}},
	OP_SFL:  [][]int{{VT_INT, VT_INT}},
	OP_SFR:  [][]int{{VT_INT, VT_INT}},
	OP_GE:   [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}},
	OP_GT:   [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}},
	OP_LE:   [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}},
	OP_LT:   [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}},
	OP_EQ:   [][]int{{VT_BOOL, VT_BOOL}, {VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}, {VT_STRING, VT_STRING}},
	OP_NE:   [][]int{{VT_BOOL, VT_BOOL}, {VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_INT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}, {VT_STRING, VT_STRING}},
	OP_BND:  [][]int{{VT_INT, VT_INT}},
	OP_BXR:  [][]int{{VT_INT, VT_INT}},
	OP_BOR:  [][]int{{VT_INT, VT_INT}},
	OP_AND:  [][]int{{VT_BOOL, VT_BOOL}},
	OP_OR:   [][]int{{VT_BOOL, VT_BOOL}},
	OP_ASN:  [][]int{{VT_BOOL, VT_BOOL}, {VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_CHAR, VT_CHAR}, {VT_STRING, VT_STRING}},
	OP_MLN:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_FLOAT, VT_INT}},
	OP_DVN:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_FLOAT, VT_INT}},
	OP_MDN:  [][]int{{VT_INT, VT_INT}},
	OP_ADN:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_FLOAT, VT_INT}, {VT_STRING, VT_STRING}},
	OP_SBN:  [][]int{{VT_INT, VT_INT}, {VT_FLOAT, VT_FLOAT}, {VT_FLOAT, VT_INT}},
	OP_SLN:  [][]int{{VT_INT, VT_INT}},
	OP_SRN:  [][]int{{VT_INT, VT_INT}},
	OP_BDN:  [][]int{{VT_INT, VT_INT}},
	OP_BXN:  [][]int{{VT_INT, VT_INT}},
	OP_BON:  [][]int{{VT_INT, VT_INT}},
	OP_CMA:  [][]int{},
	OP_IF:   [][]int{{VT_BOOL, VT_BOOL, VT_BOOL}, {VT_BOOL, VT_INT, VT_INT}, {VT_BOOL, VT_FLOAT, VT_FLOAT}, {VT_BOOL, VT_INT, VT_FLOAT}, {VT_BOOL, VT_FLOAT, VT_INT}, {VT_BOOL, VT_CHAR, VT_CHAR}, {VT_BOOL, VT_STRING, VT_STRING}},
}

//判断是否符合规范语义表
func IsSemNorm(op int, vts []int) bool {
	if _, ok := g_mapSemNorm[op]; !ok {
		return false
	}
	norms := g_mapSemNorm[op]
	if len(norms) == 0 {
		return true
	}
	for _, norm := range norms {
		if len(norm) != len(vts) {
			continue
		}
		found := true
		for i := 0; i < len(norm); i++ {
			if norm[i] != vts[i] {
				found = false
				break
			}
		}
		if found {
			return true
		}
	}
	return false
}

//获取语义类型，由低到高转换，例如int和float，则往float转换
func GetSemType(op int, vts []int) int {
	if op == OP_CMA {
		return vts[0]
	}
	if op >= OP_GE && op <= OP_NE {
		return VT_BOOL
	}
	ret := VT_NONE
	for _, v := range vts {
		if ret == VT_NONE || v > ret {
			ret = v
		}
	}
	return ret
}

//控制语句的语义
func IsSemCtl(s string, ctl string) bool {
	if ctl == S_CTL_ELSE {
		return s[len(s)-1] == '{' && len(s) > len(ctl) && s[:len(ctl)] == ctl
	} else if ctl == S_CTL_IF || ctl == S_CTL_ELSEIF || ctl == S_CTL_WHILE {
		return s[len(s)-1] == '{' && len(s) > len(ctl) && s[:len(ctl)] == ctl && (s[len(ctl)] == ' ' || s[len(ctl)] == '	')
	} else if ctl == S_CTL_RETURN {
		return s == ctl || len(s) > len(ctl) && s[:len(ctl)] == ctl && (s[len(ctl)] == ' ' || s[len(ctl)] == '	')
	} else if ctl == S_CTL_CONTINUE || ctl == S_CTL_BREAK {
		return s == ctl
	} else {
		return false
	}
}

//函数声明的语义
func IsSemFunc(s string) bool {
	return s[len(s)-1] == '{' && len(s) > len(S_FUNC) && s[:len(S_FUNC)] == S_FUNC && (s[len(S_FUNC)] == ' ' || s[len(S_FUNC)] == '	')
}
