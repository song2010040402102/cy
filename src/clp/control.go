package clp

import (
//"fmt"
)

const (
	S_CTL_IF       = "if"
	S_CTL_ELSEIF   = "else if"
	S_CTL_ELSE     = "else"
	S_CTL_WHILE    = "while"
	S_CTL_RETURN   = "return"
	S_CTL_BREAK    = "break"
	S_CTL_CONTINUE = "continue"
)

const (
	CTL_RETURN   uint32 = 0xffffffff
	CTL_BREAK    uint32 = 0xfffffffe
	CTL_CONTINUE uint32 = 0xfffffffd
)

type OperCtlIf struct {
	a, b, c uint32
}

func (o *OperCtlIf) Run() {
	va := GetData(o.a)
	if va.(*ValBool).val {
		g_blocks[o.b].Exec()
		if g_regPC < CTL_CONTINUE {
			g_regPC = o.c
		}
	}
}

type OperCtlIfElse struct {
	a, b, c uint32
}

func (o *OperCtlIfElse) Run() {
	prePC := g_regPC
	va := GetData(o.a)
	if va.(*ValBool).val {
		g_blocks[o.b].Exec()
	} else {
		g_blocks[o.c].Exec()
	}
	if g_regPC < CTL_CONTINUE {
		g_regPC = prePC
	}
}

type OperCtlWhile struct {
	a, b, c uint32
}

func (o *OperCtlWhile) Run() {
	va := GetData(o.a)
	if va.(*ValBool).val {
		prePC := g_regPC
		g_blocks[o.b].Exec()
		if g_regPC < CTL_BREAK {
			g_regPC = o.c
		} else if g_regPC == CTL_BREAK {
			g_regPC = prePC
		}
	}
}

type OperCtlReturn struct {
}

func (o *OperCtlReturn) Run() {
	g_regPC = CTL_RETURN
}

type OperCtlReturn2 struct {
	a uint32
}

func (o *OperCtlReturn2) Run() {
	MoveValue(g_dataS[g_curStackId][g_regEAX-STACK_OFFSET], GetData(o.a))
	g_regPC = CTL_RETURN
}

type OperCtlBreak struct {
}

func (o *OperCtlBreak) Run() {
	g_regPC = CTL_BREAK
}

type OperCtlContinue struct {
}

func (o *OperCtlContinue) Run() {
	g_regPC = CTL_CONTINUE
}
