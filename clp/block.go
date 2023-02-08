package clp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const IGNORE_CASE bool = false

const (
	BK_TYPE_GLOBAL int = iota
	BK_TYPE_FUNC
	BK_TYPE_IF
	BK_TYPE_WHILE
)

type Express struct {
	exp    string
	etype  int
	opera  int
	childs []*Express
}

func NewExpress() *Express {
	e := &Express{
		etype: VT_NONE,
		opera: OP_NONE,
	}
	return e
}

func (e *Express) IsEmpty() bool {
	return len(e.exp) == 0 && e.opera == OP_NONE
}

func (e *Express) IsOpera() bool {
	return e.opera > OP_NONE && e.opera < OP_FUNC && len(e.childs) == 0
}

type Block struct {
	btype   int
	parent  *Block
	pfunc   *Function
	symbTbl *SymbolTbl
	offset  uint32
	opers   []IOperator
}

func NewBlock(btype int, parent *Block, pfunc *Function, symbT *SymbolTbl) *Block {
	b := &Block{
		btype:   btype,
		parent:  parent,
		pfunc:   pfunc,
		symbTbl: symbT,
		offset:  GetOperaOffset(),
	}
	return b
}

func (b *Block) Parse(statements []string) error {
	for i := 0; i < len(statements); {
		next, err := b.parseControl(statements, i)
		if err != nil {
			return err
		}
		if next == i {
			s, err := b.defineVar(statements[i])
			if err != nil {
				return err
			}
			if len(s) > 0 {
				if exp, err := b.parseExpress(s); err == nil {
					b.createOpera(exp)
				} else {
					return err
				}
			}
			i++
		} else {
			i = next
		}
	}
	return nil
}

func (b *Block) Exec() {
	g_regPC = b.offset
	for {
		i := g_regPC - b.offset
		if i >= uint32(len(b.opers)) {
			break
		}
		b.opers[i].Run()
		if g_regPC < CTL_CONTINUE {
			g_regPC++
		}
	}
}

func (b *Block) parseControl(statements []string, cur int) (int, error) {
	if b.btype == BK_TYPE_GLOBAL {
		return cur, nil
	}
	if IsSemCtl(statements[cur], S_CTL_IF) {
		return b.parseCTLIf(statements, cur)
	}
	if IsSemCtl(statements[cur], S_CTL_WHILE) {
		return b.parseCTLWhile(statements, cur)
	}
	if IsSemCtl(statements[cur], S_CTL_RETURN) {
		return cur + 1, b.parseCTLReturn(statements[cur])
	}
	if IsSemCtl(statements[cur], S_CTL_BREAK) {
		return cur + 1, b.parseCTLBreak(statements[cur])
	}
	if IsSemCtl(statements[cur], S_CTL_CONTINUE) {
		return cur + 1, b.parseCTLContinue(statements[cur])
	}
	return cur, nil
}

func (b *Block) parseCTLIf(statements []string, cur int) (int, error) {
	ops := []int{}
	nxt := cur
	s := statements[cur]
	for {
		ctl := ""
		if IsSemCtl(s, S_CTL_IF) {
			ctl = S_CTL_IF
		} else if IsSemCtl(s, S_CTL_ELSEIF) {
			ctl = S_CTL_ELSEIF
		} else if IsSemCtl(s, S_CTL_ELSE) {
			ctl = S_CTL_ELSE
		} else {
			return cur, errors.New(s + "illegal!")
		}
		ci := uint32(0)
		if ctl != S_CTL_ELSE {
			exp, err := b.parseExpress(s[len(ctl) : len(s)-1])
			if err != nil {
				return cur, err
			}
			if exp.etype != VT_BOOL {
				return cur, errors.New("Condition type not bool!")
			}
			ci = b.exp2addr(exp)
		}

		var err error
		nxt, err = getBraceRange(statements, cur)
		if err != nil {
			return cur, err
		}

		symbT := NewSymbolTbl(b.symbTbl)
		block := NewBlock(BK_TYPE_IF, b, b.pfunc, symbT)
		err = block.Parse(statements[cur+1 : nxt])
		if err != nil {
			return cur, err
		}
		b.symbTbl.AddSymbolTbl(symbT)
		g_blocks = append(g_blocks, block)
		bi := uint32(len(g_blocks) - 1)

		if ctl == S_CTL_IF || ctl == S_CTL_ELSEIF {
			b.opers = append(b.opers, &OperCtlIf{ci, bi, b.offset + uint32(len(b.opers))})
			ops = append(ops, len(b.opers)-1)
		} else {
			oif := b.opers[len(b.opers)-1].(*OperCtlIf)
			b.opers[len(b.opers)-1] = &OperCtlIfElse{oif.a, oif.b, bi}
		}

		s = removeLeftBlank(statements[nxt][1:])
		if len(s) == 0 {
			break
		} else if ctl == S_CTL_ELSE {
			return cur, errors.New("else must be end!")
		}
		cur = nxt
	}
	for i := 0; i < len(ops)-1; i++ {
		b.opers[ops[i]].(*OperCtlIf).c = b.offset + uint32(len(b.opers)-1)
	}
	return nxt + 1, nil
}

func (b *Block) parseCTLWhile(statements []string, cur int) (int, error) {
	s := statements[cur]
	exp, err := b.parseExpress(s[len(S_CTL_WHILE) : len(s)-1])
	if err != nil {
		return cur, err
	}
	if exp.etype != VT_BOOL {
		return cur, errors.New("While condition type not bool!")
	}
	pi := b.offset + uint32(len(b.opers)-1)
	ci := b.exp2addr(exp)

	nxt, err := getBraceRange(statements, cur)
	if err != nil {
		return cur, err
	}

	symbT := NewSymbolTbl(b.symbTbl)
	block := NewBlock(BK_TYPE_WHILE, b, b.pfunc, symbT)
	err = block.Parse(statements[cur+1 : nxt])
	if err != nil {
		return cur, err
	}
	b.symbTbl.AddSymbolTbl(symbT)
	g_blocks = append(g_blocks, block)
	bi := uint32(len(g_blocks) - 1)

	b.opers = append(b.opers, &OperCtlWhile{ci, bi, pi})
	return nxt + 1, nil
}

func (b *Block) parseCTLReturn(s string) error {
	retype := b.pfunc.ftype
	if len(s) == len(S_CTL_RETURN) && retype != VT_NONE {
		return errors.New("Return cannot be empty!")
	}
	if len(s) > len(S_CTL_RETURN) {
		exp, err := b.parseExpress(s[len(S_CTL_RETURN):])
		if err != nil {
			return err
		}
		if retype != exp.etype && (retype != VT_FLOAT || exp.etype != VT_INT) {
			return errors.New("Return type not match!")
		}
		addr := b.exp2addr(exp)
		if retype != VT_NONE {
			b.opers = append(b.opers, &OperCtlReturn2{addr})
		}
	}
	if retype == VT_NONE {
		b.opers = append(b.opers, &OperCtlReturn{})
	}
	return nil
}

func (b *Block) parseCTLBreak(s string) error {
	if !b.inWhile() {
		return errors.New("Break not in while!")
	}
	b.opers = append(b.opers, &OperCtlBreak{})
	return nil
}

func (b *Block) parseCTLContinue(s string) error {
	if !b.inWhile() {
		return errors.New("Continue not in while!")
	}
	b.opers = append(b.opers, &OperCtlContinue{})
	return nil
}

func (b *Block) inWhile() bool {
	p := b
	for {
		if p == nil {
			break
		}
		if p.btype == BK_TYPE_WHILE {
			return true
		}
		p = b.parent
	}
	return false
}

func (b *Block) defineVar(s string) (string, error) {
	p := -1
	asn := false
	vars := []string{}
	for i := 0; i < len(s); i++ {
		if i+len(S_ASN) < len(s) && s[i:i+len(S_ASN)] == S_ASN {
			vars = append(vars, s[i:len(s)])
			asn = true
			break
		}
		if p != -1 {
			if s[i] == ' ' || s[i] == '	' {
				vars = append(vars, s[p:i])
				p = -1
			}
		} else {
			if s[i] != ' ' && s[i] != '	' {
				p = i
			}
		}
		if i == len(s)-1 && p != -1 {
			vars = append(vars, s[p:i+1])
		}
	}
	if len(vars) == 0 {
		return "", nil
	}

	var vtype int
	vtype, ok := g_mapType[vars[0]]
	if !ok {
		return s, nil
	} else if len(vars) == 1 {
		return "", errors.New("Cannot define empty variables!")
	} else if asn && len(vars) != 3 {
		return "", errors.New("The number of variables is illegal for assignment!")
	}

	for i := 1; i < len(vars); i++ {
		if asn && i == len(vars)-1 {
			break
		}
		v := vars[i]
		if IGNORE_CASE {
			v = strings.ToLower(v)
		}
		if vtype, _ := b.findSymbol(v, true); vtype != VT_NONE {
			return "", errors.New(fmt.Sprintf("%s redefined!", v))
		}
		b.symbTbl.AddSymbol(v, vtype)
	}
	if asn {
		return vars[1] + vars[2], nil
	}
	return "", nil
}

func (b *Block) findSymbol(s string, cur bool) (int, uint32) {
	p := b.symbTbl
	for {
		if p == nil || p.list == nil {
			break
		}
		addr := p.list.offset
		for _, v := range p.list.childs {
			if v.list != nil {
				addr += v.list.count
			} else if v.symbol.name == s {
				return v.symbol.vtype, addr
			} else {
				addr++
			}
		}
		if cur {
			break
		}
		p = p.list.parent
	}
	return VT_NONE, 0
}

func (b *Block) parseExpress(exp string) (*Express, error) {
	exp = removeBlank(exp)
	return b.parseExp(exp)
}

func (b *Block) lexAnalysis(exp string) ([]*Express, error) {
	if len(exp) == 0 {
		return nil, errors.New("Express cannot empty!")
	}
	last := 0
	exps := []*Express{}
	for i := 0; i <= len(exp); {
		s, e, exp, err := handleMatch(exp, i)
		if err != nil {
			return exps, err
		} else {
			i = e
		}
		op := ""
		for _, v := range g_strOpera {
			if i <= len(exp)-len(v) && exp[i:i+len(v)] == v {
				op = v
				break
			}
		}
		if op == "" && i < len(exp) {
			i++
		} else {
			var ep *Express
			if s < len(exp) && exp[s] == '(' {
				if s == last {
					ep, err = b.parseExp(exp[s+1 : e-1])
				} else {
					ep, err = b.parseFunc(exp[last:s], exp[s+1:e-1])
				}
				if ep == nil {
					return exps, err
				}
			} else if i > last {
				ep = NewExpress()
				ep.exp = exp[last:i]
			}
			if ep != nil {
				exps = append(exps, ep)
			}
			if op != "" {
				eop := NewExpress()
				eop.opera = g_mapOpera[op]
				exps = append(exps, eop)
			}
			i += len(op)
			last = i
			if last == len(exp) {
				break
			}
		}
	}
	return exps, nil
}

func (b *Block) syntaxAnalysis(exps []*Express) (*Express, error) {
	for i := 0; i < len(exps); i++ {
		if !exps[i].IsOpera() {
			continue
		}
		if exps[i].opera == OP_SUB {
			if i == 0 || exps[i-1].IsOpera() && (OP_UNARY(exps[i-1].opera) != 1 || !g_opCombine[exps[i-1].opera]) {
				exps[i].opera = OP_CONVERT(exps[i].opera)
			}
		} else if exps[i].opera == OP_ADS1 || exps[i].opera == OP_SBS1 {
			if i > 0 && !exps[i-1].IsOpera() {
				exps[i].opera = OP_CONVERT(exps[i].opera)
			}
		}
	}
	for {
		min := -1
		for i := 0; i < len(exps); i++ {
			if exps[i].IsOpera() { //取最高优先级运算符
				if min == -1 || g_opPrior[exps[i].opera] < g_opPrior[exps[min].opera] {
					min = i
				}
			}
		}
		var new *Express
		if min != -1 {
			new = NewExpress()
			new.opera = exps[min].opera
			un := OP_UNARY(exps[min].opera)
			if un == 1 {
				if g_opCombine[new.opera] { //左结合
					if min < 1 {
						return nil, errors.New("Syntax error!")
					}
					new.childs = append(new.childs, exps[min-1])
					exps = updateExps(exps, new, min-1, min)
				} else { //右结合
					if min >= len(exps)-1 {
						return nil, errors.New("Syntax error!")
					}
					new.childs = append(new.childs, exps[min+1])
					exps = updateExps(exps, new, min, min+1)
				}
			} else if un == 2 {
				if min < 1 || min >= len(exps)-1 {
					return nil, errors.New("Syntax error!")
				}
				new.childs = append(new.childs, exps[min-1])
				new.childs = append(new.childs, exps[min+1])
				exps = updateExps(exps, new, min-1, min+1)
			} else if un > 2 {
				pairs := g_mapPair[new.opera]
				if un-2 != len(pairs) {
					return nil, errors.New("Internal error!")
				}
				if len(pairs)*2+1+min >= len(exps) {
					return nil, errors.New("Syntax error!")
				}
				for i, pair := range pairs {
					if !exps[min+(i+1)*2].IsOpera() || exps[min+(i+1)*2].opera != pair {
						return nil, errors.New("Syntax error!")
					}
				}
				for i := 0; i < un; i++ {
					new.childs = append(new.childs, exps[min+i*2-1])
				}
				exps = updateExps(exps, new, min-1, min+un*2-3)
			} else {
				return nil, errors.New("Unknown unary!")
			}
		} else if len(exps) != 1 {
			return nil, errors.New("Syntax error!")
		} else {
			new = exps[0]
		}
		err := b.semanticAnalysis(new)
		if err != nil {
			return nil, err
		}
		if len(exps) < 2 {
			return new, nil
		}
	}
	return nil, errors.New("Unknown error!")
}

func (b *Block) semanticAnalysis(exp *Express) error {
	if exp.opera >= OP_FUNC {
		return nil
	}
	err := b.semantic(exp)
	if err != nil {
		return err
	}
	if len(exp.childs) > 0 {
		vts := []int{}
		for _, child := range exp.childs {
			if err = b.semantic(child); err != nil {
				return err
			}
			vts = append(vts, child.etype)
		}
		if !IsSemNorm(exp.opera, vts) {
			return errors.New("Type and operation not matched!")
		}
		exp.etype = GetSemType(exp.opera, vts)
	}
	return nil
}

func (b *Block) semantic(exp *Express) error {
	if exp == nil || exp.IsEmpty() {
		return errors.New("Express empty!")
	}
	if exp.IsOpera() {
		return errors.New("Cannot be operation!")
	}
	if len(exp.exp) > 0 {
		exp.etype = b.getExpType(exp.exp)
		if exp.etype == VT_NONE {
			return errors.New(fmt.Sprintf("%s undefined!", exp.exp))
		}
	}
	return nil
}

func (b *Block) parseExp(exp string) (*Express, error) {
	exps, err := b.lexAnalysis(exp)
	if err != nil {
		return nil, err
	}
	return b.syntaxAnalysis(exps)
}

func (b *Block) parseFunc(name string, paras string) (*Express, error) {
	index, sys := FindFuncIndex(name)
	if index == -1 {
		return nil, errors.New(fmt.Sprintf("%s undefined!", name))
	}
	if !sys && b.btype == BK_TYPE_GLOBAL {
		return nil, errors.New("Global scope cannot call method!")
	}
	ret := &Express{
		etype: g_funcs[index].ftype,
	}
	if sys {
		ret.opera = index + OP_FUNC
	} else {
		ret.opera = index + OP_MTHD
	}
	if len(paras) > 0 {
		exp, err := b.parseExp(paras)
		if exp == nil {
			return exp, err
		}
		ret.childs = b.spliteCommaExp(exp)
	}
	if g_funcs[index].ptype != nil { //函数的语义单独分析
		if len(ret.childs) != len(g_funcs[index].ptype) {
			return nil, errors.New("parameter num diff!")
		}
		for i, v := range ret.childs {
			if v.etype != g_funcs[index].ptype[i] {
				return nil, errors.New("parameter type diff!")
			}
		}
	}
	return ret, nil
}

func (b *Block) spliteCommaExp(exp *Express) []*Express {
	exps := []*Express{}
	p := exp
	for {
		if p.opera == OP_CMA {
			exps = append(exps, p.childs[1])
		} else {
			exps = append(exps, p)
			break
		}
		p = p.childs[0]
	}
	ret := make([]*Express, 0, len(exps))
	for i := len(exps) - 1; i >= 0; i-- {
		ret = append(ret, exps[i])
	}
	return ret
}

func (b *Block) createOpera(exp *Express) {
	if exp.opera == OP_NONE {
		return
	}
	addrs := make([]uint32, len(exp.childs)+1)
	for i, child := range exp.childs {
		addrs[i] = b.exp2addr(child)
	}
	b.symbTbl.AddSymbol("", exp.etype)
	addrs[len(exp.childs)] = b.getCurAddr()

	switch exp.opera {
	case OP_NOT:
		b.opers = append(b.opers, &OperNot{addrs[0], addrs[1]})
	case OP_NEG:
		b.opers = append(b.opers, &OperNeg{addrs[0], addrs[1]})
	case OP_LDI:
		b.opers = append(b.opers, &OperLDI{addrs[0], addrs[1]})
	case OP_ADS1:
		b.opers = append(b.opers, &OperADS1{addrs[0], addrs[1]})
	case OP_ADS2:
		b.opers = append(b.opers, &OperADS2{addrs[0], addrs[1]})
	case OP_SBS1:
		b.opers = append(b.opers, &OperSBS1{addrs[0], addrs[1]})
	case OP_SBS2:
		b.opers = append(b.opers, &OperSBS2{addrs[0], addrs[1]})
	case OP_MUL:
		b.opers = append(b.opers, &OperMul{addrs[0], addrs[1], addrs[2]})
	case OP_DIV:
		b.opers = append(b.opers, &OperDiv{addrs[0], addrs[1], addrs[2]})
	case OP_MOD:
		b.opers = append(b.opers, &OperMod{addrs[0], addrs[1], addrs[2]})
	case OP_ADD:
		b.opers = append(b.opers, &OperAdd{addrs[0], addrs[1], addrs[2]})
	case OP_SUB:
		b.opers = append(b.opers, &OperSub{addrs[0], addrs[1], addrs[2]})
	case OP_SFL:
		b.opers = append(b.opers, &OperSFL{addrs[0], addrs[1], addrs[2]})
	case OP_SFR:
		b.opers = append(b.opers, &OperSFR{addrs[0], addrs[1], addrs[2]})
	case OP_EQ:
		b.opers = append(b.opers, &OperEQ{addrs[0], addrs[1], addrs[2]})
	case OP_GE:
		b.opers = append(b.opers, &OperGE{addrs[0], addrs[1], addrs[2]})
	case OP_GT:
		b.opers = append(b.opers, &OperGT{addrs[0], addrs[1], addrs[2]})
	case OP_LE:
		b.opers = append(b.opers, &OperLE{addrs[0], addrs[1], addrs[2]})
	case OP_LT:
		b.opers = append(b.opers, &OperLT{addrs[0], addrs[1], addrs[2]})
	case OP_NE:
		b.opers = append(b.opers, &OperNE{addrs[0], addrs[1], addrs[2]})
	case OP_BND:
		b.opers = append(b.opers, &OperBND{addrs[0], addrs[1], addrs[2]})
	case OP_BXR:
		b.opers = append(b.opers, &OperBXR{addrs[0], addrs[1], addrs[2]})
	case OP_BOR:
		b.opers = append(b.opers, &OperBOR{addrs[0], addrs[1], addrs[2]})
	case OP_AND:
		b.opers = append(b.opers, &OperAnd{addrs[0], addrs[1], addrs[2]})
	case OP_OR:
		b.opers = append(b.opers, &OperOr{addrs[0], addrs[1], addrs[2]})
	case OP_ASN:
		b.opers = append(b.opers, &OperASN{addrs[0], addrs[1], addrs[2]})
	case OP_MLN:
		b.opers = append(b.opers, &OperMLN{addrs[0], addrs[1], addrs[2]})
	case OP_DVN:
		b.opers = append(b.opers, &OperDVN{addrs[0], addrs[1], addrs[2]})
	case OP_MDN:
		b.opers = append(b.opers, &OperMDN{addrs[0], addrs[1], addrs[2]})
	case OP_ADN:
		b.opers = append(b.opers, &OperADN{addrs[0], addrs[1], addrs[2]})
	case OP_SBN:
		b.opers = append(b.opers, &OperSBN{addrs[0], addrs[1], addrs[2]})
	case OP_SLN:
		b.opers = append(b.opers, &OperSLN{addrs[0], addrs[1], addrs[2]})
	case OP_SRN:
		b.opers = append(b.opers, &OperSRN{addrs[0], addrs[1], addrs[2]})
	case OP_BDN:
		b.opers = append(b.opers, &OperBDN{addrs[0], addrs[1], addrs[2]})
	case OP_BXN:
		b.opers = append(b.opers, &OperBXN{addrs[0], addrs[1], addrs[2]})
	case OP_BON:
		b.opers = append(b.opers, &OperBON{addrs[0], addrs[1], addrs[2]})
	case OP_CMA:
		b.opers = append(b.opers, &OperCMA{addrs[0], addrs[1], addrs[2]})
	case OP_IF:
		b.opers = append(b.opers, &OperIF{addrs[0], addrs[1], addrs[2], addrs[3]})
	default:
		if exp.opera >= OP_MTHD {
			b.opers = append(b.opers, &OperMthd{uint32(exp.opera - OP_MTHD), addrs})
		} else if exp.opera >= OP_FUNC {
			b.opers = append(b.opers, &OperFunc{uint32(exp.opera - OP_FUNC), addrs})
		}
	}
}

func (b *Block) exp2addr(exp *Express) uint32 {
	if exp.opera == OP_NONE {
		if index, err := b.getExpIndex(exp.exp); err != nil {
			updateConstSymbol(exp.exp)
			return g_rootBlock.getCurAddr()
		} else {
			return index
		}
	} else {
		b.createOpera(exp)
		return b.getCurAddr()
	}
}

func (b *Block) getExpType(s string) int {
	if s == S_TRUE || s == S_FALSE {
		return VT_BOOL
	}
	if len(s) > 2 {
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return VT_CHAR
		}
		if s[0] == '"' && s[len(s)-1] == '"' {
			return VT_STRING
		}
	}
	_, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return VT_INT
	}
	_, err = strconv.ParseFloat(s, 32)
	if err == nil {
		return VT_FLOAT
	}

	if IGNORE_CASE {
		s = strings.ToLower(s)
	}
	if vtype, _ := b.findSymbol(s, false); vtype != VT_NONE {
		return vtype
	}
	if b.btype != BK_TYPE_GLOBAL {
		if vtype, _ := g_rootBlock.findSymbol(s, true); vtype != VT_NONE {
			return vtype
		}
	}
	return VT_NONE
}

func (b *Block) getExpIndex(s string) (uint32, error) {
	if IGNORE_CASE {
		s = strings.ToLower(s)
	}
	if vtype, index := b.findSymbol(s, false); vtype != VT_NONE {
		return b.convertIndex(index), nil
	} else if b.btype != BK_TYPE_GLOBAL {
		if vtype, index := g_rootBlock.findSymbol(s, true); vtype != VT_NONE {
			return index, nil
		}
	}
	return 0, errors.New(s + " not found!")
}

func (b *Block) convertIndex(index uint32) uint32 {
	if b.btype == BK_TYPE_GLOBAL {
		return index
	} else {
		return index + STACK_OFFSET
	}
}

func (b *Block) getCurAddr() uint32 {
	if b.btype == BK_TYPE_GLOBAL {
		return b.symbTbl.list.offset + b.symbTbl.list.count - 1
	} else {
		return b.symbTbl.list.offset + b.symbTbl.list.count - 1 + STACK_OFFSET
	}
}

func updateConstSymbol(s string) {
	vtype := VT_NONE
	if s == S_TRUE || s == S_FALSE {
		vtype = VT_BOOL
	} else if s[0] == '\'' {
		vtype = VT_CHAR
	} else if s[0] == '"' {
		vtype = VT_STRING
	} else if _, err := strconv.ParseInt(s, 10, 32); err == nil {
		vtype = VT_INT
	} else {
		vtype = VT_FLOAT
	}
	g_rootBlock.symbTbl.AddSymbol(s, vtype)
}

func updateExps(exps []*Express, new *Express, start, end int) []*Express {
	ret := make([]*Express, 0, len(exps)-end+start)
	if start <= 0 && end >= len(exps)-1 {
		ret = append(ret, new)
	} else if start <= 0 {
		ret = append([]*Express{new}, exps[end+1:]...)
	} else if end >= len(exps)-1 {
		ret = append(exps[:start], new)
	} else {
		ret = append(exps[:start], append([]*Express{new}, exps[end+1:]...)...)
	}
	return ret
}

func handleMatch(exp string, cur int) (int, int, string, error) {
	if cur < len(exp) {
		var c byte = exp[cur]
		var c2 byte
		if (cur < 1 || exp[cur-1] != '\\') && (c == '\'' || c == '"') {
			c2 = c
		} else if c == '(' {
			c2 = ')'
		} else {
			return cur, cur, exp, nil
		}
		if c == c2 {
			for i := cur + 1; i < len(exp); i++ {
				if exp[i] == c {
					if i < 1 || exp[i-1] != '\\' {
						return cur, i + 1, exp, nil
					} else {
						exp = exp[:i-1] + exp[i:]
						i--
					}
				}
			}
		} else {
			n := 1
			var cq byte
			for i := cur + 1; i < len(exp); i++ {
				if (i < 1 || exp[i-1] != '\\') && (exp[i] == '\'' || exp[i] == '"') {
					if cq == 0 {
						cq = exp[i]
					} else if cq == exp[i] {
						cq = 0
					} else {
						return 0, 0, "", errors.New(fmt.Sprintf("%s %d %d %c not match quote!", exp, cur, cq, exp[i]))
					}
				}
				if cq == 0 {
					if exp[i] == c {
						n++
					} else if exp[i] == c2 {
						n--
					}
					if n == 0 {
						return cur, i + 1, exp, nil
					}
				}
			}
			if cq != 0 {
				return 0, 0, "", errors.New(fmt.Sprintf("%s %d %d %c not match quote!", exp, cur, cq))
			}
		}
		return 0, 0, "", errors.New(fmt.Sprintf("%s %d %d %c not match!", exp, cur, c))
	}
	return cur, cur, exp, nil
}

func removeBlank(s string) string {
	var c byte
	for i := 0; i < len(s); i++ {
		if (i < 1 || s[i-1] != '\\') && (s[i] == '\'' || s[i] == '"') {
			if c == 0 {
				c = s[i]
			} else if s[i] == c {
				c = 0
			}
		} else if c == 0 && (s[i] == ' ' || s[i] == '\t') {
			s = s[:i] + s[i+1:]
			i--
		}
	}
	return s
}
