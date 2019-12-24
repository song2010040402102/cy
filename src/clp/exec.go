package clp

import (
	"fmt"
	"runtime/debug"
	"strconv"
)

func Exec() {
	initGlobalData()
	if g_rootBlock != nil {
		g_rootBlock.Exec()
	}
}

func initGlobalData() {
	if g_rootBlock == nil || g_rootBlock.symbTbl == nil || g_rootBlock.symbTbl.list == nil {
		return
	}
	for _, v := range g_rootBlock.symbTbl.list.childs {
		if v != nil {
			g_dataG = append(g_dataG, Symbol2Value(v.symbol))
		}
	}
	g_dataS = make([][]IValue, 1)
}

func initStackData(symbTbl *SymbolTbl) {
	if symbTbl.list.parent == nil {
		g_regEBP = g_regESP
		g_regESP += symbTbl.list.count
	}
	for _, v := range symbTbl.list.childs {
		if v.symbol != nil {
			g_dataS[g_curStackId] = append(g_dataS[g_curStackId], Symbol2Value(v.symbol))
		} else {
			initStackData(v)
		}
	}
}

func freeStackData() {
	g_dataS[g_curStackId] = g_dataS[g_curStackId][:len(g_dataS[g_curStackId])-int(g_regESP-g_regEBP)]
	g_regESP = g_regEBP
}

func Symbol2Value(symbol *Symbol) IValue {
	if symbol == nil {
		return nil
	}
	switch symbol.vtype {
	case VT_NONE:
		return &ValNone{}
	case VT_BOOL:
		if symbol.name == S_TRUE {
			return &ValBool{true}
		}
		return &ValBool{false}
	case VT_INT:
		if v, err := strconv.ParseInt(symbol.name, 10, 32); err == nil {
			return &ValInt{v}
		}
		return &ValInt{0}
	case VT_FLOAT:
		if v, err := strconv.ParseFloat(symbol.name, 32); err == nil {
			return &ValFloat{v}
		}
		return &ValFloat{0}
	case VT_CHAR:
		if symbol.name != "" {
			return &ValChar{symbol.name[1]}
		}
		return &ValChar{0}
	case VT_STRING:
		if symbol.name != "" {
			return &ValString{symbol.name[1 : len(symbol.name)-1]}
		}
		return &ValString{""}
	}
	return nil
}

func PrintStack() {
	fmt.Println(string(debug.Stack()))
}
