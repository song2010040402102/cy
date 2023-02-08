package clp

import (
//"fmt"
)

const (
	S_TRUE  = "true"
	S_FALSE = "false"
)

const (
	S_NOT = "!"
	S_LDI = "~"
	S_ADS = "++"
	S_SBS = "--"
	S_MUL = "*"
	S_DIV = "/"
	S_MOD = "%"
	S_ADD = "+"
	S_SUB = "-"
	S_SFL = "<<"
	S_SFR = ">>"
	S_GE  = ">="
	S_GT  = ">"
	S_LE  = "<="
	S_LT  = "<"
	S_EQ  = "=="
	S_NE  = "!="
	S_BND = "&"
	S_BXR = "^"
	S_BOR = "|"
	S_AND = "&&"
	S_OR  = "||"
	S_ASN = "="
	S_MLN = "*="
	S_DVN = "/="
	S_MDN = "%="
	S_ADN = "+="
	S_SBN = "-="
	S_SLN = "<<="
	S_SRN = ">>="
	S_BDN = "&="
	S_BXN = "^="
	S_BON = "|="
	S_CMA = ","
	S_IF  = "?"
	S_SEL = ":"
)

const (
	OP_NONE int = iota
	OP_NOT
	OP_NEG
	OP_LDI
	OP_ADS1
	OP_ADS2
	OP_SBS1
	OP_SBS2
	OP_MUL
	OP_DIV
	OP_MOD
	OP_ADD
	OP_SUB
	OP_SFL
	OP_SFR
	OP_GE
	OP_GT
	OP_LE
	OP_LT
	OP_EQ
	OP_NE
	OP_BND
	OP_BXR
	OP_BOR
	OP_AND
	OP_OR
	OP_ASN
	OP_MLN
	OP_DVN
	OP_MDN
	OP_ADN
	OP_SBN
	OP_SLN
	OP_SRN
	OP_BDN
	OP_BXN
	OP_BON
	OP_CMA
	OP_IF
	OP_SEL
	OP_FUNC = 100
	OP_MTHD = 1000
)

var g_strOpera []string = []string{S_LDI, S_ADS, S_SBS, S_MLN, S_DVN, S_MDN, S_ADN, S_SBN, S_MUL, S_DIV, S_MOD, S_ADD, S_SUB,
	S_SLN, S_SRN, S_SFL, S_SFR, S_GE, S_GT, S_LE, S_LT, S_EQ, S_NE, S_AND, S_OR, S_NOT, S_BDN, S_BXN,
	S_BON, S_BND, S_BXR, S_BOR, S_ASN, S_CMA, S_IF, S_SEL}

var g_mapOpera map[string]int = map[string]int{
	S_NOT: OP_NOT,
	S_LDI: OP_LDI,
	S_ADS: OP_ADS1,
	S_SBS: OP_SBS1,
	S_MUL: OP_MUL,
	S_DIV: OP_DIV,
	S_MOD: OP_MOD,
	S_ADD: OP_ADD,
	S_SUB: OP_SUB,
	S_SFL: OP_SFL,
	S_SFR: OP_SFR,
	S_GE:  OP_GE,
	S_GT:  OP_GT,
	S_LE:  OP_LE,
	S_LT:  OP_LT,
	S_EQ:  OP_EQ,
	S_NE:  OP_NE,
	S_BND: OP_BND,
	S_BXR: OP_BXR,
	S_BOR: OP_BOR,
	S_AND: OP_AND,
	S_OR:  OP_OR,
	S_ASN: OP_ASN,
	S_MLN: OP_MLN,
	S_DVN: OP_DVN,
	S_MDN: OP_MDN,
	S_ADN: OP_ADN,
	S_SBN: OP_SBN,
	S_SLN: OP_SLN,
	S_SRN: OP_SRN,
	S_BDN: OP_BDN,
	S_BXN: OP_BXN,
	S_BON: OP_BON,
	S_CMA: OP_CMA,
	S_IF:  OP_IF,
	S_SEL: OP_SEL,
}

var g_mapPair map[int][]int = map[int][]int{ //多元运算符配对关系
	OP_IF: []int{OP_SEL},
}

var g_opPrior []int = []int{0, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 3, 3, 4, 4, 5, 5, 5, 5, 6, 6, 7, 8, 9, 10, 11, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 14, 12, 12} //运算符优先级

var g_opCombine []bool = []bool{false, false, false, false, false, true, false, true} //一元运算符结合性，true左结合，false右结合

func OP_UNARY(op int) int { //获取运算符元数
	if op >= OP_NOT && op <= OP_SBS2 {
		return 1
	} else if op >= OP_MUL && op <= OP_CMA {
		return 2
	} else if op >= OP_IF && op <= OP_SEL {
		return 3
	}
	return 0
}

func OP_CONVERT(op int) int { //运算符同词转换
	if op == OP_SUB {
		return OP_NEG
	} else if op == OP_ADS1 {
		return OP_ADS2
	} else if op == OP_SBS1 {
		return OP_SBS2
	} else {
		return op
	}
}

type IOperator interface {
	Run()
}

type OperNot struct {
	a, b uint32
}

func (o *OperNot) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vb.(*ValBool).val = !va.(*ValBool).val
}

type OperNeg struct {
	a, b uint32
}

func (o *OperNeg) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	if va.GetType() == VT_INT {
		vb.(*ValInt).val = -va.(*ValInt).val
	} else {
		vb.(*ValFloat).val = -va.(*ValFloat).val
	}
}

type OperLDI struct {
	a, b uint32
}

func (o *OperLDI) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vb.(*ValInt).val = ^va.(*ValInt).val
}

type OperADS1 struct {
	a, b uint32
}

func (o *OperADS1) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	va.(*ValInt).val++
	vb.(*ValInt).val = va.(*ValInt).val
}

type OperADS2 struct {
	a, b uint32
}

func (o *OperADS2) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vb.(*ValInt).val = va.(*ValInt).val
	va.(*ValInt).val++
}

type OperSBS1 struct {
	a, b uint32
}

func (o *OperSBS1) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	va.(*ValInt).val--
	vb.(*ValInt).val = va.(*ValInt).val
}

type OperSBS2 struct {
	a, b uint32
}

func (o *OperSBS2) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vb.(*ValInt).val = va.(*ValInt).val
	va.(*ValInt).val--
}

type OperMul struct {
	a, b, c uint32
}

func (o *OperMul) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValInt).val = va.(*ValInt).val * vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValFloat).val = float64(va.(*ValInt).val) * vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValFloat).val = va.(*ValFloat).val * float64(vb.(*ValInt).val)
	} else {
		vc.(*ValFloat).val = va.(*ValFloat).val * vb.(*ValFloat).val
	}
}

type OperDiv struct {
	a, b, c uint32
}

func (o *OperDiv) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValInt).val = va.(*ValInt).val / vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValFloat).val = float64(va.(*ValInt).val) / vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValFloat).val = va.(*ValFloat).val / float64(vb.(*ValInt).val)
	} else {
		vc.(*ValFloat).val = va.(*ValFloat).val / vb.(*ValFloat).val
	}
}

type OperMod struct {
	a, b, c uint32
}

func (o *OperMod) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValInt).val = va.(*ValInt).val % vb.(*ValInt).val
}

type OperAdd struct {
	a, b, c uint32
}

func (o *OperAdd) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_STRING && vb.GetType() == VT_STRING {
		vc.(*ValString).val = va.(*ValString).val + vb.(*ValString).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValInt).val = va.(*ValInt).val + vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValFloat).val = float64(va.(*ValInt).val) + vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValFloat).val = va.(*ValFloat).val + float64(vb.(*ValInt).val)
	} else {
		vc.(*ValFloat).val = va.(*ValFloat).val + vb.(*ValFloat).val
	}
}

type OperSub struct {
	a, b, c uint32
}

func (o *OperSub) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValInt).val = va.(*ValInt).val - vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValFloat).val = float64(va.(*ValInt).val) - vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValFloat).val = va.(*ValFloat).val - float64(vb.(*ValInt).val)
	} else {
		vc.(*ValFloat).val = va.(*ValFloat).val - vb.(*ValFloat).val
	}
}

type OperSFL struct {
	a, b, c uint32
}

func (o *OperSFL) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValInt).val = va.(*ValInt).val << uint64(vb.(*ValInt).val)
}

type OperSFR struct {
	a, b, c uint32
}

func (o *OperSFR) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValInt).val = va.(*ValInt).val >> uint64(vb.(*ValInt).val)
}

type OperEQ struct {
	a, b, c uint32
}

func (o *OperEQ) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_BOOL {
		vc.(*ValBool).val = va.(*ValBool).val == vb.(*ValBool).val
	} else if va.GetType() == VT_CHAR {
		vc.(*ValBool).val = va.(*ValChar).val == vb.(*ValChar).val
	} else if va.GetType() == VT_STRING {
		vc.(*ValBool).val = va.(*ValString).val == vb.(*ValString).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValInt).val == vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValBool).val = float64(va.(*ValInt).val) == vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValFloat).val == float64(vb.(*ValInt).val)
	} else {
		vc.(*ValBool).val = va.(*ValFloat).val == vb.(*ValFloat).val
	}
}

type OperGE struct {
	a, b, c uint32
}

func (o *OperGE) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_CHAR {
		vc.(*ValBool).val = va.(*ValChar).val >= vb.(*ValChar).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValInt).val >= vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValBool).val = float64(va.(*ValInt).val) >= vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValFloat).val >= float64(vb.(*ValInt).val)
	} else {
		vc.(*ValBool).val = va.(*ValFloat).val >= vb.(*ValFloat).val
	}
}

type OperGT struct {
	a, b, c uint32
}

func (o *OperGT) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_CHAR {
		vc.(*ValBool).val = va.(*ValChar).val > vb.(*ValChar).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValInt).val > vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValBool).val = float64(va.(*ValInt).val) > vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValFloat).val > float64(vb.(*ValInt).val)
	} else {
		vc.(*ValBool).val = va.(*ValFloat).val > vb.(*ValFloat).val
	}
}

type OperLE struct {
	a, b, c uint32
}

func (o *OperLE) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_CHAR {
		vc.(*ValBool).val = va.(*ValChar).val <= vb.(*ValChar).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValInt).val <= vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValBool).val = float64(va.(*ValInt).val) <= vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValFloat).val <= float64(vb.(*ValInt).val)
	} else {
		vc.(*ValBool).val = va.(*ValFloat).val <= vb.(*ValFloat).val
	}
}

type OperLT struct {
	a, b, c uint32
}

func (o *OperLT) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_CHAR {
		vc.(*ValBool).val = va.(*ValChar).val < vb.(*ValChar).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValInt).val < vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValBool).val = float64(va.(*ValInt).val) < vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValFloat).val < float64(vb.(*ValInt).val)
	} else {
		vc.(*ValBool).val = va.(*ValFloat).val < vb.(*ValFloat).val
	}
}

type OperNE struct {
	a, b, c uint32
}

func (o *OperNE) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	if va.GetType() == VT_BOOL {
		vc.(*ValBool).val = va.(*ValBool).val != vb.(*ValBool).val
	} else if va.GetType() == VT_CHAR {
		vc.(*ValBool).val = va.(*ValChar).val != vb.(*ValChar).val
	} else if va.GetType() == VT_STRING {
		vc.(*ValBool).val = va.(*ValString).val != vb.(*ValString).val
	} else if va.GetType() == VT_INT && vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValInt).val != vb.(*ValInt).val
	} else if va.GetType() == VT_INT {
		vc.(*ValBool).val = float64(va.(*ValInt).val) != vb.(*ValFloat).val
	} else if vb.GetType() == VT_INT {
		vc.(*ValBool).val = va.(*ValFloat).val != float64(vb.(*ValInt).val)
	} else {
		vc.(*ValBool).val = va.(*ValFloat).val != vb.(*ValFloat).val
	}
}

type OperBND struct {
	a, b, c uint32
}

func (o *OperBND) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValInt).val = va.(*ValInt).val & vb.(*ValInt).val
}

type OperBXR struct {
	a, b, c uint32
}

func (o *OperBXR) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValInt).val = va.(*ValInt).val ^ vb.(*ValInt).val
}

type OperBOR struct {
	a, b, c uint32
}

func (o *OperBOR) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValInt).val = va.(*ValInt).val | vb.(*ValInt).val
}

type OperAnd struct {
	a, b, c uint32
}

func (o *OperAnd) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValBool).val = va.(*ValBool).val && vb.(*ValBool).val
}

type OperOr struct {
	a, b, c uint32
}

func (o *OperOr) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vc.(*ValBool).val = va.(*ValBool).val || vb.(*ValBool).val
}

type OperASN struct {
	a, b, c uint32
}

func (o *OperASN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	switch vc.GetType() {
	case VT_BOOL:
		va.(*ValBool).val = vb.(*ValBool).val
		vc.(*ValBool).val = va.(*ValBool).val
	case VT_INT:
		va.(*ValInt).val = vb.(*ValInt).val
		vc.(*ValInt).val = va.(*ValInt).val
	case VT_FLOAT:
		if vb.GetType() == VT_INT {
			va.(*ValFloat).val = float64(vb.(*ValInt).val)
		} else {
			va.(*ValFloat).val = vb.(*ValFloat).val
		}
		vc.(*ValFloat).val = va.(*ValFloat).val
	case VT_CHAR:
		va.(*ValChar).val = vb.(*ValChar).val
		vc.(*ValChar).val = va.(*ValChar).val
	case VT_STRING:
		va.(*ValString).val = vb.(*ValString).val
		vc.(*ValString).val = va.(*ValString).val
	}
}

type OperMLN struct {
	a, b, c uint32
}

func (o *OperMLN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	switch vc.GetType() {
	case VT_INT:
		va.(*ValInt).val *= vb.(*ValInt).val
		vc.(*ValInt).val = va.(*ValInt).val
	case VT_FLOAT:
		if vb.GetType() == VT_INT {
			va.(*ValFloat).val *= float64(vb.(*ValInt).val)
		} else {
			va.(*ValFloat).val *= vb.(*ValFloat).val
		}
		vc.(*ValFloat).val = va.(*ValFloat).val
	}
}

type OperDVN struct {
	a, b, c uint32
}

func (o *OperDVN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	switch vc.GetType() {
	case VT_INT:
		va.(*ValInt).val /= vb.(*ValInt).val
		vc.(*ValInt).val = va.(*ValInt).val
	case VT_FLOAT:
		if vb.GetType() == VT_INT {
			va.(*ValFloat).val /= float64(vb.(*ValInt).val)
		} else {
			va.(*ValFloat).val /= vb.(*ValFloat).val
		}
		vc.(*ValFloat).val = va.(*ValFloat).val
	}
}

type OperMDN struct {
	a, b, c uint32
}

func (o *OperMDN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	va.(*ValInt).val %= vb.(*ValInt).val
	vc.(*ValInt).val = va.(*ValInt).val
}

type OperADN struct {
	a, b, c uint32
}

func (o *OperADN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	switch vc.GetType() {
	case VT_INT:
		va.(*ValInt).val += vb.(*ValInt).val
		vc.(*ValInt).val = va.(*ValInt).val
	case VT_FLOAT:
		if vb.GetType() == VT_INT {
			va.(*ValFloat).val += float64(vb.(*ValInt).val)
		} else {
			va.(*ValFloat).val += vb.(*ValFloat).val
		}
		vc.(*ValFloat).val = va.(*ValFloat).val
	case VT_STRING:
		va.(*ValString).val += vb.(*ValString).val
		vc.(*ValString).val = va.(*ValString).val
	}
}

type OperSBN struct {
	a, b, c uint32
}

func (o *OperSBN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	switch vc.GetType() {
	case VT_INT:
		va.(*ValInt).val -= vb.(*ValInt).val
		vc.(*ValInt).val = va.(*ValInt).val
	case VT_FLOAT:
		if vb.GetType() == VT_INT {
			va.(*ValFloat).val -= float64(vb.(*ValInt).val)
		} else {
			va.(*ValFloat).val -= vb.(*ValFloat).val
		}
		vc.(*ValFloat).val = va.(*ValFloat).val
	}
}

type OperSLN struct {
	a, b, c uint32
}

func (o *OperSLN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	va.(*ValInt).val <<= uint64(vb.(*ValInt).val)
	vc.(*ValInt).val = va.(*ValInt).val
}

type OperSRN struct {
	a, b, c uint32
}

func (o *OperSRN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	va.(*ValInt).val >>= uint64(vb.(*ValInt).val)
	vc.(*ValInt).val = va.(*ValInt).val
}

type OperBDN struct {
	a, b, c uint32
}

func (o *OperBDN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	va.(*ValInt).val &= vb.(*ValInt).val
	vc.(*ValInt).val = va.(*ValInt).val
}

type OperBXN struct {
	a, b, c uint32
}

func (o *OperBXN) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	va.(*ValInt).val ^= vb.(*ValInt).val
	vc.(*ValInt).val = va.(*ValInt).val
}

type OperBON struct {
	a, b, c uint32
}

func (o *OperBON) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	va.(*ValInt).val |= vb.(*ValInt).val
	vc.(*ValInt).val = va.(*ValInt).val
}

type OperCMA struct {
	a, b, c uint32
}

func (o *OperCMA) Run() {
	SetData(o.c, GetData(o.a))
}

type OperIF struct {
	a, b, c, d uint32
}

func (o *OperIF) Run() {
	va := GetData(o.a)
	vb := GetData(o.b)
	vc := GetData(o.c)
	vd := GetData(o.d)
	if va.(*ValBool).val {
		switch vd.GetType() {
		case VT_BOOL:
			vd.(*ValBool).val = vb.(*ValBool).val
		case VT_INT:
			vd.(*ValInt).val = vb.(*ValInt).val
		case VT_FLOAT:
			if vb.GetType() == VT_INT {
				vd.(*ValFloat).val = float64(vb.(*ValInt).val)
			} else {
				vd.(*ValFloat).val = vb.(*ValFloat).val
			}
		case VT_CHAR:
			vd.(*ValChar).val = vb.(*ValChar).val
		case VT_STRING:
			vd.(*ValString).val = vb.(*ValString).val
		}
	} else {
		switch vd.GetType() {
		case VT_BOOL:
			vd.(*ValBool).val = vc.(*ValBool).val
		case VT_INT:
			vd.(*ValInt).val = vc.(*ValInt).val
		case VT_FLOAT:
			if vc.GetType() == VT_INT {
				vd.(*ValFloat).val = float64(vc.(*ValInt).val)
			} else {
				vd.(*ValFloat).val = vc.(*ValFloat).val
			}
		case VT_CHAR:
			vd.(*ValChar).val = vc.(*ValChar).val
		case VT_STRING:
			vd.(*ValString).val = vc.(*ValString).val
		}
	}
}

type OperFunc struct {
	f    uint32
	addr []uint32
}

func (o *OperFunc) Run() {
	paras := make([]IValue, len(o.addr)-1)
	for i := 0; i < len(o.addr)-1; i++ {
		paras[i] = GetData(o.addr[i])
	}
	SetData(o.addr[len(o.addr)-1], g_funcs[o.f].sysF(paras...))
}

type OperMthd struct {
	f    uint32
	addr []uint32
}

func (o *OperMthd) Run() {
	preEBP := g_regEBP
	preEAX := g_regEAX
	prePC := g_regPC
	if len(o.addr) > 0 {
		g_regEAX = o.addr[len(o.addr)-1] + g_regEBP
	}
	block := g_blocks[g_funcs[o.f].indexF]
	initStackData(block.symbTbl)
	for i := 0; i < len(o.addr)-1; i++ {
		var fp, ap IValue
		fp = GetData(STACK_OFFSET + uint32(i))
		if o.addr[i] < STACK_OFFSET {
			ap = GetData(o.addr[i])
		} else {
			ap = g_dataS[g_curStackId][o.addr[i]-STACK_OFFSET+preEBP]
		}
		MoveValue(fp, ap)
	}
	block.Exec()
	freeStackData()
	g_regEBP = preEBP
	g_regEAX = preEAX
	g_regPC = prePC
}
