package main

import (
	"strconv"
	"strings"
)

type Masu int

func newMasu(suji byte, dan byte) Masu {
	return Masu(int(suji)*10 + int(dan))
}

// 7g -> 77
func str2Masu(str string) Masu {
	int_x, _ := strconv.Atoi(str[0:1])
	char_y := str[1:2]
	int_y := strings.Index("0abcdefghi", char_y)
	return Masu(int_x*10 + int_y)
}
