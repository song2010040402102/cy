package main

import (
	"clp"
	"fmt"
	"os"
)

func main() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Println("filename cannot empty!")
		return
	}
	filename := os.Args[1]
	err := clp.ParseFile(filename)
	if err != nil {
		fmt.Println(filename, err)
		return
	}
	clp.Exec()
}
