package clp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ParseFile(filename string) error {
	lines, _ := readFile(filename)
	if len(lines) == 0 {
		return nil
	}
	gLines := []string{}
	fparas := [][]string{}
	flStart, flEnd := []int{}, []int{}
	for i := 0; i < len(lines); {
		next, paras, err := parseFuncDeclar(lines, i)
		if err != nil {
			return err
		}
		if next == i {
			gLines = append(gLines, lines[i])
			i++
		} else {
			fparas = append(fparas, paras)
			flStart = append(flStart, i+1)
			flEnd = append(flEnd, next-1)
			i = next
		}
	}
	g_rootBlock = NewBlock(BK_TYPE_GLOBAL, nil, nil, NewSymbolTbl(nil))
	err := g_rootBlock.Parse(gLines)
	if err != nil {
		return err
	}
	index, _ := FindFuncIndex(S_MAIN)
	if index == -1 {
		return errors.New("Function main not found!")
	}
	g_rootBlock.opers = append(g_rootBlock.opers, &OperMthd{uint32(index), nil})
	if len(flStart) > 0 {
		sysN := GetSysFuncNum()
		for i := 0; i < len(flStart); i++ {
			block := NewBlock(BK_TYPE_FUNC, nil, g_funcs[sysN+i], NewSymbolTbl(nil))
			err := block.Parse(append(fparas[i], lines[flStart[i]:flEnd[i]]...))
			if err != nil {
				return err
			}
			g_blocks = append(g_blocks, block)
			g_funcs[sysN+i].indexF = len(g_blocks) - 1
		}
	}
	return nil
}

func parseFuncDeclar(lines []string, cur int) (int, []string, error) {
	s := lines[cur]
	if !IsSemFunc(s) {
		return cur, nil, nil
	}
	f := &Function{}
	s = removeSideBlank(s[len(S_FUNC) : len(s)-1])
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '	' || s[i] == '(' {
			f.name = s[:i]
			s = removeLeftBlank(s[i:])
			break
		}
	}
	if f.name == "" {
		return cur, nil, errors.New("Function name cannot empty!")
	}
	for _, v := range g_defType {
		if len(s) > len(v) && s[len(s)-len(v):] == v {
			f.ftype = g_mapType[v]
			s = removeRightBlank(s[:len(s)-len(v)])
			break
		}
	}
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return cur, nil, errors.New(s + " lack brackets!")
	}
	var paras []string
	s = s[1 : len(s)-1]
	if len(s) > 0 {
		paras = strings.Split(s, ",")
		for _, para := range paras {
			para = removeLeftBlank(para)
			for _, v := range g_defType {
				if len(para) > len(v) && para[:len(v)] == v {
					f.ptype = append(f.ptype, g_mapType[v])
					break
				}
			}
		}
	}
	g_funcs = append(g_funcs, f)
	cur, err := getBraceRange(lines, cur)
	if err != nil {
		return cur, nil, err
	}
	if len(lines[cur]) != 1 {
		return cur, nil, errors.New(lines[cur][1:] + " should be removed!")
	}
	return cur + 1, paras, nil
}

func readFile(filename string) ([]string, error) {
	var lines []string
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(filename, "Open error!", err)
		return lines, err
	}
	defer f.Close()

	count := 0
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
		}
		if index := strings.Index(line, "//"); index >= 0 {
			line = line[:index]
		} else if index = strings.Index(line, "/*"); index >= 0 {
			count++
			line = line[:index]
		} else if index = strings.Index(line, "*/"); index >= 0 && count > 0 {
			count--
			if index < len(line)-2 {
				line = line[index+2:]
			} else {
				line = ""
			}
		}
		line = removeSideBlank(line)
		if count == 0 && len(line) > 0 {
			lines = append(lines, line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(filename, "Read error!", err)
			return []string{}, err
		}
	}
	return lines, nil
}

func removeLeftBlank(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '	' {
			s = s[i:]
			break
		}
	}
	return s
}

func removeRightBlank(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != ' ' && s[i] != '	' {
			s = s[:i+1]
			break
		}
	}
	return s
}

func removeSideBlank(s string) string {
	return removeRightBlank(removeLeftBlank(s))
}

func getBraceRange(lines []string, cur int) (int, error) {
	count := 1
	for cur++; cur < len(lines); cur++ {
		if lines[cur][0] == '}' {
			count--
			if count == 0 {
				break
			}
		}
		if lines[cur][len(lines[cur])-1] == '{' {
			count++
		}
	}
	if count != 0 {
		return cur, errors.New("} missing!")
	}
	return cur, nil
}
