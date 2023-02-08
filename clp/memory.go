package clp

import (
//"fmt"
)

const STACK_OFFSET uint32 = 0x10000000 //栈偏移，用于区分全局数据区和栈数据区

var g_regEBP uint32 = 0 //栈底寄存器
var g_regESP uint32 = 0 //栈顶寄存器
var g_regEAX uint32 = 0 //存储返回值寄存器
var g_regPC uint32 = 0  //程序计数器

type Symbol struct {
	name  string //符号名称
	vtype int    //推断类型
}

type SymbolList struct {
	offset uint32 //符号起始偏移量
	count  uint32 //当前符号个数
	parent *SymbolTbl
	childs []*SymbolTbl
}

type SymbolTbl struct {
	symbol *Symbol
	list   *SymbolList
}

func NewSymbolT(name string, vtype int) *SymbolTbl {
	symb := &Symbol{
		name:  name,
		vtype: vtype,
	}
	return &SymbolTbl{
		symbol: symb,
	}
}

func NewSymbolTbl(parent *SymbolTbl) *SymbolTbl {
	offset := uint32(0)
	if parent != nil {
		offset = parent.list.offset + parent.list.count
	}
	list := &SymbolList{
		offset: offset,
		parent: parent,
	}
	return &SymbolTbl{
		list: list,
	}
}

func (st *SymbolTbl) AddSymbol(name string, vtype int) {
	st.list.childs = append(st.list.childs, NewSymbolT(name, vtype))
	st.list.count++
}

func (st *SymbolTbl) AddSymbolTbl(symbT *SymbolTbl) {
	st.list.childs = append(st.list.childs, symbT)
	st.list.count += symbT.list.count
}

var g_dataG []IValue   //全局数据区，为了简化逻辑，把常量区、全局变量区、静态变量区合并
var g_dataS [][]IValue //栈数据区

var g_rootBlock *Block //全局启动块
var g_blocks []*Block  //函数、条件、循环等语句块

var g_curStackId uint32 = 0 //当前栈id

func GetData(addr uint32) IValue {
	if addr < STACK_OFFSET {
		return g_dataG[addr]
	}
	return g_dataS[g_curStackId][addr-STACK_OFFSET+g_regEBP]
}

func SetData(addr uint32, v IValue) {
	if addr < STACK_OFFSET {
		g_dataG[addr] = v
	} else {
		g_dataS[g_curStackId][addr-STACK_OFFSET+g_regEBP] = v
	}
}

func MoveData(addr1, addr2 uint32) {
	MoveValue(GetData(addr1), GetData(addr2))
}

func MoveValue(val1, val2 IValue) {
	switch val1.GetType() {
	case VT_BOOL:
		val1.(*ValBool).val = val2.(*ValBool).val
	case VT_INT:
		val1.(*ValInt).val = val2.(*ValInt).val
	case VT_FLOAT:
		if val2.GetType() == VT_INT {
			val1.(*ValFloat).val = float64(val2.(*ValInt).val)
		} else {
			val1.(*ValFloat).val = val2.(*ValFloat).val
		}
	case VT_CHAR:
		val1.(*ValChar).val = val2.(*ValChar).val
	case VT_STRING:
		val1.(*ValString).val = val2.(*ValString).val
	}
}

//获取当前指令的偏移量
func GetOperaOffset() uint32 {
	if len(g_blocks) == 0 {
		if g_rootBlock != nil {
			return uint32(len(g_rootBlock.opers))
		}
	} else {
		return g_blocks[len(g_blocks)-1].offset + uint32(len(g_blocks[len(g_blocks)-1].opers))
	}
	return 0
}
