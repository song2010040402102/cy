package clp

const (
	DEF_BOOL   string = "bool"
	DEF_INT    string = "int"
	DEF_FLOAT  string = "float"
	DEF_CHAR   string = "char"
	DEF_STRING string = "string"
)

var g_defType []string = []string{DEF_BOOL, DEF_INT, DEF_FLOAT, DEF_CHAR, DEF_STRING}

const (
	VT_NONE int = iota
	VT_BOOL
	VT_INT
	VT_FLOAT
	VT_CHAR
	VT_STRING
)

var g_mapType map[string]int = map[string]int{
	DEF_BOOL:   VT_BOOL,
	DEF_INT:    VT_INT,
	DEF_FLOAT:  VT_FLOAT,
	DEF_CHAR:   VT_CHAR,
	DEF_STRING: VT_STRING,
}

type IValue interface {
	GetType() int
	GetValue() interface{}
}

type ValNone struct {
}

func (v *ValNone) GetType() int {
	return VT_NONE
}

func (v *ValNone) GetValue() interface{} {
	return nil
}

type ValBool struct {
	val bool
}

func (v *ValBool) GetType() int {
	return VT_BOOL
}

func (v *ValBool) GetValue() interface{} {
	return v.val
}

type ValInt struct {
	val int64
}

func (v *ValInt) GetType() int {
	return VT_INT
}

func (v *ValInt) GetValue() interface{} {
	return v.val
}

type ValFloat struct {
	val float64
}

func (v *ValFloat) GetType() int {
	return VT_FLOAT
}

func (v *ValFloat) GetValue() interface{} {
	return v.val
}

type ValChar struct {
	val byte
}

func (v *ValChar) GetType() int {
	return VT_CHAR
}

func (v *ValChar) GetValue() interface{} {
	return v.val
}

type ValString struct {
	val string
}

func (v *ValString) GetType() int {
	return VT_STRING
}

func (v *ValString) GetValue() interface{} {
	return v.val
}
