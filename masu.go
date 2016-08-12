package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Masu int

const (
	MU      Masu = 0
	KOMADAI Masu = 100
)

func newMasu(suji int, dan int) Masu {
	return Masu(suji*10 + dan)
}

func (masu Masu) suji() int {
	return int(masu) / 10
}

func (masu Masu) dan() int {
	return int(masu) % 10
}

// test ok
func str2Masu(str string) Masu {
	// 7g -> 77
	int_x, _ := strconv.Atoi(str[0:1])
	char_y := str[1:2]
	int_y := strings.Index("0abcdefghi", char_y)
	return Masu(int_x*10 + int_y)
}

const DAN2STR string = "0abcdefghi"

func (masu Masu) masu2Str() string {
	// 77 -> 7g
	suji := masu.suji()
	dan := masu.dan()
	return fmt.Sprint(suji) + DAN2STR[dan:dan+1]
}
