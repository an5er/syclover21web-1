package model

import (
	"bufio"
	"fmt"
	"os"
	"watch01/core"
)

func Get_urllist() (lines []string) {
	file, err := os.Open(*core.Urllist)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
